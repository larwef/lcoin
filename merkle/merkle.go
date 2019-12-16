package merkle

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
)

var ErrNodeNotFound = errors.New("couldn't find node")

type Tree struct {
	Root  *Node
	Depth int
}

type Node struct {
	Index int
	Left  *Node
	Right *Node
	Hash  [32]byte
}

func NewTree(hashes [][32]byte) *Tree {
	var nodes []*Node

	depth := -1
	// Make leaf nodes
	for i, elem := range hashes {
		node := Node{
			Index: i,
			Hash:  elem,
		}
		nodes = append(nodes, &node)
	}
	depth++

	for len(nodes) > 1 {
		var tmp []*Node
		for i := 0; i < len(nodes)-1; i = i + 2 {
			tmp = append(tmp, newNode(nodes[i], nodes[i+1]))
		}

		if len(nodes)%2 != 0 {
			tmp = append(tmp, newNode(nodes[len(nodes)-1], nodes[len(nodes)-1]))
		}

		nodes = tmp
		depth++
	}

	return &Tree{
		Root:  nodes[0],
		Depth: depth,
	}
}

func (t *Tree) ProofSearch(hash [32]byte) (*Proof, error) {
	leaf := t.findLeaf(hash)
	if leaf == nil {
		return nil, ErrNodeNotFound
	}

	return t.Proof(leaf.Index), nil
}

func (t *Tree) Proof(index int) *Proof {
	result := &Proof{
		Index: index,
		Depth: t.Depth,
	}

	path := 1
	exponent := t.Depth
	for exponent != 0 {
		path *= 2
		exponent -= 1
	}
	path += index

	pathStr := fmt.Sprintf("%b", path)

	current := t.Root
	for _, elem := range pathStr[1:] {
		if string(elem) == "0" {
			result.Hashes = append(result.Hashes, current.Right.Hash)
			current = current.Left
		} else {
			result.Hashes = append(result.Hashes, current.Left.Hash)
			current = current.Right
		}
	}

	return result
}

func (t *Tree) findLeaf(hash [32]byte) *Node {
	var stack []*Node
	stack = append(stack, t.Root)
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1] // Pop of stack

		if n.Left != nil {
			stack = append(stack, n.Left)
			stack = append(stack, n.Right)
		}

		if bytes.Equal(n.Hash[:], hash[:]) {
			return n
		}
	}

	return nil
}

func newNode(left, right *Node) *Node {
	h1 := sha256.Sum256(append(left.Hash[:], right.Hash[:]...))
	node := &Node{
		Left:  left,
		Right: right,
		Hash:  sha256.Sum256(h1[:]),
	}

	return node
}

type Proof struct {
	Index  int        `json:"index"`
	Depth  int        `json:"depth"`
	Hashes [][32]byte `json:"hashes"`
}

func (p *Proof) Root(leaf [32]byte) [32]byte {
	var tmp [32]byte

	path := 1
	exponent := p.Depth
	for exponent != 0 {
		path *= 2
		exponent -= 1
	}
	path += p.Index

	pathStr := fmt.Sprintf("%b", path)

	copy(tmp[:], leaf[:])
	for i := len(pathStr) - 1; i > 0; i-- {
		if string(pathStr[i]) == "0" {
			tmp = sha256.Sum256(append(tmp[:], p.Hashes[i-1][:]...))
		} else {
			tmp = sha256.Sum256(append(p.Hashes[i-1][:], tmp[:]...))
		}
		tmp = sha256.Sum256(tmp[:])
	}

	return tmp
}
