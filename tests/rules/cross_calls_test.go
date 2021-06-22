package rules

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestCrossCalls(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color highload
 */
function f1() { f2(); }
function f2() { f3(); f4(); f6(); }
function f3() { echo 1; }
function f4() { f5(); }

/**
 * @color no-highload
 */
function f5() { echo 1; }

/**
 * @color highload
 */
function f6() { echo 1; }

f1();
`)

	suite.Expect = []string{
		`
highload no-highload => Calling no-highload function from highload function
  This color rule is broken, call chain:
f1@highload -> f2 -> f4 -> f5@no-highload
`,
	}

	suite.RunAndMatch()
}
