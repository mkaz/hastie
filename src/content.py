"""Hastie resources, pages and categories."""

import frontmatter
from markdown import markdown
import os
from pathlib import Path
from typing import Dict, List
import sys

from config import config
import utils


def get_page(filename: os.PathLike) -> Dict:
    """Read page in from file system, parse front matter and render markdown."""
    page = {}
    exts = ["codehilite", "fenced_code", "tables"]
    with open(filename, "r") as f:
        # check for frontmatter format
        fm = config.get("frontmatter", "yaml")
        match fm:
            case "toml":
                handler = frontmatter.default_handlers.TOMLHandler()
            case "json":
                handler = frontmatter.default_handlers.JSONHandler()
            case _:
                handler = frontmatter.default_handlers.YAMLHandler()

        try:
            page = frontmatter.load(f, handler=handler).to_dict()
        except Exception as err:
            print("Error parsing frontmatter")
            print(f"    Filename: {filename}")
            print(err)
            sys.exit(1)

    try:
        page["content"] = markdown(page.get("content", ""), extensions=exts)
    except Exception as err:
        print("Error converting markdown")
        print(f"    Filename: {filename}")
        print(err)
        sys.exit(1)

    page["filename"] = filename
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
            page["parent"] = get_parent_name(parent_path, content_dir)
            page["category"] = os.path.relpath(category_path, parent_path)

        page["name"] = Path(f.parent, f.stem)
        page["url"] = base_url + os.path.relpath(page["name"], start=content_dir)
        pages.append(page)
    return pages


def get_parent_name(p: os.PathLike, c: os.PathLike) -> str:
    name = os.path.relpath(p, start=c)
    if name == ".":
        name = ""
    return name


def gather_categories(content_dir: os.PathLike, base_url="/") -> List:
    categories = []

    paths = content_dir.glob("**/*")
    for p in paths:
        if not p.is_dir():
            continue  # skip

        parent_name = get_parent_name(p.parent, content_dir)
        name = os.path.relpath(p, start=p.parent)

        if parent_name == ".." or name == ".":
            continue  # skip

        index = Path(p, "index.md")
        if not index.is_file():
            continue  # skip - it's not a category without a index.md

        page = get_page(index)
        page["category"] = name

        category = {
            "name": name,
            "parent": parent_name,
            "page": page,
            "url": utils.urljoin([base_url, parent_name, name]),
        }
        categories.append(category)

    return categories
