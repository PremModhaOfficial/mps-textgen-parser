#!/usr/bin/env python3
"""
parse_textgen.py — Parse TextGen pseudo-code into MPS .mps XML

Every segment (constants AND expressions) becomes ConstantStringAppendPart.
Expressions get ▶ ◀ markers so you can find/replace them in MPS.
Control flow (foreach/if/for/var) becomes XML comments as a guide.

Parsing rule:
  } <space>  = segment boundary
  }          = literal content (no space after = part of Go code)

Usage:
  python3 parse_textgen.py appends.md -o GoTextGen.textGen.mps
  python3 parse_textgen.py appends.md --concept-name Code --concept-id TODO
"""

import re
import sys
import argparse
from xml.sax.saxutils import escape as xml_escape


def xml_comment_safe(s):
    """Sanitize string for use inside XML comments (-- is illegal)."""
    return s.replace('--', '~~')


# ── ID Generator ──────────────────────────────────────────────

class IdGen:
    """Generate MPS-compatible node IDs (alphanumeric, like 2aSUx1NZjnH)."""
    CHARS = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'

    def __init__(self):
        self.n = 0

    def next(self, prefix="n"):
        self.n += 1
        # Generate a base62 ID from counter
        num = self.n + 100000  # offset to get consistent length
        result = []
        while num > 0:
            result.append(self.CHARS[num % 62])
            num //= 62
        return '7tgPrs' + ''.join(reversed(result))


# ── Tokenizer: split append content into segments ────────────

def tokenize_append(raw):
    """Split raw append content (after 'append ', before ' ;') into segments.

    Returns list of ('const', text) or ('expr', text).

    Rule: } followed by space = boundary.  } followed by anything else = literal.
    """
    segs = []
    i = 0
    while i < len(raw):
        if i + 1 < len(raw) and raw[i] == '$' and raw[i + 1] == '{':
            # ── Expression: ${...} ──
            j = raw.index('}', i + 2)
            segs.append(('expr', raw[i + 2:j]))
            i = j + 1
            if i < len(raw) and raw[i] == ' ':
                i += 1  # skip separator space
        elif raw[i] == '{':
            # ── Constant: {...} ──
            j = i + 1
            while j < len(raw):
                if raw[j] == '}':
                    # Boundary: } at end, or } followed by space
                    if j + 1 >= len(raw) or raw[j + 1] == ' ':
                        segs.append(('const', raw[i + 1:j]))
                        j += 1
                        if j < len(raw) and raw[j] == ' ':
                            j += 1  # skip space
                        break
                    else:
                        j += 1  # } is content (e.g. }{ in Go code)
                else:
                    j += 1
            i = j
        else:
            i += 1
    return segs


# ── Escape processor ─────────────────────────────────────────

def expand_escapes(text):
    r"""Split constant text on \n into parts. Convert \t to real tab.

    Returns [('text', str) | ('newline', None), ...]
    """
    parts = []
    chunks = text.split('\\n')
    for idx, chunk in enumerate(chunks):
        if chunk:
            parts.append(('text', chunk.replace('\\t', '\t')))
        if idx < len(chunks) - 1:
            parts.append(('newline', None))
    return parts


# ── Flatten segments into XML-ready parts ─────────────────────

def flatten(segments):
    """Convert tokenized segments into flat parts list."""
    parts = []
    for kind, val in segments:
        if kind == 'const':
            parts.extend(expand_escapes(val))
        else:  # expr
            parts.append(('expr', val))
    return parts


# ── XML node emitters ─────────────────────────────────────────

def emit_const(value, ids):
    v = xml_escape(value, {'"': '&quot;'})
    return (
        f'        <node concept="lc7rE" id="{ids.next("a")}" role="3cqZAp">\n'
        f'          <node concept="la8eA" id="{ids.next("c")}" role="lcghm">\n'
        f'            <property role="lacIc" value="{v}" />\n'
        f'          </node>\n'
        f'        </node>'
    )

def emit_expr(expr, ids):
    v = xml_escape(f'\u25b6{expr}\u25c0', {'"': '&quot;'})
    return (
        f'        <node concept="lc7rE" id="{ids.next("a")}" role="3cqZAp">\n'
        f'          <node concept="la8eA" id="{ids.next("x")}" role="lcghm">\n'
        f'            <property role="lacIc" value="{v}" />\n'
        f'          </node>\n'
        f'        </node>'
    )

def emit_newline(ids):
    return (
        f'        <node concept="lc7rE" id="{ids.next("a")}" role="3cqZAp">\n'
        f'          <node concept="l8MVK" id="{ids.next("nl")}" role="lcghm" />\n'
        f'        </node>'
    )

def emit_blank(ids):
    return f'        <node concept="3clFbH" id="{ids.next("s")}" role="3cqZAp" />'


