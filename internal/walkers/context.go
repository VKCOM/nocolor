package walkers

import (
	"encoding/gob"
	"io"
	"log"

	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"

	"github.com/vkcom/nocolor/internal/symbols"
)

type GlobalContext struct {
	Info *meta.Info

	Functions *symbols.Functions
}

func NewGlobalContext(info *meta.Info) *GlobalContext {
	return &GlobalContext{
		Info:      info,
		Functions: symbols.NewFunctions(),
	}
}

// Version returns the current version of the cache.
func (ctx *GlobalContext) Version() string {
	return "1.0.0"
}

// Encode caches the data of one rootWalker of one file.
func (ctx *GlobalContext) Encode(writer io.Writer, checker linter.RootChecker) error {
	if ctx.Info.IsLoadingStubs() {
		return nil
	}

	ind := checker.(*RootIndexer)

	enc := gob.NewEncoder(writer)
	if err := enc.Encode(&ind.meta); err != nil {
		log.Printf("cache error: encode %s: %v", ind.ctx.Filename(), err)
		return err
	}

	return nil
}

// Decode recovers data from cache.
func (ctx *GlobalContext) Decode(r io.Reader, filename string) error {
	if ctx.Info.IsLoadingStubs() {
		return nil
	}

	var m FileMeta

	dec := gob.NewDecoder(r)
	if err := dec.Decode(&m); err != nil {
		log.Printf("cache error: decode %s: %v", filename, err)
		return err
	}

	ctx.UpdateMeta(&m, filename)

	return nil
}

// UpdateMeta recovers data by collecting it from each file.
func (ctx *GlobalContext) UpdateMeta(f *FileMeta, filename string) {
	for _, fun := range f.Functions.Raw() {
		ctx.Functions.Add(fun)
	}
}
