package lcoin

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"log"
)

// Message is ...
type Message struct {
	SenderAddress   [32]byte
	ReceiverAddress [32]byte
	Payload         string
}

// Hash is ...
func (m Message) Hash() [32]byte {
	h := sha256.New()
	h.Write(m.SenderAddress[:])
	h.Write(m.ReceiverAddress[:])
	h.Write([]byte(m.Payload))

	return sha256.Sum256(h.Sum(nil))
}

// Claim is ...
func (m Message) Claim(signature []byte, pubkey *ecdsa.PublicKey) bool {
	pubkeyHash := sha256.Sum256(elliptic.MarshalCompressed(elliptic.P256(), pubkey.X, pubkey.Y))
	address := sha256.Sum256(pubkeyHash[:])
	if !bytes.Equal(address[:], m.ReceiverAddress[:]) {
		return false
	}

	messageHash := m.Hash()
	return ecdsa.VerifyASN1(pubkey, messageHash[:], signature)
}

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

// NewUserFromKey returns a new user with an existing keyin on SEC 1, ASN.1 DER form.
func NewUserFromKey(name string, der []byte) (User, error) {
	privKey, err := x509.ParseECPrivateKey(der)
	if err != nil {
		return User{}, err
	}

	return User{
		Name:       name,
		PrivateKey: privKey,
	}, nil
}

// Address returns the address of the user.
func (u User) Address() [32]byte {
	// Rmember: u.PrivateKey.X == u.PrivateKey.PublicKey.X
	h1 := sha256.Sum256(elliptic.MarshalCompressed(elliptic.P256(), u.PrivateKey.X, u.PrivateKey.Y))
	return sha256.Sum256(h1[:])
}
