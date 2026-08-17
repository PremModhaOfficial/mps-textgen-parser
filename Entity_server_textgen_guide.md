# TextGen Fixup Guide

Source: `Entity_server_textgen.md`
Output: `Entity_server_textgen.mps`

## File Metadata (add manually in MPS)

- Line 2: `text gen component for concept Entity {`
- Line 3: `file name :`
- Line 5: `file path :`
- Line 6: `extension :`

## Expressions (91 total)

Search `{???-` in MPS. Each is a ConstantStringAppendPart — replace content with actual expression.

| # | Expression | Line | Type |
|---|---|---|---|
| 1 | `name` | 34 | D: local var |
| 2 | `field.name.capitalize()` | 36 | C: chained |
| 3 | `field.goType()` | 36 | B: method() |
| 4 | `field.jsonName()` | 36 | B: method() |
| 5 | `field.dbName()` | 36 | B: method() |
| 6 | `name` | 45 | D: local var |
| 7 | `name` | 46 | D: local var |
| 8 | `name` | 46 | D: local var |
| 9 | `nameLower` | 46 | D: local var |
| 10 | `name` | 53 | D: local var |
| 11 | `name` | 54 | D: local var |
| 12 | `name` | 54 | D: local var |
| 13 | `nameLower` | 54 | D: local var |
| 14 | `name` | 61 | D: local var |
| 15 | `name` | 62 | D: local var |
| 16 | `nameLower` | 62 | D: local var |
| 17 | `name` | 69 | D: local var |
| 18 | `name` | 78 | D: local var |
| 19 | `name` | 79 | D: local var |
| 20 | `nameLower` | 79 | D: local var |
| 21 | `name` | 88 | D: local var |
| 22 | `name` | 93 | D: local var |
| 23 | `name` | 93 | D: local var |
| 24 | `name` | 94 | D: local var |
| 25 | `name` | 106 | D: local var |
| 26 | `opName` | 106 | D: local var |
| 27 | `name` | 108 | D: local var |
| 28 | `opName` | 108 | D: local var |
| 29 | `name` | 115 | D: local var |
| 30 | `name` | 118 | D: local var |
| 31 | `name` | 121 | D: local var |
| 32 | `name` | 124 | D: local var |
| 33 | `name` | 127 | D: local var |
| 34 | `name` | 143 | D: local var |
| 35 | `field.name` | 143 | A: property |
| 36 | `name` | 146 | D: local var |
| 37 | `field.name` | 146 | A: property |
| 38 | `nameLower` | 152 | D: local var |
| 39 | `name` | 160 | D: local var |
| 40 | `pkField` | 160 | D: local var |
| 41 | `nameLower` | 161 | D: local var |
| 42 | `name` | 169 | D: local var |
| 43 | `nameLower` | 170 | D: local var |
| 44 | `name` | 178 | D: local var |
| 45 | `nameLower` | 179 | D: local var |
| 46 | `nameLower` | 199 | D: local var |
| 47 | `name` | 199 | D: local var |
| 48 | `pkField` | 199 | D: local var |
| 49 | `nameLower` | 202 | D: local var |
| 50 | `name` | 202 | D: local var |
| 51 | `pkField` | 202 | D: local var |
| 52 | `nameLower` | 205 | D: local var |
| 53 | `name` | 205 | D: local var |
| 54 | `nameLower` | 208 | D: local var |
| 55 | `name` | 208 | D: local var |
| 56 | `opName` | 215 | D: local var |
| 57 | `nameLower` | 223 | D: local var |
| 58 | `opKind` | 223 | D: local var |
| 59 | `name` | 238 | D: local var |
| 60 | `opKind` | 238 | D: local var |
| 61 | `opName` | 242 | D: local var |
| 62 | `name` | 254 | D: local var |
| 63 | `hookName` | 254 | D: local var |
| 64 | `name` | 254 | D: local var |
| 65 | `name` | 258 | D: local var |
| 66 | `hookName` | 258 | D: local var |
| 67 | `name` | 258 | D: local var |
| 68 | `name` | 265 | D: local var |
| 69 | `hookName` | 265 | D: local var |
| 70 | `name` | 265 | D: local var |
| 71 | `name` | 269 | D: local var |
| 72 | `hookName` | 269 | D: local var |
| 73 | `name` | 269 | D: local var |
| 74 | `name` | 276 | D: local var |
| 75 | `hookName` | 276 | D: local var |
| 76 | `name` | 276 | D: local var |
| 77 | `name` | 280 | D: local var |
| 78 | `hookName` | 280 | D: local var |
| 79 | `name` | 280 | D: local var |
| 80 | `name` | 287 | D: local var |
| 81 | `hookName` | 287 | D: local var |
| 82 | `name` | 287 | D: local var |
| 83 | `name` | 291 | D: local var |
| 84 | `hookName` | 291 | D: local var |
| 85 | `name` | 291 | D: local var |
| 86 | `name` | 298 | D: local var |
| 87 | `hookName` | 298 | D: local var |
| 88 | `name` | 298 | D: local var |
| 89 | `name` | 302 | D: local var |
| 90 | `hookName` | 302 | D: local var |
| 91 | `name` | 302 | D: local var |

