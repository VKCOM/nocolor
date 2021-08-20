package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestPHP8NullsafeMethodCall(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
<?php

class Foo {
    /**
     * @color red
     */
    public function f() {}
}

/**
 * @color green
 */
function f($cond) {
    $a = new Foo;
    if ($cond) {
        $a = null;
    }
    $a?->f();
}
`)

	suite.Expect = []string{
		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
f@green -> Foo::f@red
`,
	}

	suite.RunAndMatch()
}
