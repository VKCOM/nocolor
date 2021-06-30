package tests

import (
	"testing"

	"github.com/i582/cfmt/cmd/cfmt"

	"github.com/vkcom/nocolor/internal/linttest"
)

const defaultPalette = `
red group:
  - red green: red green mixing
  - green red: calling red from green is prohibited
internals group:
  - internals: call internals
  - module internals: ""
`

func init() {
	cfmt.DisableColors()
}

func TestEdges(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php

#region Simple

/** @color red */
function f1() { 
  f2(); // error
  f10(); // undefined func, skip
}

/** @color green */
function f2() { echo 1; }

#endregion

#region Class

class Foo {
  /** @color red */
  public function method1() {
    $this->method2(); // error
    self::staticMethod1(); // error
    static::staticMethod1(); // error, equal with prev

    $this->method4(); // undefined method, skip
    self::staticMethod4(); // undefined method, skip
    static::staticMethod4(); // undefined method, skip
  }

  /** @color green */
  public function method2() { echo 1; }

  /** @color green */
  public function method3() { echo 1; }

  /** @color green */
  public static function staticMethod1() { echo 1; }

  /** @color green */
  public static function staticMethod2() { echo 1; }

  /** @color green */
  public static function staticMethod3() { echo 1; }
}

/** @color red */
function f3() {
  $foo = new Foo();
  $foo->method3(); // error
  Foo::staticMethod2(); // error
  $foo::staticMethod3(); // error

  $foo->method4(); // undefined method, skip
  Foo::staticMethod4(); // undefined method, skip
  $foo::staticMethod4(); // undefined method, skip
}

#endregion

#region ClassConstructor

class Boo {
  /** @color green */
  public function __construct() { echo 1; }
}

/** @color red */
function f4() {
  $_ = new Boo(); // error
  $_ = new Woo(); // undefined class, skip
}

#endregion

#region GlobalClassObjects

class Goo {
  /** @color internals */
  public function __construct() { echo 1; }

  /** @color internals */
  public function method1() { echo 1; }

  /** @color internals */
  public static function staticMethod1() { echo 1; }

  /** @color internals */
  public static function staticMethod2() { echo 1; }
}

$g = new Goo(); // error
$g->method1(); // error
Goo::staticMethod1(); // error
$g::staticMethod2(); // error

$g->method2(); // undefined method, skip
Goo::staticMethod3(); // undefined method, skip
$g::staticMethod3(); // undefined method, skip

#endregion


#region CallInArgs

class Doo {
  /** @color red */
  public static function staticMethod() {
    self::staticMethod2(self::staticMethod3()); // error
  }

  public static function staticMethod2($a) { echo 1; }

  /** @color green */
  public static function staticMethod3() { return 1; }

  /** @color red */
  public function method() {
    $this->method2(self::method3()); // error
  }

  public function method2($a) { echo 1; }

  /** @color green */
  public function method3() { return 1; }
}

/** @color red */
function f5() {
  f7(f6());
}

/** @color green */
function f6() {
  return 1;
}

function f7($a) {
  echo $a;
}

#endregion

`)

	suite.Expect = []string{
		`
red green => red green mixing
  This color rule is broken, call chain:
f1@red -> f2@green
`,

		`
red green => red green mixing
  This color rule is broken, call chain:
Foo::method1@red -> Foo::staticMethod1@green
`,

		`
red green => red green mixing
  This color rule is broken, call chain:
Foo::method1@red -> Foo::method2@green
`,

		`
red green => red green mixing
  This color rule is broken, call chain:
f4@red -> Boo::__construct@green
`,

		`
red green => red green mixing
  This color rule is broken, call chain:
f3@red -> Foo::method3@green
`,

		`
red green => red green mixing
  This color rule is broken, call chain:
f3@red -> Foo::staticMethod2@green
`,

		`
red green => red green mixing
  This color rule is broken, call chain:
f3@red -> Foo::staticMethod3@green
`,

		`
internals => call internals
  This color rule is broken, call chain:
file '_file0.php' scope -> Goo::__construct@internals
`,

		`
internals => call internals
  This color rule is broken, call chain:
file '_file0.php' scope -> Goo::method1@internals
`,

		`
internals => call internals
  This color rule is broken, call chain:
file '_file0.php' scope -> Goo::staticMethod1@internals
`,

		`
internals => call internals
  This color rule is broken, call chain:
file '_file0.php' scope -> Goo::staticMethod2@internals
`,

		`
red green => red green mixing
  This color rule is broken, call chain:
Doo::staticMethod@red -> Doo::staticMethod3@green
`,

		`
red green => red green mixing
  This color rule is broken, call chain:
Doo::method@red -> Doo::method3@green
`,

		`
red green => red green mixing
  This color rule is broken, call chain:
f5@red -> f6@green
`,
	}

	suite.RunAndMatch()
}

func TestMethodWithSeveralClasses(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
class Foo {
  /**
   * @color red
   */
  public function some() {}
}

class Zoo {}

class Goo {
  /**
   * @color red
   */
  public function some() {}
}

/**
 * @color green
 * @param Foo|Goo|Zoo|int $inst
 */
function f($inst) {
  $inst->some();
}
`)

	suite.Expect = []string{
		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
f@green -> Foo::some@red
`,
		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
f@green -> Goo::some@red
`,
	}

	suite.RunAndMatch()
}

func TestStaticMethodWithSeveralClasses(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
class Foo {
  /**
   * @color red
   */
  public static function someStatic() {}
}

class Zoo {}

class Goo {
  /**
   * @color red
   */
  public static function someStatic() {}
}

/**
 * @color green
 * @param Foo|Goo|Zoo|int $inst
 */
function f($inst) {
  $inst::someStatic();
}
`)

	suite.Expect = []string{
		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
f@green -> Foo::someStatic@red
`,
		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
f@green -> Goo::someStatic@red
`,
	}

	suite.RunAndMatch()
}

func TestConstructorWithSeveralClasses(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color red
 */
class Foo {
  public function __construct() {}
}

class Zoo {}

/**
 * @color red
 */
class Goo {}

/**
 * @color green
 * @param Foo|Goo|Zoo|int $inst
 */
function f($inst) {
  $a = new $inst;
}
`)

	suite.Expect = []string{
		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
f@green -> Foo::__construct@red
`,
		`
green red => calling red from green is prohibited
  This color rule is broken, call chain:
f@green -> Goo::__construct (default autogenerated)@red
`,
	}

	suite.RunAndMatch()
}
