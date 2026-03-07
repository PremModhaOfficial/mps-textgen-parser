# TextGen Fixup Guide

Source: `appends.md`
Output: `GoTextGen.textGen.mps`

## File Metadata (add manually in MPS)

- Line 70: `text gen component for concept Code {`
- Line 71: `file name : (node)->string { node.model.name; }`
- Line 557: `file path : (node)->string { "go/"; }`
- Line 558: `extension : (node)->string { "go"; }`

## Expressions (172 total)

Search `▶` in MPS. Each is a ConstantStringAppendPart — replace content with actual expression.

| # | Expression | Line | Type |
|---|---|---|---|
| 1 | `infra.modulePath` | 77 | A: property |
| 2 | `model.name` | 77 | A: property |
| 3 | `model.name` | 77 | A: property |
| 4 | `infra.port` | 77 | A: property |
| 5 | `schema.structName()` | 83 | B: method() |
| 6 | `f.pascalName()` | 87 | B: method() |
| 7 | `f.goType()` | 87 | B: method() |
| 8 | `f.name` | 87 | A: property |
| 9 | `f.name` | 87 | A: property |
| 10 | `schema.structName()` | 103 | B: method() |
| 11 | `schema.structName()` | 103 | B: method() |
| 12 | `f.pascalName()` | 108 | B: method() |
| 13 | `f.name` | 108 | A: property |
| 14 | `f.pascalName()` | 117 | B: method() |
| 15 | `schema.createStructName()` | 126 | B: method() |
| 16 | `f.pascalName()` | 131 | B: method() |
| 17 | `f.goType()` | 131 | B: method() |
| 18 | `f.name` | 131 | A: property |
| 19 | `schema.structName()` | 142 | B: method() |
| 20 | `fr.pascalName()` | 146 | B: method() |
| 21 | `fr.name` | 146 | A: property |
| 22 | `fr.name` | 146 | A: property |
| 23 | `f.pascalName()` | 150 | B: method() |
| 24 | `f.goType()` | 150 | B: method() |
| 25 | `f.name` | 150 | A: property |
| 26 | `f.name` | 150 | A: property |
| 27 | `fr.target_schema.structName()` | 160 | C: chained |
| 28 | `fr.pascalName()` | 160 | B: method() |
| 29 | `fr.name` | 160 | A: property |
| 30 | `rn` | 180 | D: local var |
| 31 | `rn` | 183 | D: local var |
| 32 | `sn` | 183 | D: local var |
| 33 | `tn` | 183 | D: local var |
| 34 | `f.name` | 190 | A: property |
| 35 | `i` | 198 | D: local var |
| 36 | `f.name` | 204 | A: property |
| 37 | `f.pascalName()` | 213 | B: method() |
| 38 | `f.pascalName()` | 222 | B: method() |
| 39 | `rn` | 228 | D: local var |
| 40 | `sn` | 228 | D: local var |
| 41 | `sn` | 228 | D: local var |
| 42 | `f.name` | 232 | A: property |
| 43 | `tn` | 235 | D: local var |
| 44 | `f.pascalName()` | 239 | B: method() |
| 45 | `rn` | 245 | D: local var |
| 46 | `sn` | 245 | D: local var |
| 47 | `f.name` | 249 | A: property |
| 48 | `tn` | 252 | D: local var |
| 49 | `sn` | 252 | D: local var |
| 50 | `sn` | 252 | D: local var |
| 51 | `f.pascalName()` | 256 | B: method() |
| 52 | `rn` | 262 | D: local var |
| 53 | `sn` | 262 | D: local var |
| 54 | `tn` | 262 | D: local var |
| 55 | `f.name` | 270 | A: property |
| 56 | `uidx` | 270 | D: local var |
| 57 | `nextParam` | 275 | D: local var |
| 58 | `f.pascalName()` | 279 | B: method() |
| 59 | `rn` | 285 | D: local var |
| 60 | `tn` | 285 | D: local var |
| 61 | `rn` | 299 | D: local var |
| 62 | `rn` | 302 | D: local var |
| 63 | `fr.name` | 308 | A: property |
| 64 | `sn` | 313 | D: local var |
| 65 | `sn` | 313 | D: local var |
| 66 | `tn` | 313 | D: local var |
| 67 | `fr.name` | 319 | A: property |
| 68 | `i` | 326 | D: local var |
| 69 | `fr.name` | 334 | A: property |
| 70 | `fr.name` | 344 | A: property |
| 71 | `f.name` | 348 | A: property |
| 72 | `fr.name` | 356 | A: property |
| 73 | `fr.pascalName()` | 365 | B: method() |
| 74 | `f.pascalName()` | 369 | B: method() |
| 75 | `rn` | 376 | D: local var |
| 76 | `fr.name` | 382 | A: property |
| 77 | `tn` | 386 | D: local var |
| 78 | `fr.name` | 393 | A: property |
| 79 | `wIdx` | 393 | D: local var |
| 80 | `fr.name` | 400 | A: property |
| 81 | `rn` | 415 | D: local var |
| 82 | `ts` | 415 | D: local var |
| 83 | `frA.target_schema.structName()` | 415 | C: chained |
| 84 | `frA.name` | 415 | A: property |
| 85 | `ts` | 415 | D: local var |
| 86 | `f.name` | 419 | A: property |
| 87 | `tt` | 422 | D: local var |
| 88 | `tn` | 422 | D: local var |
| 89 | `frB.name` | 422 | A: property |
| 90 | `frA.name` | 422 | A: property |
| 91 | `frA.name` | 422 | A: property |
| 92 | `ts` | 422 | D: local var |
| 93 | `ts` | 422 | D: local var |
| 94 | `f.pascalName()` | 426 | B: method() |
| 95 | `sn` | 446 | D: local var |
| 96 | `sn` | 449 | D: local var |
| 97 | `rn` | 449 | D: local var |
| 98 | `sn` | 449 | D: local var |
| 99 | `sn` | 452 | D: local var |
| 100 | `rn` | 452 | D: local var |
| 101 | `sn` | 455 | D: local var |
| 102 | `rn` | 455 | D: local var |
| 103 | `sn` | 458 | D: local var |
| 104 | `rn` | 458 | D: local var |
| 105 | `sn` | 458 | D: local var |
| 106 | `sn` | 461 | D: local var |
| 107 | `rn` | 461 | D: local var |
| 108 | `sn` | 474 | D: local var |
| 109 | `secondRef.target_schema.structName()` | 487 | C: chained |
| 110 | `rn` | 487 | D: local var |
| 111 | `firstRef.name` | 487 | A: property |
| 112 | `secondRef.target_schema.structName()` | 487 | C: chained |
| 113 | `firstRef.name` | 487 | A: property |
| 114 | `secondRef.pascalName()` | 487 | B: method() |
| 115 | `secondRef.target_schema.structName()` | 490 | C: chained |
| 116 | `rn` | 490 | D: local var |
| 117 | `firstRef.name` | 490 | A: property |
| 118 | `secondRef.name` | 490 | A: property |
| 119 | `secondRef.name` | 490 | A: property |
| 120 | `firstRef.name` | 490 | A: property |
| 121 | `secondRef.name` | 490 | A: property |
| 122 | `frA.target_schema.structName()` | 499 | C: chained |
| 123 | `frB.target_schema.structName()` | 499 | C: chained |
| 124 | `rn` | 499 | D: local var |
| 125 | `frB.target_schema.structName()` | 499 | C: chained |
| 126 | `frA.target_schema.structName()` | 499 | C: chained |
| 127 | `infra.dbUser` | 511 | A: property |
| 128 | `infra.dbPass` | 511 | A: property |
| 129 | `infra.dbName` | 511 | A: property |
| 130 | `schema.singularName()` | 516 | B: method() |
| 131 | `schema.repoName()` | 516 | B: method() |
| 132 | `schema.singularName()` | 521 | B: method() |
| 133 | `schema.repoName()` | 521 | B: method() |
| 134 | `sn` | 533 | D: local var |
| 135 | `tn` | 533 | D: local var |
| 136 | `sn` | 533 | D: local var |
| 137 | `vr` | 533 | D: local var |
| 138 | `tn` | 533 | D: local var |
| 139 | `sn` | 533 | D: local var |
| 140 | `vr` | 533 | D: local var |
| 141 | `tn` | 533 | D: local var |
| 142 | `sn` | 533 | D: local var |
| 143 | `vr` | 533 | D: local var |
| 144 | `tn` | 533 | D: local var |
| 145 | `sn` | 533 | D: local var |
| 146 | `vr` | 533 | D: local var |
| 147 | `tn` | 533 | D: local var |
| 148 | `sn` | 533 | D: local var |
| 149 | `vr` | 533 | D: local var |
| 150 | `schema.structName()` | 550 | B: method() |
| 151 | `fRef.target_schema.name` | 550 | C: chained |
| 152 | `sRef.target_schema.name` | 550 | C: chained |
| 153 | `sRef.target_schema.structName()` | 550 | C: chained |
| 154 | `vr` | 550 | D: local var |
| 155 | `fRef.target_schema.name` | 550 | C: chained |
| 156 | `sRef.target_schema.name` | 550 | C: chained |
| 157 | `sRef.name` | 550 | A: property |
| 158 | `sRef.target_schema.structName()` | 550 | C: chained |
| 159 | `vr` | 550 | D: local var |
| 160 | `fRef.target_schema.name` | 550 | C: chained |
| 161 | `sRef.target_schema.name` | 550 | C: chained |
| 162 | `fRef.target_schema.structName()` | 550 | C: chained |
| 163 | `sRef.target_schema.structName()` | 550 | C: chained |
| 164 | `vr` | 550 | D: local var |
| 165 | `sRef.target_schema.name` | 550 | C: chained |
| 166 | `fRef.target_schema.name` | 550 | C: chained |
| 167 | `sRef.target_schema.structName()` | 550 | C: chained |
| 168 | `fRef.target_schema.structName()` | 550 | C: chained |
| 169 | `vr` | 550 | D: local var |
| 170 | `infra.port` | 554 | A: property |
| 171 | `infra.port` | 554 | A: property |
| 172 | `infra.port` | 554 | A: property |

