package lcoin

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	ErrPrevHashNotEqual = errors.New("previous hash not equal")
	ErrHashNotEqual     = errors.New("hash not equal")
)

type Error struct {
	err   error
	Index int
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error on Index %d: %v", e.Index, e.err)
}

func (e *Error) Unwrap() error {
	return e.err
}

type Chain struct {
	Blocks []*Block `json:"blocks"`
}

func (c *Chain) Add(b *Block) {
	if len(c.Blocks) == 0 {
		b.PrevHash = []byte{0}
	} else {
		b.PrevHash = c.Blocks[len(c.Blocks)-1].Hash
	}

	h := sha256.New()
	h.Write(b.PrevHash)
	h.Write([]byte(b.Payload))
	b.Hash = h.Sum(nil)

	c.Blocks = append(c.Blocks, b)
}

func (c *Chain) Validate() error {
	h := sha256.New()
	for i := range c.Blocks {
		if i < len(c.Blocks)-1 && !isEqual(c.Blocks[i].Hash, c.Blocks[i+1].PrevHash) {
			return &Error{
				err:   ErrPrevHashNotEqual,
				Index: i + 1,
			}
		}

		h.Reset()
		h.Write(c.Blocks[i].PrevHash)
		h.Write([]byte(c.Blocks[i].Payload))
		elemHash := h.Sum(nil)

		if !isEqual(c.Blocks[i].Hash, elemHash) {
			return &Error{
				err:   ErrHashNotEqual,
				Index: i,
			}
		}
	}

	return nil
}

func isEqual(x, y []byte) bool {
	if len(x) != len(y) {
		return false
	}

	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			return false
		}
	}

	return true
}

func (c *Chain) Print() {
	for _, elem := range c.Blocks {
		elem.print()
		fmt.Println("")
	}
}

func (c *Chain) ToJson(filename string) error {
	bytes, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *Chain) FromJson(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	return json.NewDecoder(file).Decode(c)
}
