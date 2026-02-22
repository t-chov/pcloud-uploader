---
name: version-bump
description: Bump version number across Makefile and README.md
disable-model-invocation: true
argument-hint: "[version-number]"
---

# Version Bump

Bump the application version to `$0`.

## Steps

1. Read `Makefile` and update the `VERSION` variable to `$0`
2. Read `README.md` and update the version string in the usage block to `$0`
3. Verify both files contain the new version by grepping for `$0`
4. Show a summary of the changes (do NOT commit automatically)
