package rules

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestAllowed(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = defaultPalette
	suite.AddFile(`<?php
/**
 * @color ssr
 */
function someSsr() {
    hasDbAccess();
}

/**
 * @color ssr-allow-db
 * @color has-db-access
 */
function hasDbAccess() {
    echo 1;
}

someSsr();


/**
 * @color api
 */
function api() {
    hasCurl();
}

/**
 * @color api-allow-curl
 * @color has-curl
 */
function hasCurl() {
    echo 1;
}

api();


/**
 * @color api
 */
function apiCallApiCallback() {
    apiCallback();
}

/**
 * @color api-callback
 */
function apiCallback() {
    hasCurl();
}

apiCallApiCallback();


/**
 * @color api
 * @color api-callback
 */
function apiWithApiCallback() {
    hasCurl();
}

apiWithApiCallback();


/**
 * @color message-internals
 */
function messageInternals() {
    echo 1;
}

/**
 * @color message-module
 */
function messageModule() {
    messageInternals();
}

messageModule();
`)

	suite.RunAndMatch()
}
