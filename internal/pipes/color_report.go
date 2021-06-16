package pipes

import (
	"github.com/vkcom/nocolor/internal/callgraph"
	"github.com/vkcom/nocolor/internal/palette"
)

// ColorReport is a structure for storing color mixing error information.
type ColorReport struct {
	Rule      *palette.Rule
	CallChain callgraph.Nodes
	Message   string

	Palette *palette.Palette
}
