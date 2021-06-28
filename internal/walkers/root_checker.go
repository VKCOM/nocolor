package walkers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/VKCOM/noverify/src/constfold"
	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"
	"github.com/VKCOM/noverify/src/phpdoc"
	"github.com/VKCOM/noverify/src/solver"
	"github.com/VKCOM/noverify/src/types"

	"github.com/vkcom/nocolor/internal/palette"
	"github.com/vkcom/nocolor/internal/symbols"
)

// RootChecker is a walker that collects information about the
// colors of functions and checks them for correctness.
type RootChecker struct {
	linter.RootCheckerDefaults

	ctx       *linter.RootContext
	palette   *palette.Palette
	globalCtx *GlobalContext

	fileFunction *symbols.Function

	colorTag string
}

// NewRootChecker returns a new walker.
func NewRootChecker(palette *palette.Palette, globalCtx *GlobalContext, ctx *linter.RootContext, colorTag string) *RootChecker {
	return &RootChecker{
		ctx:       ctx,
		palette:   palette,
		globalCtx: globalCtx,
		colorTag:  colorTag,
	}
}

// EnterNode is method to use RootChecker in the Walk method of AST nodes.
func (r *RootChecker) EnterNode(n ir.Node) bool {
	r.BeforeEnterNode(n)
	return true
}

// LeaveNode is method to use RootChecker in the Walk method of AST nodes.
func (r *RootChecker) LeaveNode(ir.Node) {}

// BeforeEnterFile sets the current function of the file.
func (r *RootChecker) BeforeEnterFile() {
	fun, ok := r.globalCtx.Functions.Get(generateFileFuncName(r.ctx.Filename()))
	if !ok {
		r.fileFunction = &symbols.Function{
			Name:     generateFileFuncName(r.ctx.Filename()),
			Type:     symbols.MainFunc,
			Called:   symbols.NewFunctions(),
			CalledBy: symbols.NewFunctions(),
		}
		return
	}

	r.fileFunction = fun
}

// AfterEnterNode
func (r *RootChecker) AfterEnterNode(n ir.Node) {
	switch n := n.(type) {
	case *ir.NewExpr:
		r.handleNew(n)
	case *ir.FunctionCallExpr:
		r.handleFunctionCall(n, r)
	case *ir.StaticCallExpr:
		r.handleStaticCall(n, nil)
	case *ir.MethodCallExpr:
		r.handleMethodCall(n, nil, r)

	case *ir.ClassStmt:
		r.handleClassStmt(n.ClassName, n.Doc)
	case *ir.InterfaceStmt:
		r.handleClassStmt(n.InterfaceName, n.Doc)
	case *ir.TraitStmt:
		r.handleClassStmt(n.TraitName, n.Doc)

	case *ir.ClassMethodStmt:
		r.handleClassMethodStmt(n)
	case *ir.FunctionStmt:
		r.handleFunctionStmt(n)
	case *ir.ImportExpr:
		r.handleImportExpr(n)
	}
}

func (r *RootChecker) handleClassStmt(className ir.Node, doc phpdoc.Comment) {
	errs := r.checkPhpDocColors(doc)
	for _, err := range errs {
		r.ctx.Report(className, linter.LevelError, "errorColor", err)
	}
}

func (r *RootChecker) handleImportExpr(n *ir.ImportExpr) {
	pathValue := constfold.Eval(r.ctx.ClassParseState(), n.Expr)
	if !pathValue.IsValid() {
		return
	}

	path, ok := pathValue.ToString()
	if !ok {
		return
	}

	path, ok = r.getImportAbsPath(path)
	if !ok {
		return
	}

	curFunc, ok := r.getCurrentFunc()
	if !ok {
		return
	}

	fileFunc, ok := r.globalCtx.Functions.Get(generateFileFuncName(path))
	if !ok {
		return
	}

	curFunc.Called.Add(fileFunc)
	fileFunc.CalledBy.Add(curFunc)
}

func (r *RootChecker) getImportAbsPath(path string) (string, bool) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), true
	}

	// If relative path.
	if strings.HasPrefix(path, ".") || strings.HasPrefix(path, "..") {
		currentFilePath := r.ctx.Filename()
		dir := filepath.Dir(currentFilePath)

		absPath := filepath.Clean(filepath.Join(dir, path))
		return absPath, true
	}

	return "", false
}

func (r *RootChecker) handleFunctionStmt(n *ir.FunctionStmt) {
	errs := r.checkPhpDocColors(n.Doc)
	for _, err := range errs {
		r.ctx.Report(n.FunctionName, linter.LevelError, "errorColor", err)
	}
}

func (r *RootChecker) handleClassMethodStmt(n *ir.ClassMethodStmt) {
	errs := r.checkPhpDocColors(n.Doc)
	for _, err := range errs {
		r.ctx.Report(n.MethodName, linter.LevelError, "errorColor", err)
	}
}

