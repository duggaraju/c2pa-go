package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/duggaraju/c2pa-go/lib"
)

const DEFAULT_MANIFEST = "{}"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	// Support a global version flag as first argument: -v or --version
	if os.Args[1] == "-v" || os.Args[1] == "--version" {
		fmt.Printf("version: %s\n", lib.CpaVersion())
		return
	}

	sub := os.Args[1]

	readCmd := flag.NewFlagSet("read", flag.ExitOnError)
	readIn := readCmd.String("i", "", "input file (required)")

	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
	signIn := signCmd.String("i", "", "input file (required)")
	signOut := signCmd.String("o", "", "output file (required)")
	signManifest := signCmd.String("m", "", "manifest file (required)")
	certificates := signCmd.String("c", "", "certificate file (required)")
	key := signCmd.String("k", "", "key file (required)")

	switch sub {
	case "read":
		readCmd.Parse(os.Args[2:])
		if *readIn == "" {
			log.Fatalf("read: -i is required")
		}
		handleRead(*readIn)
	case "sign":
		signCmd.Parse(os.Args[2:])
		if *signIn == "" {
			log.Fatalf("sign: -i is required")
		}
		if *signOut == "" {
			log.Fatalf("sign: -o is required")
		}

		handleSign(*signIn, *signOut, *signManifest, *certificates, *key)
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

func handleRead(path string) {
	r, err := lib.ReaderFromFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open reader: %v", err)
	}
	defer r.Close()

	json := r.Json()
	fmt.Println(json)
}

func handleSign(input, output, manifest, certificates, key string) {
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

	builder, err := lib.BuilderFromJson(manifest)
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

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not an RSA private key")
	}

	return &TestSigner{
		certificates: string(certificates),
		key:          rsaKey,
	}, nil
}
func (s *TestSigner) Sign(input []byte, output []byte) (int, error) {

	// create sha256 hash of input
	hash := sha256.Sum256(input)
	// sign using rsa algorithm
	_, err := rsa.SignPKCS1v15(rand.Reader, s.key, crypto.SHA256, hash[:])
	if err != nil {
		return 0, err
	}
	return len(input), nil
}

func (s *TestSigner) Alg() lib.SigningAlg {
	return lib.SigningAlgPs256
}

func (s *TestSigner) TimeStampUrl() string {
	return ""
}

func (s *TestSigner) Certificates() string {
	return s.certificates
}
