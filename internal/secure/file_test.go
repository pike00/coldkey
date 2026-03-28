package secure

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileCreatesWithRestrictedPerms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.key")
	data := []byte("secret-key-data")

	if err := WriteFile(path, data); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("permissions = %o, want 0600", perm)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("content = %q, want %q", got, data)
	}
}

func TestWriteFileRefusesOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.key")

	if err := WriteFile(path, []byte("first")); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	err := WriteFile(path, []byte("second"))
	if err == nil {
		t.Fatal("WriteFile should refuse to overwrite existing file")
	}
}

func TestWriteFileForceOverwrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.key")

	if err := WriteFile(path, []byte("first")); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := WriteFileForce(path, []byte("second")); err != nil {
		t.Fatalf("WriteFileForce: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "second" {
		t.Errorf("content = %q, want %q", got, "second")
	}
}

func TestShred(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "shred-me")

	if err := os.WriteFile(path, []byte("sensitive data here"), 0600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := Shred(path); err != nil {
		t.Fatalf("Shred: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("Shred should remove the file")
	}
}

func TestShredNonexistent(t *testing.T) {
	err := Shred("/nonexistent/path/to/file")
	if err == nil {
		t.Fatal("Shred should return error for nonexistent file")
	}
}
