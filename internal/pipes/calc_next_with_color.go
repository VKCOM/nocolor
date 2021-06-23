package pipes

import (
	"github.com/vkcom/nocolor/internal/callgraph"
	"github.com/vkcom/nocolor/internal/palette"
)

type visitedMap map[*callgraph.Node]int

// CalcNextWithColor function to calculate node.NextWithColors for every function.
// We'll use it for perform dfs only for colored nodes of a call graph.
//
// Note: this precalculation is needed, because dfs for a whole call graph is too
// heavy on a large code base.
func CalcNextWithColor(graph *callgraph.Graph) {
	nodes := graph.Functions
	nodesCount := len(nodes)

	visited := make(visitedMap, nodesCount)
	topSorted := make(callgraph.Nodes, 0, nodesCount)

	for _, fun := range nodes {
		topSortedDFS(fun, graph, visited, &topSorted)
	}

	for i := range visited {
		visited[i] = 0
	}

	wasColors := make([]int, nodesCount+1)
	currentColor := 0

	for i := len(topSorted) - 1; i >= 0; i-- {
		node := topSorted[i]
		if visited[node] == 1 {
			continue
		}

		var component callgraph.Nodes
		var edges callgraph.Nodes

		currentColor++
		componentDFS(node, graph, currentColor, visited, &wasColors, &component, &edges)

		eachComponent(component, edges)
	}
}

func eachComponent(component, edges callgraph.Nodes) {
	nextColoredUniq := map[*callgraph.Node]struct{}{}

	// If an edge is colored, append this edge.
	// If not, append NextWithColors from this edge.
	for _, fun := range edges {
		if fun.Function.HasColors() {
			nextColoredUniq[fun] = struct{}{}
		} else if fun.NextWithColors != nil {
			for _, node := range *fun.NextWithColors {
				nextColoredUniq[node] = struct{}{}
			}
		}
	}

	// For a recursive bunch of functions (e.g. f1 -> f2 -> f1) assume that
	// "every functions is reachable from any of them" — so, append all colored
	// functions from a recurse. This is not true in all cases, but enough for
	// practical usage.
	//
	// Also, to reduce a number of permutations for next colored path searching,
	// choose only one representative of each color.
	hasColorsInComponent := false
	for _, node := range component {
		if node.Function.HasColors() {
			hasColorsInComponent = true
			break
		}
	}
	if hasColorsInComponent && len(component) > 1 {
		added := palette.NewEmptyColorMasks()
		for _, node := range component {
			for _, color := range node.Function.Colors.Colors {
				if !added.Contains(color) {
					added = added.Add(color)
					nextColoredUniq[node] = struct{}{}
				}
			}
		}
	}

	if len(nextColoredUniq) != 0 {
		nextWithColors := make(callgraph.Nodes, 0, len(nextColoredUniq))
		for node := range nextColoredUniq {
			nextWithColors = append(nextWithColors, node)
		}

		for _, node := range component {
			node.NextWithColors = &nextWithColors
		}
	}
}

func topSortedDFS(fun *callgraph.Node, graph *callgraph.Graph, visited visitedMap, topSorted *callgraph.Nodes) {
	if visited[fun] == 1 {
		return
	}
	visited[fun] = 1

	for _, prev := range graph.RevGraph[fun] {
		topSortedDFS(prev, graph, visited, topSorted)
	}

	*topSorted = append(*topSorted, fun)
}

func componentDFS(fun *callgraph.Node, graph *callgraph.Graph, color int, visited visitedMap, wasColors *[]int, component, edges *callgraph.Nodes) {
	otherColor := visited[fun]
	if otherColor == color {
		return
	}

	if otherColor != 0 {
		if (*wasColors)[otherColor] != color {
			(*wasColors)[otherColor] = color
			*edges = append(*edges, fun)
		}
		return
	}

	visited[fun] = color
	*component = append(*component, fun)

	for _, next := range graph.Graph[fun] {
		componentDFS(next, graph, color, visited, wasColors, component, edges)
	}
}
