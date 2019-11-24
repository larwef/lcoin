package lcoin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
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

	b.Index = *big.NewInt(int64(len(c.Blocks)))

	proof := big.NewInt(0)
	var hsh []byte
	for {
		b.Proof = *proof
		hsh = b.HashBlock()
		if hsh[0] == 0 && hsh[1] == 0 {
			break
		}

		proof = proof.Add(proof, big.NewInt(1))
	}

	b.Hash = hsh

	c.Blocks = append(c.Blocks, b)
}

func (c *Chain) Validate() error {
	for i := range c.Blocks {
		if i < len(c.Blocks)-1 && !bytes.Equal(c.Blocks[i].Hash, c.Blocks[i+1].PrevHash) {
			return &Error{
				err:   ErrPrevHashNotEqual,
				Index: i + 1,
			}
		}

		elemHash := c.Blocks[i].HashBlock()

		if !bytes.Equal(c.Blocks[i].Hash, elemHash) {
			return &Error{
				err:   ErrHashNotEqual,
				Index: i,
			}
		}
	}

	return nil
}

func (c *Chain) Print() {
	for _, elem := range c.Blocks {
		elem.print()
		fmt.Println("")
	}
}

func (c *Chain) ToJson(filename string) error {
	btes, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, btes, 0644)
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
