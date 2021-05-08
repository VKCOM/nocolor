package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestLongPath(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/** @color highload */
function fHigh() {
     shortAllowedPath();
     longDisallowedPath();
}
fHigh();

/** @color no-highload */
function fLow() { [1]; }

function shortAllowedPath() {
    s1();
}

/** @color allow-hh */
function s1() { fLow(); }

function longDisallowedPath() {
    l1();
}

function l1() { l2(); }
function l2() { l3(); }
function l3() { l4(); }
function l4() { fLow(); }
`)

	suite.Expect = []string{
		`
highload no-highload => Calling no-highload function from highload function
  This color rule is broken, call chain:
fHigh@highload -> longDisallowedPath -> l1 -> l2 -> l3 -> l4 -> fLow@no-highload
`,
	}

	suite.RunAndMatch()
}
