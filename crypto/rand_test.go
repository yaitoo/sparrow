package crypto

import (
	"testing"
)

func TestRandBytes(t *testing.T) {
	bytes, err := RandBytes(4)
	if err != nil {
		// Serve an appropriately vague error to the
		// user, but log the details internally.
		t.Error(err)
	}

	if len(bytes) != 4 {
		t.Error("Error on length")
	}
}

func TestRandString(t *testing.T) {
	_, err := RandString(4)
	if err != nil {
		// Serve an appropriately vague error to the
		// user, but log the details internally.
		t.Error(err)
	}

	// if len(token) != 4 {
	// 	t.Error("Error on length")
	// }
}

func TestRandInt(t *testing.T) {
	n, err := RandInt(4)
	if err != nil {
		// Serve an appropriately vague error to the
		// user, but log the details internally.
		t.Error(err)
	}

	t.Log(n)

	// if len(bytes) != 4 {
	// 	t.Error("Error on length")
	// }
}
