package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
	"github.com/vkcom/nocolor/internal/palette"
)

func TestMasks(t *testing.T) {
	suite := linttest.NewSuite(t)

	defer func(count int) {
		palette.MaxColorsInMask = count
	}(palette.MaxColorsInMask)
	palette.MaxColorsInMask = 3

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
 * @color api
 */
function f3() { f4(); }
/**
 * @color has-curl
 */
function f4() { echo 1; }

/**
 * @color message-module
 */
function f5() { f6(); }
/**
 * @color message-internals
 */
function f6() { echo 1; }

/**
 * @color message-internals
 */
function f7() { echo 1; }

/**
 * @color highload
 */
function f8() {
    f9();
}

/**
 * @color ssr
 */
function f9() {
    f10();
}

/**
 * @color no-highload
 */
function f10() {
    echo 1;
}

f1();
f3();
f5();
f7();
f8();
`)

	suite.Expect = []string{
		`
highload no-highload => Calling no-highload function from highload function
  This color rule is broken, call chain:
f1@highload -> f2@no-highload
`,

		`
api has-curl => Calling curl function from API functions
  This color rule is broken, call chain:
f3@api -> f4@has-curl
`,

		`
message-internals => Calling function marked as internal outside of functions with the color message-module
  This color rule is broken, call chain:
file '_file0.php' scope -> f7@message-internals
`,

		`
highload no-highload => Calling no-highload function from highload function
  This color rule is broken, call chain:
f8@highload -> f9 -> f10@no-highload
`,
	}

	suite.RunAndMatch()
}
