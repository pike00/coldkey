package secure

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// LockMemory calls mlockall to prevent key material from being swapped to disk.
// Prints a warning to stderr if insufficient privileges (no CAP_IPC_LOCK).
func LockMemory() {
	if err := unix.Mlockall(unix.MCL_CURRENT | unix.MCL_FUTURE); err != nil {
		fmt.Fprintf(os.Stderr, "coldkey: warning: could not lock memory (keys may be swapped to disk): %v\n", err)
	}
}

// Zero overwrites a byte slice with zeros to erase key material from memory.
// Uses the Go built-in clear() which is guaranteed not to be optimized away.
// Best-effort in Go due to GC; combined with mlockall this prevents swap exposure.
func Zero(b []byte) {
	clear(b)
}
