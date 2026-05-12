package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/duggaraju/c2pa-go/c2pa"
)

const DEFAULT_MANIFEST = "{}"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	// Support a global version flag as first argument: -v or --version
	if os.Args[1] == "-v" || os.Args[1] == "--version" {
		fmt.Printf("version: %s\n", c2pa.CpaVersion())
		return
	}

	sub := os.Args[1]

	readCmd := flag.NewFlagSet("read", flag.ExitOnError)
	readIn := readCmd.String("i", "", "input file (required)")
	readSettings := readCmd.String("s", "settings.toml", "settings file (TOML)")

	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
	signIn := signCmd.String("i", "", "input file (required)")
	signOut := signCmd.String("o", "", "output file (required)")
	signManifest := signCmd.String("m", "", "manifest file (required)")
	certificates := signCmd.String("c", "", "certificate file (required)")
	key := signCmd.String("k", "", "key file (required)")
	signSettings := signCmd.String("s", "settings.toml", "settings file (TOML)")

	switch sub {
	case "read":
		readCmd.Parse(os.Args[2:])
		if *readIn == "" {
			log.Fatalf("read: -i is required")
		}
		ctx, err := createContext(*readSettings)
		if err != nil {
			log.Fatalf("failed to create context: %v", err)
		}
		defer ctx.Close()
		handleRead(ctx, *readIn)
	case "sign":
		signCmd.Parse(os.Args[2:])
		if *signIn == "" {
			log.Fatalf("sign: -i is required")
		}
		if *signOut == "" {
			log.Fatalf("sign: -o is required")
		}
		ctx, err := createContext(*signSettings)
		if err != nil {
			log.Fatalf("failed to create context: %v", err)
		}
		defer ctx.Close()
		handleSign(ctx, *signIn, *signOut, *signManifest, *certificates, *key)
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s <command> [options]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  read  -i <file>           Read and print reader JSON")
	fmt.Fprintln(os.Stderr, "  sign  -i <file> -o <file> [-m <manifest>]   Sign the input file (placeholder)")
}

// createContext builds a *c2pa.Context using settings loaded from the given
// TOML file. If the file does not exist, defaults are used.
func createContext(settingsPath string) (*c2pa.Context, error) {
	builder, err := c2pa.NewContextBuilder()
	if err != nil {
		return nil, err
	}
	defer builder.Close()
	resolver, err := c2pa.NewHttpResolver(&c2pa.DefaultHttpResolver{})
	if err != nil {
		return nil, err
	}
	defer resolver.Close()
	if err := builder.SetHttpResolver(resolver); err != nil {
		return nil, err
	}

	if settingsPath != "" {
		if content, err := os.ReadFile(settingsPath); err == nil {
			settings, err := c2pa.NewSettings()
			if err != nil {
				return nil, err
			}
			defer settings.Close()
			if err := settings.UpdateFromString(string(content), "toml"); err != nil {
				return nil, fmt.Errorf("failed to load settings %s: %w", settingsPath, err)
			}
			if err := builder.SetSettings(settings); err != nil {
				return nil, err
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("failed to read settings %s: %w", settingsPath, err)
		}
	}

	return builder.Build()
}

func handleRead(ctx *c2pa.Context, path string) {
	r, err := c2pa.ReaderFromFile(ctx, path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open reader: %v", err)
		return
	}
	defer r.Close()

	json := r.Json()
	fmt.Println(json)
}

func handleSign(ctx *c2pa.Context, input, output, manifest, certificates, key string) {
	_, err := os.Stat(manifest)
	if errors.Is(err, os.ErrNotExist) {
		manifest = DEFAULT_MANIFEST
	} else {
		content, err := os.ReadFile(manifest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read manifest file: %v", err)
		}
		manifest = string(content)
	}
	builder, err := c2pa.BuilderFromJson(ctx, manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create builder: %v", err)
	}
	defer builder.Close()

	signer, err := CreateTestSigner(certificates, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create signer: %v", err)
		return
	}

	bytes, err := builder.Sign(input, output, signer)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to sign file: %v", err)
	} else {
		fmt.Printf("Signed file %s, manifest bytes: %d\n", output, len(bytes))
	}
}

type TestSigner struct {
	certificates string
	key          *rsa.PrivateKey
}

func CreateTestSigner(cert string, key string) (*TestSigner, error) {
	certificates, err := os.ReadFile(cert)
	if err != nil {
		return nil, err
	}

	keyBytes, err := os.ReadFile(key)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	rsaKey, err := parseRSAPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &TestSigner{
		certificates: string(certificates),
		key:          rsaKey,
	}, nil
}

// parseRSAPrivateKey accepts a PKCS#1 or PKCS#8 encoded RSA key. PKCS#8 keys
// using the RSASSA-PSS OID (1.2.840.113549.1.1.10) are also accepted; the
// inner PKCS#1 RSAPrivateKey is parsed directly.
func parseRSAPrivateKey(der []byte) (*rsa.PrivateKey, error) {
	if k, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return k, nil
	}
	if k, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		if rsaKey, ok := k.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
		return nil, fmt.Errorf("key is not an RSA private key")
	}
	// Fallback: unwrap PKCS#8 manually for keys with RSASSA-PSS OID.
	var pkcs8 struct {
		Version    int
		Algo       asn1.RawValue
		PrivateKey []byte
	}
	if _, err := asn1.Unmarshal(der, &pkcs8); err != nil {
		return nil, fmt.Errorf("failed to parse PKCS#8: %w", err)
	}
	return x509.ParsePKCS1PrivateKey(pkcs8.PrivateKey)
}
func (s *TestSigner) Sign(input []byte, output []byte) (int, error) {
	hash := sha256.Sum256(input)
	sig, err := rsa.SignPSS(rand.Reader, s.key, crypto.SHA256, hash[:], &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
	})
	if err != nil {
		return -1, err
	}
	if len(sig) > len(output) {
		return -1, fmt.Errorf("output buffer too small: need %d, have %d", len(sig), len(output))
	}
	return copy(output, sig), nil
}

func (s *TestSigner) Alg() c2pa.SigningAlg {
	return c2pa.SigningAlgPs256
}

func (s *TestSigner) TimeStampUrl() string {
	return "http://timestamp.digicert.com"
}

func (s *TestSigner) Certificates() string {
	return s.certificates
}
