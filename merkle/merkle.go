package merkle

import (
	"bytes"
	"crypto/sha256"
)

type Tree struct {
	Root *Node
}

type Node struct {
	Left      *Node
	Right     *Node
	Hash      [32]byte
}

func NewTree(hashes [][32]byte) *Tree {
	var nodes []*Node

	// Make leaf nodes
	for _, elem := range hashes {
		node := Node{Hash: elem}
		nodes = append(nodes, &node)
	}

	for len(nodes) > 1 {
		var tmp []*Node
		for i := 0; i < len(nodes)-1; i = i + 2 {
			tmp = append(tmp, newNode(nodes[i], nodes[i+1]))
		}

		if len(nodes)%2 != 0 {
			tmp = append(tmp, newNode(nodes[len(nodes)-1], nodes[len(nodes)-1]))
		}

		nodes = tmp
	}

	return &Tree{
		Root: nodes[0],
	}
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
		Left:   left,
		Right:  right,
		Hash:   sha256.Sum256(h1[:]),
	}

	return node
}
