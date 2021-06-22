1. **Specify the type:**

   `bugfix` / `feature` / `refactoring` / `docs update`

2. **Describe the changes:**
   A clear and concise description of what has changed.

3. **Specify the issue number that this PR solves.** If the issue does not exist, then create it.

   > For example:
   > Fixes #999

---

Example:

**Type**: `bugfix`

Fixed a bug when the user flew to the moon due to overflow. The type for fields storing the ship's speed has been changed to `int64`.

```go
// Optional code to explain the changes.
package main

func main() {
    return
}
```

Fixes #999
