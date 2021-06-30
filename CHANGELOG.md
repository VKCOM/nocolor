# Changelog

All notable changes to this project will be documented in this file, in reverse chronological order by release.

## `1.0.3` 2021-01-07

### Added

- [#14](https://github.com/VKCOM/nocolor/pull/14): Added `cache-clear` command.

### Fixed

- [`1a7ac`](https://github.com/VKCOM/nocolor/commit/1a7ac0f04f1abd89b272e2222a155af485f24524): Fixed panic if no arguments or commands were passed.

## `1.0.2` 2021-30-06

### Added

- [#8](https://github.com/VKCOM/nocolor/pull/11): Added support for the following magic methods:
  - `__clone`
  - `__invoke`
  - `__call`
  - `__callStatic`
  - `__get`
  - `__set`


### Changed

- [#11](https://github.com/VKCOM/nocolor/pull/11): Changed behavior, if the new operator is called with a variable that has a class type, then we assume that the constructor of this class is called;

- [#12](https://github.com/VKCOM/nocolor/pull/12): Changed behavior, if a method is called from a variable with several possible classes, then a connection will be created with the methods of all classes;

- [`3bf46`](https://github.com/VKCOM/nocolor/commit/3bf46ab1fcd773fc780873fa8dc6a9cdc0d7a937): Improved the output for the `version` command.

## `1.0.1` 2021-28-06

### Changed

- [#1](https://github.com/VKCOM/nocolor/pull/1): Changed function for recalculating masks for colors, which allowed to increase the speed by 2-3%.

### Fixed

- [#2](https://github.com/VKCOM/nocolor/issues/2): Fixed a bug when calling the new operator for a class that does not have an explicit constructor;
- [#3](https://github.com/VKCOM/nocolor/issues/3): Fixed color mixing for classes and methods.

## `1.0.0` 2021-27-06

First stable version.
