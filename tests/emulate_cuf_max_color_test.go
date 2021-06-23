package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
	"github.com/vkcom/nocolor/internal/palette"
)

func TestEmulateCufMaxColor(t *testing.T) {
	suite := linttest.NewSuite(t)

	defer func(count int) {
		palette.MaxColorsInMask = count
	}(palette.MaxColorsInMask)
	palette.MaxColorsInMask = 2

	suite.Palette = `
-
  - "api has-curl": "dont call curl from api"
  - "api has-curl please": ""
`
	suite.AddFile(`<?php
function cuf($name) {
	if(0) api_1();
	if(0) api_2();
	api_3(); api_4(); api_5(); api_6();
	f1(); f2(); f3(); f4(); f5(); f6();
}

/** @color api */
function api_1() { f1(); f2(); f3(); api_2(); }
function f1() { cuf('api_1'); }

/** @color api */
function api_2() { f2(); f1(); }
function f2() { cuf('api_2'); }

/** @color api */
function api_3() { f3(); f4(); f1();  }
function f3() { cuf('api_3'); }

/** @color api */
function api_4() { f4();  }
/** @color has-curl */
function f4() { cuf('api_3'); }

/** @color api */
function api_5() { f5();  }
function f5() { cuf('api_3'); }

/** @color api */
function api_6() { f6();  }
function f6() { cuf('api_3'); }

cuf('api_1');
`)

	suite.Expect = []string{
		`
api has-curl => dont call curl from api
  This color rule is broken, call chain:

`,
	}

	suite.RunAndMatch()
}