def parts_to_xml(parts, ids):
    nodes = []
    for kind, val in parts:
        if kind == 'text':
            if val:
                nodes.append(emit_const(val, ids))
        elif kind == 'expr':
            nodes.append(emit_expr(val, ids))
        elif kind == 'newline':
            nodes.append(emit_newline(ids))
    return nodes


# ── Line classifier ──────────────────────────────────────────

def classify(line):
    """Returns (type, payload)."""
    s = line.strip()
    if not s:
        return 'blank', None

    # Skip markdown / structural
    if s in ('```', '```java'):
        return 'skip', None
    if s.startswith('text gen component'):
        return 'meta', s
    if s.startswith('file name :'):
        return 'meta', s
    if s.startswith('file path :'):
        return 'meta', s
    if s.startswith('extension :'):
        return 'meta', s
    if s.startswith('(node)->'):
        return 'body_start', s

    # ── Inline if/else-if with append ──
    # Pattern: if (...) { append ... ; }   OR   else if (...) { ... ; }
    if s.startswith(('if ', 'if(', 'else if')) and s.endswith('; }'):
        bp = s.find(') {')
        if bp >= 0:
            cond = s[:bp + 1]
            inner = s[bp + 3:-1].strip()  # between { and }
            if inner.startswith('append '):
                return 'inline_if_append', (cond, inner)
            else:
                return 'inline_if_ctrl', (cond, inner)

    # ── Append ──
    if s.startswith('append '):
        return 'append', s

    # ── Control flow (block openers) ──
    if s.startswith('foreach '):
        return 'foreach', s
    if s.startswith('if ') or s.startswith('if('):
        return 'if', s
    if s.startswith('else if'):
        return 'else_if', s
    if s.startswith('for (') or s.startswith('for('):
        return 'for', s

    # ── Block close ──
    if s == '}':
        return 'end', None

    # ── Variable / assignment ──
    if re.match(r'^(node<|string |int |boolean )', s):
        return 'var', s
    if re.match(r'^[a-zA-Z_]\w*\s*=\s*.+;$', s):
        return 'assign', s

    # ── Comments ──
    if s.startswith('//'):
        return 'comment', s

    return 'unknown', s


# ── Parse append line ─────────────────────────────────────────

def parse_append(line):
    """Extract raw content from append statement, tokenize, flatten."""
    s = line.strip()
    s = s[7:]  # strip 'append '
    if s.endswith(' ;'):
        s = s[:-2]
    segs = tokenize_append(s)
    return flatten(segs), segs


# ── Main ──────────────────────────────────────────────────────

