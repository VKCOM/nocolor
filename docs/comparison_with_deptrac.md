# Comparison with Deptrac

[**Deptrac**](https://github.com/qossmic/deptrac) is a popular tool for checking architecture through the concept of layers and rules. Each layer consists of a set of classes, which are formed using different collectors. Rules are a whitelist of allowed dependencies, that is, a list of layers that can depend on other layers.

**NoColor** takes a slightly different path. **Deptrac** works with classes, while **NoColor** works with functions. At the same time, **NoColor** uses tag in PHPDoc instead of various collectors in the config to mark functions.

This approach is more flexible, but more time-consuming since you need to mark all the colors manually, however, **NoColor** is not aimed at ensuring that as many functions as possible are marked, rather, on the contrary, colors should be used locally for strictly defined places.

Due to the fact that **NoColor** works at the function level, it is a candidate for testing legacy projects that have a lot of functions and few classes.

### Is it possible to replace Deptrac with NoColor?

Yes, it is theoretically possible.

However, to do this, you need to mark all classes with colors, which can be too expensive for huge projects and if you need to change some rules, then changing colors will take more costs than changing the config in **Deptrac**.

## Differences

### Number of levels of nesting when checking

When analyzing **Deptrac**, look only at one level.

In the case of **NoColor**, the maximum nesting level is 50 colored functions (in fact, if the depth is 100, but among these 100 functions only 49 are colored, then the error can still be found).

For example:

```php
<?php

class Service {
    public static function method() {
        Controller::method(); // allowed in config
    }
}

class Controller {
    public static function method() {
        Repository::method(); // allowed in config
    }
}

class Repository {
    public static function method() {
        echo 1;
    }
}
```

Here the `Service` class depends on the `Repository` at a depth of 2 and therefore **Deptrac** does not find an error, although this dependency is not allowed in the config.

Let's take a look at the same example for **NoColor**.

```php
// palette.yaml
// -
//   - "controller": ""
// -
//   - "service repository": "Service -> Repository dependency not allowed."

<?php

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
```

**NoColor** can find the following error here:

```
Error at the stage of checking colors
   test.php:25  in function Repository::method

service repository => Service -> Repository dependency not allowed.
  This color rule is broken, call chain:
Service::method@service -> Controller::method -> Repository::method@repository
```

### Type inference and return types

**Deptrac** does not do type inference, which is why quite a few dependencies can be lost.

**NoColor** makes type inference, which can cover 90% of possible dependencies.

> Please note that since PHP is a dynamically typed language, it is impossible to infer types always and everywhere.

For example:

```php
<?php

class Something {
    private Repository $repository;

    public function __construct() {
        $this->repository = new Repository;
    }

    public function getRepository(): Repository {
        return $this->repository;
    }
}

/**
 * @color repository
 */
class Repository {
    public function do() {
        echo 1;
    }
}

/**
 * @color service
 */
class Service {
    public static function method() {
        $something = new Something;
        $repository = $something->getRepository();
        $repository->do();
    }
}
```

**Deptrac** thinks this code is correct, but the `Service` actually depends on the `Repository` because the `Something::getRepository` method returns an object of the `Repository` class.

**NoColor** found the following error:

```
Error at the stage of checking colors
   script.php:20  in function Repository::do

service repository => Service -> Repository dependency not allowed.
  This color rule is broken, call chain:
Service::method@service -> Repository::do@repository
```

### Whitelist vs Blacklist with Exclusions

**Deptrac** uses a whitelist that allows one layer to be dependent on another.

**NoColor** uses a blacklist, which actually consists of rules that prohibit calling one function from another. Also, in these rules, you can write exclusion rules for specific cases.

### Performance

**Deptrac** runs slightly faster on personal computers thanks to the above cut corners.

**NoColor**, due to the fact that it does a lot of additional things, loses somewhat in speed. However, unlike **Deptrac**, **NoColor** is easy to parallelize, which allows you to use all the powers provided. On development servers, **NoColor** can run faster than **Deptrac**.

## Conclusions

1. **Deptrac** works with classes and **NoColor** works with functions;
2. In **Deptrac**, all configuration takes place in the config using various selectors, and in **NoColor**, using annotations over functions, and a palette in the form of a config file;
3. **Deptrac** only supports one level of nesting, and **NoColor** can check in depth for 50 color functions;
4. In **Deptrac**, due to the lack of type inference, not all dependencies can be found, and in **NoColor** there is type inference, so all possible calls will be found;
5. **NoColor** is better suited for analyzing legacy projects where many functions and few classes are used;
6. **Deptrac** works faster on personal computers, but in one thread, and **NoColor** is easily parallelized and can work at all the capacities provided;
7. **Deptrac** provides various nice and visual views like Graphviz, **NoColor** doesn't;
8. **Deptrac** provides many collectors, which can greatly simplify the work with the tool when **NoColor** has only one way to set the color — in PHPDoc.

As a result, **NoColor** is not a replacement for **Deptrac**, but its addition aimed at deeper validation of function calls.

## Next steps

- [How to contribute](https://github.com/vkcom/nocolor/blob/master/docs/how_to_contribute.md)
