package rules

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestNotAllowedMessageInternals(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color message-internals
 */
function messageInternals() {
    echo 1;
}

function someFunction() {
    messageInternals();
}

someFunction();
`)

	suite.Expect = []string{
		`
message-internals => Calling function marked as internal outside of functions with the color message-module
  This color rule is broken, call chain:
file '_file0.php' scope -> someFunction -> messageInternals@message-internals
`,
	}

	suite.RunAndMatch()
}
