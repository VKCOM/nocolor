package rules

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestTwoErrorInSubTree(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color highload
 */
function f1() { f2(); /* error 1 */ }
/**
 * @color no-highload
 */
function f2() { f3(); }
/**
 * @color highload
 */
function f3() { f4(); /* error 2 (not displayed until error 1 is fixed) */ }
/**
 * @color no-highload
 */
function f4() { echo 1; }

f1();
`)

	suite.Expect = []string{
		`
highload no-highload => Calling no-highload function from highload function
  This color rule is broken, call chain:
f1@highload -> f2@no-highload
`,
	}

	suite.RunAndMatch()
}
