package pipes

import (
	"sort"

	"github.com/i582/cfmt/cmd/cfmt"

	"github.com/vkcom/nocolor/internal/callgraph"
	"github.com/vkcom/nocolor/internal/palette"
)

// CheckColorsInGraph check all palette rules like "api has-curl" or
// "messages-module messages-internals".
//
// Start from the root, make dfs expanding NextWithColors
// and check all rules for each dfs chain.
//
// These chains are left-to-right, e.g. colored f1 -> f2 -> f3,
// we find and check f1, then f1 -> f2, then f1 -> f2 -> f3.
func CheckColorsInGraph(graph *callgraph.Graph, palette *palette.Palette) []*ColorReport {
	roots := FindRootNodes(graph)

	var reports []*ColorReport
	checker := newCheckerFunctionsColors(graph, palette)
	for _, root := range roots {
		callstack := callgraph.NewCallstackOfColoredFunctions()
		reports = append(reports,
			checker.checkFuncDFS(callstack, graph, root)...,
		)
	}

	return reports
}

// FindRootNodes finds the nodes with the minimum number
// of input edges and has NextWithColors.
func FindRootNodes(graph *callgraph.Graph) callgraph.Nodes {
	funcs := graph.Functions

	if len(funcs) == 0 {
		return nil
	}

	sort.Slice(funcs, func(i, j int) bool {
		return len(funcs[i].Prev) < len(funcs[j].Prev)
	})

	minIndex := 0
	minEdges := len(funcs[0].Prev)
	for i, node := range funcs {
		if len(node.Prev) > minEdges {
			break
		}

		minIndex = i
	}
	if minIndex == len(funcs)-1 {
		return callgraph.Nodes{funcs[0]}
	}

	return funcs[:minIndex+1]
}

const maxShownErrorsCount = 10

type checkerFunctionsColors struct {
	callGraph *callgraph.Graph
	palette   *palette.Palette

	// shownErrors is used to prevent duplicate errors from being shown.
	shownErrors map[string]struct{}
}

func newCheckerFunctionsColors(callGraph *callgraph.Graph, palette *palette.Palette) *checkerFunctionsColors {
	return &checkerFunctionsColors{
		callGraph:   callGraph,
		palette:     palette,
		shownErrors: map[string]struct{}{},
	}
}

func (c *checkerFunctionsColors) checkFuncDFS(callstack *callgraph.CallstackOfColoredFunctions, graph *callgraph.Graph, node *callgraph.Node) (reports []*ColorReport) {
	callstack.Append(node)

	wasAnyError := false

	for _, ruleset := range c.palette.Rulesets {
		for i := len(ruleset) - 1; i >= 0; i-- {
			rule := ruleset[i]

			if !matchRule(callstack, rule) {
				continue
			}

			if rule.IsError() {
				report := c.errorOnRuleBroken(callstack, rule)
				if report != nil {
					reports = append(reports, report)
				}
				wasAnyError = true
			}
			break
		}
	}

	if !wasAnyError && node.NextWithColors != nil {
		for _, next := range *node.NextWithColors {
			if callstack.Size() < 50 && !callstack.Contains(next) {
				reps := c.checkFuncDFS(callstack, graph, next)
				reports = append(reports, reps...)
			}
		}
	}

	callstack.PopBack()
	return reports
}

func matchRule(callstack *callgraph.CallstackOfColoredFunctions, rule *palette.Rule) bool {
	matchMasks := callstack.ColorsMasks
	if len(matchMasks) == 0 || !rule.ContainsIn(matchMasks) {
		return false
	}

	return matchTwoVectors(rule.Colors, callstack.ColorsChain)
}

func matchTwoVectors(ruleChain, actualChain []palette.Color) bool {
	ruleIndex := len(ruleChain) - 1
	actualIndex := len(actualChain) - 1

	for {
		rightmostMatched := ruleChain[ruleIndex] == actualChain[actualIndex]

		if rightmostMatched {
			if ruleIndex == 0 {
				return true
			}
			if actualIndex == 0 {
				return false
			}
			ruleIndex--
			actualIndex--
		} else {
			if actualIndex == 0 {
				return false
			}
			actualIndex--
		}
	}
}

