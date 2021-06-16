package pipes

import (
	"sort"
	"sync"

	"github.com/vkcom/nocolor/internal/callgraph"
)

// Async starts the passed number of workers to process the graphs passed to
// the channel, each of which is processed in the passed callback function.
func Async(workers int, input chan *callgraph.Graph, output chan []*ColorReport, cb func(*callgraph.Graph) []*ColorReport) {
	go func() {
		var wg sync.WaitGroup

		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func(id int) {
				var reports []*ColorReport

				for graph := range input {
					reports = append(reports, cb(graph)...)
				}

				output <- reports

				wg.Done()
			}(i)
		}
		wg.Wait()
		close(output)
	}()
}

// WriteGraphsAsync asynchronously writes the given graphs to the transferred channel.
func WriteGraphsAsync(graphs []*callgraph.Graph, graphsCh chan *callgraph.Graph) {
	go func() {
		for _, graph := range graphs {
			graphsCh <- graph
		}
		close(graphsCh)
	}()
}

// ReadReportsSync synchronously reads all reports from the channel,
// blocking the stream until the channel is closed.
func ReadReportsSync(output chan []*ColorReport) []*ColorReport {
	var allReports []*ColorReport
	for reports := range output {
		allReports = append(allReports, reports...)
	}

	sort.Slice(allReports, func(i, j int) bool {
		return allReports[i].Message < allReports[j].Message
	})

	return allReports
}
