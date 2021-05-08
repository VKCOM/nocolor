package tests

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestResumableFail(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = `
-
  - "api has-curl": "curl from api"
`
	suite.AddFile(`<?php
function init() {
    $ii = fork(api());
    wait($ii);
}

/** @color api */
function api() {
    sched_yield();
    $i = fork(resum());
    wait($i);
    return null;
}

function resum() {
    sched_yield();
    resum2();
    return null;
}

function init2() {
    $i = fork(resum2());
    wait($i);
}

function resum2() {
    sched_yield();
    callurl();
    return null;
}

/** @color has-curl */
function callurl() { sched_yield(); 1; }

init();
init2();
`)

	suite.Expect = []string{
		`
api has-curl => curl from api
  This color rule is broken, call chain:
api@api -> resum -> resum2 -> callurl@has-curl
`,
	}

	suite.RunAndMatch()
}
