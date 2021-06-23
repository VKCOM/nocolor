package callgraph

import (
	"github.com/vkcom/nocolor/internal/palette"
)

// CallstackOfColoredFunctions is a structure for storing a stack of called
// colored functions with a quick check for the presence of a certain node.
type CallstackOfColoredFunctions struct {
	// Stack is functions placed in order, only colored functions.
	Stack Nodes
	// ColorsChain is blended colors of stacked functions, one-by-one.
	ColorsChain []palette.Color
	// IndexSet is quick index for Contains(), has the same elements as stack.
	IndexSet map[*Node]struct{}
	// ColorsMask is mask of all ColorsChain.
	ColorsMasks palette.ColorMasks
}

// NewCallstackOfColoredFunctions creates a new callstack.
func NewCallstackOfColoredFunctions() *CallstackOfColoredFunctions {
	return &CallstackOfColoredFunctions{
		Stack:    make(Nodes, 0, 10),
		IndexSet: make(map[*Node]struct{}, 10),
	}
}

// Size returns the number of functions.
func (c *CallstackOfColoredFunctions) Size() int {
	return len(c.Stack)
}

// AsVector returns a slice of the functions that are on the stack.
func (c *CallstackOfColoredFunctions) AsVector() Nodes {
	return c.Stack
}

// Append adds the passed node onto the stack.
func (c *CallstackOfColoredFunctions) Append(fun *Node) {
	c.Stack = append(c.Stack, fun)
	c.IndexSet[fun] = struct{}{}
	c.ColorsChain = append(c.ColorsChain, fun.Function.Colors.Colors...)
	c.recalcMask()
}

// PopBack removes the last function from the stack.
func (c *CallstackOfColoredFunctions) PopBack() {
	if len(c.Stack) == 0 {
		return
	}

	back := c.Stack[len(c.Stack)-1]
	c.Stack = c.Stack[:len(c.Stack)-1]
	delete(c.IndexSet, back)
	c.ColorsChain = c.ColorsChain[:len(c.ColorsChain)-len(back.Function.Colors.Colors)]
	c.recalcMask()
}

// Contains checks for the existence of the passed node.
func (c *CallstackOfColoredFunctions) Contains(fun *Node) bool {
	_, ok := c.IndexSet[fun]
	return ok
}

func (c *CallstackOfColoredFunctions) recalcMask() {
	c.ColorsMasks = palette.NewColorMasks(c.ColorsChain)
}
