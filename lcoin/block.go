package lcoin

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

type Chain struct {
	Root *Block
}

func (c *Chain) Add(b *Block) {

}

type Header struct {
	PrevBlockHash string `json:"prevHash"`
	MerkleRoot    string `json:"merkleRoot"`
	Time          int64  `json:"time"`
	Nonce         int64  `json:"nonce"`
}

func (h *Header) Hash() ([32]byte, error) {
	hsh := sha256.New()
	prev, err := hex.DecodeString(h.PrevBlockHash)
	if err != nil {
		return [32]byte{}, err
	}
	root, err := hex.DecodeString(h.MerkleRoot)
	if err != nil {
		return [32]byte{}, err
	}

	var t []byte
	binary.BigEndian.PutUint64(t, uint64(h.Time))

	var nonce []byte
	binary.BigEndian.PutUint64(nonce, uint64(h.Nonce))

	hsh.Write(prev)
	hsh.Write(root)
	hsh.Write(t)
	hsh.Write(nonce)

	return sha256.Sum256(hsh.Sum(nil)), nil
}

type Input struct {
	PrevTX    string `json:"prevTX"`
	Index     int    `json:"index"`
	ScriptSig string `json:"scriptSig"`
}

type Output struct {
	Value           int    `json:"value"`
	ScriptPublicKey string `json:"scriptPublicKey"`
}

type Transaction struct {
	Inputs  []*Input  `json:"inputs"`
	Outputs []*Output `json:"outputs"`
}

func (t *Transaction) Hash() [32]byte {
	h := sha256.New()

	for _, elem := range t.Inputs {
		h.Write([]byte(elem.PrevTX))
		var index []byte
		binary.BigEndian.PutUint64(index, uint64(elem.Index))
		h.Write(index)
		h.Write([]byte(elem.ScriptSig))
	}

	for _, elem := range t.Outputs {
		var value []byte
		binary.BigEndian.PutUint64(value, uint64(elem.Value))
		h.Write(value)
		h.Write([]byte(elem.ScriptPublicKey))
	}

	return sha256.Sum256(h.Sum(nil))
}

type Block struct {
	Header       Header         `json:"header"`
	Transactions []*Transaction `json:"transactions"`
}
