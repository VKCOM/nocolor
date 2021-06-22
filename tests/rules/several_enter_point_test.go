package rules

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestSeveralEnterPoint(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color highload
 */
function f1() { f2(); }

/**
 * @color no-highload
 */
function f2() { echo 1; }

/**
 * @color highload
 */
function f3() { f2(); }


function f4() { f3(); }
`)

	suite.Expect = []string{
		`
highload no-highload => Calling no-highload function from highload function
  This color rule is broken, call chain:
f1@highload -> f2@no-highload
`,

		`
highload no-highload => Calling no-highload function from highload function
  This color rule is broken, call chain:
f3@highload -> f2@no-highload
`,
	}

	suite.RunAndMatch()
}
