# TextGen Fixup Guide

Source: `Relation_server_textgen.md`
Output: `throwed.textGen.mps`

## File Metadata (add manually in MPS)

- Line 2: `text gen component for concept Relation {`
- Line 3: `file name :`
- Line 5: `file path :`
- Line 6: `extension :`

## Expressions (51 total)

Search `{???-` in MPS. Each is a ConstantStringAppendPart — replace content with actual expression.

| # | Expression | Line | Type |
|---|---|---|---|
| 1 | `node.from.name` | 29 | C: chained |
| 2 | `node.to.name` | 29 | C: chained |
| 3 | `node.from.name` | 30 | C: chained |
| 4 | `node.from.name.toLowerCase()` | 30 | C: chained |
| 5 | `node.to.name` | 31 | C: chained |
| 6 | `node.to.name.toLowerCase()` | 31 | C: chained |
| 7 | `node.from.name` | 38 | C: chained |
| 8 | `node.to.name` | 38 | C: chained |
| 9 | `node.from.name` | 39 | C: chained |
| 10 | `node.from.name.toLowerCase()` | 39 | C: chained |
| 11 | `node.to.name` | 40 | C: chained |
| 12 | `node.to.name.toLowerCase()` | 40 | C: chained |
| 13 | `node.from.name` | 47 | C: chained |
| 14 | `node.to.name` | 47 | C: chained |
| 15 | `node.from.name` | 48 | C: chained |
| 16 | `node.from.name.toLowerCase()` | 48 | C: chained |
| 17 | `node.from.name` | 59 | C: chained |
| 18 | `node.to.name` | 59 | C: chained |
| 19 | `node.from.name` | 64 | C: chained |
| 20 | `node.to.name` | 64 | C: chained |
| 21 | `node.from.name` | 64 | C: chained |
| 22 | `node.to.name` | 64 | C: chained |
| 23 | `node.from.name` | 65 | C: chained |
| 24 | `node.to.name` | 65 | C: chained |
| 25 | `node.from.name` | 75 | C: chained |
| 26 | `node.to.name` | 75 | C: chained |
| 27 | `op.capitalizedName()` | 75 | B: method() |
| 28 | `node.from.name` | 76 | C: chained |
| 29 | `node.to.name` | 76 | C: chained |
| 30 | `op.capitalizedName()` | 76 | B: method() |
| 31 | `node.from.name` | 83 | C: chained |
| 32 | `node.to.name` | 83 | C: chained |
| 33 | `node.from.name` | 86 | C: chained |
| 34 | `node.to.name` | 86 | C: chained |
| 35 | `node.from.name` | 89 | C: chained |
| 36 | `node.to.name` | 89 | C: chained |
| 37 | `node.from.name` | 100 | C: chained |
| 38 | `node.to.name` | 100 | C: chained |
| 39 | `node.from.name.toLowerCase()` | 101 | C: chained |
| 40 | `node.to.name.toLowerCase()` | 101 | C: chained |
| 41 | `node.from.name` | 107 | C: chained |
| 42 | `node.from.name.toLowerCase()` | 108 | C: chained |
| 43 | `node.from.name.toLowerCase()` | 118 | C: chained |
| 44 | `node.from.name` | 118 | C: chained |
| 45 | `node.to.name.toLowerCase()` | 119 | C: chained |
| 46 | `node.to.name` | 119 | C: chained |
| 47 | `node.from.name.toLowerCase()` | 122 | C: chained |
| 48 | `node.from.name` | 122 | C: chained |
| 49 | `node.from.name.toLowerCase()` | 130 | C: chained |
| 50 | `node.to.name.toLowerCase()` | 130 | C: chained |
| 51 | `op.relationOperation.name` | 130 | C: chained |

## Control Flow (12 items)

These are XML comments. Wrap enclosed nodes in MPS.

| Type | Code | Line |
|---|---|---|
| FOREACH | `foreach op in node.operations` | 26 |
| IF | `if (op.relationOperation == RelationOperation:assign)` | 28 |
| IF | `if (op.relationOperation == RelationOperation:remove)` | 37 |
| IF | `if (op.relationOperation == RelationOperation:list)` | 46 |
| FOREACH | `foreach op in node.operations` | 73 |
| IF | `if (op.relationOperation == RelationOperation:assign)` | 82 |
| IF | `if (op.relationOperation == RelationOperation:remove)` | 85 |
| IF | `if (op.relationOperation == RelationOperation:list)` | 88 |
| IF | `if (op.relationOperation == RelationOperation:assign || op.relationOperation == RelationOperation:remove)` | 99 |
| IF | `if (op.relationOperation == RelationOperation:list)` | 106 |
| IF | `if (op.relationOperation == RelationOperation:assign || op.relationOperation == RelationOperation:remove)` | 117 |
| IF | `if (op.relationOperation == RelationOperation:list)` | 121 |
