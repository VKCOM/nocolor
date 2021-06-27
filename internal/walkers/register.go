package walkers

import (
	"github.com/VKCOM/noverify/src/linter"

	"github.com/vkcom/nocolor/internal/palette"
)

// Register registers custom walkers to collect information about functions.
func Register(config *linter.Config, globalCtx *GlobalContext, pal *palette.Palette, colorTag string) {
	config.Checkers.AddBlockChecker(func(ctx *linter.BlockContext) linter.BlockChecker {
		if ctx.ClassParseState().Info.IsIndexingComplete() {
			return NewBlockChecker(ctx, ctx.RootState()["lints-root"].(*RootChecker))
		}

		return &BlockIndexer{}
	})

	config.Checkers.AddRootCheckerWithCacher(globalCtx, func(ctx *linter.RootContext) linter.RootChecker {
		if ctx.ClassParseState().Info.IsIndexingComplete() {
			checker := NewRootChecker(pal, globalCtx, ctx, colorTag)

			ctx.State()["lints-root"] = checker
			return checker
		}

		indexer := NewRootIndexer(pal, globalCtx, ctx, colorTag)
		ctx.State()["lints-root"] = indexer
		return indexer
	})
}
