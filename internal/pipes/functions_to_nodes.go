package pipes

import (
	"github.com/vkcom/nocolor/internal/callgraph"
	"github.com/vkcom/nocolor/internal/symbols"
)

func FunctionsToNodes(funcs *symbols.Functions) callgraph.Nodes {
	nodes := make(callgraph.Nodes, 0, funcs.Len())
	visited := make(map[*symbols.Function]*callgraph.Node, funcs.Len())

	for _, fun := range funcs.Functions {
		nodes = append(nodes, functionToNode(fun, visited))
	}

	return nodes
}

func functionToNode(fun *symbols.Function, visited map[*symbols.Function]*callgraph.Node) *callgraph.Node {
	if node, ok := visited[fun]; ok {
		return node
	}

	var node callgraph.Node
	visited[fun] = &node

	node.Function = fun

	for _, called := range fun.Called.Functions {
		node.Next = append(node.Next, functionToNode(called, visited))
	}

	for _, calledBy := range fun.CalledBy.Functions {
		node.Prev = append(node.Prev, functionToNode(calledBy, visited))
	}

	return &node
}
