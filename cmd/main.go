package main

import (
	"encoding/hex"
	"fmt"

	"github.com/larwef/lcoin"
)

func main() {
	u1 := lcoin.NewUser("User1")
	u2 := lcoin.NewUser("User2")
	u3 := lcoin.NewUser("User3")

	address1 := u1.Address()
	fmt.Println(hex.EncodeToString(address1[:]))

	address2 := u2.Address()
	fmt.Println(hex.EncodeToString(address2[:]))

	address3 := u3.Address()
	fmt.Println(hex.EncodeToString(address3[:]))
}
