package pipes

import (
	"github.com/vkcom/nocolor/internal/callgraph"
	"github.com/vkcom/nocolor/internal/palette"
)

// EraseNodesWithRemoverColor is function for drop '@color remover' functions from a
// call graph.
//
// To perform calculating NextWithColors and checking color rules from the palette,
// we need to drop '@color remover' functions from a call graph completely, like
// they don't exist at all, this special color is for manual cutting connectivity
// rules, allowing to explicitly separate recursively-joint components
func EraseNodesWithRemoverColor(graph *callgraph.Graph) {
	removerColor := palette.NewColor(palette.SpecialColorRemover, 0)

	for _, fun := range graph.Functions {
		if fun.Function.HasColors() && fun.Function.Colors.Contains(removerColor) {
			graph.Remove(fun)
			graph.Functions.Remove(fun)
		}
	}
}
