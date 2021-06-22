package rules

import (
	"testing"

	"github.com/i582/cfmt/cmd/cfmt"

	"github.com/vkcom/nocolor/internal/linttest"
)

const defaultPalette = `
-
  - "highload no-highload": "Calling no-highload function from highload function"
  - "highload allow-hh no-highload": ""
-
  - "ssr has-db-access": "Calling function working with the database in the server side rendering function"
  - "ssr ssr-allow-db has-db-access": ""
-
  - "api has-curl": "Calling curl function from API functions"
  - "api api-callback has-curl": ""
  - "api api-allow-curl has-curl": ""
-
  - "api2 has-curl": "Calling curl function from API2 functions"
  - "api2 api-callback has-curl": ""
  - "api2 api-allow-curl has-curl": ""
-
  - "message-internals": "Calling function marked as internal outside of functions with the color message-module"
  - "message-module message-internals": ""
`

func init() {
	cfmt.DisableColors()
}

func TestHighload(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color highload - optional comment
 * @return int
 */
function someHighload() {
    someNoHighload();
    return 2;
}

/**
 * @color no-highload
 */
function someA() {
    echo 1;
}

/**
 * @color no-highload
 */
function someNoHighload() {
    echo 1;
}

someHighload();
`)

	suite.Expect = []string{
		`
highload no-highload => Calling no-highload function from highload function
  This color rule is broken, call chain:
someHighload@highload -> someNoHighload@no-highload
`,
	}

	suite.RunAndMatch()
}
