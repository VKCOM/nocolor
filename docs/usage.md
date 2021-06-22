# Getting started

> Here and then we mean that the **NoColor** binary is available by name.

To analyze the project, just run the following command:

```sh
$ nocolor check --index-only-files='./vendor,./tests' ./src
```

Pay attention to the `--index-only-files` flag, it sets paths to **files that do not need to be analyzed**, but they are **needed for correct work** (for example, the `vendor` folder is needed so that the analyzer can find function definitions).

### Flags

- `--palette` — path to the file with the palette. The default is `palette.yaml`;
- `--tag` — the tag to be used to set the color in PHPDoc; The default is `color`;
- `--index-only-files` — comma-separated list of paths to files, which should be **indexed**, but **not analyzed**. It is used to specify folders that contain *definitions of functions, classes, etc.*, which is **important for correct analysis**, but the files themselves do not need to be analyzed;
- `--output` — path to the file where the errors will be written in JSON format;
- `--php-exts` — comma-separated list of PHP extensions to be analyzed. The default is `php, inc, php5, phtml`;
- `--cache-dir` — path to the directory where the cache will be saved. The cache allows you to analyze the project **much faster**, since a large stage of data collection is skipped. The default is `$TMP/nocolor-cache`;
- `--disable-cache` — flag to disable caching. The default is `false`;
- `--cores` — maximum number of threads for analysis. By default, it is equal to the number of processor threads;
- And some others for debugging.

To get all possible flags, run the command:

```sh
$ nocolor check help
```

After the flags, you can specify several folders or files for analysis.

For example:

```sh
$ nocolor check --cores=4 --output="reports.json" ./folder1 ./folder2 ./folder3/file.php
```

If no folders or files are specified, the current folder is analyzed. This is identical with the call with `./`:

```sh
$ nocolor check
# identically with
$ nocolor check ./
```

Usually, for analysis it is enough to pass the `src` folder and add the `vendor` folder to the `--index-only-files` flag:

```sh
$ nocolor check --index-only-files='./vendor' ./src
```

## Example run

To check that everything is working correctly, create a file `test.php` with the following content:

```php
<?php

/** 
 * @color highload
 */
function f1() {
  f2();
}

/** 
 * @color no-highload
 */
function f2() {
  f1();
}

f1();
```

And next to it is the `palette.yaml` file to define the palette:

```yaml
-
  - "highload no-highload": "Don't call a no-highload function from a highload one"
  - "highload allow-hh no-highload": ""
```

And in the folder with these files, run the following command:

```sh
$ nocolor check ./test.php
```

You should see an error:

```
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Error at the stage of checking colors
   test.php:13  in function f2

highload no-highload => Don't call a no-highload function from a highload one
  This color rule is broken, call chain:
f1@highload -> f2@no-highload
```

If so, then everything is working correctly.

Read the [colors concept description](https://github.com/vkcom/nocolor/blob/master/docs/concept_of_colors.md) to understand what happened here.

## Next steps

- [Description of the color concept](https://github.com/vkcom/nocolor/blob/master/docs/concept_of_colors.md)

