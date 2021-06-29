package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestCloneMagicMethod(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color red
 */
class FooCloneable {
  public function __clone() {}
}

/**
 * @color red
 */
class FooNotCloneable {}

/**
 * @color green
 */
function f() {
  $obj1 = new FooCloneable;
  $obj2 = new FooNotCloneable;
  $obj11 = clone $obj1;
  $obj22 = clone $obj2;
}
`)

	suite.Expect = []string{
		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
f@green -> FooCloneable::__construct (default autogenerated)@red
`,
		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
f@green -> FooNotCloneable::__construct (default autogenerated)@red
`,
		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
f@green -> FooCloneable::__clone@red
`,
	}

	suite.RunAndMatch()
}

func TestInvokeMagicMethod(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
class Foo {
  /**
   * @color red
   */
  public function __invoke() {}
}

/**
 * @return Foo[]
 */
function getFooArray(): array {
  return [new Foo];
}

/**
 * @color green
 */
function f() {
  $obj = new Foo;
  echo $obj();
}

/**
 * @color green
 */
function g() {
  echo getFooArray()[0]();
}
`)

	suite.Expect = []string{
		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
f@green -> Foo::__invoke@red
`,

		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
g@green -> Foo::__invoke@red
`,
	}

	suite.RunAndMatch()
}
