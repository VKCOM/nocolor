package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestComplexOk(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color ssr
 * @color message-module
 */
function functionWithSsrAndMessageModule() {
    functionWithMessageModule();
}

/**
 * @color message-module
 */
function functionWithMessageModule() {
    functionWithAllowDbAccessAndMessageInternals(); // ok, has allow color
}

/**
 * @color ssr-allow-db
 * @color has-db-access
 * @color message-internals
 */
function functionWithAllowDbAccessAndMessageInternals() {
    functionWithMessageInternals();
}

/**
 * @color message-internals
 */
function functionWithMessageInternals() {
    echo 1;
}

functionWithSsrAndMessageModule();
`)

	suite.RunAndMatch()
}
