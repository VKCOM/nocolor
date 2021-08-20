# NoColor config and options

This page is dedicated to some technical details.

## How to parse as PHP 7

It looks like this:

```bash
nocolor check --php7 ./src
```

By default, all code is parsed as PHP 8, however, some projects use names that have become reserved in PHP 8, so they need to be parsed as PHP 7.

## How to exclude some folders from checking

It looks like this:

```bash
nocolor check --index-only-files='./src/tests' ./src
```

The `--index-only-files` option sets paths that won't be analyzed, they will be just indexed.


## How to include the `vendor` folder

>  Since version 1.1, the `vendor` folder is added by default if it exists.

Using the same option:

```bash
nocolor check --index-only-files='./vendor' ./src
```

This option will add a folder to be indexed. This could be useful for the analyzer to find external function declarations for type inferring, for example.


## Console options

A full launch command line is
```bash
nocolor check --option1=xxx --option2=yyy ... [folder_or_file] [folder_or_file] ...
```

When no options are specified, their default values are used.

When no folders and files are specified, the current directory `.` is assumed.

Below is the list of available options, which can also be found with
```bash
nocolor check help
```

- `--palette` — a path to the file with the palette; by default, `palette.yaml`
- `--tag` — a PHPDoc color tag name; by default, `color`
- `--index-only-files` — a comma-separated list of paths to files, which should be **indexed**, but **not analyzed**, see the section above; by default, empty
- `--output` — a path to the file where the errors will be written in JSON format instead of printing a human-readable message to console; by default, empty
- `--php-exts` — a comma-separated list of PHP extensions to be analyzed; by default, `php, inc, php5, phtml`
- `--cache-dir` — a path to the directory where the cache will be saved between runs; by default, `$TMP/nocolor-cache`
- `--disable-cache` — a flag to disable caching; by default, `false`
- `--cores` — the maximum number of threads for analysis; by default, the number of CPUs


## Format of the `palette.yaml` file

**Tip**. The `nocolor init` command creates a `palette.yaml` file with some examples to do by analogy with.

The format looks like this:
```yaml
first group title:
- rule
- rule
# other rules in this group

# optionally, there may be many groups
```

Rules are key-value pairs `color pattern: error or an empty string`.

Here is a working example. Two rulesets, with two rules in each:
```yaml
finding performance leaks:
- fast slow: potential performance leak 
- fast slow-ignore slow: ""

# you can add optional comments in yaml
preventing data fetching from ssr:  
- ssr db: don't fetch data from templates 
- ssr allow-db db: ""
```

In .yaml syntax, it's a map from a string key (ruleset description) to a list (rules of that ruleset). 

On the top-level, there are one or more rulesets. Each ruleset is a list of color mixing rules. Each rule is a string key, and a string value. A key is a color pattern to be matched against every call chain. A value is an error message, or an empty string, meaning there is no error.

The necessity of having multiple rulesets comes from [conflict resolution](/docs/introducing_colors.md#conflict-resolution-if-many-rules-are-matched). Making every ruleset a list is also important, because the order of rules is meaningful (in case of a sub-object syntax (without dashes), that order would be missed).

Since it's a .yaml format, you can use quotes if you prefer. 

*Transparent* color and wildcard can't occur in selectors.

**You can't use colors missing in the palette**. This restriction is on purpose: it prevents you from occasional misprints in `@color` doc tags. It means, that before using a new color, you must add a rule with it. If you aren't ready to define a sensible rule yet, you should at least write a "declaration rule" like `color-for-the-future: ""`
