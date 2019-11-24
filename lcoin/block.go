package lcoin

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
)

type MyHash []byte

func (m *MyHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(*m))
}

func (m *MyHash) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	hexBytes, err := hex.DecodeString(v)
	if err != nil {
		return err
	}

	*m = hexBytes

	return nil
}

type Block struct {
	Index    big.Int `json:"index"`
	Payload  string  `json:"payload"`
	PrevHash MyHash  `json:"prevHash"`
	Hash     MyHash  `json:"hash"`
	Proof    big.Int `json:"proof"`
}

func (b *Block) HashBlock() []byte {
	h := sha256.New()
	h.Write(b.Index.Bytes())
	h.Write([]byte(b.Payload))
	h.Write(b.Proof.Bytes())
	h.Write(b.PrevHash)

	return h.Sum(nil)
}

func (b *Block) print() {
	fmt.Printf("Payload:\t%s\nPrev Hash:\t%x\nHash:\t\t%x\n", b.Payload, b.PrevHash, b.Hash)
}
