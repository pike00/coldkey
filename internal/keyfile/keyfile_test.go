package keyfile

import (
	"os"
	"path/filepath"
	"testing"
)

const validPQKey = `# created: 2026-01-15T10:30:00Z
# public key: age1pq1example-public-key
AGE-SECRET-KEY-PQ-1EXAMPLEKEYDATA
`

const classicKey = `# created: 2026-01-15T10:30:00Z
# public key: age1examplepublickey
AGE-SECRET-KEY-1EXAMPLEKEYDATA
`

func TestParseValidPQKey(t *testing.T) {
	ki, err := Parse([]byte(validPQKey))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if ki.SecretKey != "AGE-SECRET-KEY-PQ-1EXAMPLEKEYDATA" {
		t.Errorf("SecretKey = %q, want AGE-SECRET-KEY-PQ-1EXAMPLEKEYDATA", ki.SecretKey)
	}
	if ki.PublicKey != "age1pq1example-public-key" {
		t.Errorf("PublicKey = %q, want age1pq1example-public-key", ki.PublicKey)
	}
	if ki.CreatedAt.IsZero() {
		t.Error("CreatedAt should be parsed")
	}
	if ki.SHA256 == "" {
		t.Error("SHA256 should be set")
	}
}

func TestParseRejectsClassicKey(t *testing.T) {
	_, err := Parse([]byte(classicKey))
	if err == nil {
		t.Fatal("Parse should reject classic X25519 keys")
	}
}

func TestParseRejectsEmptyFile(t *testing.T) {
	_, err := Parse([]byte(""))
	if err == nil {
		t.Fatal("Parse should reject empty file")
	}
}

func TestParseRejectsNoKey(t *testing.T) {
	_, err := Parse([]byte("# just a comment\n"))
	if err == nil {
		t.Fatal("Parse should reject file with no key")
	}
}

func TestReadFileSizeLimit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "big.key")

	// Create a file larger than maxKeyFileSize
	data := make([]byte, maxKeyFileSize+1)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := Read(path)
	if err == nil {
		t.Fatal("Read should reject files larger than maxKeyFileSize")
	}
}

func TestReadNonexistent(t *testing.T) {
	_, err := Read("/nonexistent/key.txt")
	if err == nil {
		t.Fatal("Read should return error for nonexistent file")
	}
}
