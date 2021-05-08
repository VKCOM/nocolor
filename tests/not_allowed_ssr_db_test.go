package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestNotAllowedSsrDb(t *testing.T) {
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
`)

	suite.Expect = []string{
		`
ssr has-db-access => Calling function working with the database in the server side rendering function
  This color rule is broken, call chain:
someSsr@ssr -> hasDbAccess@has-db-access
`,
	}

	suite.RunAndMatch()
}
