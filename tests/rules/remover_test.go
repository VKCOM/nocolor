package rules

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestRemover(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/** @color remover */
function callAct() {
    if(0) f1();
    if(0) f2();
    if(0) f3();
    if(0) f4();
}

/** @color highload */
function f1() { g1(); }
function g1() { callAct(); }

/** @color highload */
function f2() { g2(); }
function g2() { callAct(); }

function f3() { g3(); }
/** @color no-highload */
function g3() { callAct(); }

function f4() { g4(); }
/** @color no-highload */
function g4() { callAct(); }

callAct();
`)

	suite.RunAndMatch()
}
