"""Hastie resources, pages and categories."""

import frontmatter
from markdown import markdown
import os
from pathlib import Path
from typing import Dict, List


def get_page(filename: os.PathLike) -> Dict:
    """Read page in from file system, parse front matter and render markdown."""
    page = {}
    exts = ["fenced_code", "tables"]
    with open(filename, "r") as f:
        page = frontmatter.load(f)

    page["filename"] = filename
    page.content = markdown(page.content, extensions=exts)
    return page


def gather_pages(content_dir: os.PathLike, base_url="/") -> List:
    pages = []
    files = content_dir.glob("**/*.md")
    for f in files:
        if f.name == "index.md":
            continue

        page = get_page(f)
        page["category"] = ""  # default
        page["parent"] = ""  # default

        if f.parent != content_dir:
            # in a subdirectory, thus a category
            category_path = f.parent
            parent_path = category_path.parent
            page["parent"] = os.path.relpath(parent_path, content_dir)
            page["category"] = os.path.relpath(category_path, parent_path)

        page["name"] = f.with_suffix(".html")
        page["url"] = base_url + os.path.relpath(page["name"], start=content_dir)
        pages.append(page)
    return pages


def gather_categories(content_dir: os.PathLike, base_url="/") -> List:
    categories = []

    paths = content_dir.glob("**/*")
    for p in paths:
        if not p.is_dir():
            continue  # skip

        parent = os.path.relpath(p.parent, start=content_dir)
        name = os.path.relpath(p, start=p.parent)

        if parent == "..":
            continue

        if parent == ".":
            parent = ""

        if name == ".":
            continue  # skip

        index = Path(p, "index.md")
        if not index.is_file():
            continue  # skip - it's not a category without a index.md

        page = get_page(index)
        page["category"] = name
        if parent:
            url = base_url + parent + "/" + name + "/"
        else:
            url = base_url + name + "/"

        category = {
            "name": name,
            "parent": parent,
            "page": page,
            "url": url,
        }
        categories.append(category)

    return categories
