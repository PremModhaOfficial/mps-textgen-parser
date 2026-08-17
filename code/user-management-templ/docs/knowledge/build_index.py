#!/usr/bin/env python3
"""
build_index.py — Generate _index.json from knowledge graph markdown files.

Walks docs/knowledge/**/*.md, extracts YAML frontmatter + [[wiki-links]],
and produces _index.json with nodes[], edges[], categories[], projects{}.

Also reports broken links (references to nodes that don't exist).

Usage: python3 docs/knowledge/build_index.py
"""

import json
import re
import sys
from datetime import datetime
from pathlib import Path

KNOWLEDGE_DIR = Path(__file__).parent
WIKI_LINK_RE = re.compile(r'\[\[([^\]]+)\]\]')
SKIP_FILES = {'_template.md', 'README.md', '_index.json'}

PROJECTS = {
    "v1": {"path": "/home/prem-modha/MPSProjects/UserManagement", "status": "failed-prototype"},
    "v2-uman": {"path": "/home/prem-modha/MPSProjects/UMAN", "status": "working-abandoned"},
    "v3": {"path": "/home/prem-modha/UserManagmentMps", "status": "production"},
    "experiments": {"path": "/home/prem-modha/MPSProjects", "status": "various"},
    "dls": {"path": "/home/prem-modha/projects/DLS", "status": "supplementary"},
}


def parse_frontmatter(content: str) -> dict:
    """Extract YAML frontmatter from markdown content."""
    if not content.startswith('---'):
        return {}
    end = content.find('---', 3)
    if end == -1:
        return {}
    fm_text = content[3:end].strip()
    result = {}
    for line in fm_text.split('\n'):
        if ':' in line:
            key, _, value = line.partition(':')
            key = key.strip()
            value = value.strip().strip('"').strip("'")
            if value.startswith('[') and value.endswith(']'):
                value = [v.strip().strip('"').strip("'") for v in value[1:-1].split(',') if v.strip()]
            result[key] = value
    return result


def extract_links(content: str) -> list[str]:
    """Extract all [[wiki-links]] from content."""
    return WIKI_LINK_RE.findall(content)


def extract_link_types(content: str) -> list[dict]:
    """Extract links with their type annotations from the Links section."""
    edges = []
    for line in content.split('\n'):
        links = WIKI_LINK_RE.findall(line)
        for link in links:
            link_type = "related"
            type_match = re.search(r'type:\s*([\w-]+)', line)
            if type_match:
                link_type = type_match.group(1)
            edges.append({"target": link, "type": link_type})
    return edges


def main():
    md_files = sorted(KNOWLEDGE_DIR.rglob('*.md'))

    nodes = []
    all_edges = []
    categories = set()
    broken_links = []
    node_files = set()

    for md_file in md_files:
        rel_path = md_file.relative_to(KNOWLEDGE_DIR)
        if rel_path.name in SKIP_FILES:
            continue
        if str(rel_path).startswith('images/'):
            continue

        content = md_file.read_text()
        fm = parse_frontmatter(content)

        if not fm.get('title'):
            continue

        file_key = str(rel_path).removesuffix('.md')
        node_files.add(file_key)

        category = fm.get('category', '')
        if category:
            categories.add(category)

        links = extract_links(content)
        typed_edges = extract_link_types(content)

        node = {
            "file": str(rel_path),
            "title": fm.get('title', rel_path.stem),
            "category": category,
            "project": fm.get('project', ''),
            "depth": fm.get('depth', ''),
            "tags": fm.get('tags', []),
            "links": links,
        }
        nodes.append(node)

        for edge in typed_edges:
            all_edges.append({
                "from": str(rel_path),
                "to": edge["target"],
                "type": edge["type"],
            })

    # Check for broken links
    for edge in all_edges:
        target = edge["to"]
        if target not in node_files and f"{target}" not in node_files:
            broken_links.append(f"{edge['from']} -> [[{target}]] (not found)")

    index = {
        "version": 1,
        "generated": datetime.now().isoformat(),
        "description": "DSL development learnings knowledge graph",
        "node_count": len(nodes),
        "categories": sorted(categories),
        "projects": PROJECTS,
        "nodes": nodes,
        "edges": all_edges,
    }

    out_path = KNOWLEDGE_DIR / '_index.json'
    out_path.write_text(json.dumps(index, indent=2) + '\n')
    print(f"Generated _index.json: {len(nodes)} nodes, {len(all_edges)} edges, {len(categories)} categories")

    if broken_links:
        print(f"\nBroken links ({len(broken_links)}):")
        for bl in broken_links:
            print(f"  - {bl}")
    else:
        print("No broken links.")


if __name__ == "__main__":
    main()
