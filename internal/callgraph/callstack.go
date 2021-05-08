package callgraph

import (
	"github.com/vkcom/nocolor/internal/palette"
)

type CallstackOfColoredFunctions struct {
	// Stack is functions placed in order, only colored functions.
	Stack Nodes
	// ColorsChain is blended colors of stacked functions, one-by-one.
	ColorsChain []palette.Color
	// IndexSet is quick index for contains(), has the same elements as stack.
	IndexSet map[*Node]struct{}
	// ColorsMask is mask of all colors_chain.
	ColorsMask palette.ColorMask
}

func NewCallstackOfColoredFunctions() *CallstackOfColoredFunctions {
	return &CallstackOfColoredFunctions{
		Stack:    make(Nodes, 0, 10),
		IndexSet: make(map[*Node]struct{}, 10),
	}
}

func (c *CallstackOfColoredFunctions) Size() int {
	return len(c.Stack)
}

func (c *CallstackOfColoredFunctions) AsVector() Nodes {
	return c.Stack
}

func (c *CallstackOfColoredFunctions) Append(fun *Node) {
	c.Stack = append(c.Stack, fun)
	c.IndexSet[fun] = struct{}{}
	c.ColorsChain = append(c.ColorsChain, fun.Function.Colors.Colors...)
	c.recalcMask()
}

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

func (c *CallstackOfColoredFunctions) Contains(fun *Node) bool {
	_, ok := c.IndexSet[fun]
	return ok
}

func (c *CallstackOfColoredFunctions) recalcMask() {
	c.ColorsMask = 0
	for _, color := range c.ColorsChain {
		c.ColorsMask |= color
	}
}
