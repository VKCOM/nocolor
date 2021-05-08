package walkers

import (
	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/linter"
)

type BlockIndexer struct {
	linter.BlockCheckerDefaults
}

type BlockChecker struct {
	linter.BlockCheckerDefaults
	ctx  *linter.BlockContext
	root *RootChecker
}

func NewBlockChecker(ctx *linter.BlockContext, root *RootChecker) *BlockChecker {
	return &BlockChecker{
		ctx:  ctx,
		root: root,
	}
}

func (b *BlockChecker) EnterNode(n ir.Node) bool {
	b.BeforeEnterNode(n)
	return true
}

func (b *BlockChecker) LeaveNode(n ir.Node) {}

func (b *BlockChecker) BeforeEnterNode(n ir.Node) {
	switch n := n.(type) {
	case *ir.NewExpr:
		b.root.handleNew(n)
	case *ir.FunctionCallExpr:
		b.root.handleFunctionCall(n, b)
	case *ir.StaticCallExpr:
		b.root.handleStaticCall(n, b.ctx.Scope())
	case *ir.MethodCallExpr:
		b.root.handleMethodCall(n, b.ctx.Scope(), b)
	case *ir.ImportExpr:
		b.root.handleImportExpr(n)
	}
}
