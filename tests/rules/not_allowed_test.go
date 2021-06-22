package rules

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestNotAllowed(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color ssr
 */
function someSsr() {
    hasDbAccess();
}

/**
 * @color has-db-access
 */
function hasDbAccess() {
    echo 1;
}

someSsr();


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
ssr has-db-access => Calling function working with the database in the server side rendering function
  This color rule is broken, call chain:
someSsr@ssr -> hasDbAccess@has-db-access
`,
		`
api has-curl => Calling curl function from API functions
  This color rule is broken, call chain:
api@api -> hasCurl@has-curl
`,
	}

	suite.RunAndMatch()
}
