# Concept of function colors

The concept of function colors is a way of describing rules for function calls from other functions based on the analysis of function call chains.

The color concept consists of two parts: **colors** and a **mixing palette**.

- **Colors** define the **properties** of functions.
- **Palette** describes the **rules** for mixing colors. 

Color mixing here means the **combination of colors of two or more functions that are called in a chain**.

So, for example, if there is a function `f1` with `red` color and a function `f2` with `green` color, then if the function `f1` calls `f2`, then we say that two colors are mixed: `red` and `green`.

```php
/**
 * @color red
 */
function f1() { 
  f2(); // color mixing: [red, green]
}

/**
 * @color green
 */
function f2() { echo 1; }
```

Thanks to the color palette, we can set rules that prohibit certain color mixing. In fact, this means that we forbid calling some functions from others.

Let's say we want to prohibit calling `green` functions from `red`.

To do this, we will add the following rule to the palette:

```php
"red green" => "it is forbidden to call green functions from red ones"
```

> Please note that on the left we write colors separated by a space, the mixing of which we want to prohibit, and on the right, an error message.

Thus, if we now try to call the `green` function from the `red` one, we will receive an **error**.

> Please note, however, that if we call the `red` function from the `green` one, there will be **no error**!

### Transparent color

The color concept has a special `transparent` color, which implicitly has all the functions without colors.

Transparent functions are skipped when checking, that is, this means that if there are 50 transparent functions between two color functions, then this does not affect anything, and if mixing is prohibited for these two color functions, then this will be found.

For example:

```php
/**
 * @color red
 */
function f1() { 
  f3(); // error
}

function f3() { // function has no colors, so it is transparent
  f4();
}

function f4() { // function has no colors, so it is transparent
  f2();
}

/**
 * @color green
 */
function f2() { echo 1; }
```

Despite the two transparent functions, the `red` function calls the `green` function, which is prohibited by the rule above.

### Exclusion rules

Even though some mixing is prohibited, we may need it, and we need a way that says that for some mixing it is not necessary to throw an error, although it is prohibited.

To do this, you need to write a **more specific rule**.

```php
"red allow-rg green" => ""
```

> Please note that there is an empty line on the right, this is a sign of an exclusion rule.

This approach is very similar to CSS.

Let's imagine we have a style:

```css
/* rule 1 */
.red .green {
    text: "it is forbidden to call green functions from red ones"
}
```

But we can make a more specific selector that overrides the previous one:

```css
/* rule 2 */
.red .allow-rg .green {
    text: ""
}
```

In HTML:

```html
<div class="red">
    <div class="green"></div> <!-- rule 1 is used -->
</div>

<div class="red allow-rg">
    <div class="green"></div> <!-- rule 2 is used -->
</div>

<!-- or -->

<div class="red">
    <div class="allow-rg green"></div> <!-- rule 2 is used -->
</div>

<!-- or -->

<div class="red">
    <div class="allow-rg">
        <div class="green"></div> <!-- rule 2 is used -->
    </div>
</div>
```

As you can see, there are three ways to apply a more specific selector. Now imagine that each **div** is a **function**, and its **classes** are its **colors**.

Then, to apply a specific rule for functions, you need to either add an `allow-rg` color for the `red` function, or the `green` one or add another function between them with the `allow-rg` color.

```php
/**
 * @color red
 * @color allow-rg
 */
function f1() { 
  f2(); // ok
}

/**
 * @color green
 */
function f2() { echo 1; }

// or --------------------------------

/**
 * @color red
 */
function f1() { 
  f2(); // ok
}

/**
 * @color allow-rg
 * @color green
 */
function f2() { echo 1; }

// or --------------------------------

/**
 * @color red
 */
function f1() { 
  allower(); // ok
}

/**
 * @color allow-rg
 */
function allower() { 
  f2();
}

/**
 * @color green
 */
function f2() { echo 1; }
```

> Please note that order is important, i.e.
>
> ```
> @color allow-rg
> @color green
> ```
>
> and
>
> ```
> @color green
> @color allow-rg
> ```
>
> are different cases.
>
> When colors are blended, function colors are mixed from top to bottom.
>
> ```php
> /**
>  * @color red
>  */
> function f1() { 
>   f2(); // ok
> }
> 
> /**
>  * @color green
>  * @color allow-rg
>  */
> function f2() { echo 1; }
> ```
>
> This will be the color chain `red green allow-rg`, not the `red allow-rg green` required for the exclusion rule.

### Groups of rules

All colors are divided into logical groups, which usually contain one rule leading to an error and several exclusion rules.

The discussed rules above are combined into one group:

```php
"red green"          => "it is forbidden to call green functions from red ones"
"red allow-rg green" => ""
```

The color palette can contain any number of such groups.

> Please note that the exclusion rules must follow strictly the error rules.

If within one group you need to write an even more specific rule that leads to an error, then it is written at the very end:

```php
"red green"               => "it is forbidden to call green functions from red ones"
"red allow-rg green"      => ""
"red allow-rg blue green" => "some error"
```

Exclusion rules are also written after the rule resulting in the error.

```php
"red green"                         => "it is forbidden to call green functions from red ones"
"red allow-rg green"                => ""
"red allow-rg blue green"           => "some error"
"red allow-rg allow-rbg blue green" => ""    
```

### Special color `@remover`

In large projects like ours, some functions call many other functions, and therefore, because of them, the analysis process can be very long, since many functions are reachable from them. To exclude such functions from the call graph, thereby manually decoupling the large connectivity component, there is a special color `@remover` that removes the function from the call graph.

For example:

```php
/**
 * @color remover
 */
function router() { 
  f1(); f2(); f3(); f4(); f5(); f6();
}

/**
 * @color green
 */
function f1() { echo 1; } // same for f2, f3, ...
```

Here the `router` function calls many functions, and because of this, the call graph may become too large for analysis, or become recursive, that is, each function will be reachable from any other.

To avoid this, we added the `remover` color, and during the analysis, the function will be removed from the call graph, that is, if some function was reachable only through the `router` function, then now it will be unreachable since the `router` function does not exist within the analysis.

## Colors and palette in code

### Colors

The PHPDoc tag `@color` is used to set colors, each tag **must contain only one color** (all others will be ignored). Use multiple tags to define multiple colors. 

> Please note, If you want to use a different tag, then use the `--tag` flag, which accepts the string value of the tag that will be used to set the color.

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

You can also set colors for all methods of a class/interface/trait using the tag in PHPDoc above the class/interface/trait:

```php
/**
 * @color service
 */
class Service {
  public function method() {}  // contains the 'service' color implicitly
  public function method2() {} // contains the 'service' color implicitly
}
```

### Palette

The `yaml` config is used to define the palette. The structure of the `yaml` config is as follows:

```yaml
# In short, the structure consists of an array of groups, 
# where each group is an array of map<colors, error value or empty string> objects.

# first group
-
  # array of rules
  # each rule is a separate object
  # the first rule must be the rule leading to an error
  - "highload no-highload": "it is forbidden to call green functions from red ones"
  - "highload allow-hh no-highload": ""

# second group
-
  - "ssr rpc": "don't perform rpc queries from SSR, pass already prepared data"
  - "ssr ssr-allow-rpc rpc": ""
```

> Please note that if you want to use a different name or location for the palette file, use the `--palette` flag.

## Examples close to reality

If you want, [here](https://github.com/vkcom/nocolor/blob/master/docs/examples_close_to_reality.md) you can see examples of using the concept, close to reality.

## Next steps

- [Comparison with Deptrac](https://github.com/vkcom/nocolor/blob/master/docs/comparison_with_deptrac.md)