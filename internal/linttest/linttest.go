package linttest

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"

	"github.com/VKCOM/noverify/src/cmd"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/linttest"
	"github.com/VKCOM/noverify/src/workspace"

	cmdp "github.com/vkcom/nocolor/cmd"
	"github.com/vkcom/nocolor/internal/palette"
	"github.com/vkcom/nocolor/internal/pipes"
	"github.com/vkcom/nocolor/internal/walkers"
)

// Suite is a configurable test runner for linter.
//
// Use NewSuite to create usable instance.
type Suite struct {
	t testing.TB

	Palette string
	Files   []linttest.TestFile
	Expect  []string

	config *linter.Config
	linter *linter.Linter
}

// NewSuite returns a new linter test suite for t.
func NewSuite(t testing.TB) *Suite {
	conf := linter.NewConfig("8.1")
	return &Suite{
		t:      t,
		config: conf,
		linter: linter.NewLinter(conf),
	}
}

// AddFile adds a file to a suite file list.
// File gets an auto-generated name. If custom name is important,
// use AddNamedFile.
func (s *Suite) AddFile(contents string) {
	s.Files = append(s.Files, linttest.TestFile{
		Name: fmt.Sprintf("_file%d.php", len(s.Files)),
		Data: []byte(contents),
	})
}

// AddNamedFile adds a file with a specific name to a suite file list.
func (s *Suite) AddNamedFile(name, contents string) {
	s.Files = append(s.Files, linttest.TestFile{
		Name: name,
		Data: []byte(contents),
	})
}

// RunLinter executes linter over s Files and returns all issue reports
// that were produced during that.
func (s *Suite) RunLinter() []*pipes.ColorReport {
	s.t.Helper()

	globalContext := walkers.NewGlobalContext(s.linter.MetaInfo())
	pal := palette.NewPalette()
	walkers.Register(s.config, globalContext, pal, "color")

	var err error
	paletteFromFile, err := palette.ReadPaletteFileYAML("palette.yaml", []byte(s.Palette))
	if err != nil {
		s.t.Fatalf("%v", err)
	}

	*pal = *paletteFromFile

	indexing := s.linter.NewIndexingWorker(0)

	shuffleFiles(s.Files)
	for _, f := range s.Files {
		parseTestFile(s.t, indexing, f)
	}

	s.linter.MetaInfo().SetIndexingComplete(true)

	linting := s.linter.NewLintingWorker(0)

	shuffleFiles(s.Files)
	for _, f := range s.Files {
		if f.Nolint {
			// Mostly used to add builtin definitions
			// and for other kind of stub code that was
			// inserted to make actual testing easier (or possible, even).
			continue
		}

		parseTestFile(s.t, linting, f)
	}

	reports := cmdp.HandleFunctions(&cmd.AppContext{
		ParsedFlags: cmd.ParsedFlags{
			MaxConcurrency: 1,
		},
	}, globalContext.Functions, pal)

	return reports
}

// RunAndMatch calls Match with the results of RunLinter.
//
// This is a recommended way to use the Suite, but if
// reports slice is needed, one can use RunLinter directly.
func (s *Suite) RunAndMatch() {
	s.t.Helper()
	s.Match(s.RunLinter())
}

// Match tries to match every report against Expect list of s.
//
// If expect slice is nil or empty, only nil (or empty) reports
// slice would match it.
func (s *Suite) Match(reports []*pipes.ColorReport) {
	expect := s.Expect
	t := s.t

	for i := range expect {
		expect[i] = strings.TrimPrefix(expect[i], "\n")
		expect[i] = strings.TrimSuffix(expect[i], "\n")
	}

	t.Helper()

	if len(reports) != len(expect) {
		t.Errorf("unexpected number of reports: expected %d, got %d",
			len(expect), len(reports))
	}

	matchedReports := map[*pipes.ColorReport]bool{}
	usedMatchers := map[int]bool{}
	for _, r := range reports {
		have := r.Message
		for i, want := range expect {
			if usedMatchers[i] {
				continue
			}
			if strings.Contains(have, want) {
				matchedReports[r] = true
				usedMatchers[i] = true
				break
			}
		}
	}
	for i, r := range reports {
		if matchedReports[r] {
			continue
		}
		t.Errorf("unexpected report %d: %s", i, r.Message)
	}
	for i, want := range expect {
		if usedMatchers[i] {
			continue
		}
		t.Errorf("pattern %d matched nothing: %s", i, want)
	}

	// Only print all reports if test failed.
	if t.Failed() {
		t.Log(">>> issues reported:")
		for _, r := range reports {
			t.Log(r.Message)
		}
		t.Log("<<<")
	}
}

func init() {
	var once sync.Once
	once.Do(func() { go linter.MemoryLimiterThread(0) })
}

func shuffleFiles(files []linttest.TestFile) {
	rand.Shuffle(len(files), func(i, j int) {
		files[i], files[j] = files[j], files[i]
	})
}

func parseTestFile(t testing.TB, worker *linter.Worker, f linttest.TestFile) {
	file := workspace.FileInfo{
		Name:     f.Name,
		Contents: f.Data,
	}

	var err error
	if worker.MetaInfo().IsIndexingComplete() {
		_, err = worker.ParseContents(file)
	} else {
		err = worker.IndexFile(file)
	}
	if err != nil {
		t.Fatalf("could not parse %s: %v", f.Name, err.Error())
	}
}
