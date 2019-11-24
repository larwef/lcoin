package main

import (
	"fmt"
	"github.com/larwef/lcoin/lcoin"
	"log"
)

func main() {

	var chn lcoin.Chain

	b1 := &lcoin.Block{Payload: "This is the first Block"}
	b2 := &lcoin.Block{Payload: "This is the second Block"}
	b3 := &lcoin.Block{Payload: "This is the third Block"}
	b4 := &lcoin.Block{Payload: "This is the fourth Block"}
	b5 := &lcoin.Block{Payload: "This is the fifth Block"}

	chn.Add(b1)
	chn.Add(b2)
	chn.Add(b3)
	chn.Add(b4)
	chn.Add(b5)

	chn.Print()

	if err := chn.Validate(); err != nil {
		fmt.Printf("Chain not valid: %v\n", chn.Validate())
	} else {
		fmt.Println("Chain valid")
	}

	if err := chn.ToJson("chain.json"); err != nil {
		log.Fatal(err)
	}

}