def main():
    ap = argparse.ArgumentParser(description='Parse TextGen pseudo-code → MPS XML')
    ap.add_argument('input', help='Input file (e.g. appends.md)')
    ap.add_argument('--model-uuid', default='TODO-MODEL-UUID')
    ap.add_argument('--model-name', default='YourLang.textGen')
    ap.add_argument('--structure-uuid', default='TODO-STRUCTURE-UUID')
    ap.add_argument('--structure-name', default='YourLang.structure')
    ap.add_argument('--concept-name', default='Code')
    ap.add_argument('--concept-id', default='TODO-CONCEPT-ID')
    ap.add_argument('-o', '--output', default='output.textGen.mps')
    args = ap.parse_args()

    with open(args.input) as f:
        lines = f.readlines()

    ids = IdGen()
    body = []           # XML node strings
    exprs = []          # (num, expr_text, src_line) for guide
    ctrls = []          # (type, text, src_line) for guide
    meta_lines = []     # file metadata
    scope = []          # stack for END tracking
    expr_n = 0
    in_body = False
    in_fence = False

    for ln, raw in enumerate(lines, 1):
        s = raw.rstrip('\n')
        stripped = s.strip()

        # Track code fence
        if stripped in ('```', '```java'):
            in_fence = not in_fence
            continue
        if not in_fence:
            continue

        kind, payload = classify(s)

        if kind == 'meta':
            meta_lines.append((ln, payload))
            continue
        if kind == 'body_start':
            in_body = True
            continue
        if kind == 'skip':
            continue
        if not in_body:
            continue

        # ── Blank line → empty statement ──
        if kind == 'blank':
            body.append(emit_blank(ids))

        # ── Append ──
        elif kind == 'append':
            parts, segs = parse_append(payload)
            for pk, pv in parts:
                if pk == 'expr':
                    expr_n += 1
                    exprs.append((expr_n, pv, ln))
            body.extend(parts_to_xml(parts, ids))

        # ── Inline if + append ──
        elif kind == 'inline_if_append':
            cond, inner_append = payload
            body.append(f'        <!-- IF: {xml_comment_safe(xml_escape(cond))} -->')
            ctrls.append(('IF (inline)', cond, ln))

            parts, segs = parse_append(inner_append)
            for pk, pv in parts:
                if pk == 'expr':
                    expr_n += 1
                    exprs.append((expr_n, pv, ln))
            body.extend(parts_to_xml(parts, ids))

            body.append(f'        <!-- END IF -->')

        # ── Inline if + control (assignment etc.) ──
        elif kind == 'inline_if_ctrl':
            cond, inner = payload
            body.append(f'        <!-- IF: {xml_comment_safe(xml_escape(cond))} {{ {xml_comment_safe(xml_escape(inner))} }} -->')
            ctrls.append(('IF+STMT (inline)', f'{cond} {{ {inner} }}', ln))

        # ── Foreach ──
        elif kind == 'foreach':
            # strip trailing {
            txt = stripped.rstrip(' {').rstrip('{')
            body.append(f'        <!-- ═══ FOREACH: {xml_comment_safe(xml_escape(txt))} ═══ -->')
            ctrls.append(('FOREACH', txt, ln))
            scope.append('FOREACH')

        # ── If ──
        elif kind == 'if':
            txt = stripped.rstrip(' {').rstrip('{')
            body.append(f'        <!-- ─── IF: {xml_comment_safe(xml_escape(txt))} ─── -->')
            ctrls.append(('IF', txt, ln))
            scope.append('IF')

        # ── Else if ──
        elif kind == 'else_if':
            txt = stripped.rstrip(' {').rstrip('{')
            body.append(f'        <!-- ─── ELSE IF: {xml_comment_safe(xml_escape(txt))} ─── -->')
            ctrls.append(('ELSE IF', txt, ln))
            # pop IF, push ELSE_IF (same level)
            if scope and scope[-1] in ('IF', 'ELSE_IF'):
                scope[-1] = 'ELSE_IF'
            else:
                scope.append('ELSE_IF')

        # ── For loop ──
        elif kind == 'for':
            txt = stripped.rstrip(' {').rstrip('{')
            body.append(f'        <!-- ═══ FOR: {xml_comment_safe(xml_escape(txt))} ═══ -->')
            ctrls.append(('FOR', txt, ln))
            scope.append('FOR')

        # ── End block ──
        elif kind == 'end':
            tag = scope.pop() if scope else '?'
            body.append(f'        <!-- END {tag} -->')

        # ── Variable declaration ──
        elif kind == 'var':
            body.append(f'        <!-- VAR: {xml_comment_safe(xml_escape(stripped))} -->')
            ctrls.append(('VAR', stripped, ln))

        # ── Assignment ──
        elif kind == 'assign':
            body.append(f'        <!-- ASSIGN: {xml_comment_safe(xml_escape(stripped))} -->')
            ctrls.append(('ASSIGN', stripped, ln))

        # ── Comment ──
        elif kind == 'comment':
            body.append(f'        <!-- {xml_comment_safe(xml_escape(stripped))} -->')

        # ── Unknown ──
        elif kind == 'unknown':
            body.append(f'        <!-- ?? {xml_comment_safe(xml_escape(stripped))} -->')

    # ── Assemble full XML ─────────────────────────────────────
    body_xml = '\n'.join(body)

    tg_id = ids.next("tg")
    gt_id = ids.next("gt")
    sl_id = ids.next("sl")
    fn_id = ids.next("fn")
    fn_sl = ids.next("fns")
    fn_ret = ids.next("fnr")
    fn_str = ids.next("fns2")

    xml = f'''<?xml version="1.0" encoding="UTF-8"?>
<model ref="{args.model_uuid}">
  <persistence version="9" />
  <languages>
    <use id="b83431fe-5c8f-40bc-8a36-65e25f4dd253" name="jetbrains.mps.lang.textGen" version="1" />
    <devkit ref="fa73d85a-ac7f-447b-846c-fcdc41caa600(jetbrains.mps.devkit.aspect.textgen)" />
  </languages>
  <imports>
    <import index="s1" ref="{args.structure_uuid}" implicit="true" />
  </imports>
  <registry>
    <language id="f3061a53-9226-4cc5-a443-f952ceaf5816" name="jetbrains.mps.baseLanguage">
      <concept id="1137021947720" name="jetbrains.mps.baseLanguage.structure.ConceptFunction" flags="in" index="2VMwT0">
        <child id="1137022507850" name="body" index="2VODD2" />
      </concept>
      <concept id="1070475926800" name="jetbrains.mps.baseLanguage.structure.StringLiteral" flags="nn" index="Xl_RD">
        <property id="1070475926801" name="value" index="Xl_RC" />
      </concept>
      <concept id="1068580123157" name="jetbrains.mps.baseLanguage.structure.Statement" flags="nn" index="3clFbH" />
      <concept id="1068580123136" name="jetbrains.mps.baseLanguage.structure.StatementList" flags="sn" stub="5293379017992965193" index="3clFbS">
        <child id="1068581517665" name="statement" index="3cqZAp" />
      </concept>
      <concept id="1068581242878" name="jetbrains.mps.baseLanguage.structure.ReturnStatement" flags="nn" index="3cpWs6">
        <child id="1068581517676" name="expression" index="3cqZAk" />
      </concept>
    </language>
    <language id="b83431fe-5c8f-40bc-8a36-65e25f4dd253" name="jetbrains.mps.lang.textGen">
      <concept id="45307784116571022" name="jetbrains.mps.lang.textGen.structure.FilenameFunction" flags="ig" index="29tfMY" />
      <concept id="1237305208784" name="jetbrains.mps.lang.textGen.structure.NewLineAppendPart" flags="ng" index="l8MVK" />
      <concept id="1237305557638" name="jetbrains.mps.lang.textGen.structure.ConstantStringAppendPart" flags="ng" index="la8eA">
        <property id="1237305576108" name="value" index="lacIc" />
      </concept>
      <concept id="1237306079178" name="jetbrains.mps.lang.textGen.structure.AppendOperation" flags="nn" index="lc7rE">
        <child id="1237306115446" name="part" index="lcghm" />
      </concept>
      <concept id="1233670071145" name="jetbrains.mps.lang.textGen.structure.ConceptTextGenDeclaration" flags="ig" index="WtQ9Q">
        <reference id="1233670257997" name="conceptDeclaration" index="WuzLi" />
        <child id="45307784116711884" name="filename" index="29tGrW" />
        <child id="1233749296504" name="textGenBlock" index="11c4hB" />
      </concept>
      <concept id="1233748055915" name="jetbrains.mps.lang.textGen.structure.NodeParameter" flags="nn" index="117lpO" />
      <concept id="1233749247888" name="jetbrains.mps.lang.textGen.structure.GenerateTextDeclaration" flags="in" index="11bSqf" />
    </language>
  </registry>
  <node concept="WtQ9Q" id="{tg_id}">
    <ref role="WuzLi" to="s1:{args.concept_id}" resolve="{args.concept_name}" />
    <node concept="29tfMY" id="{fn_id}" role="29tGrW">
      <node concept="3clFbS" id="{fn_sl}" role="2VODD2">
        <node concept="3cpWs6" id="{fn_ret}" role="3cqZAp">
          <node concept="Xl_RD" id="{fn_str}" role="3cqZAk">
            <property role="Xl_RC" value="generated" />
          </node>
        </node>
      </node>
    </node>
    <node concept="11bSqf" id="{gt_id}" role="11c4hB">
      <node concept="3clFbS" id="{sl_id}" role="2VODD2">
{body_xml}
      </node>
    </node>
  </node>
</model>
'''

    with open(args.output, 'w') as f:
        f.write(xml)

    # ── Generate guide ────────────────────────────────────────
    guide = args.output.replace('.mps', '_guide.md')
    with open(guide, 'w') as f:
        f.write(f'# TextGen Fixup Guide\n\n')
        f.write(f'Source: `{args.input}`\n')
        f.write(f'Output: `{args.output}`\n\n')

        if meta_lines:
            f.write(f'## File Metadata (add manually in MPS)\n\n')
            for mln, mtxt in meta_lines:
                f.write(f'- Line {mln}: `{mtxt}`\n')
            f.write('\n')

        f.write(f'## Expressions ({expr_n} total)\n\n')
        f.write(f'Search `▶` in MPS. Each is a ConstantStringAppendPart — replace content with actual expression.\n\n')
        f.write(f'| # | Expression | Line | Type |\n')
        f.write(f'|---|---|---|---|\n')
        for n, expr, ln in exprs:
            if '.' not in expr:
                t = 'D: local var'
            elif '(' in expr:
                t = 'C: chained' if expr.count('.') > 1 else 'B: method()'
            else:
                t = 'C: chained' if expr.count('.') > 1 else 'A: property'
            f.write(f'| {n} | `{expr}` | {ln} | {t} |\n')

        f.write(f'\n## Control Flow ({len(ctrls)} items)\n\n')
        f.write(f'These are XML comments. Wrap enclosed nodes in MPS.\n\n')
        f.write(f'| Type | Code | Line |\n')
        f.write(f'|---|---|---|\n')
        for ct, cc, cl in ctrls:
            f.write(f'| {ct} | `{cc}` | {cl} |\n')

    # ── Summary ───────────────────────────────────────────────
    print(f'✓ {args.output}')
    print(f'  {expr_n} expression placeholders (▶...◀)')
    print(f'  {len(ctrls)} control flow comments')
    print(f'  {ids.n} total XML nodes')
    print(f'✓ {guide}')


if __name__ == '__main__':
    main()
