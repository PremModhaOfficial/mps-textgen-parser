# TextGen Parser — MPS TextGen from Pseudo-Code

Converts TextGen pseudo-code (the `append {...} ${expr} {...} ;` format used in MPS) into a valid `.mps` XML file that MPS can load directly.

## What It Does

- **Constants** (`{literal text}`) become `ConstantStringAppendPart` nodes — fully automated
- **Expressions** (`${schema.structName()}`) become `ConstantStringAppendPart` with `▶expr◀` markers — you swap them manually in MPS
- **Newlines** (`\n`) become `NewLineAppendPart` nodes — fully automated
- **Tabs** (`\t`) are preserved as literal tabs in constants
- **Control flow** (`foreach`, `if`, `for`, variable declarations) become XML comments as a guide for manual wrapping in MPS

## Installation

No dependencies. Python 3.6+.

```bash
python3 parse_textgen.py --help
```

## Usage

### Step 1: Create a dummy MPS language project

Create a throwaway MPS project. Add a Structure concept (can be empty) and a TextGen aspect. This gives you the UUIDs.

### Step 2: Get the UUIDs

```bash
grep "model ref" /path/to/YourLang.structure.mps
# r:40b69eb5-3c02-4fef-92ca-f6d48dc1801a(YourLang.structure)

grep "model ref" /path/to/YourLang.textGen.mps
# r:16369af0-98d1-4f43-8cee-8de735066b4d(YourLang.textGen)
```

Get the concept ID from the structure file:
```bash
grep 'concept="1TIwiD"' /path/to/YourLang.structure.mps
# id="5Hvr0stUi4w" — this is your concept-id
```

### Step 3: Run the parser

```bash
python3 parse_textgen.py appends.md \
  --model-uuid "r:16369af0-...(YourLang.textGen)" \
  --structure-uuid "r:40b69eb5-...(YourLang.structure)" \
  --concept-id "5Hvr0stUi4w" \
  --concept-name "Code" \
  -o /path/to/YourLang/models/YourLang.textGen.mps
```

### Step 4: Open in MPS

Reload the project. All constants appear in the TextGen editor.

### Step 5: Fix manually in MPS

1. Search for `▶` — each is an expression placeholder. Replace the ConstantStringAppendPart with a real NodeAppendPart expression.
2. Use the generated `_guide.md` file as a checklist for expressions and control flow.
3. Wrap sections in `foreach`/`if` blocks using MPS intentions (select nodes → wrap).

## CLI Flags

| Flag | Default | Description |
|---|---|---|
| `input` | (required) | Input file with TextGen pseudo-code |
| `--model-uuid` | `TODO-MODEL-UUID` | Full model ref from `YourLang.textGen.mps` line 2 |
| `--structure-uuid` | `TODO-STRUCTURE-UUID` | Full model ref from `YourLang.structure.mps` line 2 |
| `--concept-name` | `Code` | Name of the concept this TextGen is for |
| `--concept-id` | `TODO-CONCEPT-ID` | Node ID of the concept in the structure model |
| `-o`, `--output` | `output.textGen.mps` | Output `.mps` file path |

## Input Format

The parser reads any text file and looks for code inside triple-backtick fences. The pseudo-code syntax:

### Append statements

```
append {package main} ;                          → constant "package main"
append {type } ${schema.structName()} { struct} ; → constant "type " + expr + constant " struct"
append ${f.name} ;                                → expr only
```

**Parsing rule:** `} ` (brace + space) = segment boundary. `}` without space = literal Go brace content.

### Escape sequences

| Sequence | Result |
|---|---|
| `\n` | Separate `NewLineAppendPart` node |
| `\t` | Literal tab character in constant |

### Control flow (becomes XML comments)

```
foreach schema in model.schemas {     → <!-- ═══ FOREACH: ... ═══ -->
if (!(schema.hasReferences())) {      → <!-- ─── IF: ... ─── -->
for (int i = 1; i <= idx; i++) {      → <!-- ═══ FOR: ... ═══ -->
}                                     → <!-- END FOREACH/IF/FOR -->
string sn = schema.structName();      → <!-- VAR: ... -->
idx = idx + 1;                        → <!-- ASSIGN: ... -->
```

### Inline if + append

```
if (idx > 0) { append {, } ; }       → <!-- IF: ... --> constant "," <!-- END IF -->
```

## Output Format

### 1. `.mps` file (MPS XML)

Valid MPS Persistence v9 XML with:
- Full registry section (baseLanguage + textGen concepts)
- `ConceptTextGenDeclaration` with filename function and body
- Every constant as `ConstantStringAppendPart`
- Every expression as `ConstantStringAppendPart` with `▶expr◀` markers
- Every `\n` as `NewLineAppendPart`
- Control flow as XML comments

### 2. `_guide.md` file (fixup checklist)

Generated alongside the `.mps` file. Contains:
- **Expressions table** — every `▶expr◀` placeholder with type classification:
  - Type A: property access (`schema.name`)
  - Type B: method call (`schema.structName()`)
  - Type C: chained access (`fr.target_schema.structName()`)
  - Type D: local variable reference (`sn`)
- **Control flow table** — every foreach/if/for/var with source line numbers

## Workflow: Dummy Project Buffer

To avoid corrupting existing TextGen declarations:

1. Create a dummy MPS project
2. Run the parser against the dummy project's textGen file
3. Open in MPS — verify constants render
4. Copy the body content from the dummy project
5. Paste into your real project's TextGen
6. Fix expressions and control flow in the real project

This way you never touch the real project's file directly.

## Limitations

- **No expression trees** — all `${...}` become literal text placeholders, not real MPS NodeAppendPart expressions
- **No structure refs** — the concept-id must be provided manually (from the structure model)
- **Single concept** — generates one `ConceptTextGenDeclaration` per run
- **Filename hardcoded** — the filename function returns `"generated"` — change it in MPS after import
- **Control flow is comments only** — you manually wrap nodes in foreach/if/for in MPS's editor

## Expression Types Reference

When fixing `▶expr◀` placeholders in MPS:

| Type | Example | MPS Action |
|---|---|---|
| A: property | `▶schema.name◀` | NodeAppendPart → `schema.name` |
| B: method | `▶schema.structName()◀` | NodeAppendPart → `schema.structName()` |
| C: chained | `▶fr.target_schema.structName()◀` | NodeAppendPart → `fr.target_schema.structName()` |
| D: local var | `▶sn◀` | NodeAppendPart → `sn` |
