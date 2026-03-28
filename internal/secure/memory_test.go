package secure

import "testing"

func TestZero(t *testing.T) {
	b := []byte{0xff, 0xfe, 0xfd, 0xfc, 0xfb}
	Zero(b)
	for i, v := range b {
		if v != 0 {
			t.Errorf("Zero: byte %d = %#x, want 0", i, v)
		}
	}
}

func TestZeroEmpty(t *testing.T) {
	b := []byte{}
	Zero(b) // should not panic
}

func TestZeroNil(t *testing.T) {
	Zero(nil) // should not panic
}
