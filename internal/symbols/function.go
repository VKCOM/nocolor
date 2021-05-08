package symbols

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/VKCOM/noverify/src/meta"

	"github.com/vkcom/nocolor/internal/palette"
)

type FunctionType int64

const (
	MainFunc FunctionType = iota
	LocalFunc
	ExternFunc
)

type Function struct {
	Name   string
	Type   FunctionType
	Pos    meta.ElementPosition
	Colors palette.ColorContainer

	Called   *Functions
	CalledBy *Functions
}

func (f *Function) HumanReadableName() string {
	if f.Type == MainFunc {
		name := f.Name
		path := name[strings.LastIndex(name, "$")+1:]

		wd, err := os.Getwd()
		if err == nil {
			relPath, err := filepath.Rel(wd, path)
			if err == nil {
				path = relPath
			}
		}

		path = filepath.ToSlash(path)

		return fmt.Sprintf("file '%s' scope", path)
	}

	return strings.TrimPrefix(f.Name, `\`)
}

func (f *Function) String() string {
	return f.Name
}

func (f *Function) HasColors() bool {
	return !f.Colors.Empty()
}

type Functions struct {
	mtx       sync.Mutex
	Functions map[string]*Function
}

func NewFunctions() *Functions {
	return &Functions{Functions: map[string]*Function{}}
}

func (f *Functions) Get(name string) (*Function, bool) {
	fun, ok := f.Functions[name]
	return fun, ok
}

func (f *Functions) Raw() map[string]*Function {
	return f.Functions
}

func (f *Functions) Len() int {
	return len(f.Functions)
}

func (f *Functions) Add(fun *Function) {
	f.mtx.Lock()
	f.Functions[fun.Name] = fun
	f.mtx.Unlock()
}
