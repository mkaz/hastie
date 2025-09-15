"""
Hastie resources, pages and categories.
"""

import sys
from pathlib import Path
from typing import Any

import frontmatter
from markdown import markdown

import hastie.utils as utils


def get_page(filename: Path, config: dict[str, Any]) -> dict[str, Any]:
    """Read page in from file system, parse frontmatter and render markdown."""
    try:
        page = read_page(filename, config)
        page["content"] = process_markdown(page.get("content", ""))
    except Exception as err:
        print(f"Error reading page {filename}")
        print(err)
        sys.exit(1)
    return page


def read_page(filename: Path, config: dict[str, Any] | None = None) -> dict[str, Any]:
    if config is None:
        config = {}
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


def gather_pages(content_dir: Path, config: dict[str, Any]) -> list[dict[str, Any]]:
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
            page["name"] = f.parent.relative_to(content_dir).as_posix()
        else:
            page["name"] = (f.parent / f.stem).relative_to(content_dir).as_posix()

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
        return category_path.relative_to(content_dir).as_posix()

    # Option 3 - Grandparent is category
    return category_path.parent.relative_to(content_dir).as_posix()


def gather_categories(
    content_dir: Path, config: dict[str, Any]
) -> list[dict[str, Any]]:
    """Build list of categories from the filesystem."""
    categories = []
    baseurl = config["site"]["baseurl"]

    ## do not recurse, categories are top level directories
    paths = content_dir.glob("*")
    for p in paths:
        if not p.is_dir():
            continue  # skip files

        name = p.relative_to(content_dir).as_posix()

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


def gather_subpages(filepath: Path, config: dict[str, Any]) -> list[dict[str, Any]]:
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
                page["name"] = f.parent.relative_to(content_dir).as_posix()
            else:
                page["name"] = (f.parent / f.stem).relative_to(content_dir).as_posix()

            page["url"] = utils.urljoin([baseurl, page["name"]])

            subpages.append(page)

        utils.human_sort(subpages, "title")
    return subpages


def filter_category_pages(
    category: str, pages: list[dict[str, Any]], include_archived: bool = False
) -> list[dict[str, Any]]:
    category_pages = []

    for page in pages:
        if page["category"] != category:
            continue

        if "draft" in page:
            continue

        if "archive" in page and not include_archived:
            continue

        category_pages.append(page)

    return category_pages
