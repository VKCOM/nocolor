package pipes

import (
	"os"
	"path/filepath"

	"github.com/i582/cfmt/cmd/cfmt"

	"github.com/vkcom/nocolor/internal/callgraph"
	"github.com/vkcom/nocolor/internal/palette"
)

type Report struct {
	Rule      *palette.Rule
	CallChain callgraph.Nodes
	Message   string
}

func (r *Report) String() string {
	last := r.CallChain[len(r.CallChain)-1].Function

	path := last.Pos.Filename
	wd, err := os.Getwd()
	if err == nil {
		path, err = filepath.Rel(wd, path)
		if err != nil {
			path = last.Pos.Filename
		}
	}

	path = filepath.ToSlash(path)

	return cfmt.Sprintf(`~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
{{Error}}::red at the stage of checking colors
   %s:%d  in function {{%s}}::yellow

%s

`, path, last.Pos.Line, last.HumanReadableName(), r.Message)

}
