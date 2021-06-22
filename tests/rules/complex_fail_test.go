package rules

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestComplexFail(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color ssr
 * @color message-module
 * @color highload
 */
function functionWithSsrAndMessageModuleAndHighload() {
    functionWithMessageModule();
}

/**
 * @color message-module
 */
function functionWithMessageModule() {
    functionWithAllowDbAccessAndNoHighload(); // error, call no-highload in highload
    functionWithDbAccessWithoutAllow(); // error, db is not allowed
}

/**
 * @color has-db-access
 * @color ssr-allow-db
 * @color no-highload
 */
function functionWithAllowDbAccessAndNoHighload() {
    echo 1;
}

/**
 * @color has-db-access
 */
function functionWithDbAccessWithoutAllow() {
    echo 1;
}

function dangerZoneCallingMessageInternals() {
    messageInternals();
}

/**
 * @color message-internals
 */
function messageInternals() {
    echo 1;
}


function main() {
    dangerZoneCallingMessageInternals();
    functionWithSsrAndMessageModuleAndHighload();
}

main();
`)

	suite.Expect = []string{
		`
highload no-highload => Calling no-highload function from highload function
  This color rule is broken, call chain:
functionWithSsrAndMessageModuleAndHighload@highload -> functionWithMessageModule -> functionWithAllowDbAccessAndNoHighload@no-highload
`,

		`
message-internals => Calling function marked as internal outside of functions with the color message-module
  This color rule is broken, call chain:
file '_file0.php' scope -> main -> dangerZoneCallingMessageInternals -> messageInternals@message-internals
`,

		`
ssr has-db-access => Calling function working with the database in the server side rendering function
  This color rule is broken, call chain:
functionWithSsrAndMessageModuleAndHighload@ssr -> functionWithMessageModule -> functionWithAllowDbAccessAndNoHighload@has-db-access
`,

		`
ssr has-db-access => Calling function working with the database in the server side rendering function
  This color rule is broken, call chain:
functionWithSsrAndMessageModuleAndHighload@ssr -> functionWithMessageModule -> functionWithDbAccessWithoutAllow@has-db-access
`,
	}

	suite.RunAndMatch()
}
