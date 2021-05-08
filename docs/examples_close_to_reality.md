# Examples close to reality

Disallowing green functions to be called from red ones **may seem meaningless**, so let's look at some examples that are **close to reality**.

### Highly loaded places

There are function `getDataFromDatabase` that go to a **slow database** for data, and there are function `getDataFromMemcached` that go to **fast memcached**. We have a highly loaded place `highloadFunction` where the data is **accessed many times** and goes into the database for the data to turn into a disaster. We *don't want this to happen*.

```php
<?php

/**
 * @color highload
 */
function highloadFunction() {
  for ($i = 0; $i < 1000; $i++) {
    echo getDataFromMemcached(); // ok 
    // echo getDataFromDatabase(); // bad
  }
}

/**
 * @color no-highload
 * @color db-access
 * @return int
 */
function getDataFromDatabase() {
  return 1;
}

/**
 * @color highload
 * @return int
 */
function getDataFromMemcached() {
  return 1;
}
```

First of all, we need to determine the colors that we need.

We can use the `highload` color to indicate functions that are both **highly loaded** and functions that **can work under high load**.

For the function that goes to the database, we can come up with several options. It is either the color `db-access` or `no-highload`. The first option **gives more information**, while the second is **better suited** to the concept of `highload / no-highload` functions.

One way or another, with colors, we must write **rules for mixing colors**, that is, rules that **prohibit or allow** calling some functions from others.

We want to prevent calling `no-highload` or `db-access` functions from `highload` functions.

The palette will look like this:

```yaml
# The first group of rules will prohibit the first option.
-
  - "highload no-highload": "Calling no-highload function from highload function"

# And the second, respectively, the second.
-
  - "highload db-access": "Calling function with access to DB from highload function"
```

In this case, if someone tries to call the `no-highload` or `db-access` function, they will receive an error message.

#### Afterword

Of these two rules, the first one sounds quite unambiguous, cannot call functions in this way, and most likely there should be no exceptions. But the second option is too harsh. A situation is possible when it is necessary when the developer understands what he is doing and what it can lead to.

Exception rules have already been mentioned in the description of the concept, so let's write another rule that will allow using the database functions in `highload` functions.

```yaml
-
  - "highload no-highload": "Calling no-highload function from highload function"

-
  - "highload db-access": "Calling function with access to DB from highload function"
  - "highload allow-highload-db db-access": ""
```

Now, just add the color  `allow-highload-db` on the function `highloadFunction` that calls or on the function `getDataFromDatabase` that is called. In the first case, which is preferable, in the `highloadFunction` function it will be possible to call the `db-access` functions.

```php
/**
 * @color highload
 * @color allow-highload-db
 */
function highloadFunction() {
  for ($i = 0; $i < 1000; $i++) {
    echo getDataFromMemcached(); // ok 
    // echo getDataFromDatabase(); // will be ok
  }
}
```

The `getDataFromDatabase` function can be called from any `highload` function in the second case.

```php
/**
 * @color no-highload
 * @color db-access
 * @color allow-highload-db
 * @return int
 */
function getDataFromDatabase() {
  return 1;
}

```

### Internals

Suppose we have some folder with some module. This module contains some API and you would not want anyone to use something outside of this API.

We want to prevent some functions from being called from this module because they are **internal**.

Or you already have such function calls and **want to find them**.

Let's describe the rules that **prohibit doing this outside of the module**.

```yaml
-
  - "some-internals": "Calling function marked as internal outside of module functions"
  - "some-module some-internals": ""
```

The first rule will **disallow any function calls** that are marked as internal.

The second rule will **add an exception that allows** such functions to be called from functions with the `some-module` color.

Now it is enough to add the color `some-internals` to the functions that you do not want to be called outside the module, but for the places in the module where they are called add the color `some-module`. Thus, all calls outside the module will give errors.

Note that if there are **other functions between** the function with the `some-module` color and the `some-internal` color in the chain, then this **will not give an error either**.

#### Afterword

This pattern will help keep **track of where functions can be called**.

If you need **several allowed places**, then it is enough to **add a new exclusion rule**.

```yaml
-
  - "some-internals": "Calling function marked as internal outside of module functions"
  - "some-module some-internals": ""
  - "some-module-2 some-internals": ""
```

## Next steps

- [Comparison with Deptrac](https://github.com/vkcom/nocolor/blob/master/docs/nocolor_vs_deptrac.md)