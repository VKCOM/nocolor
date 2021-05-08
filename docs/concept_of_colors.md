# Concept of function colors

The color concept consists of two parts: *colors* and a *mixing palette*.

Colors define the *properties* of functions, and a palette describes the *rules for mixing* colors. Color mixing here
means the **combination of colors of two or more functions that are called in a chain**.

So, for example, if there is a function `f1@red` with red color and a function `f2@green` with green color, then if the
function `f1` calls `f2`, then we say that two colors are mixed: `red` and `green`.

```php
/**
 * @color red
 */
function f1() { 
  f2(); // [red, green]
}

/**
 * @color green
 */
function f2() { echo 1; }
```

The blending palette allows us to set rules that *describe situations that we do not want to allow*. Let's say we want
to prohibit calling `green` functions from `red`, for this, the following rule is added to the palette:

```
"red green" => "it is forbidden to call green functions from red ones"
```

However, sometimes we need it, we do it with full understanding and we need a way to do it and not get an error.

For this, there are so-called **exclusion rules**, these rules specify the rule describing the error.

So, for example, we want to allow calling `green` functions from `red` ones, for this we will create the following rule
in addition to the first one:

```
"red allow-rg green" => ""
```

This rule contains an empty string on the right, which means that the rule describes an exception.

Now, if we want to allow calling `green` functions from `red` in someplace, it is enough to add the color `allow-rg` to
the function with `red` color. By this action, we **allow to do what the first rule prohibits**. Now we have
three `red allow-rg green` colors mixed, they match the exclusion rule, so there will be no error.

As a result, we get the following rules:

```
"red green"          => "it is forbidden to call green functions from red ones
"red allow-rg green" => ""
```

These two rules form a group of rules. Each rule group has one error rule, which **necessarily comes first**, and any
number of exclusion rules.

There can be several groups depending on the need. However, there cannot be more than **63** unique colors due to
implementation peculiarities.

## Special color `@remover`

In large projects like ours, there are functions that call a lot of functions and therefore, because of them, the check process can be very long, despite the optimization of the algorithm. In order to exclude such functions from the call graph, there is a special color `@remover`, which removes a function from the call graph.

## Colors and palette in code

### Colors

The phpdoc annotation `@color` is used to set colors, each annotation **must contain only one color** (all others will be ignored). Use multiple annotations to define multiple colors:

```php
/** 
 * @color red 
 * @color allow-rg - Optional comment.
 */
function f1() {}

/** 
 * @color green 
 */
function f2() {}
```

You can also set colors for all methods of a class/interface/trait using the annotation above the class/interface/trait:

```php
/**
 * @color service
 */
class Service {
  public function method() {}
  public function method2() {}
}
```

### Palette

The `yaml` or `json` config is used to set the palette. It is preferable to use `yaml`, its structure will be explained further. An example `json` can be viewed [here](https://github.com/vkcom/nocolor/blob/master/palette.json).

The structure of the `yaml` config is as follows:

```yaml
# In short, the structure consists of an array of groups, 
# where each group is an array of <colors, error value or empty string> objects.

# first group
-
  # array of rules
  # each rule is a separate object
  # the first rule must be the rule leading to an error
  - "highload no-highload": "Calling no-highload function from highload function"
  - "highload allow-hh no-highload": ""

# second group
-
  - "ssr has-db-access": "Calling function working with the database in the server side rendering function"
  - "ssr ssr-allow-db has-db-access": ""
```

## Examples close to reality

You can see examples of using the concept close to reality [here](https://github.com/vkcom/nocolor/blob/master/docs/examples_close_to_reality.md).

## Next steps

- [Comparison with Deptrac](https://github.com/vkcom/nocolor/blob/master/docs/nocolor_vs_deptrac.md)