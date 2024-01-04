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
        # Check if it is an index file
        if f.name == "index.md":
            # We don't want category pages, they are special
            if f.parent.parent == content_dir:
                continue

        page = get_page(f, config)
        page["category"] = determine_category_from_path(f.parent, content_dir)

        # determine name different for directory page
        if f.name == "index.md":
            page["name"] = os.path.relpath(f.parent, start=content_dir)
        else:
            page["name"] = os.path.relpath(Path(f.parent, f.stem), start=content_dir)

        page["url"] = utils.urljoin([baseurl, page["name"]])

        if "subpages" in page and page["subpages"] == "skip":
            page["subpages"] = []
        else:
            page["subpages"] = gather_subpages(f, config)

        pages.append(page)

    return utils.human_sort(pages, "title")


def determine_category_from_path(file_parent: Path, content_dir: Path) -> str:
    """Get the page category from the files parent path.

    Option 1 - Top level page
        - the file parent is the same as content dir

    Option 2 - a single markdown file as page
        - the file parent is the category (grandparent is content dir)

    Option 3 - a directory page with index.md
        - the grandparent is category
    """

    # Option 1- Top level page
    if file_parent == content_dir:
        return ""
    # Option 2 - Parent is category
    category_path = file_parent
    if category_path.parent == content_dir:
        return os.path.relpath(category_path, start=content_dir)

    # Option 3 - Grandparent is category
    return os.path.relpath(category_path.parent, start=content_dir)


def gather_categories(content_dir: Path, config: Dict) -> List:
    """Build list of categories from the filesystem."""
    categories = []
    baseurl = config["site"]["baseurl"]

    ## do not recurse, categories are top level directories
    paths = content_dir.glob("*")
    for p in paths:
        if not p.is_dir():
            continue  # skip files

        name = os.path.relpath(p, start=content_dir)

        index = Path(p, "index.md")
        if not index.is_file():
            continue  # skip - it's not a category without a index.md

        page = get_page(index, config)
        page["category"] = name

        category = {
            "name": name,
            "page": page,
            "url": utils.urljoin([baseurl, name]),
        }
        categories.append(category)

    utils.human_sort(categories, "name")

    return categories


def gather_subpages(filepath: Path, config: Dict) -> List:
    """Build the list of subpages from page system."""
    subpages = []
    content_dir = config["content_dir"]
    baseurl = config["site"]["baseurl"]

    dirs = []
    for en in filepath.parent.iterdir():
        if en.is_dir():
            dirs.append(en)

    for d in dirs:
        files = d.glob("*.md")
        for f in files:
            page = get_page(f, config)

            # determine name different for directory page
            if f.name == "index.md":
                page["name"] = os.path.relpath(f.parent, start=content_dir)
            else:
                page["name"] = os.path.relpath(
                    Path(f.parent, f.stem), start=content_dir
                )

            page["url"] = utils.urljoin([baseurl, page["name"]])

            subpages.append(page)

        utils.human_sort(subpages, "title")
    return subpages


def filter_category_pages(category: str, pages: List) -> List:
    category_pages = []

    for page in pages:
        if page["category"] != category:
            continue

        if "draft" in page:
            continue

        if "archive" in page:
            continue

        category_pages.append(page)

    return category_pages
