package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestNotAllowedMore(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = `
-
  - "ssr": ""
  - "api": ""
-
  - "red green blue": "dont allow rgb nesting"
`
	suite.AddFile(`<?php

function init() { r(); }
init();

/**
 * @color red
 * @color api
 */
function r() { r2(); }
/** @color ssr */
function r2() { g(); b(); }

/** @color green */
function g() { g2(); }
function g2() { ssr1(); }

/**
  * @color blue
  */
function b() { b2(); }
function b2() { [1]; }

/** @color ssr */
function ssr1() { [1]; ssr2(); }
/** @color ssr */
function ssr2() { [1]; b(); }
`)

	suite.Expect = []string{
		`
red green blue => dont allow rgb nesting
  This color rule is broken, call chain:
r@red -> r2 -> g@green -> g2 -> ssr1 -> ssr2 -> b@blue
`,
	}

	suite.RunAndMatch()
}
