package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestRecursive(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = `
-
  - "api has-curl": "dont call curl from api"
  - "api has-curl please": ""
`
	suite.AddFile(`<?php
function f1() { f2(); f3(); }
/**
 * @color has-curl
 */
function f2() { f10(); }
/** @color api */
function f3() { f10(); }
function f10() { f11(); f100(); }
function f11() { [1]; f12(); f13(); }
function f12() { [2]; f1(); }
function f13() { [3]; }
function f100() { [100]; f101(); }
function f101() { [101]; f102(); }
function f102() { [101]; f101(); }
f1();
	`)

	suite.Expect = []string{
		`
api has-curl => dont call curl from api
  This color rule is broken, call chain:
f3@api -> f10 -> f11 -> f12 -> f1 -> f2@has-curl
`,
	}

	suite.RunAndMatch()
}
