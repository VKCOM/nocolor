package callgraph

import (
	"github.com/vkcom/nocolor/internal/symbols"
)

type RawGraph map[*Node]Nodes

type Nodes []*Node

func (n Nodes) Remove(node *Node) Nodes {
	index := -1
	for i, fun := range n {
		if fun == node {
			index = i
			continue
		}
	}

	if index != -1 {
		return removeHelper(n, index)
	}

	return n
}

func removeHelper(s Nodes, i int) Nodes {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

type Graph struct {
	Functions Nodes

	Graph    RawGraph
	RevGraph RawGraph
}

func (g *Graph) Remove(node *Node) {
	callers := g.RevGraph[node]
	for _, caller := range callers {
		g.Graph[caller] = g.Graph[caller].Remove(node)
	}
}

type Node struct {
	Function *symbols.Function

	// Next is an array of functions that are called from the current one.
	// Prev is an array of functions that call the current one.
	//
	// All functions are always contained here, in contrast to Graph,
	// where some functions can be deleted.
	// Use these fields only if you need to know if there is a connection,
	// but for other, use the Graph.
	Next Nodes
	Prev Nodes

	// Pointer to a slice containing the
	// following nodes that have colors.
	NextWithColors *Nodes
}

// String method for debugging.
func (n *Node) String() string {
	return n.Function.HumanReadableName()
}
