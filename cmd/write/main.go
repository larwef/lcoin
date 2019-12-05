package main

import (
	"fmt"
	"github.com/larwef/lcoin/lcoin"
)

func main() {

	var transactions []string
	for i := 0; i < 16; i++ {
		transactions = append(transactions, fmt.Sprintf("Transaction %d", i))
	}

	b1 := &lcoin.Block{
		Transactions: transactions,
	}

	b1.MerkleRoot()
}
