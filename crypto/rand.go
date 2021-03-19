//Generating Secure Random Numbers Using crypto/rand

package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

// RandBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func RandBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// RandString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func RandString(s int) (string, error) {
	b, err := RandBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// RandInt returns an int64
// securely generated random int64.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func RandInt(max int64) (int64, error) {
	num, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return -1, err
	}

	return num.Int64(), nil
}
