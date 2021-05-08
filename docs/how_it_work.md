# How it works internally

Let's look at how colors are checked and what **problems** can arise if we talk about **large projects**.

First of all, remember that our concept consists of **colors** and **mixing rules**, where the rules are grouped, in each group the **rule leading to an error is the first rule**, followed by **optional exception rules**.

```yaml
-
  - "highload no-highload": "Calling no-highload function from highload function"

-
  - "highload db-access": "Calling function with access to DB from highload function"
  - "highload allow-highload-db db-access": ""
```

Each function can have **any number of colors** and the **order** of the colors in the function is **important**.

```php
/**
 * @color no-highload
 */
function f1() {}

/**
 * @color highload
 */
function f2() {}
```

Let's see what happens when checking if the function `f2` calls `f1`:

### Naive approach

1. First of all, the entire code is analyzed and the **call graph** is **built**. Thus, we get a **directed graph** in which the **vertices denote functions**, and the **edges denote the fact of calling** a function from another function.

   In our case, we will get a simple graph with two vertices and one edge:

   ```
   f2@highload -> f1@no-highload
   ```

   > After the `@` sign, we will denote the colors of functions.

2. The next step, using **breadth-first traversal**, we traverse this graph, while **collecting the function call stack**. 

   That is, first we will have a stack of **one element**:
   
   ```
   f2@highload
   ```

   Then of the **two**:
   
   ```
   f2@highload, f1@no-highload
   ```

   And at **every step**, we **call a check**.

   Checking is comparing the **received call stack** with the **rules** that are in the palette.

   In the first step, we will have **one function** in the **stack** with one color `highload`. We need to **find a rule** that **matches this chain** (in our case, size 1) of colors.

   Since our **rules on the right consist** of a **set of colors**, it is enough to **take all the colors from our call stack** and compare them with **each right side of the rules** that are in the palette.

   So, for example, we have a rule:
   
   ```yaml
   "highload no-highload": "Calling no-highload function from highload function"
   ```

   On the **right** we have a set of two colors:
   
   ```
   highload, no-highload
   ```

   And in the current call stack, there is **one function** with the color `highload`. So all we have to do is **take and compare these two sets**.

   As we can see, they **do not match**, since their **length is not equal**, so this rule does **not match** and we **move on to the next one**.

   And so each rule is checked and **if the rule matches**, and **if this rule describes an error**, then we **throw an error**.

   Very important here is the fact that each **group contains only one rule that leads to an error**, and **all the others should be exception rules**. Therefore, when checking the rules, the rules are **compared from the end**, that is, the **exclusion rules are checked first**, and at the **very end, the rule is an error**. If the **current color chain** matches **some exception rule**, then the rules of this group are **not checked further**.

   Let's consider a case when we have two functions in the call stack:
   
   ```
   f2@highload, f1@no-highload
   ```

   Here the chain of colors will be two in length:
   
   ```
   highload, no-highload
   ```

   Next, we go **through each group** and **compare the color chains from the last rule**.

   In the first group, the first and last rule is as follows:
   
   ```
   "highload no-highload": "Calling no-highload function from highload function"
   ```
   
   On the **right**, the **color chain matches with the current color chain**, which means that the **current rule describes the current situation,** and since this **rule leads to an error**, we **throw it**.

This approach works well on a **small codebase** when the **number of functions is small** and the **call graph connectivity is not high**. However, **we have over 200k functions** in our codebase, and this is a **huge monolith** with **strong connectivity**. A **naive algorithm cannot cope** with such a huge graph, so we had to come up with ways to **optimize it**.

### Optimize

First of all, let's consider this situation, we have 100 functions that are linked into a call graph. If we put colors on the first function in the graph and on the last, then to check this chain of calls we will need to put all 100 functions on the call stack, and at the same time, if the graph is strongly connected, then to pass to the last function, we will need to make a very large traversal, the value which will increase with the number of functions and connectivity.

Let's take the following graph as an example:

```
f1@highload -> f2  -> f3  -> ... -> f20 -> f21@no-highload
           |-> f22 -> f23 -> ... -> f50 -> f51@no-highload
```

Here, the naive algorithm will have to go through a lot of functions and only at the end check the color chain for matching the rules.

But we can do the following.

Add a preparatory pass through the graph, which forms for each color function a list of color functions that are reachable from the current one. And also this list will be at the root node.

