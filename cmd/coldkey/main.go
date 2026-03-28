package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pike00/coldkey/internal/backup"
	"github.com/pike00/coldkey/internal/keygen"
	"github.com/pike00/coldkey/internal/keyfile"
	"github.com/pike00/coldkey/internal/prompt"
	"github.com/pike00/coldkey/internal/secure"
)

var version = "dev"

func main() {
	secure.LockMemory()

	if len(os.Args) < 2 {
		interactive()
		return
	}

	switch os.Args[1] {
	case "generate":
		cmdGenerate(os.Args[2:])
	case "backup":
		cmdBackup(os.Args[2:])
	case "version":
		fmt.Printf("coldkey %s\n", version)
	case "-v", "--version":
		fmt.Printf("coldkey %s\n", version)
	case "-h", "--help", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "coldkey: unknown command %q\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`coldkey — Post-quantum age key generation and paper backup

Usage:
  coldkey                          Interactive mode
  coldkey generate [flags]         Generate a new PQ age key
  coldkey backup [flags] KEYFILE   Create paper backup from existing key
  coldkey version                  Print version

Generate flags:
  -o PATH    Key file output path (default: stdout)
  -f         Overwrite existing file
  --no-backup  Skip HTML backup generation

Backup flags:
  -o PATH    HTML output path (default: <keyfile>-backup.html)

Docker (via just):
  just docker          Build the image
  just docker-run      Interactive mode with security hardening
  just docker-backup KEYFILE   Backup an existing key
`)
}

func cmdGenerate(args []string) {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	output := fs.String("o", "", "Key file output path (default: stdout)")
	force := fs.Bool("f", false, "Overwrite existing file")
	noBackup := fs.Bool("no-backup", false, "Skip HTML backup generation")
	fs.Parse(args)

	fmt.Fprintln(os.Stderr, "coldkey: generating post-quantum age key (ML-KEM-768 + X25519)...")

	kp, err := keygen.Generate()
	if err != nil {
		fatalf("key generation failed: %v", err)
	}
	keyData := keygen.FormatKeyFile(kp)
	defer secure.Zero(keyData)

	if *output == "" {
		// Write key to stdout
		os.Stdout.Write(keyData)
		if !*noBackup {
			fmt.Fprintln(os.Stderr, "coldkey: hint: use -o PATH to also generate an HTML backup")
		}
		return
	}

	// Ensure parent directory exists
	if dir := filepath.Dir(*output); dir != "." {
		if err := os.MkdirAll(dir, 0700); err != nil {
			fatalf("creating output directory: %v", err)
		}
	}

	writeFunc := secure.WriteFile
	if *force {
		writeFunc = secure.WriteFileForce
	}
	if err := writeFunc(*output, keyData); err != nil {
		fatalf("writing key file: %v (use -f to overwrite)", err)
	}
	fmt.Fprintf(os.Stderr, "coldkey: key written to %s\n", *output)
	fmt.Fprintf(os.Stderr, "coldkey: public key: %s\n", kp.PublicKey)

	if *noBackup {
		return
	}

	// Generate backup
	ki, err := keyfile.Read(*output)
	if err != nil {
		fatalf("reading back key file: %v", err)
	}
	backupPath := *output + "-backup.html"
	if err := backup.WriteHTML(ki, backupPath, version); err != nil {
		fatalf("generating backup: %v", err)
	}
	fmt.Fprintf(os.Stderr, "coldkey: backup written to %s\n", backupPath)
	fmt.Fprintln(os.Stderr, "coldkey: open in a browser to print, then shred the HTML file")
}

func cmdBackup(args []string) {
	fs := flag.NewFlagSet("backup", flag.ExitOnError)
	output := fs.String("o", "", "HTML output path")
	fs.Parse(args)

	if fs.NArg() < 1 {
		fatalf("usage: coldkey backup [-o PATH] KEYFILE")
	}
	keyPath := fs.Arg(0)

	ki, err := keyfile.Read(keyPath)
	if err != nil {
		fatalf("reading key file: %v", err)
	}

	backupPath := *output
	if backupPath == "" {
		backupPath = keyPath + "-backup.html"
	}

	if err := backup.WriteHTML(ki, backupPath, version); err != nil {
		fatalf("generating backup: %v", err)
	}
	fmt.Fprintf(os.Stderr, "coldkey: backup written to %s\n", backupPath)
	fmt.Fprintln(os.Stderr, "coldkey: open in a browser to print, then shred the HTML file")
}

func interactive() {
	fmt.Println("coldkey — Post-Quantum Age Key Tool")
	fmt.Printf("version %s\n", version)

	choice, err := prompt.Choice("What would you like to do?", []string{
		"Generate a new post-quantum key",
		"Create backup from existing key file",
	})
	if err != nil {
		fatalf("reading input: %v", err)
	}

	switch choice {
	case 0:
		interactiveGenerate()
	case 1:
		interactiveBackup()
	}
}

func interactiveGenerate() {
	output, err := prompt.String("Key file output path", "/out/keys.txt")
	if err != nil {
		fatalf("reading input: %v", err)
	}

	// Check for existing file
	writeFunc := secure.WriteFile
	if _, err := os.Stat(output); err == nil {
		overwrite, err := prompt.Confirm(fmt.Sprintf("File %s already exists. Overwrite?", output))
		if err != nil {
			fatalf("reading input: %v", err)
		}
		if !overwrite {
			fmt.Println("Aborted.")
			return
		}
		writeFunc = secure.WriteFileForce
	}

	fmt.Println()
	fmt.Println("Generating post-quantum age key (ML-KEM-768 + X25519)...")

	kp, err := keygen.Generate()
	if err != nil {
		fatalf("key generation failed: %v", err)
	}
	keyData := keygen.FormatKeyFile(kp)
	defer secure.Zero(keyData)

	// Ensure parent directory exists
	if dir := filepath.Dir(output); dir != "." {
		if err := os.MkdirAll(dir, 0700); err != nil {
			fatalf("creating output directory: %v", err)
		}
	}

	if err := writeFunc(output, keyData); err != nil {
		fatalf("writing key file: %v", err)
	}
	fmt.Printf("Key written to %s\n", output)
	fmt.Printf("Public key: %s\n", kp.PublicKey)

	// Backup
	fmt.Println()
	doBackup, err := prompt.Confirm("Generate printable HTML backup?")
	if err != nil {
		fatalf("reading input: %v", err)
	}
	if !doBackup {
		return
	}

	backupPath := output + "-backup.html"
	backupPath, err = prompt.String("Backup output path", backupPath)
	if err != nil {
		fatalf("reading input: %v", err)
	}

	ki, err := keyfile.Read(output)
	if err != nil {
		fatalf("reading key file: %v", err)
	}

	if err := backup.WriteHTML(ki, backupPath, version); err != nil {
		fatalf("generating backup: %v", err)
	}
	fmt.Printf("Backup written to %s\n", backupPath)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Open the HTML file in a browser and print it")
	fmt.Println("  2. Store the printout in a secure location")
	fmt.Println("  3. Shred the HTML file: shred -u " + backupPath)
}

func interactiveBackup() {
	keyPath, err := prompt.String("Path to existing key file", "")
	if err != nil {
		fatalf("reading input: %v", err)
	}
	if keyPath == "" {
		fatalf("key file path is required")
	}

	ki, err := keyfile.Read(keyPath)
	if err != nil {
		fatalf("reading key file: %v", err)
	}

	defaultBackup := keyPath + "-backup.html"
	backupPath, err := prompt.String("Backup output path", defaultBackup)
	if err != nil {
		fatalf("reading input: %v", err)
	}

	if err := backup.WriteHTML(ki, backupPath, version); err != nil {
		fatalf("generating backup: %v", err)
	}
	fmt.Printf("Backup written to %s\n", backupPath)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Open the HTML file in a browser and print it")
	fmt.Println("  2. Store the printout in a secure location")
	fmt.Println("  3. Shred the HTML file: shred -u " + backupPath)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "coldkey: error: "+format+"\n", args...)
	os.Exit(1)
}
