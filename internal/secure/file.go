package secure

import (
	"crypto/rand"
	"os"
)

// WriteFile creates a new file at path with mode 0600 and calls fsync to ensure durability.
// It refuses to overwrite an existing file; use WriteFileForce for that.
func WriteFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	return syncAndClose(f, data)
}

// WriteFileForce writes data to path with mode 0600, overwriting any existing file.
func WriteFileForce(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	return syncAndClose(f, data)
}

func syncAndClose(f *os.File, data []byte) error {
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// Shred overwrites a file with 3 passes of random data, then removes it.
func Shred(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	size := info.Size()

	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	buf := make([]byte, size)
	defer Zero(buf)
	for pass := 0; pass < 3; pass++ {
		if _, err := rand.Read(buf); err != nil {
			f.Close()
			return err
		}
		if _, err := f.WriteAt(buf, 0); err != nil {
			f.Close()
			return err
		}
		if err := f.Sync(); err != nil {
			f.Close()
			return err
		}
	}
	f.Close()
	return os.Remove(path)
}
