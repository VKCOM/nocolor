# How to contribute

We are very glad that you want to participate in the development of this tool. There are several ways to help out:

- Create an [issue](https://github.com/vkcom/nocolor/issues/) on GitHub, if you have found a bug or have an idea for a feature;
- Write test cases for open bug issues;
- Write patches for open bug/feature issues.

There are a few guidelines that we need contributors to follow, so that we have a chance of keeping on top of things.

- The code must follow the Go coding standard. A static analyzer will help you with this, which can be launched with one command (see the testing section below).

- All commits should have the following structure:

  ```
  pkgs: short
  
  More complete description, if needed.
  ```

  Where: 

  - `pkgs` is the name of the package or packages separated by commas in which the change occurred;
  - `short` is a short description of changes starting with a lowercase letter and no punctuation marks at the end.

- All code changes should be covered by unit tests.

## A short description if you'd like to contribute by writing code

Below we describe how to build a project and test it.

### Build

The first thing you need is [Golang](https://golang.org/dl/). The project is using Go version 1.16, so make sure that the version you have installed is 1.16 or higher.

Next, clone the repository to the desired folder:

```sh
git clone https://github.com/vkcom/nocolor
```

Go to the project folder:

```sh
cd nocolor
```

And run the following command:

```sh
make build
```

It will build all sources into a single binary file and place it in the `build` folder.

### Testing

The project uses standard tests provided by Go. Run the following command to run the tests:

```sh
make test
```

It will run tests from the `./tests` folder.

All tests in the project are located in the `./tests` folder and are divided into two parts:

- The `rules` folder contains tests for rules that check the correctness of the errors for the code;

- The `edges` folder contains tests that check the correctness of the constructed call graphs.

#### Linter for source code

To keep the code clean and correct, we use the static analyzer (linter) [`golangci-lint`](https://github.com/golangci/golangci-lint). Its configuration file is located at the root of the project with `.golangci.yml` name.

To start the linter run the following command:

```sh
make lint
```

It will install the `golangci-lint` linter and run the analysis.

>  For convenience, there is a command `make check`, which runs the linter first, and then the tests.

### Release

We do not use complicated methods for releases. Each release is created manually:

- Update the version in the `cmd/main.go` file;
- Run the `make release` command, which creates archives with executable files for Linux, Windows, and macOS;
- Create a new release in GitHub, which describes the changes and uploads the archives;
- Release.

## Brief description of project folders

- `cmd` — folder with CLI-related code;
- `docs` — folder with documentation;
- `internal` — folder with implementation;
  - `callgraph` — description of the call graph for functions, as well as the function call stack;
  - `checkers` — description of checkers for PHP code, in fact, there is only one checker for checking PHPDoc tags to define colors for functions or classes;
  - `linttest` — functions for testing;
  - `palette` — description of the palette, colors, and config for the palette;
  - `pipes` — set of steps for analyzing call graphs;
  - `symbols` — description of the structure of the function for storage;
  - `walkers` — description of the walkers that bypass files, classes, functions, etc. in the form of AST and which collect all the information about which functions call which and vice versa, as well as the colors of these functions;
- `tests` — folder with tests;
- `_scripts` — folder with scripts for creating a release.

## Next steps

- [Description of the color concept](https://github.com/vkcom/nocolor/blob/master/docs/concept_of_colors.md)