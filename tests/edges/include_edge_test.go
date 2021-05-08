package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestIncludeEdgesRelativePaths(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddNamedFile(`file1.php`, `<?php

require_once "./file3.php";
f3(); // error

/**
 * @color green
 */
function f1() {
    echo 1;
}

f1();

`)

	suite.AddNamedFile(`dir/file2.php`, `<?php
/**
 * @color red
 */
function f2() {
    require_once "../file1.php"; // error
}

f2();
`)

	suite.AddNamedFile(`file3.php`, `<?php
/**
 * @color internals
 */
function f3() {
    echo 1;
}
`)

	suite.Expect = []string{
		`
internals => call internals
  This color rule is broken, call chain:
file 'dir/file2.php' scope -> f2 -> file 'file1.php' scope -> f3@internals
`,
		`
red green => red green mixing
  This color rule is broken, call chain:
f2@red -> file 'file1.php' scope -> f1@green
`,
	}

	suite.RunAndMatch()
}
