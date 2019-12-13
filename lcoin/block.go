package lcoin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/larwef/lcoin/merkle"
)

type Header struct {
	PrevHash   string `json:"prevHash"`
	MerkleRoot string `json:"merkleRoot"`
	Time       int64  `json:"time"`
	Nonce      int64  `json:"nonce"`
}

type Block struct {
	Header       Header   `json:"header"`
	Transactions []string `json:"transactions"`
}

func (b *Block) AddTransaction(transaction string) {
	b.Transactions = append(b.Transactions, transaction)
}

func (b *Block) MerkleRoot() string {
	var hashes [][32]byte
	for _, elem := range b.Transactions {
		hashes = append(hashes, sha256.Sum256([]byte(elem)))
	}

	_ = merkle.NewTree(hashes)

	return ""
}

func printHasehs(hashes [][]byte) {
	for _, elem := range hashes {
		fmt.Println(hex.EncodeToString(elem))
	}
}
