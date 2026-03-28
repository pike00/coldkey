package keygen

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	kp, err := Generate()
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if !strings.HasPrefix(kp.SecretKey, "AGE-SECRET-KEY-PQ-") {
		t.Errorf("SecretKey should start with AGE-SECRET-KEY-PQ-, got prefix %q", kp.SecretKey[:20])
	}
	if kp.PublicKey == "" {
		t.Error("PublicKey should not be empty")
	}
	if kp.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}

func TestFormatKeyFile(t *testing.T) {
	kp, err := Generate()
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	data := FormatKeyFile(kp)
	content := string(data)

	if !strings.Contains(content, "# created: ") {
		t.Error("FormatKeyFile should contain created timestamp")
	}
	if !strings.Contains(content, "# public key: ") {
		t.Error("FormatKeyFile should contain public key")
	}
	if !strings.Contains(content, "AGE-SECRET-KEY-PQ-") {
		t.Error("FormatKeyFile should contain the secret key")
	}
	if !strings.HasSuffix(content, "\n") {
		t.Error("FormatKeyFile should end with newline")
	}
}