## Control Flow (40 items)

These are XML comments. Wrap enclosed nodes in MPS.

| Type | Code | Line |
|---|---|---|
| VAR | `string name = node.name;` | 11 |
| VAR | `string nameLower = node.name.toLowerCase();` | 12 |
| VAR | `string pkField = node.primaryKeyField().name;` | 13 |
| FOREACH | `foreach field in node.fields` | 35 |
| FOREACH | `foreach op in node.operations` | 42 |
| IF | `if (op.entityOperation == EntityOperation:create)` | 44 |
| IF | `if (op.entityOperation == EntityOperation:update)` | 52 |
| IF | `if (op.entityOperation == EntityOperation:delete)` | 60 |
| IF | `if (op.entityOperation == EntityOperation:list)` | 68 |
| IF | `if (op.entityOperation == EntityOperation:get)` | 77 |
| FOREACH | `foreach op in node.operations` | 102 |
| VAR | `string opName = op.capitalizedName();` | 103 |
| VAR | `string opKind = op.entityOperation.name;` | 104 |
| IF | `if (op.entityOperation == EntityOperation:create)` | 114 |
| IF | `if (op.entityOperation == EntityOperation:update)` | 117 |
| IF | `if (op.entityOperation == EntityOperation:delete)` | 120 |
| IF | `if (op.entityOperation == EntityOperation:list)` | 123 |
| IF | `if (op.entityOperation == EntityOperation:get)` | 126 |
| IF | `if (op.entityOperation == EntityOperation:create)` | 138 |
| VAR | `int valIdx = 0;` | 139 |
| FOREACH | `foreach field in node.fields` | 140 |
| IF | `if (!(field.hasAnnotation(FieldAnnotation:primaryKey)) && !(field.hasAnnotation(FieldAnnotation:auto)) && !(field.hasAnnotation(FieldAnnotation:hidden)) && !(field.hasAnnotation(FieldAnnotation:nullable)))` | 141 |
| IF | `if (valIdx == 0)` | 142 |
| IF | `if (valIdx > 0)` | 145 |
| ASSIGN | `valIdx = valIdx + 1;` | 148 |
| IF | `if (op.entityOperation == EntityOperation:update)` | 159 |
| IF | `if (op.entityOperation == EntityOperation:delete)` | 168 |
| IF | `if (op.entityOperation == EntityOperation:get)` | 177 |
| IF | `if (op.entityOperation == EntityOperation:list)` | 186 |
| IF | `if (op.entityOperation == EntityOperation:create)` | 198 |
| IF | `if (op.entityOperation == EntityOperation:update)` | 201 |
| IF | `if (op.entityOperation == EntityOperation:delete)` | 204 |
| IF | `if (op.entityOperation == EntityOperation:get)` | 207 |
| FOREACH | `foreach op in node.operations` | 250 |
| VAR | `string hookName = op.capitalizedName();` | 251 |
| IF | `if (op.entityOperation == EntityOperation:create)` | 253 |
| IF | `if (op.entityOperation == EntityOperation:update)` | 264 |
| IF | `if (op.entityOperation == EntityOperation:delete)` | 275 |
| IF | `if (op.entityOperation == EntityOperation:get)` | 286 |
| IF | `if (op.entityOperation == EntityOperation:list)` | 297 |
