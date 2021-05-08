package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestDifferentPathsOk(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/** @color api */
function api1() { allow1(); }
/** @color api-allow-curl */
function allow1() { common(); }

/** @color api-allow-curl */
function allow2() { common(); }
/** @color api */
function api2() { allow2(); }

function common() { callurl(); }
/** @color has-curl */
function callurl() { [1]; }

api2();
api1();
`)

	suite.RunAndMatch()
}
