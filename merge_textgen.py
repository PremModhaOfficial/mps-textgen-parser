#!/usr/bin/env python3
"""
merge_textgen.py — Merges multiple generated .textGen.mps files into one.

Usage:
  python3 merge_textgen.py --target <base.textGen.mps> --sources <a.mps> <b.mps> ...

Extracts the top-level <node concept="WtQ9Q"> block from each source and
appends them into the target before </model>.
Deletes source files after merging.
"""
import sys
import re
import os
import argparse

def extract_wtq9q_node(content):
    """Extract the top-level <node concept="WtQ9Q">...</node> block."""
    # Find start
    start = content.find('<node concept="WtQ9Q"')
    if start == -1:
        return None
    # Walk forward counting open/close tags to find matching </node>
    depth = 0
    i = start
    while i < len(content):
        open_m = content.find('<node', i)
        close_m = content.find('</node>', i)
        selfclose_m = content.find('/>', i)
        # find next event
        candidates = [(p, t) for p, t in [
            (open_m, 'open'), (close_m, 'close'), (selfclose_m, 'self')
        ] if p != -1]
        if not candidates:
            break
        pos, typ = min(candidates, key=lambda x: x[0])
        if typ == 'open':
            # only count if not self-closing on same position
            sc = content.find('/>', pos)
            nl = content.find('>', pos)
            if sc != -1 and sc == nl:
                # self-closing, don't increment depth
                i = sc + 2
            else:
                depth += 1
                i = nl + 1
        elif typ == 'self':
            i = pos + 2
        elif typ == 'close':
            depth -= 1
            if depth == 0:
                end = pos + len('</node>')
                return content[start:end]
            i = pos + len('</node>')
    return None

def merge(target_path, source_paths, delete_sources=True):
    with open(target_path, 'r', encoding='utf-8') as f:
        target = f.read()

    # Find insertion point — just before </model>
    insert_pos = target.rfind('</model>')
    if insert_pos == -1:
        print(f"ERROR: </model> not found in {target_path}", file=sys.stderr)
        sys.exit(1)

    # Check which concept IDs already exist to avoid duplicates
    existing_ids = set(re.findall(r'<node concept="WtQ9Q" id="([^"]+)"', target))

    nodes_to_insert = []
    for src in source_paths:
        with open(src, 'r', encoding='utf-8') as f:
            content = f.read()
        node = extract_wtq9q_node(content)
        if node is None:
            print(f"WARNING: No WtQ9Q node found in {src}", file=sys.stderr)
            continue
        # Check for duplicate ID
        id_match = re.search(r'id="([^"]+)"', node)
        if id_match and id_match.group(1) in existing_ids:
            print(f"SKIP: {src} — id {id_match.group(1)} already in target")
            continue
        nodes_to_insert.append((src, node))
        print(f"  + Extracted WtQ9Q from {os.path.basename(src)}")

    if not nodes_to_insert:
        print("Nothing to merge.")
        return

    # Also merge imports — collect imports from sources not already in target
    existing_imports = set(re.findall(r'<import[^/]*/>', target))
    extra_imports = []
    for src, _ in nodes_to_insert:
        with open(src, 'r', encoding='utf-8') as f:
            content = f.read()
        for imp in re.findall(r'<import[^/]*/>', content):
            if imp not in existing_imports:
                existing_imports.add(imp)
                extra_imports.append(imp)

    # Insert extra imports before </imports>
    if extra_imports:
        imports_end = target.find('</imports>')
        if imports_end != -1:
            imp_block = '\n'.join(f'  {i}' for i in extra_imports) + '\n'
            target = target[:imports_end] + imp_block + target[imports_end:]
            # Recalculate insert_pos after modification
            insert_pos = target.rfind('</model>')

    # Insert all WtQ9Q nodes before </model>
    insertion = '\n' + '\n'.join(f'  {node}' for _, node in nodes_to_insert) + '\n'
    target = target[:insert_pos] + insertion + target[insert_pos:]

    with open(target_path, 'w', encoding='utf-8') as f:
        f.write(target)

    print(f"✓ Merged {len(nodes_to_insert)} textgen component(s) into {os.path.basename(target_path)}")

    if delete_sources:
        for src, _ in nodes_to_insert:
            os.remove(src)
            print(f"  - Deleted {os.path.basename(src)}")

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Merge .textGen.mps files into one')
    parser.add_argument('--target', required=True, help='Base .textGen.mps to merge into')
    parser.add_argument('--sources', nargs='+', required=True, help='Source .mps files to merge from')
    parser.add_argument('--keep', action='store_true', help='Keep source files after merge')
    args = parser.parse_args()

    merge(args.target, args.sources, delete_sources=not args.keep)
