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
	h := sha256.New()
	var hashes [][]byte
	for _, elem := range b.Transactions {
		h.Write([]byte(elem))
		hashes = append(hashes, h.Sum(nil))
		h.Reset()
	}

	_ = merkle.NewTree(hashes)

	return ""
}

func printHasehs(hashes [][]byte) {
	for _, elem := range hashes {
		fmt.Println(hex.EncodeToString(elem))
	}
}
