package walkers

import (
	"github.com/vkcom/nocolor/internal/symbols"
)

// FileMeta describes the data to be cached.
type FileMeta struct {
	Functions *symbols.Functions
	Classes   *symbols.Classes
}

// NewFileMeta returns a new FileMeta instance with pre-allocated fields.
func NewFileMeta() FileMeta {
	return FileMeta{
		Functions: symbols.NewFunctions(),
		Classes:   symbols.NewClasses(),
	}
}
