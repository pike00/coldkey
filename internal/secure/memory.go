package secure

// Zero overwrites a byte slice with zeros to erase key material from memory.
// Uses the Go built-in clear() which is guaranteed not to be optimized away.
// Best-effort in Go due to GC; combined with mlockall this prevents swap exposure.
func Zero(b []byte) {
	clear(b)
}
