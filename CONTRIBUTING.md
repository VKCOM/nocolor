# How to contribute

There are several ways to help out:

- create an [issue](https://github.com/vkcom/nocolor/issues/) on GitHub in case you have found a bug or have a feature request
- write test cases for open bug issues
- write patches for open bug/feature issues

There are a few guidelines that we ask contributors to observe:

- The code must follow the Go coding standard (checked by a linter, see below).

- All commits messages should be formatted as
  ```
  pkgs: short desc
  
  A more detailed description.
  ```
  where `pkgs` is the name of the package or comma-separated packages in which the change occurred.  

- All code changes should be covered by unit tests.


## A short description if you'd like to contribute by writing code

Below you'll find how to build a project and test it.

### Building

Make sure you have [Go](https://golang.org/dl/) version 1.16 or higher installed.  
Clone this repository and run `make build`:
```bash
git clone https://github.com/vkcom/nocolor
cd nocolor
make build
```

A resulting binary will be placed in the `./build` folder.

### Testing

The project uses standard tests provided by Go:
```bash
make test
```

It will run all tests from the `./tests` folder. 

Tests in the `./tests/rules` folder check the correctness of errors for PHP code.  
Tests in the `./tests/edges` folder check the call graph building process.

### Linting

We use [golangci-lint](https://github.com/golangci/golangci-lint). Its configuration file is located at `/.golangci.yml`.
```bash
make lint
```

This command will install the `golangci-lint` linter and run the analysis.

>  For convenience, there is a command `make check`, which runs the linter first, and then runs the tests.

### Releasing

We do not use complicated methods for releases. Each release is created manually:

- update the version the `Makefile`
- prepend changes to the `CHANGELOG.md`
- run the `make release` command, it will create archives with executable files in `release-v[version]`
- create a new release in GitHub with description, and upload the archives


## A brief description of project folders

- `cmd` — CLI-related code;
- `docs` — documentation;
- `internal` — implementation:
  - `callgraph` — description of the call graph for functions, as well as the function call stack;
  - `checkers` — description of checkers for PHP code, in fact, there is only one checker for checking PHPDoc tags to define colors for functions or classes;
  - `linttest` — functions for testing;
  - `palette` — description of the palette, colors, and config for the palette;
  - `pipes` — set of steps for analyzing call graphs;
  - `symbols` — description of the structure of the function for storage;
  - `walkers` — description of the walkers that traverse files, classes, functions, etc. in the form of AST and which collect all the information about which functions call which and vice versa, as well as the colors of these functions;
- `tests` — unit tests;
- `_scripts` — scripts for creating a release.

