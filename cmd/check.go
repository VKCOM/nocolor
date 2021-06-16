package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/VKCOM/noverify/src/cmd"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/vkcom/nocolor/internal/callgraph"
	"github.com/vkcom/nocolor/internal/palette"
	"github.com/vkcom/nocolor/internal/pipes"
	"github.com/vkcom/nocolor/internal/symbols"
	"github.com/vkcom/nocolor/internal/walkers"
)

type extraCheckFlags struct {
	PaletteSrc string
	ColorTag   string
	Output     string
}

// Check is the function that starts the analysis of the project.
func Check(ctx *cmd.AppContext, globalContext *walkers.GlobalContext) (status int, err error) {
	flags := ctx.CustomFlags.(*extraCheckFlags)

	pal, err := palette.OpenPaletteFromFile(flags.PaletteSrc)
	if err != nil {
		return 1, fmt.Errorf("error open palette file '%s': %v", flags.PaletteSrc, err)
	}

	// Registering custom walkers for collecting the call graph.
	walkers.Register(ctx.MainConfig.LinterConfig, globalContext, pal, flags.ColorTag)

	// If there are no arguments, then we interpret this as
	// an analysis of the current directory.
	if len(ctx.ParsedArgs) == 0 {
		ctx.ParsedArgs = append(ctx.ParsedArgs, "./")
	}

	// The main function for analyzing in NoVerify,
	// in it we collect all the functions of the project.
	_, err = cmd.Check(ctx)
	if len(LinterReports) != 0 {
		HandleShowLinterReports(ctx, LinterReports)
		return 2, err
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
		HandleShowColorReports(ctx, reports)
		return 2, nil
	}

	log.Printf("No critical issues found. Your code is perfect.")
	return 0, nil
}

// HandleFunctions is a function that starts checking colors.
func HandleFunctions(ctx *cmd.AppContext, funcs *symbols.Functions, palette *palette.Palette) []*pipes.ColorReport {
	workers := ctx.ParsedFlags.MaxConcurrency
	reportsCh := make(chan []*pipes.ColorReport, 10)
	graphsCh := make(chan *callgraph.Graph, 10)

	nodes := pipes.FunctionsToNodes(funcs)
	graphs := pipes.NodesToGraphs(nodes)

	pipes.WriteGraphsAsync(graphs, graphsCh)

	pipes.Async(workers, graphsCh, reportsCh, func(graph *callgraph.Graph) []*pipes.ColorReport {
		pipes.EraseNodesWithRemoverColor(graph)
		pipes.CalcNextWithColor(graph)

		reports := pipes.CheckColorsInGraph(graph, palette)
		return reports
	})

	return pipes.ReadReportsSync(reportsCh)
}

func HandleShowColorReports(ctx *cmd.AppContext, reports []*pipes.ColorReport) {
	generalReports := make([]*pipes.GeneralReport, 0, len(reports))
	for _, report := range reports {
		generalReports = append(generalReports, pipes.NewGeneralReportFromColorReport(report))
	}
	handleShowReports(ctx, generalReports)
}

func HandleShowLinterReports(ctx *cmd.AppContext, reports []*linter.Report) {
	generalReports := make([]*pipes.GeneralReport, 0, len(reports))
	for _, report := range reports {
		generalReports = append(generalReports, pipes.NewGeneralReportFromLinterReport(report))
	}
	handleShowReports(ctx, generalReports)
}

func handleShowReports(ctx *cmd.AppContext, reports []*pipes.GeneralReport) {
	flags := ctx.CustomFlags.(*extraCheckFlags)
	toJSON := flags.Output != ""
	if toJSON {
		fileName := flags.Output
		data, err := json.Marshal(reports)
		if err != nil {
			log.Printf("Error marshal json: %v", err)
			return
		}
		err = ioutil.WriteFile(fileName, data, 0644)
		if err != nil {
			log.Printf("Error write json: %v", err)
			return
		}

		log.Printf("Found %d critical reports\n", len(reports))
		log.Printf("Reports are written to the '%s' file\n", fileName)
		return
	}

	for _, report := range reports {
		fmt.Println(report)
	}
	log.Printf("Found %d critical reports\n", len(reports))
}
