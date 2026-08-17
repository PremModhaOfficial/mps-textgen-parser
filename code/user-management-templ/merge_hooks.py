#!/usr/bin/env python3
"""
merge_hooks.py — Extract hook stubs from generated .go files and merge into userDefinedHooks.go

MPS generates .go files with hook stubs between #HOOKS_START / #HOOKS_END markers.
This script:
  1. Finds all .go files in src/ with #HOOKS_START markers
  2. Extracts the hook stubs from between the markers
  3. Strips the hooks section from the generated file
  4. Appends only NEW stubs to userDefinedHooks.go (preserves existing user implementations)
"""

import re
import sys
import os
from pathlib import Path

SRC_DIR = Path(__file__).parent / "src"
HOOKS_FILE = SRC_DIR / "userDefinedHooks.go"

HOOKS_HEADER = """\
// userDefinedHooks.go
// AUTO-GENERATED STUBS — safe to edit.
// New hooks are appended automatically. Your implementations are preserved.
package main

import (
\t"context"

\t"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
)
"""

# Regex to match a Go method: func (receiver) name(params) [returns] { ... }
# Captures: (1) receiver type (e.g. "*UserHandler"), (2) function name
FUNC_SIG_RE = re.compile(r'^func\s+\(\w+\s+(\*?\w+)\)\s+(\w+)\s*\(', re.MULTILINE)


def extract_hooks_section(content: str) -> tuple[str, str]:
    """Extract content between ALL #HOOKS_START and #HOOKS_END pairs.
    Returns (cleaned_content_without_hooks, combined_hooks_section)."""

    start_marker = "// #HOOKS_START"
    end_marker = "// #HOOKS_END"

    all_hooks = []
    cleaned = content

    while True:
        start_idx = cleaned.find(start_marker)
        end_idx = cleaned.find(end_marker)

        if start_idx == -1 or end_idx == -1:
            break

        hooks_section = cleaned[start_idx + len(start_marker):end_idx].strip()
        if hooks_section:
            all_hooks.append(hooks_section)

        cleaned = cleaned[:start_idx].rstrip() + "\n" + cleaned[end_idx + len(end_marker):].lstrip("\n")

    return cleaned, "\n\n".join(all_hooks)


def make_func_key(receiver_type: str, func_name: str) -> str:
    """Create a unique key from receiver type + function name.
    e.g. '*UserHandler.preCreateD'"""
    return f"{receiver_type}.{func_name}"


def split_into_functions(hooks_text: str) -> list[tuple[str, str]]:
    """Split a block of Go code into individual functions.
    Returns list of (func_key, full_func_text) where func_key is 'ReceiverType.funcName'."""

    functions = []
    lines = hooks_text.split("\n")
    current_func = []
    current_key = ""
    brace_depth = 0
    in_func = False

    for line in lines:
        # Check if this line starts a new function
        match = FUNC_SIG_RE.match(line)
        if match and not in_func:
            current_key = make_func_key(match.group(1), match.group(2))
            current_func = [line]
            brace_depth = line.count("{") - line.count("}")
            in_func = True
            if brace_depth == 0:
                # Single-line function
                functions.append((current_key, "\n".join(current_func)))
                current_func = []
                current_key = ""
                in_func = False
            continue

        if in_func:
            current_func.append(line)
            brace_depth += line.count("{") - line.count("}")
            if brace_depth <= 0:
                functions.append((current_key, "\n".join(current_func)))
                current_func = []
                current_key = ""
                in_func = False

    return functions


def normalize_signature(sig_line: str) -> str:
    """Normalize a Go function signature line for comparison.
    Strips whitespace variations so 'func (s *X) foo() error {' matches itself."""
    return " ".join(sig_line.split())


def get_existing_functions(filepath: Path) -> dict[str, tuple[str, str]]:
    """Read a Go file and return dict[func_name → (signature_line, full_func_text)]."""
    if not filepath.exists():
        return {}

    content = filepath.read_text()
    funcs = split_into_functions(content)
    result = {}
    for name, text in funcs:
        first_line = text.split("\n")[0]
        result[name] = (first_line, text)
    return result


def replace_function_in_file(filepath: Path, old_func_text: str, new_func_text: str):
    """Replace an existing function in a file with new text."""
    content = filepath.read_text()
    content = content.replace(old_func_text, new_func_text)
    filepath.write_text(content)


def main():
    # Find all .go files with #HOOKS_START
    go_files = list(SRC_DIR.glob("*.go"))

    all_new_stubs = []
    files_processed = 0

    for go_file in go_files:
        if go_file.name == "userDefinedHooks.go":
            continue

        content = go_file.read_text()
        if "// #HOOKS_START" not in content:
            continue

        files_processed += 1
        cleaned, hooks_section = extract_hooks_section(content)

        # Write back the cleaned file (without hooks)
        go_file.write_text(cleaned)
        print(f"  {go_file.name}: stripped hooks section")

        if hooks_section:
            funcs = split_into_functions(hooks_section)
            all_new_stubs.extend(funcs)

    if files_processed == 0:
        print("No files with #HOOKS_START found. Nothing to do.")
        return

    # Create hooks file if it doesn't exist
    if not HOOKS_FILE.exists():
        HOOKS_FILE.write_text(HOOKS_HEADER)
        print(f"  Created {HOOKS_FILE.name}")

    # Read existing hooks file
    existing_funcs = get_existing_functions(HOOKS_FILE)

    # Merge stubs: add new, preserve matching, replace on signature mismatch
    added = 0
    skipped = 0
    replaced = 0

    with open(HOOKS_FILE, "a") as f:
        for func_key, func_text in all_new_stubs:
            if func_key not in existing_funcs:
                f.write(f"\n{func_text}\n")
                existing_funcs[func_key] = (func_text.split("\n")[0], func_text)
                added += 1
                continue

            old_sig, old_text = existing_funcs[func_key]
            new_sig = func_text.split("\n")[0]

            if normalize_signature(old_sig) == normalize_signature(new_sig):
                skipped += 1
            else:
                # Signature changed — replace with new stub, print old body as warning
                print(f"  ⚠ SIGNATURE CHANGED: {func_key}")
                print(f"    old: {old_sig.strip()}")
                print(f"    new: {new_sig.strip()}")
                print(f"    Old implementation (copy back if needed):")
                for line in old_text.split("\n"):
                    print(f"      {line}")
                replace_function_in_file(HOOKS_FILE, old_text, func_text)
                existing_funcs[func_key] = (new_sig, func_text)
                replaced += 1

    print(f"  Hooks: {added} added, {skipped} preserved, {replaced} replaced (signature changed)")


if __name__ == "__main__":
    main()
