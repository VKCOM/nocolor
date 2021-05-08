package cmd

import (
	"fmt"
	"log"

	"github.com/VKCOM/noverify/src/cmd"

	"github.com/vkcom/nocolor/internal/callgraph"
	"github.com/vkcom/nocolor/internal/palette"
	"github.com/vkcom/nocolor/internal/pipes"
	"github.com/vkcom/nocolor/internal/symbols"
	"github.com/vkcom/nocolor/internal/walkers"
)

type ExtraCheckFlags struct {
	PaletteSrc string
}

func Check(ctx *cmd.AppContext, globalContext *walkers.GlobalContext, pal *palette.Palette) (status int, err error) {
	flags := ctx.CustomFlags.(*ExtraCheckFlags)

	// We need to open the palette at the very beginning,
	// since during the analysis it should already be initialized,
	// since phpdoc colors is parsed during the function traversal.
	paletteFromFile, err := palette.OpenPaletteFromFile(flags.PaletteSrc)
	if err != nil {
		return 1, err
	}

	*pal = *paletteFromFile

	// If there are no arguments, then we interpret this as
	// an analysis of the current directory.
	if len(ctx.ParsedArgs) == 0 {
		ctx.ParsedArgs = append(ctx.ParsedArgs, "./")
	}

	// The main function for analyzing in NoVerify,
	// in it we collect all the functions of the project.
	status, err = cmd.Check(ctx)
	if err != nil {
		return status, err
	}

	// If the status is not zero, it means that there are
	// some errors at the stage of data collection.
	//
	// No further analysis is needed.
	if status != 0 {
		return status, nil
	}

	// Function that starts checking colors.
	reports := HandleFunctions(ctx, globalContext.Functions, pal)

	if len(reports) != 0 {
		for _, report := range reports {
			fmt.Println(report)
		}
		log.Printf("Found %d critical reports\n", len(reports))
		return 2, nil
	}

	log.Printf("No critical issues found. Your code is perfect.")
	return 0, nil
}

func HandleFunctions(ctx *cmd.AppContext, funcs *symbols.Functions, palette *palette.Palette) []*pipes.Report {
	workers := ctx.ParsedFlags.MaxConcurrency
	reportsCh := make(chan []*pipes.Report, 10)
	graphsCh := make(chan *callgraph.Graph, 10)

	nodes := pipes.FunctionsToNodes(funcs)
	graphs := pipes.NodesToGraphs(nodes)

	pipes.WriteGraphsAsync(graphs, graphsCh)

	pipes.Async(workers, graphsCh, reportsCh, func(graph *callgraph.Graph) []*pipes.Report {
		pipes.EraseNodesWithRemoverColor(graph)
		pipes.CalcNextWithColor(graph)

		reports := pipes.CheckColorsInGraph(graph, palette)
		return reports
	})

	return pipes.ReadReportsSync(reportsCh)
}