## Control Flow (202 items)

These are XML comments. Wrap enclosed nodes in MPS.

| Type | Code | Line |
|---|---|---|
| VAR | `node<Code> n = node;` | 73 |
| VAR | `node<Models> model = n.model;` | 74 |
| VAR | `node<Infra> infra = n.infra;` | 75 |
| FOREACH | `foreach schema in model.schemas` | 80 |
| IF | `if (!(schema.hasReferences()))` | 81 |
| FOREACH | `foreach field in schema.Fields` | 84 |
| IF | `if (field.isInstanceOf(Field))` | 85 |
| VAR | `node<Field> f = (Field) field;` | 86 |
| VAR | `boolean hasSensitive = false;` | 93 |
| FOREACH | `foreach field in schema.Fields` | 94 |
| IF | `if (field.isInstanceOf(Field))` | 95 |
| VAR | `node<Field> f = (Field) field;` | 96 |
| IF+STMT (inline) | `if (f.Sensitive) { hasSensitive = true; }` | 97 |
| IF | `if (hasSensitive)` | 102 |
| FOREACH | `foreach field in schema.Fields` | 104 |
| IF | `if (field.isInstanceOf(Field))` | 105 |
| VAR | `node<Field> f = (Field) field;` | 106 |
| IF | `if (f.Sensitive)` | 107 |
| FOREACH | `foreach field in schema.Fields` | 113 |
| IF | `if (field.isInstanceOf(Field))` | 114 |
| VAR | `node<Field> f = (Field) field;` | 115 |
| IF | `if (f.Sensitive)` | 116 |
| FOREACH | `foreach field in schema.Fields` | 127 |
| IF | `if (field.isInstanceOf(Field))` | 128 |
| VAR | `node<Field> f = (Field) field;` | 129 |
| IF | `if (f.dataType != DataType.timestamp)` | 130 |
| FOREACH | `foreach schema in model.schemas` | 140 |
| IF | `if (schema.hasReferences())` | 141 |
| FOREACH | `foreach field in schema.Fields` | 143 |
| IF | `if (field.isInstanceOf(FieldRefrence))` | 144 |
| VAR | `node<FieldRefrence> fr = (FieldRefrence) field;` | 145 |
| IF | `if (field.isInstanceOf(Field))` | 148 |
| VAR | `node<Field> f = (Field) field;` | 149 |
| VAR | `int refCount = 0;` | 155 |
| FOREACH | `foreach field in schema.Fields` | 156 |
| IF | `if (field.isInstanceOf(FieldRefrence))` | 157 |
| VAR | `node<FieldRefrence> fr = (FieldRefrence) field;` | 158 |
| IF | `if (refCount == 1)` | 159 |
| ASSIGN | `refCount = refCount + 1;` | 162 |
| FOREACH | `foreach schema in model.schemas` | 174 |
| IF | `if (!(schema.hasReferences()))` | 175 |
| VAR | `string sn = schema.structName();` | 176 |
| VAR | `string rn = schema.repoName();` | 177 |
| VAR | `string tn = schema.name;` | 178 |
| VAR | `int idx = 0;` | 184 |
| FOREACH | `foreach field in schema.Fields` | 185 |
| IF | `if (field.isInstanceOf(Field))` | 186 |
| VAR | `node<Field> f = (Field) field;` | 187 |
| IF | `if (f.dataType != DataType.timestamp)` | 188 |
| IF (inline) | `if (idx > 0)` | 189 |
| ASSIGN | `idx = idx + 1;` | 191 |
| FOR | `for (int i = 1; i <= idx; i++)` | 196 |
| IF (inline) | `if (i > 1)` | 197 |
| FOREACH | `foreach field in schema.Fields` | 201 |
| IF | `if (field.isInstanceOf(Field))` | 202 |
| VAR | `node<Field> f = (Field) field;` | 203 |
| IF (inline) | `if (f.dataType == DataType.timestamp)` | 204 |
| FOREACH | `foreach field in schema.Fields` | 210 |
| IF | `if (field.isInstanceOf(Field))` | 211 |
| VAR | `node<Field> f = (Field) field;` | 212 |
| IF (inline) | `if (f.dataType != DataType.timestamp)` | 213 |
| FOREACH | `foreach field in schema.Fields` | 219 |
| IF | `if (field.isInstanceOf(Field))` | 220 |
| VAR | `node<Field> f = (Field) field;` | 221 |
| IF (inline) | `if (f.dataType == DataType.timestamp)` | 222 |
| FOREACH | `foreach field in schema.Fields` | 229 |
| IF | `if (field.isInstanceOf(Field))` | 230 |
| VAR | `node<Field> f = (Field) field;` | 231 |
| FOREACH | `foreach field in schema.Fields` | 236 |
| IF | `if (field.isInstanceOf(Field))` | 237 |
| VAR | `node<Field> f = (Field) field;` | 238 |
| FOREACH | `foreach field in schema.Fields` | 246 |
| IF | `if (field.isInstanceOf(Field))` | 247 |
| VAR | `node<Field> f = (Field) field;` | 248 |
| FOREACH | `foreach field in schema.Fields` | 253 |
| IF | `if (field.isInstanceOf(Field))` | 254 |
| VAR | `node<Field> f = (Field) field;` | 255 |
| VAR | `int uidx = 0;` | 263 |
| FOREACH | `foreach field in schema.Fields` | 264 |
| IF | `if (field.isInstanceOf(Field))` | 265 |
| VAR | `node<Field> f = (Field) field;` | 266 |
| IF | `if (f.dataType != DataType.timestamp)` | 267 |
| IF (inline) | `if (uidx > 0)` | 268 |
| ASSIGN | `uidx = uidx + 1;` | 269 |
| VAR | `int nextParam = uidx + 1;` | 274 |
| FOREACH | `foreach field in schema.Fields` | 276 |
| IF | `if (field.isInstanceOf(Field))` | 277 |
| VAR | `node<Field> f = (Field) field;` | 278 |
| IF (inline) | `if (f.dataType != DataType.timestamp)` | 279 |
| FOREACH | `foreach schema in model.schemas` | 293 |
| IF | `if (schema.hasReferences())` | 294 |
| VAR | `string sn = schema.structName();` | 295 |
| VAR | `string rn = schema.repoName();` | 296 |
| VAR | `string tn = schema.name;` | 297 |
| VAR | `int argIdx = 0;` | 303 |
| FOREACH | `foreach field in schema.Fields` | 304 |
| IF | `if (field.isInstanceOf(FieldRefrence))` | 305 |
| VAR | `node<FieldRefrence> fr = (FieldRefrence) field;` | 306 |
| IF (inline) | `if (argIdx > 0)` | 307 |
| ASSIGN | `argIdx = argIdx + 1;` | 309 |
| VAR | `int fkIdx = 0;` | 314 |
| FOREACH | `foreach field in schema.Fields` | 315 |
| IF | `if (field.isInstanceOf(FieldRefrence))` | 316 |
| VAR | `node<FieldRefrence> fr = (FieldRefrence) field;` | 317 |
| IF (inline) | `if (fkIdx > 0)` | 318 |
| ASSIGN | `fkIdx = fkIdx + 1;` | 320 |
| FOR | `for (int i = 1; i <= fkIdx; i++)` | 324 |
| IF (inline) | `if (i > 1)` | 325 |
| VAR | `int ckIdx = 0;` | 329 |
| FOREACH | `foreach field in schema.Fields` | 330 |
| IF | `if (field.isInstanceOf(FieldRefrence))` | 331 |
| VAR | `node<FieldRefrence> fr = (FieldRefrence) field;` | 332 |
| IF (inline) | `if (ckIdx > 0)` | 333 |
| ASSIGN | `ckIdx = ckIdx + 1;` | 335 |
| VAR | `int retIdx = 0;` | 339 |
| FOREACH | `foreach field in schema.Fields` | 340 |
| IF (inline) | `if (retIdx > 0)` | 341 |
| IF | `if (field.isInstanceOf(FieldRefrence))` | 342 |
| VAR | `node<FieldRefrence> fr = (FieldRefrence) field;` | 343 |
| IF | `if (field.isInstanceOf(Field))` | 346 |
| VAR | `node<Field> f = (Field) field;` | 347 |
| ASSIGN | `retIdx = retIdx + 1;` | 350 |
| FOREACH | `foreach field in schema.Fields` | 353 |
| IF | `if (field.isInstanceOf(FieldRefrence))` | 354 |
| VAR | `node<FieldRefrence> fr = (FieldRefrence) field;` | 355 |
| VAR | `int scanIdx = 0;` | 360 |
| FOREACH | `foreach field in schema.Fields` | 361 |
| IF (inline) | `if (scanIdx > 0)` | 362 |
| IF | `if (field.isInstanceOf(FieldRefrence))` | 363 |
| VAR | `node<FieldRefrence> fr = (FieldRefrence) field;` | 364 |
| IF | `if (field.isInstanceOf(Field))` | 367 |
| VAR | `node<Field> f = (Field) field;` | 368 |
| ASSIGN | `scanIdx = scanIdx + 1;` | 371 |
| VAR | `int rmIdx = 0;` | 377 |
| FOREACH | `foreach field in schema.Fields` | 378 |
| IF | `if (field.isInstanceOf(FieldRefrence))` | 379 |
| VAR | `node<FieldRefrence> fr = (FieldRefrence) field;` | 380 |
| IF (inline) | `if (rmIdx > 0)` | 381 |
| ASSIGN | `rmIdx = rmIdx + 1;` | 383 |
| VAR | `int wIdx = 0;` | 387 |
| FOREACH | `foreach field in schema.Fields` | 388 |
| IF | `if (field.isInstanceOf(FieldRefrence))` | 389 |
| VAR | `node<FieldRefrence> fr = (FieldRefrence) field;` | 390 |
| IF (inline) | `if (wIdx > 0)` | 391 |
| ASSIGN | `wIdx = wIdx + 1;` | 392 |
| FOREACH | `foreach field in schema.Fields` | 397 |
| IF | `if (field.isInstanceOf(FieldRefrence))` | 398 |
| VAR | `node<FieldRefrence> fr = (FieldRefrence) field;` | 399 |
| FOREACH | `foreach refA in schema.Fields` | 406 |
| IF | `if (refA.isInstanceOf(FieldRefrence))` | 407 |
| VAR | `node<FieldRefrence> frA = (FieldRefrence) refA;` | 408 |
| FOREACH | `foreach refB in schema.Fields` | 409 |
| IF | `if (refB.isInstanceOf(FieldRefrence) && refB != refA)` | 410 |
| VAR | `node<FieldRefrence> frB = (FieldRefrence) refB;` | 411 |
| VAR | `string ts = frB.target_schema.structName();` | 412 |
| VAR | `string tt = frB.target_schema.name;` | 413 |
| FOREACH | `foreach tf in frB.target_schema.Fields` | 416 |
| IF | `if (tf.isInstanceOf(Field))` | 417 |
| VAR | `node<Field> f = (Field) tf;` | 418 |
| FOREACH | `foreach tf in frB.target_schema.Fields` | 423 |
| IF | `if (tf.isInstanceOf(Field))` | 424 |
| VAR | `node<Field> f = (Field) tf;` | 425 |
| FOREACH | `foreach schema in model.schemas` | 441 |
| IF | `if (!(schema.hasReferences()))` | 442 |
| VAR | `string sn = schema.structName();` | 443 |
| VAR | `string rn = schema.repoName();` | 444 |
| FOREACH | `foreach schema in model.schemas` | 469 |
| IF | `if (schema.hasReferences())` | 470 |
| VAR | `string sn = schema.structName();` | 471 |
| VAR | `string rn = schema.repoName();` | 472 |
| VAR | `node<FieldRefrence> firstRef = null;` | 476 |
| VAR | `node<FieldRefrence> secondRef = null;` | 477 |
| FOREACH | `foreach field in schema.Fields` | 478 |
| IF | `if (field.isInstanceOf(FieldRefrence))` | 479 |
| VAR | `node<FieldRefrence> fr = (FieldRefrence) field;` | 480 |
| IF+STMT (inline) | `if (firstRef == null) { firstRef = fr; }` | 481 |
| IF+STMT (inline) | `else if (secondRef == null) { secondRef = fr; }` | 482 |
| FOREACH | `foreach refA in schema.Fields` | 493 |
| IF | `if (refA.isInstanceOf(FieldRefrence))` | 494 |
| VAR | `node<FieldRefrence> frA = (FieldRefrence) refA;` | 495 |
| FOREACH | `foreach refB in schema.Fields` | 496 |
| IF | `if (refB.isInstanceOf(FieldRefrence) && refB != refA)` | 497 |
| VAR | `node<FieldRefrence> frB = (FieldRefrence) refB;` | 498 |
| FOREACH | `foreach schema in model.schemas` | 514 |
| IF | `if (!(schema.hasReferences()))` | 515 |
| FOREACH | `foreach schema in model.schemas` | 519 |
| IF | `if (schema.hasReferences())` | 520 |
| FOREACH | `foreach schema in model.schemas` | 528 |
| IF | `if (!(schema.hasReferences()))` | 529 |
| VAR | `string vr = schema.singularName() + "Repo";` | 530 |
| VAR | `string sn = schema.structName();` | 531 |
| VAR | `string tn = schema.name;` | 532 |
| FOREACH | `foreach schema in model.schemas` | 538 |
| IF | `if (schema.hasReferences())` | 539 |
| VAR | `string vr = schema.singularName() + "Repo";` | 540 |
| VAR | `node<FieldRefrence> fRef = null;` | 541 |
| VAR | `node<FieldRefrence> sRef = null;` | 542 |
| FOREACH | `foreach field in schema.Fields` | 543 |
| IF | `if (field.isInstanceOf(FieldRefrence))` | 544 |
| VAR | `node<FieldRefrence> fr = (FieldRefrence) field;` | 545 |
| IF+STMT (inline) | `if (fRef == null) { fRef = fr; }` | 546 |
| IF+STMT (inline) | `else if (sRef == null) { sRef = fr; }` | 547 |
