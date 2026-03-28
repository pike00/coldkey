package keygen

import (
	"fmt"
	"time"

	"filippo.io/age"
)

// KeyPair holds a generated post-quantum age identity and its formatted output.
type KeyPair struct {
	SecretKey string
	PublicKey string
	CreatedAt time.Time
}

// Generate creates a new ML-KEM-768 + X25519 hybrid age key pair.
func Generate() (*KeyPair, error) {
	identity, err := age.GenerateHybridIdentity()
	if err != nil {
		return nil, fmt.Errorf("generating PQ identity: %w", err)
	}
	return &KeyPair{
		SecretKey: identity.String(),
		PublicKey: identity.Recipient().String(),
		CreatedAt: time.Now().UTC(),
	}, nil
}

// FormatKeyFile produces the standard age key file content:
//
//	# created: <RFC3339>
//	# public key: <recipient>
//	AGE-SECRET-KEY-PQ-1...
func FormatKeyFile(kp *KeyPair) []byte {
	return []byte(fmt.Sprintf(
		"# created: %s\n# public key: %s\n%s\n",
		kp.CreatedAt.Format(time.RFC3339),
		kp.PublicKey,
		kp.SecretKey,
	))
}
