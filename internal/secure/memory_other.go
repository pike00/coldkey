//go:build !linux

package secure

import (
	"fmt"
	"os"
)

// LockMemory is a no-op on non-Linux platforms where mlockall is unavailable.
func LockMemory() {
	fmt.Fprintln(os.Stderr, "coldkey: warning: memory locking not supported on this platform")
}
