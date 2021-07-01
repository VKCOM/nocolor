package symbols

import (
	"strings"
	"sync"

	"github.com/VKCOM/noverify/src/meta"

	"github.com/vkcom/nocolor/internal/palette"
)

type ClassType uint8

const (
	PlainClass ClassType = iota
	Interface
	Trait
)

// Class is a structure for storing information about a class.
type Class struct {
	Name   string
	Type   ClassType
	Pos    meta.ElementPosition
	Colors *palette.ColorContainer

	WithExplicitConstructor bool
}

// HumanReadableName returns a string with a name that is understandable.
func (f *Class) HumanReadableName() string {
	return strings.TrimPrefix(f.Name, `\`)
}

func (f *Class) String() string {
	return f.Name
}

func (f *Class) HasColors() bool {
	return !f.Colors.Empty()
}

type Classes struct {
	mtx     sync.Mutex
	Classes map[string]*Class
}

func NewClasses() *Classes {
	return &Classes{Classes: map[string]*Class{}}
}

func (f *Classes) Get(name string) (*Class, bool) {
	class, ok := f.Classes[name]
	return class, ok
}

func (f *Classes) Raw() map[string]*Class {
	return f.Classes
}

func (f *Classes) Len() int {
	return len(f.Classes)
}

func (f *Classes) Add(class *Class) {
	f.mtx.Lock()
	f.Classes[class.Name] = class
	f.mtx.Unlock()
}
