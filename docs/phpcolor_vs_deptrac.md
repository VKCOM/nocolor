# Comparison of NoColor and Deptrac

[**Deptrac**](https://github.com/qossmic/deptrac) is a popular tool for checking architecture through the concept of layers and rules. In this respect, it is similar to **NoColor**, which uses colors and mixing rules. However, **Deptrac** targets classes, while **NoColor** targets functions. It is possible to implement **Deptrac** capabilities using **NoColor**, but for this, you have to manually set colors to all classes, which can be very expensive. In this regard, **NoColor** is not a replacement for **Deptrac**.

The other difference is the number of nesting levels that the tools check. In the case of **Deptrac**, this is just **one level**, which sometimes may not be enough, in **NoColor** the **nesting level can reach 50** (only *functions with colors* are taken into account, that is, if the actual depth is 100, but there are only 49 color functions in this chain, then the error will still be found).

In the following example:

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

**Deptrac** will not find an error, although `Service` is not specified in the config that it can depend on `Repository`.

In **NoColor**, if you set the colors:
```php
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

With the following rules:

```yaml
-
  - "controller": ""

-
  - "service repository": "Service -> Repository dependency not allowed."
```

An error will be found:

```
Error at the stage of checking colors
   1.php:25  in function Repository::method

service repository => Service -> Repository dependency not allowed.
  This color rule is broken, call chain:
Service::method@service -> Controller::method -> Repository::method@repository
```

***

The check for 10 levels of nesting has its drawbacks, it is speed, **Deptrac**, although written in PHP, but with a cache, it outperforms **NoColor** with a cache in analysis speed.

This speed is also explained by the fact that **Deptrac** cannot track everything. For example, if some function returns a class and an instance of this class is used inside the class for which the dependency is prohibited, then **Deptrac** cannot find such an error, this comes from — because **Deptrac** does not do type inference during analysis, which significantly reduces the work time. **NoColor** does type inference, so it can find almost all function and method calls and check them.

### Whitelist vs Rules

It should also be noted that the **Deptrac** config specifies classes that **may depend** on each other, **NoColor** goes the other way around, it **prohibits** calling some functions from some others. Also, due to the fact that **NoColor** checks functions, it has the ability to add **exceptions for the rules**, if necessary. This is hardly necessary for **Deptrac**, since it operates with classes where the connections are more significant.

## Conclusions

As a result, **NoColor** is not a replacement for **Deptrac**, but its addition aimed at deeper validation of function calls.

## Next steps

- [How it works](https://github.com/vkcom/nocolor/blob/master/docs/how_it_work.md)