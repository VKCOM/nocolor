<img width="350" alt="image" src="https://user-images.githubusercontent.com/51853996/122410502-a543ff00-cf8c-11eb-9b23-6b0c6e900f1e.png">

[![](https://github.com/vkcom/nocolor/workflows/Go/badge.svg)](https://github.com/vkcom/nocolor/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/VKCOM/nocolor)](https://goreportcard.com/report/github.com/vkcom/nocolor) [![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)


# NoColor — validate the architecture of your PHP project

NoColor is an **architecture validation tool** based on the [*concept of colored functions*](/docs/introducing_colors.md). This concept was originally invented for [KPHP](https://github.com/VKCOM/kphp) and later exposed as a separate tool to be used in regular PHP projects.

NoColor and [Deptrac](https://github.com/qossmic/deptrac) have similar goals: they both search for bad architectural patterns. But NoColor uses an absolutely different scenario: function coloring. Unlike Deptrac, NoColor analyzes **call chains of any depth** and supports type inferring. Continue reading [NoColor vs Deptrac](/docs/comparison_with_deptrac.md).

NoColor performs **static code analysis**, it has no runtime overhead. It supports type inferring to handle instance method calls, but dynamic invocations like `A::$method()` can't be statically analyzed.

NoColor is built on top of [NoVerify](https://github.com/VKCOM/noverify) and **written in Go**.

Optionally, install an experimental [plugin](https://github.com/i582/nocolor-phpstorm) for **PhpStorm** that brings some handy features.

Theoretically, the same concept can be implemented for almost every language. NoColor targets PHP.


## A brief intro to colors

Using the `@color` PHPDoc tag, you assign **colors** to a function. When a function calls another one, their colors combine to a chain:
```php
/** @color green */
function f1() { f2(); }

// this function has no color (it's transparent)
function f2() { f3(); }

/** @color red */
function f3() { /* ... */ }
```

<p align="center">
    <img src="/docs/img/f1-f2-f3-g-t-r.png" alt="f1 f2 f3 colored" height="75">
</p>

In **the palette**, you define rules of colors mixing, as this error rule:
```yaml
green red: calling red from green is prohibited
```

All possible call chains are pattern-matched against rules. Hence, `f1 -> f2 -> f3` will trigger the error above.

A color is anything you want: `@color api`, or `@color model`, or `@color controller`. With the palette, you define arbitrary patterns and exceptions. You can mark classes and use namespaces. You can express modularity and even simulate the `internal` keyword.

[Continue reading about colors here](/docs/introducing_colors.md)


## Getting started

The [Getting started](/docs/getting_started.md) page contains a step-by-step guide and copy-paste examples.


## The palette.yaml file

Once you call
```bash
nocolor init
```
at the root of your project, it creates an example `palette.yaml` file. It contains rules of color mixtures: 

```yaml
first group title:
- rule1 color pattern: rule1 error or nothing
- rule2 color pattern: rule2 error or nothing
# other rules in this group

# optionally, there may be many groups
```

Consider the [Configuration](/docs/configuration.md#format-of-the-paletteyaml-file) page section.


## Running NoColor

At first, [**install**](/docs/install.md) NoColor to your system. The easiest way is just to download a ready binary.

Then, execute this command **once** at the root of your project:
```bash
nocolor init
```
It will create a `palette.yaml` file with some examples.

Every time you need to **check a project**, run
```bash
nocolor check
```
to perform checking in the current directory, or
```bash
nocolor check ./src
```
to perform checks in another folder (or many).

To exclude some paths from analyzing, or to include the `./vendor` dir, consider all possible [command-line options](/docs/configuration.md).


## Limitations and speed

As for now, NoColor supports PHP 7.4 language level. It depends on a Go package [php-parser](https://github.com/z7zmey/php-parser) which is currently frozen. This restriction can be overcome in the future.

NoColor scales easily to the capacities provided. Depending on the number of cores and the speed of a hard disk, NoColor can process up to 300k lines per second.

The number of groups, colors, and selectors is unlimited, though the more functions are colored — the slower NoColor would work, as the number of colored graph paths **exponentially increases**. Typically, 99% of classes/functions are supposed to be left transparent.   

All available call chains are calculated **on a static analysis phase**, there is no runtime overhead. NoColor uses some tricky internal optimizations to avoid useless depth searching in a call graph. Every possible colored call chain is matched against all rules in the palette.

Remember, that **PHP is an interpreted language** and allows constructions that can't be statically analyzed. If you write something like `SomeClass::$any_function()` or `new $class_name`, NoColor can't do anything about it.


## Contributing

Feel free to contribute to this project. See [CONTRIBUTING.md](/CONTRIBUTING.md) for more information.


## The License

NoColor is distributed under the MIT License, on behalf of VK.com (V Kontakte LLC).
