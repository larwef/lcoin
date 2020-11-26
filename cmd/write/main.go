package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/larwef/lcoin/user"
	"log"
)

type Transaction struct {
	Message             string
	RecipientPubKeyHash [32]byte
}

func (t *Transaction) Hash() [32]byte {
	h1 := sha256.Sum256(append([]byte(t.Message), t.RecipientPubKeyHash[:]...))
	return sha256.Sum256(h1[:])
}

func (t *Transaction) Unlock(pubKey *rsa.PublicKey, signature []byte) error {
	var pubKeyPem bytes.Buffer
	if err := pem.Encode(&pubKeyPem, &pem.Block{
		Bytes: x509.MarshalPKCS1PublicKey(pubKey),
	}); err != nil {
		return err
	}

	pubKeyHash := sha256.Sum256(pubKeyPem.Bytes())
	if !bytes.Equal(t.RecipientPubKeyHash[:], pubKeyHash[:]) {
		return errors.New("not a matching public key")
	}

	t1Hash := t.Hash()
	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, t1Hash[:], signature)
}

func main() {
	user1 := user.LoadUser("user1")
	user2 := user.LoadUser("user2")
	//user3 := user.LoadUser("user3")

	var user1PublicKey bytes.Buffer
	err := pem.Encode(&user1PublicKey, &pem.Block{
		Bytes: x509.MarshalPKCS1PublicKey(&user1.PrivateKey.PublicKey),
	})
	checkErr(err)

	t1 := &Transaction{
		Message:             "someone gives 1 token",
		RecipientPubKeyHash: sha256.Sum256(user1PublicKey.Bytes()),
	}

	t1Hash := t1.Hash()

	// Unlock with user1
	sign1, err := rsa.SignPKCS1v15(rand.Reader, user1.PrivateKey, crypto.SHA256, t1Hash[:])
	checkErr(err)

	if err := t1.Unlock(&user1.PrivateKey.PublicKey, sign1); err != nil {
		fmt.Printf("Error unlocking t1 with user1: %v.\n", err)
	} else {
		fmt.Println("User1 sucessfully unlocked t1.")
	}

	// Unlock with user2 using own pubkey
	sign21, err := rsa.SignPKCS1v15(rand.Reader, user2.PrivateKey, crypto.SHA256, t1Hash[:])
	checkErr(err)

	if err := t1.Unlock(&user2.PrivateKey.PublicKey, sign21); err != nil {
		fmt.Printf("Error unlocking t1 with user2: %v.\n", err)
	} else {
		fmt.Println("User2 sucessfully unlocked t1.")
	}

	// Unlock with user2 using user1 pubkey
	sign22, err := rsa.SignPKCS1v15(rand.Reader, user2.PrivateKey, crypto.SHA256, t1Hash[:])
	checkErr(err)

	if err := t1.Unlock(&user1.PrivateKey.PublicKey, sign22); err != nil {
		fmt.Printf("Error unlocking t1 with user2: %v.\n", err)
	} else {
		fmt.Println("User2 sucessfully unlocked t1.")
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
