package rules

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestClassColorsEqualWithMethodColor(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/** @color controller */
class Controller {
  public static function act() {
    $m = new Model;
  }
}

/** @color model */
class Model {
  /** @color model */
  public function __construct() {}
  public function act() {}
}
`)

	suite.Expect = []string{
		`
controller model => restricted dependency
  This color rule is broken, call chain:
Controller::act@controller -> Model::__construct@model
`,
	}

	suite.RunAndMatch()
}
