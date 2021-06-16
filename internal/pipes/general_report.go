package pipes

import (
	"os"
	"path/filepath"

	"github.com/VKCOM/noverify/src/linter"
	"github.com/i582/cfmt/cmd/cfmt"
)

type GeneralReport struct {
	linterReport *linter.Report
	colorReport  *ColorReport
	fullMessage  string

	Rule      string   `json:"rule"`
	CallChain []string `json:"call-chain"`
	Message   string   `json:"message"`
	Context   string   `json:"context"`
	File      string   `json:"file"`
	Line      int      `json:"line"`
}

func NewGeneralReportFromLinterReport(r *linter.Report) *GeneralReport {
	return &GeneralReport{
		linterReport: r,
		Message:      r.Message,
		Context:      r.Context,
		File:         r.Filename,
		Line:         r.Line,
	}
}

func NewGeneralReportFromColorReport(r *ColorReport) *GeneralReport {
	gr := &GeneralReport{
		colorReport: r,
		fullMessage: r.Message,
		Rule:        r.Rule.String(r.Palette),
		Message:     r.Rule.Error,
	}

	last := r.CallChain[len(r.CallChain)-1].Function
	gr.File = last.Pos.Filename

	for _, node := range r.CallChain {
		gr.CallChain = append(gr.CallChain, node.Function.HumanReadableName())
	}

	return gr
}

// String returns the string representation of the GeneralReport.
func (r *GeneralReport) String() string {
	if r.colorReport != nil {
		last := r.colorReport.CallChain[len(r.CallChain)-1].Function

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

`, path, last.Pos.Line, last.HumanReadableName(), r.fullMessage)
	}

	return cfmt.Sprintf(`~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
{{Error}}::red at the stage of collecting colors
   %s:%d
     {{%s}}::yellow

%s

`, r.File, r.Line, r.Context, r.Message)
}
