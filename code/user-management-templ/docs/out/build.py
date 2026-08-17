#!/usr/bin/env python3
"""Build styled HTML + PDF from markdown/asciidoc sources with embedded images."""

import base64
import re
import subprocess
import sys
from pathlib import Path

OUT = Path(__file__).parent
CSS = (OUT / "style.css").read_text()

IMAGES_DIR = OUT.parent / "knowledge" / "images"


def embed_images(html: str, resource_dir: Path) -> str:
    """Replace <img src="images/..."> with base64-embedded images."""
    def replace_img(match):
        full = match.group(0)
        src_match = re.search(r'src="([^"]+)"', full)
        if not src_match:
            return full
        src = src_match.group(1)

        # Try multiple paths
        for candidate in [
            resource_dir / src,
            resource_dir / "images" / Path(src).name,
            IMAGES_DIR / Path(src).name,
        ]:
            if candidate.exists():
                data = base64.b64encode(candidate.read_bytes()).decode()
                ext = candidate.suffix.lstrip(".")
                mime = {"png": "image/png", "jpg": "image/jpeg", "jpeg": "image/jpeg", "gif": "image/gif", "svg": "image/svg+xml"}.get(ext, "image/png")
                return full.replace(f'src="{src}"', f'src="data:{mime};base64,{data}"')

        return full

    return re.sub(r'<img[^>]+>', replace_img, html)


def md_to_html(md_path: Path, title: str, resource_dir: Path) -> str:
    """Convert markdown to styled HTML with embedded images."""
    result = subprocess.run(
        ["pandoc", str(md_path), "--to=html5", "--no-highlight"],
        capture_output=True, text=True
    )
    body = result.stdout

    # Embed images
    body = embed_images(body, resource_dir)

    return f"""<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{title}</title>
<style>{CSS}</style>
</head>
<body>
{body}
</body>
</html>"""


ASCIIDOCTOR = Path.home() / ".local/share/gem/ruby/3.3.0/bin/asciidoctor"


def adoc_to_html(adoc_files: list[Path], title: str, resource_dir: Path) -> str:
    """Convert multiple AsciiDoc files to styled HTML using asciidoctor."""
    # Combine all adoc files into one with includes
    combined = f"= {title}\n:toc: left\n:icons: font\n:imagesdir: images\n:source-highlighter: coderay\n\n"
    for f in sorted(adoc_files):
        combined += f.read_text() + "\n\n"

    tmp = Path("/tmp/feynman-build.adoc")
    tmp.write_text(combined)

    # Use real asciidoctor for proper rendering
    result = subprocess.run(
        [str(ASCIIDOCTOR), str(tmp), "-o", "-", "-b", "html5",
         "-a", f"imagesdir={resource_dir / 'images'}",
         "-a", "nofooter"],
        capture_output=True, text=True
    )

    if result.returncode != 0:
        print(f"  asciidoctor error: {result.stderr[:300]}")

    body_html = result.stdout

    # Extract just the body content if it's a full HTML doc
    body_match = re.search(r'<body[^>]*>(.*)</body>', body_html, re.DOTALL)
    body = body_match.group(1) if body_match else body_html

    # Extract asciidoctor's own CSS
    style_match = re.search(r'<style>(.*?)</style>', body_html, re.DOTALL)
    adoc_css = style_match.group(1) if style_match else ""

    # Embed images
    body = embed_images(body, resource_dir)

    return f"""<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{title}</title>
<style>{adoc_css}
/* Custom overrides */
body {{ max-width: 960px; margin: 0 auto; padding: 2rem; }}
img {{ max-width: 100%; border: 1px solid #d0d7de; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }}
</style>
</head>
<body>
{body}
</body>
</html>"""


def html_to_pdf(html_path: Path, pdf_path: Path):
    """Convert HTML to PDF via weasyprint."""
    result = subprocess.run(
        ["weasyprint", str(html_path), str(pdf_path)],
        capture_output=True, text=True
    )
    if result.returncode != 0:
        print(f"  weasyprint warning: {result.stderr.strip()[:200]}")


def build_all():
    html_dir = OUT / "html"
    pdf_dir = OUT / "pdf"
    html_dir.mkdir(exist_ok=True)
    pdf_dir.mkdir(exist_ok=True)

    builds = [
        ("meeting", "Meeting Presentation", OUT / "meeting" / "README.md"),
        ("blog", "Blog Post — Building a Production DSL with JetBrains MPS", OUT / "blog" / "README.md"),
        ("onboarding", "Onboarding Guide — MotaDSLUM", OUT / "onboarding" / "README.md"),
    ]

    for slug, title, md_path in builds:
        print(f"Building {slug}...")
        resource_dir = md_path.parent
        html = md_to_html(md_path, title, resource_dir)
        html_path = html_dir / f"{slug}.html"
        html_path.write_text(html)
        print(f"  HTML: {html_path.name} ({len(html)//1024}K)")

        pdf_path = pdf_dir / f"{slug}.pdf"
        html_to_pdf(html_path, pdf_path)
        print(f"  PDF:  {pdf_path.name} ({pdf_path.stat().st_size//1024}K)")

    # Feynman guide (AsciiDoc)
    print("Building feynman-guide...")
    adoc_dir = OUT / "feynman-guide"
    adoc_files = list(adoc_dir.glob("0*.adoc"))
    html = adoc_to_html(adoc_files, "Feynman Technical Guide — MotaDSLUM", adoc_dir)
    html_path = html_dir / "feynman-guide.html"
    html_path.write_text(html)
    print(f"  HTML: {html_path.name} ({len(html)//1024}K)")

    pdf_path = pdf_dir / "feynman-guide.pdf"
    html_to_pdf(html_path, pdf_path)
    print(f"  PDF:  {pdf_path.name} ({pdf_path.stat().st_size//1024}K)")

    print("\nDone! All outputs in:")
    print(f"  HTML: {html_dir}/")
    print(f"  PDF:  {pdf_dir}/")


if __name__ == "__main__":
    build_all()
