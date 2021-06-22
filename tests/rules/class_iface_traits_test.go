package rules

import (
	"testing"

	"github.com/vkcom/nocolor/internal/linttest"
)

func TestClasses(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = `
-
  - "controller": ""

-
  - "service repository": "error"
`
	suite.AddFile(`<?php
/**
 * @color service
 */
class Service {
  public static function method() {
    Controller::method();
  }
}

/**
 * @color controller
 */
class Controller {
  public static function method() {
    Repository::method();
  }
}

/**
 * @color repository
 */
class Repository {
  public static function method() {
    echo 1;
  }
}
`)

	suite.Expect = []string{
		`
service repository => error
  This color rule is broken, call chain:
Service::method@service -> Controller::method -> Repository::method@repository
`,
	}

	suite.RunAndMatch()
}

func TestClassAndTrait(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = `
-
  - "controller": ""

-
  - "service repository": "error"
`
	suite.AddFile(`<?php
/**
 * @color service
 */
class Service {
  use SomeRepositoryTrait;

  public function method() {
    $this->traitMethod();
  }
}

/**
 * @color repository
 */
class Repository {}

/**
 * @color repository
 */
trait SomeRepositoryTrait {
  public function traitMethod() {
    echo 1;
  }
}
`)

	suite.Expect = []string{
		`
service repository => error
  This color rule is broken, call chain:
Service::method@service -> SomeRepositoryTrait::traitMethod@repository
`,
	}

	suite.RunAndMatch()
}

func TestClassAndInterface(t *testing.T) {
	suite := linttest.NewSuite(t)

	suite.Palette = `
-
  - "controller": ""

-
  - "service repository": "error"
`
	suite.AddFile(`<?php
/**
 * @color service
 */
class Service {
  public function method(IRepository $repository) {
    $repository->getData();
  }
}

/**
 * @color repository
 */
interface IRepository {
  public function getData();
}

class Repository implements IRepository {
  public function getData() {
    return 1;
  }
}
`)

	suite.Expect = []string{
		`
service repository => error
  This color rule is broken, call chain:
Service::method@service -> IRepository::getData@repository
`,
	}

	suite.RunAndMatch()
}
