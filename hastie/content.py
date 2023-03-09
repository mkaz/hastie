"""
Hastie resources, pages and categories.
"""

import os
import sys
from pathlib import Path
from typing import Dict, List

import frontmatter
from markdown import markdown

import hastie.utils as utils


def get_page(filename: Path, config: Dict) -> Dict:
    """Read page in from file system, parse frontmatter and render markdown."""

    try:
        page = read_page(filename, config)
        page["content"] = process_markdown(page.get("content", ""))
    except Exception as err:
        print(f"Error reading page {filename}")
        print(err)
        sys.exit(1)

    return page


def read_page(filename: Path, config: Dict = {}) -> Dict:
    """Read page using frontmatter library."""

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

        page = frontmatter.load(f, handler=handler).to_dict()
        page["filename"] = filename

    return page


def process_markdown(md: str) -> str:
    """Take markdown content and process to HTML."""
    exts = ["codehilite", "fenced_code", "tables", "toc"]
    html = markdown(md, extensions=exts)
    return html


def gather_pages(content_dir: Path, config: Dict) -> List:
    """Build the list of pages from the file system."""
    pages = []
    baseurl = config["site"]["baseurl"]
    files = content_dir.glob("**/*.md")
    for f in files:
        if f.name == "index.md":
            continue

        page = get_page(f, config)
        page |= determine_categories_from_path(f.parent, content_dir)
        page["name"] = os.path.relpath(Path(f.parent, f.stem), start=content_dir)
        page["url"] = utils.urljoin([baseurl, page["name"]])
        pages.append(page)
    return pages


def get_parent_name(p: Path, c: Path) -> str:
    """If parent name is . return empty."""
    name = os.path.relpath(p, start=c)
    if name == ".":
        name = ""
    return name


def determine_categories_from_path(file_parent: Path, content_dir: Path) -> Dict:
    """Get the page category and parent from paths.
    Returns empty for both if top level page, or empty parent if not sub-category."""
    # set defaults to be empty
    d = {"category": "", "parent": ""}

    if file_parent != content_dir:
        # Then we're in a subdirectory, thus a category
        category_path = file_parent
        parent_path = category_path.parent
        d["parent"] = get_parent_name(parent_path, content_dir)
        d["category"] = os.path.relpath(category_path, parent_path)
    return d


def gather_categories(content_dir: Path, config: Dict) -> List:
    """Build list of categories from the filesystem."""
    categories = []
    baseurl = config["site"]["baseurl"]

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

        page = get_page(index, config)
        page["category"] = name

        category = {
            "name": name,
            "parent": parent_name,
            "page": page,
            "url": utils.urljoin([baseurl, parent_name, name]),
        }
        categories.append(category)

    return categories
