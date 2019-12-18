package merkle

import (
	"bytes"
	"crypto/sha256"
	"errors"
)

var ErrNodeNotFound = errors.New("couldn't find node")

// Tree keeps track of the root of the Merkle tree.
type Tree struct {
	Root  *Node
	Depth int
}

// Node is a representation of a node in the Merkle tree.
type Node struct {
	Index int
	Left  *Node
	Right *Node
	Hash  [32]byte
}

// NewTree makes a new Merkle tree from a list of leaf hashes.
func NewTree(hashes [][32]byte) *Tree {
	var nodes []*Node

	// Make leaf nodes
	for i, elem := range hashes {
		node := Node{
			Index: i,
			Hash:  elem,
		}
		nodes = append(nodes, &node)
	}
	depth := 0

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

// ProofSearch searches for the hash and provides the proof if present.
func (t *Tree) ProofSearch(hash [32]byte) (*Proof, error) {
	leaf := t.findLeaf(hash)
	if leaf == nil {
		return nil, ErrNodeNotFound
	}

	return t.Proof(leaf.Index), nil
}

// Proof provides the proof for the provided index.
func (t *Tree) Proof(index int) *Proof {
	result := &Proof{
		Index: index,
		Depth: t.Depth,
	}

	// (2^depth)+index. 0: left, 1: right.
	path := 1<<t.Depth + index

	current := t.Root
	for i := 0; i < t.Depth; i++ {
		if path&(1<<(t.Depth-1-i)) == 0 {
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

// Root calculated the merkle root based on the Proof.
func (p *Proof) Root(leaf [32]byte) [32]byte {
	// (2^depth)+index
	path := 1<<p.Depth + p.Index

	var tmp [32]byte
	copy(tmp[:], leaf[:])
	for i := 0; i < p.Depth; i++ {
		if path&(1<<i) == 0 {
			tmp = sha256.Sum256(append(tmp[:], p.Hashes[p.Depth-i-1][:]...))
		} else {
			tmp = sha256.Sum256(append(p.Hashes[p.Depth-i-1][:], tmp[:]...))
		}
		tmp = sha256.Sum256(tmp[:])
	}

	return tmp
}
