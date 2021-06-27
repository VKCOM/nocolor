package rules

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestMixingClassAndMethodColors(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/** @color controller */
class Controller {
  /** @color allow-model-call */
  static function act() {
    Model::act();
  }
}

/** @color model */
class Model {
  static function act() {}
}
`)

	suite.Expect = []string{}
	suite.RunAndMatch()
}

func TestMixingTraitAndMethodColors(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/** @color controller */
trait Controller {
  /** @color allow-model-call */
  static function act() {
    Model::act();
  }
}

/** @color model */
class Model {
  static function act() {}
}
`)

	suite.Expect = []string{}
	suite.RunAndMatch()
}
