package lcoin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
)

// User is ...
type User struct {
	Name       string
	PrivateKey *ecdsa.PrivateKey
}

// NewUser returns a new user initialized with a freshly generated key.
func NewUser(name string) *User {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Error generating key: %v\n", err)
	}

	return &User{
		Name:       name,
		PrivateKey: privateKey,
	}
}

// Address returns the address of the user.
func (u User) Address() [32]byte {
	// Rmember: u.PrivateKey.X == u.PrivateKey.PublicKey.X
	h1 := sha256.Sum256(elliptic.MarshalCompressed(elliptic.P256(), u.PrivateKey.X, u.PrivateKey.Y))
	return sha256.Sum256(h1[:])
}
