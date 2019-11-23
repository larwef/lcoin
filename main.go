package main

import (
	"crypto/sha256"
	"fmt"
)

type chain struct {
	root *block
	last *block
}

func (c *chain) add(b *block) {
	if c.root == nil {
		b.prevHash = []byte{0}
		c.root = b
	} else {
		b.prevHash = c.last.hash
		c.last.next = b
	}

	h := sha256.New()
	h.Write(b.prevHash)
	h.Write(b.payload)
	b.hash = h.Sum(nil)

	c.last = b
}

func (c *chain) print() {
	for current := c.root; current != nil; current = current.next {
		current.print()
		fmt.Println("")
	}
}

type block struct {
	payload  []byte
	prevHash []byte
	hash     []byte
	next     *block
}

func (b *block) print() {
	fmt.Printf("Payload:\t%s\nPrev Hash:\t%x\nHash:\t\t%x\n", string(b.payload), b.prevHash, b.hash)
}

func main() {

	var chn chain

	b1 := &block{payload: []byte("This is the first block")}
	b2 := &block{payload: []byte("This is the second block")}
	b3 := &block{payload: []byte("This is the third block")}
	b4 := &block{payload: []byte("This is the fourth block")}
	b5 := &block{payload: []byte("This is the fifth block")}

	chn.add(b1)
	chn.add(b2)
	chn.add(b3)
	chn.add(b4)
	chn.add(b5)

	chn.print()
}
