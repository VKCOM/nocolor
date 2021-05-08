package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestDifferentPaths(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/** @color api */
function api1() { allow1(); }
/** @color api-allow-curl */
function allow1() { common(); }

/** @color api-allow-curl */
function allow2() { api2(); }
/** @color api */
function api2() { common(); }

function common() { callurl(); }
/** @color has-curl */
function callurl() { [1]; }

allow2();
api1();
`)

	suite.Expect = []string{
		`
api has-curl => Calling curl function from API functions
  This color rule is broken, call chain:
api2@api -> common -> callurl@has-curl
`,
	}

	suite.RunAndMatch()
}
