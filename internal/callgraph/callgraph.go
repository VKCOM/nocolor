package callgraph

import (
	"github.com/vkcom/nocolor/internal/symbols"
)

// Nodes is an alias for a slice of graph nodes for which a Remove
// function is defined for ease of interaction.
type Nodes []*Node

// Remove is a function that removes the passed node from the slice if any.
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

// RawGraph is a map that stores all the functions that are called
// in the function or the function in which the function is called.
type RawGraph map[*Node]Nodes

// Graph is a structure for storing the complete call graph, as
// well as the functions that are included in it.
type Graph struct {
	Functions Nodes

	Graph    RawGraph
	RevGraph RawGraph
}

// Remove is a function that removes the passed node from the Graph.
//
// Note: Node is not removed from the RevGraph.
func (g *Graph) Remove(node *Node) {
	callers := g.RevGraph[node]
	for _, caller := range callers {
		g.Graph[caller] = g.Graph[caller].Remove(node)
	}
}

// Node is a structure for storing information about the functions that call the
// function, the functions in which the function is called, and the reachable
// colored functions from the function.
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
