package main

import (
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
		fmt.Printf("version: %s\n", c2pa.Version())
		return
	}

	sub := os.Args[1]

	readCmd := flag.NewFlagSet("read", flag.ExitOnError)
	readIn := readCmd.String("i", "", "input file (required)")
	readSettings := readCmd.String("s", "settings.toml", "settings file (TOML)")

	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
	signIn := signCmd.String("i", "", "input file (required)")
	signOut := signCmd.String("o", "", "output file (required)")
	signManifest := signCmd.String("m", "", "manifest file")
	certificates := signCmd.String("c", "", "certificate file (required)")
	key := signCmd.String("k", "", "key file (required)")
	signSettings := signCmd.String("s", "settings.toml", "settings file (TOML)")

	switch sub {
	case "read":
		readCmd.Parse(os.Args[2:])
		if *readIn == "" {
			log.Fatalf("read: -i is required")
		}
		builder, err := createContextBuilder(*readSettings)
		if err != nil {
			log.Fatalf("failed to create context builder: %v", err)
		}
		defer builder.Close()
		handleRead(builder, *readIn)
	case "sign":
		signCmd.Parse(os.Args[2:])
		if *signIn == "" {
			log.Fatalf("sign: -i is required")
		}
		if *signOut == "" {
			log.Fatalf("sign: -o is required")
		}
		builder, err := createContextBuilder(*signSettings)
		if err != nil {
			log.Fatalf("failed to create context builder: %v", err)
		}
		defer builder.Close()
		handleSign(builder, *signIn, *signOut, *signManifest, *certificates, *key)
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

// createContextBuilder builds a *c2pa.ContextBuilder using settings loaded
// from the given TOML file. If the file does not exist, defaults are used.
// The caller owns the returned builder and must Close() it (or consume it via
// Build()).
func createContextBuilder(settingsPath string) (*c2pa.ContextBuilder, error) {
	builder, err := c2pa.NewContextBuilder()
	if err != nil {
		return nil, err
	}
	resolver, err := c2pa.NewHttpResolver(&c2pa.DefaultHttpResolver{})
	if err != nil {
		builder.Close()
		return nil, err
	}
	defer resolver.Close()
	if err := builder.SetHttpResolver(resolver); err != nil {
		builder.Close()
		return nil, err
	}

	if settingsPath != "" {
		if content, err := os.ReadFile(settingsPath); err == nil {
			settings, err := c2pa.NewSettings()
			if err != nil {
				builder.Close()
				return nil, err
			}
			defer settings.Close()
			if err := settings.UpdateFromString(string(content), "toml"); err != nil {
				builder.Close()
				return nil, fmt.Errorf("failed to load settings %s: %w", settingsPath, err)
			}
			if err := builder.SetSettings(settings); err != nil {
				builder.Close()
				return nil, err
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			builder.Close()
			return nil, fmt.Errorf("failed to read settings %s: %w", settingsPath, err)
		}
	}

	return builder, nil
}

func handleRead(builder *c2pa.ContextBuilder, path string) {
	ctx, err := builder.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build context: %v", err)
		return
	}
	defer ctx.Close()

	r, err := c2pa.NewReader(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open reader: %v", err)
		return
	}
	defer r.Close()

	if err := r.WithFile(path); err != nil {
		fmt.Fprintf(os.Stderr, "failed to configure reader: %v", err)
		return
	}

	json := r.Json()
	fmt.Println(json)
}

func handleSign(builder *c2pa.ContextBuilder, input, output, manifest, certificates, key string) {
	certBytes, err := os.ReadFile(certificates)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read certificates: %v", err)
		return
	}
	keyBytes, err := os.ReadFile(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read private key: %v", err)
		return
	}

	info := c2pa.SignerInfo{
		Alg:          "ps256",
		SignCert:     string(certBytes),
		PrivateKey:   string(keyBytes),
		TimestampUrl: "http://timestamp.digicert.com",
	}
	if err := builder.SetSignerInfo(info); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set signer on context: %v", err)
		return
	}

	ctx, err := builder.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build context: %v", err)
		return
	}
	defer ctx.Close()

	_, err = os.Stat(manifest)
	if errors.Is(err, os.ErrNotExist) {
		manifest = DEFAULT_MANIFEST
	} else {
		content, err := os.ReadFile(manifest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read manifest file: %v", err)
		}
		manifest = string(content)
	}
	b, err := c2pa.NewBuilder(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create builder: %v", err)
		return

	}

	defer b.Close()
	b, err = b.WithDefinition(manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create builder: %v", err)
		return
	}

	bytes, err := b.SignWithContext(input, output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to sign file: %v", err)
	} else {
		fmt.Printf("Signed file %s, manifest bytes: %d\n", output, len(bytes))
	}
}
