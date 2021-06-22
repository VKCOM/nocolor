package rules

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestErrorWithDifferentCaller(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color ssr
 */
function f1() { f2(); }
/**
 * @color ssr-allow-db
 */
function f2() { f6(); /* error */ }
/**
 * @color has-db-access
 */
function f3() { f4(); }

function f4() { f5(); }
function f5() { f6(); /* error */ }

/**
 * @color highload
 */
function f6() { f7(); }
/**
 * @color no-highload
 */
function f7() { echo 1; /* error that should not be displayed because the exact same error has already been displayed */ }

f1();
f3();
`)

	suite.Expect = []string{
		`
highload no-highload => Calling no-highload function from highload function
  This color rule is broken, call chain:
f6@highload -> f7@no-highload
`,
	}

	suite.RunAndMatch()
}
