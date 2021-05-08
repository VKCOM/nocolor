package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestAlmostSame(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color api
 */
function api() {
    api2();
    hasCurl();
}

/**
 * @color api2
 */
function api2() {
    hasCurl();
}

/**
 * @color has-curl
 */
function hasCurl() {
    echo 1;
}

api();
`)

	suite.Expect = []string{
		`
api has-curl => Calling curl function from API functions
  This color rule is broken, call chain:
api@api -> api2 -> hasCurl@has-curl
`,
		`
api has-curl => Calling curl function from API functions
  This color rule is broken, call chain:
api@api -> hasCurl@has-curl
`,
		`
api2 has-curl => Calling curl function from API2 functions
  This color rule is broken, call chain:
api2@api2 -> hasCurl@has-curl
`,
	}

	suite.RunAndMatch()
}
