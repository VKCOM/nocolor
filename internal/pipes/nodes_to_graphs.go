package pipes

import (
	"github.com/vkcom/nocolor/internal/callgraph"
)

// NodesToGraphs splits the graph into connectivity components.
//
// Graphs from one node are skipped.
func NodesToGraphs(nodes callgraph.Nodes) []*callgraph.Graph {
	visited := make(map[*callgraph.Node]struct{}, len(nodes))
	graphs := make([]*callgraph.Graph, 0, 10)
	queue := make(callgraph.Nodes, 0, 10)
	graphFunctions := make(callgraph.Nodes, 0, 10)

	for _, node := range nodes {
		if _, ok := visited[node]; ok {
			continue
		}

		queue = append(queue, node)

		for len(queue) != 0 {
			node := queue[0]
			queue = queue[1:]

			if _, ok := visited[node]; ok {
				continue
			}
			visited[node] = struct{}{}

			graphFunctions = append(graphFunctions, node)

			queue = append(queue, node.Next...)
			queue = append(queue, node.Prev...)
		}

		if len(graphFunctions) > 1 {
			graph, revGraph := functionsToRawGraphs(graphFunctions)

			graphs = append(graphs, &callgraph.Graph{
				Functions: graphFunctions,
				Graph:     graph,
				RevGraph:  revGraph,
			})
		}

		graphFunctions = make(callgraph.Nodes, 0, 10)
	}

	return graphs
}

func functionsToRawGraphs(funcs callgraph.Nodes) (graph, revGraph callgraph.RawGraph) {
	graph = make(callgraph.RawGraph, len(funcs))
	revGraph = make(callgraph.RawGraph, len(funcs))

	for _, node := range funcs {
		graph[node] = node.Next
		revGraph[node] = node.Prev
	}

	return graph, revGraph
}
