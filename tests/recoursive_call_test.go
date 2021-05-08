package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestRecursiveCall(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
function withoutColor(int $a) {
    if ($a) {
        withoutColorToo($a);
    }
}

function withoutColorToo(int $a) {
    if ($a) {
        withoutColor($a);
    }
}

withoutColor((int)$_POST["id"]);

/**
 * @color highload
 */
function highload(int $a) {
    if ($a) {
        highloadToo($a);
    }

    echo 1;
}

/**
 * @color highload
 */
function highloadToo(int $a) {
    if ($a) {
        highload($a);
    }
}

highload((int)$_POST["id"]);
`)

	suite.RunAndMatch()
}
