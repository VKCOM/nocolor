package rules

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestNotAllowedApiCurl(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color api
 */
function api() {
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
api@api -> hasCurl@has-curl
`,
	}

	suite.RunAndMatch()
}