func (r *RootChecker) checkPhpDocColors(comment phpdoc.Comment) (errs []string) {
	for _, part := range comment.Parsed {
		p, ok := part.(*phpdoc.RawCommentPart)
		if !ok {
			continue
		}

		if p.Name() != r.colorTag {
			continue
		}

		if len(p.Params) == 0 {
			errs = append(errs, fmt.Sprintf("An empty '@%s' tag value", p.Name()))
			continue
		}

		colorName := p.Params[0]

		if colorName == "transparent" {
			errs = append(errs, "Use of the 'transparent' color does not make sense")
			continue
		}

		if !r.palette.ColorExists(colorName) {
			errs = append(errs, fmt.Sprintf("Color '%s' missing in palette (either a misprint or a new color that needs to be added)", colorName))
			continue
		}
	}

	return errs
}

func (r *RootChecker) handleFunctionCall(n *ir.FunctionCallExpr, v ir.Visitor) {
	for _, arg := range n.Args {
		arg.Walk(v)
	}

	fqName, ok := solver.GetFuncName(r.ctx.ClassParseState(), n.Function)
	if !ok {
		return
	}

	calledFunc, ok := r.globalCtx.Functions.Get(fqName)
	if !ok {
		return
	}

	curFunc, ok := r.getCurrentFunc()
	if !ok {
		return
	}

	curFunc.Called.Add(calledFunc)
	calledFunc.CalledBy.Add(curFunc)
}

func (r *RootChecker) handleStaticCall(n *ir.StaticCallExpr, blockScope *meta.Scope) {
	method, ok := n.Call.(*ir.Identifier)
	if !ok {
		return
	}
	methodName := method.Value

	scope := blockScope
	if scope == nil {
		scope = r.ctx.Scope()
	}

	var classType types.Map

	if vr, ok := n.Class.(*ir.SimpleVar); ok {
		classType = solver.ExprType(scope, r.ctx.ClassParseState(), vr)
	} else {
		className, ok := solver.GetClassName(r.ctx.ClassParseState(), n.Class)
		if !ok {
			return
		}

		classType = types.NewMap(className)
	}

	r.handleMethod(methodName, classType)
}

func (r *RootChecker) handleMethodCall(n *ir.MethodCallExpr, blockScope *meta.Scope, v ir.Visitor) {
	method, ok := n.Method.(*ir.Identifier)
	if !ok {
		return
	}
	methodName := method.Value

	scope := blockScope
	if scope == nil {
		scope = r.ctx.Scope()
	}

	classType := solver.ExprType(scope, r.ctx.ClassParseState(), n.Variable)

	r.handleMethod(methodName, classType)

	for _, nn := range n.Args {
		nn.Walk(v)
	}
}

func (r *RootChecker) handleNew(n *ir.NewExpr) {
	className, ok := solver.GetClassName(r.ctx.ClassParseState(), n.Class)
	if !ok {
		return
	}

	classType := types.NewMap(className)

	r.handleMethod("__construct", classType)
}

func (r *RootChecker) handleMethod(name string, classType types.Map) {
	var calledMethodInfo solver.FindMethodResult
	var implicitConstructor *symbols.Function
	var ok bool

	found := classType.Find(func(typ string) bool {
		calledMethodInfo, ok = solver.FindMethod(r.ctx.ClassParseState().Info, typ, name)

		if !ok && name == "__construct" {
			constructorName := typ + "::__construct (default autogenerated)"
			implicitConstructor, ok = r.globalCtx.Functions.Get(constructorName)
			return ok
		}

		return ok
	})

	if !found && implicitConstructor == nil {
		return
	}

	calledName := calledMethodInfo.Info.Name
	fqn := calledMethodInfo.ImplName() + "::" + calledName

	var calledFunc *symbols.Function
	if implicitConstructor != nil {
		calledFunc = implicitConstructor
	} else {
		calledFunc, ok = r.globalCtx.Functions.Get(fqn)
		if !ok {
			return
		}
	}

	curFunc, ok := r.getCurrentFunc()
	if !ok {
		return
	}

	curFunc.Called.Add(calledFunc)
	calledFunc.CalledBy.Add(curFunc)
}

func (r *RootChecker) getCurrentFunc() (*symbols.Function, bool) {
	name := r.ctx.ClassParseState().CurrentFunction
	if name == "" {
		return r.fileFunction, true
	}

	if r.ctx.ClassParseState().CurrentClass != "" {
		className, ok := solver.GetClassName(r.ctx.ClassParseState(), &ir.Name{Value: r.ctx.ClassParseState().CurrentClass})
		if !ok {
			return nil, false
		}

		fn, ok := r.globalCtx.Functions.Get(className + "::" + name)
		if !ok {
			return nil, false
		}

		return fn, true
	}

	funcName, ok := solver.GetFuncName(r.ctx.ClassParseState(), &ir.Name{
		Value: name,
	})
	if !ok {
		return nil, false
	}

	fn, ok := r.globalCtx.Functions.Get(funcName)
	if !ok {
		return nil, false
	}

	return fn, true
}
