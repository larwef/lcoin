package user

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
)

type User struct {
	Name       string
	PrivateKey *rsa.PrivateKey
}

func NewUser(name string) *User {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Error generating key: %v\n", err)
	}

	return &User{
		Name:       name,
		PrivateKey: privateKey,
	}
}

func (u *User) PersistUser() {
	file, err := os.Create(u.Name + ".key")
	if err != nil {
		log.Fatalf("Error opening file %q key: %v\n", u.Name+".key", err)
	}
	defer file.Close()

	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(u.PrivateKey),
	}

	if err := pem.Encode(file, privBlock); err != nil {
		log.Fatalf("Error encodin private key: %v\n", err)
	}
}

func LoadUser(name string) *User {
	byts, err := ioutil.ReadFile(name + ".key")
	if err != nil {
		log.Fatalf("Error opening file %q key: %v\n", name+".key", err)
	}

	block, _ := pem.Decode(byts)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Error decoding private key: %v\n", err)
	}

	return &User{
		Name:       name,
		PrivateKey: privateKey,
	}
}
