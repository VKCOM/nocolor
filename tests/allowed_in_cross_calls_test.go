package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestAllowedInCrossCalls(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color ssr
 */
function f1() { f2(); f5(); }
/**
 * @color ssr-allow-db
 */
function f2() { f3(); f6(); }
/**
 * @color has-db-access
 */
function f3() { f4(); }

function f4() { echo 1; }
function f5() { f2(); }

/**
 * @color ssr
 */
function f6() { f7(); }
/**
 * @color has-db-access
 */
function f7() { echo 1; /* not an error, since above in the f2() function there was a color allowing this */ }

f1();
`)

	suite.RunAndMatch()
}