// On error (colored chain breaks some rule), we want to find an actual chain of calling.
// findCallstackBetweenTwoFunctionsBFS is launched only on error, that's why we don't care
// about performance and just use bfs.
func (c *checkerFunctionsColors) findCallstackBetweenTwoFunctionsBFS(from, target *callgraph.Node, shouldntAppear map[*callgraph.Node]struct{}) callgraph.Nodes {
	visitedLevel := map[*callgraph.Node]int{}
	var bfsQueue callgraph.Nodes

	bfsQueue = append(bfsQueue, from)
	visitedLevel[from] = 0

	for len(bfsQueue) != 0 {
		cur := bfsQueue[0]
		bfsQueue = bfsQueue[1:]
		nextLevel := visitedLevel[cur] + 1

		if cur == target {
			break
		}

		for _, called := range c.callGraph.Graph[cur] {
			if _, ok := visitedLevel[called]; ok {
				continue
			}
			if _, ok := shouldntAppear[called]; ok {
				continue
			}

			visitedLevel[called] = nextLevel
			bfsQueue = append(bfsQueue, called)
		}
	}

	// If couldn't find, just return [from, target].
	if _, ok := visitedLevel[target]; !ok {
		return callgraph.Nodes{from, target}
	}

	var callstack callgraph.Nodes
	callstack = append(callstack, target)

	for cur := target; cur != from; {
		prevLevel := visitedLevel[cur] - 1
		for _, callee := range c.callGraph.RevGraph[cur] {
			level, has := visitedLevel[callee]
			if has && level == prevLevel {
				callstack = append(callstack, callee)
				cur = callee
				break
			}
		}
	}

	revCallstack := make(callgraph.Nodes, 0, len(callstack))
	for i := len(callstack) - 1; i >= 0; i-- {
		revCallstack = append(revCallstack, callstack[i])
	}

	return revCallstack
}

func (c *checkerFunctionsColors) errorOnRuleBroken(callstack *callgraph.CallstackOfColoredFunctions, rule *palette.Rule) *ColorReport {
	if len(c.shownErrors) > maxShownErrorsCount {
		return nil
	}

	var fullCallstack callgraph.Nodes // Will be: src_main -> ... -> f1 -> ... -> f2 -> ... -> f3.
	vector := callstack.AsVector()    // f1, f2, f3 — all of them are colored, and their chain breaks the rule.
	for i := 0; i < callstack.Size()-1; i++ {
		cur := vector[i]
		next := vector[i+1]
		shouldntAppear := map[*callgraph.Node]struct{}{}

		if cur.NextWithColors == nil {
			continue
		}

		for _, exclude := range *cur.NextWithColors {
			if exclude != next && exclude != cur {
				shouldntAppear[exclude] = struct{}{}
			}
		}

		nextCallstackPart := c.findCallstackBetweenTwoFunctionsBFS(cur, next, shouldntAppear)
		fullCallstack = append(fullCallstack, nextCallstackPart[:len(nextCallstackPart)-1]...)
	}
	fullCallstack = append(fullCallstack, vector[len(vector)-1])

	// Having full callstack like "src_main -> main -> init -> apiFn@api -> ... -> curlFn@curl"
	// we want to show a slice only containing a subchain that breaks the rule.
	firstItemToShow := 0
	for !fullCallstack[firstItemToShow].Function.Colors.Contains(rule.Colors[0]) {
		firstItemToShow++
	}
	lastItemToShow := len(fullCallstack) - 1
	for !fullCallstack[lastItemToShow].Function.Colors.Contains(rule.Colors[len(rule.Colors)-1]) {
		lastItemToShow--
	}
	if firstItemToShow == lastItemToShow {
		firstItemToShow = 0
	}

	callChainToShow := fullCallstack[firstItemToShow : lastItemToShow+1]

	callstackStr := ""
	for i, node := range callChainToShow {
		callstackStr += cfmt.Sprintf("%s{{%s}}::cyan", node.Function.HumanReadableName(), node.Function.Colors.String(c.palette, rule.Masks))
		if i != len(callChainToShow)-1 {
			callstackStr += " -> "
		}
	}

	if _, ok := c.shownErrors[callstackStr]; ok {
		return nil
	}
	c.shownErrors[callstackStr] = struct{}{}

	message := cfmt.Sprintf("{{%s}}::cyan => {{%s}}::red\n  This color rule is broken, call chain:\n%s",
		rule.String(c.palette), rule.Error, callstackStr)

	return &ColorReport{
		Rule:      rule,
		CallChain: callChainToShow,
		Message:   message,
		Palette:   c.palette,
	}
}
