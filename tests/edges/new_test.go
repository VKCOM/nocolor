package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestNewWithExpr(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/** @color red */
class Foo {}

/** @color red */
class Goo {}

/**
 * @color green
 * @param Foo|Goo|int $inst
 */
function f($inst) {
  $a = new $inst;
}

/**
 * @color green
 */
function g() {
  $inst = new Foo;
  $a = new $inst;
}

/**
 * @color green
 */
function e() {
  $inst = "Foo";
  $a = new $inst; // not working
}
`)

	suite.Expect = []string{
		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
f@green -> Foo::__construct (default autogenerated)@red
`,
		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
f@green -> Goo::__construct (default autogenerated)@red
`,
		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
g@green -> Foo::__construct (default autogenerated)@red
`,
	}

	suite.RunAndMatch()
}
