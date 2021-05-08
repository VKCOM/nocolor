package walkers

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"

	"github.com/vkcom/nocolor/internal/symbols"
)

type RootIndexer struct {
	linter.RootCheckerDefaults

	ctx       *linter.RootContext
	meta      FileMeta
	globalCtx *GlobalContext
}

func NewRootIndexer(globalCtx *GlobalContext, ctx *linter.RootContext) *RootIndexer {
	return &RootIndexer{
		ctx:       ctx,
		globalCtx: globalCtx,
		meta:      NewFileMeta(),
	}
}

func generateFileFuncName(filename string) string {
	hash := md5.Sum([]byte(filename))

	return "src$" + hex.EncodeToString(hash[:]) + "$" + filename
}

func (r *RootIndexer) BeforeEnterFile() {
	external := strings.Contains(r.ctx.Filename(), "phpstorm-stubs")
	if external {
		return
	}

	r.meta.Functions.Add(&symbols.Function{
		Name:     generateFileFuncName(r.ctx.Filename()),
		Type:     symbols.MainFunc,
		Called:   symbols.NewFunctions(),
		CalledBy: symbols.NewFunctions(),
	})
}

func (r *RootIndexer) AfterLeaveFile() {
	r.globalCtx.UpdateMeta(&r.meta, "")
}

func (r *RootIndexer) BeforeEnterNode(n ir.Node) {
	external := strings.Contains(r.ctx.Filename(), "phpstorm-stubs")
	if external {
		return
	}

	switch n := n.(type) {
	case *ir.ClassMethodStmt:
		class := r.ctx.ClassParseState().CurrentClass
		methodName := class + "::" + n.MethodName.Value

		typ := symbols.LocalFunc
		if external {
			typ = symbols.ExternFunc
		}

		r.meta.Functions.Add(&symbols.Function{
			Name:     methodName,
			Type:     typ,
			Pos:      r.getElementPos(n),
			Called:   symbols.NewFunctions(),
			CalledBy: symbols.NewFunctions(),
		})

	case *ir.FunctionStmt:
		namespace := r.ctx.ClassParseState().Namespace
		funcName := n.FunctionName.Value

		if namespace != "" {
			funcName = namespace + `\` + funcName
		} else {
			funcName = `\` + funcName
		}

		typ := symbols.LocalFunc
		if external {
			typ = symbols.ExternFunc
		}

		r.meta.Functions.Add(&symbols.Function{
			Name:     funcName,
			Type:     typ,
			Pos:      r.getElementPos(n),
			Called:   symbols.NewFunctions(),
			CalledBy: symbols.NewFunctions(),
		})
	}
}

func (r *RootIndexer) getElementPos(n ir.Node) meta.ElementPosition {
	pos := ir.GetPosition(n)

	return meta.ElementPosition{
		Filename:  r.ctx.ClassParseState().CurrentFile,
		Character: int32(0),
		Line:      int32(pos.StartLine),
		EndLine:   int32(pos.EndLine),
		Length:    int32(pos.EndPos - pos.StartPos),
	}
}
