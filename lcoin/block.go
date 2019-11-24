package lcoin

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
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

	byts, err := hex.DecodeString(v)
	if err != nil {
		return err
	}

	*m = byts

	return nil
}

type Block struct {
	Payload  string `json:"payload"`
	PrevHash MyHash `json:"prevHash"`
	Hash     MyHash `json:"hash"`
}

func (b *Block) print() {
	fmt.Printf("Payload:\t%s\nPrev Hash:\t%x\nHash:\t\t%x\n", b.Payload, b.PrevHash, b.Hash)
}
