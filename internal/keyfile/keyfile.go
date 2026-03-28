package keyfile

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"time"
)

// KeyInfo holds parsed metadata from an age key file.
type KeyInfo struct {
	SecretKey  string
	PublicKey  string
	CreatedAt  time.Time
	RawContent []byte
	SHA256     string
	FilePath   string
	FileSize   int64
}

// maxKeyFileSize is the sanity limit for age key files (64KB).
// Real PQ age keys are ~2-3KB; anything larger is not a valid key file.
const maxKeyFileSize = 64 * 1024

// Read reads and validates an age PQ key file from disk.
func Read(path string) (*KeyInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("reading key file: %w", err)
	}
	if info.Size() > maxKeyFileSize {
		return nil, fmt.Errorf("file too large to be a valid age key file (%d bytes)", info.Size())
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading key file: %w", err)
	}
	ki, err := Parse(data)
	if err != nil {
		return nil, err
	}
	ki.FilePath = path
	ki.FileSize = info.Size()
	return ki, nil
}

// Parse validates and extracts metadata from age key file content.
func Parse(data []byte) (*KeyInfo, error) {
	ki := &KeyInfo{
		RawContent: data,
		SHA256:     fmt.Sprintf("%x", sha256.Sum256(data)),
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "# created: "):
			t, err := time.Parse(time.RFC3339, strings.TrimPrefix(line, "# created: "))
			if err == nil {
				ki.CreatedAt = t
			}
		case strings.HasPrefix(line, "# public key: "):
			ki.PublicKey = strings.TrimPrefix(line, "# public key: ")
		case strings.HasPrefix(line, "AGE-SECRET-KEY-"):
			ki.SecretKey = line
		}
	}

	if ki.SecretKey == "" {
		return nil, fmt.Errorf("key file does not contain an AGE-SECRET-KEY line")
	}
	if !strings.HasPrefix(ki.SecretKey, "AGE-SECRET-KEY-PQ-") {
		return nil, fmt.Errorf("key is not post-quantum (expected AGE-SECRET-KEY-PQ- prefix, got classic X25519)")
	}
	return ki, nil
}
