package merkle

import (
	"crypto/sha256"
)

type Tree struct {
	Root *Node
}

type Node struct {
	Left  *Node
	Right *Node
	Hash  []byte
}

func NewTree(hashes [][]byte) *Tree {
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

func (t *Tree) Path(hash []byte) []string {

}

func newNode(left, right *Node) *Node {
	h := sha256.New()
	h.Write(left.Hash)
	h.Write(right.Hash)
	firstHash := h.Sum(nil)
	h.Reset()

	h.Write(firstHash)

	return &Node{
		Left:  left,
		Right: right,
		Hash:  h.Sum(nil),
	}
}