Thus, in our graph, the function `f1` will have a list: `f21, f51`. With such a list, we just need to take the colors of the function `f1`, mix them first with the colors of the first functions from the list, check that they match the rules, and then with the colors of the second function and check again.

Thus, the length of the chain between two colored functions will no longer be a problem, since we have pointers to the next colored functions, thereby completely solving the problem in the case of strongly connected huge graphs that we have.

The algorithm for calculating these "jump" lists is much simpler, the process of "bubbling" occurs when we go from the function up to the root, and if we meet a colored function, then the "pop-up" function is added to the list and the current function begins to "pop up" further.

In the case of a recursive part in a graph, we set up lists with the assumption that every colored function in the recursive loop is reachable from any other colored function in that loop.

However, in case of an error, this approach contains only functions with colors in the call stack, and in the output, I would like to see the full call stack so that the entire chain can be traced. Since this is needed already at the stage when the error is found, performance is not very important to us in this case, so we use a breadth-first search between every two colored functions from the call stack and connect them into the final chain.

#### Comparison of chains

When there are a lot of rules, each comparison can be quite expensive. Therefore, it makes sense to optimize this part as well. 

The first thing to talk about is how to store colors. The storage option in its original form is the simplest, but it is inefficient in memory and execution time since it will also require string comparison at the rule checking stage.

The very first option that comes to mind is to convert colors to numbers, then we will save time on checking the equality of two colors, but here we can make it a little more interesting, which will allow us to make a couple more optimizations.

We can set colors as masks, that is, each color is expressed as:

```
1 << n, where n is the ordinal number of the color
```

Thus, each color chain can be represented as a 64-bit mask, which will allow us to speed up the comparison.

This method imposes a limitation, we cannot store more than 64 colors, which, however, should be enough for most tasks.

Now let's take a closer look at bitmasks. As already mentioned, each rule and each color chain can be represented as a 64-bit mask. This mask will contain 1 for the color indices that are in the chain.

Let's take the following values as an example:

```
Colors:
highload =          0b0001
no-highload =       0b0010
internals =         0b0100
module =            0b1000

Rules:
-
  - "highload no-highload": "Calling no-highload function from highload function"
    0b0001 + 0b0010 = 0b0011
-
  - "internals": "Calling function marked as internal outside of module functions"
    0b0100 + 0 = 0b0100
    
  - "module internals": ""
    0b1000 + 0b0100 = 0b1100
```

Now each rule has its own mask, which determines what colors it contains.

Now, if we have a chain:

```
f3@module -> f4@internals
```

Then we can calculate its mask:

```
0b1000 + 0b0100 = 0b1100
```

Now, when comparing the rules, we can check if the mask matches the rule. Let's take the first rule, its mask is `0b0011`. Let's do a bitwise AND with a chain mask:

```
0b0011 & 0b1100 = 0b0000
```

As a result, we will get a result different from the original chain mask, which means that the chain does not contain colors from the current rule, which means that the rule does automatically not match.

Let's move on to the next rule. In it, the mask is `0b1100`. A bitwise AND with this mask will give a result equal to the chain mask, which means that the chain and the rule contain the same colors. However, this is not a sufficient fact, since the order is important to us, because the following two chains will give the same masks, but only the first one will fall under the rule:

```
f3@module    -> f4@internals
f4@internals -> f3@module
```

After we have determined that the chain and the rule have the same colors, then we need to check the order and number of colors.

Our task is to understand as quickly as possible that the current chain does not match the chain from the rule, so it makes sense to go from right to left, and not vice versa. This will allow you to identify chains where the beginning will coincide, but the end will not.

Thus, when checking, we go from right to left, and if the chains match, we say that the current rule is suitable.

We also need to take into account the fact that functions may have more colors in the color chain than in the rule:

```
f3@module -> f5@highload -> f4@internals
```

Here, when comparing, we skip colors that are not in the rule's color chain, since they are not important colors.

### Special color `@remover`

The concept introduces the concept of a special color `remover`, which is intended for manual separation of strongly connected graphs. In fact, when creating a graph, if a function has this color, then it is removed from the graph before being checked.

## Next steps

- [Description of the color concept](https://github.com/vkcom/nocolor/blob/master/docs/concept_of_colors.md)