package pipes

import (
	"github.com/vkcom/nocolor/internal/callgraph"
	"github.com/vkcom/nocolor/internal/palette"
)

// To perform calculating NextWithColors and checking color rules from the palette,
// we need to drop '@color remover' functions from a call graph completely, like
// they don't exist at all, this special color is for manual cutting connectivity
// rules, allowing to explicitly separate recursively-joint components
func EraseNodesWithRemoverColor(graph *callgraph.Graph) {
	for _, fun := range graph.Functions {
		if fun.Function.HasColors() && fun.Function.Colors.Contains(palette.SpecialColorRemover) {
			graph.Remove(fun)
			graph.Functions.Remove(fun)
		}
	}
}
