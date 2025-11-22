"""
Hastie resources, pages and categories.
"""

import asyncio
import sys
from pathlib import Path
from typing import Any

import aiofiles
import frontmatter
from markdown import markdown
from markdown.extensions.toc import TocExtension

import hastie.utils as utils


# Module-level cache for processed pages
_page_cache: dict[Path, dict[str, Any]] = {}


def clear_cache():
    """Clear the page cache."""
    _page_cache.clear()


def get_page(filename: Path, config: dict[str, Any]) -> dict[str, Any]:
    """Read page in from file system, parse frontmatter and render markdown."""
    # Check cache first
    if filename in _page_cache:
        return _page_cache[filename]

    try:
        page = read_page(filename, config)
        page["content"] = process_markdown(page.get("content", ""))
    except Exception as err:
        print(f"Error reading page {filename}")
        print(err)
        sys.exit(1)

    # Cache the result
    _page_cache[filename] = page
    return page


async def get_page_async(filename: Path, config: dict[str, Any]) -> dict[str, Any]:
    """Async version: Read page, parse frontmatter and render markdown."""
    # Check cache first
    if filename in _page_cache:
        return _page_cache[filename]

    try:
        page = await read_page_async(filename, config)
        # Run markdown processing in thread pool (it's CPU-bound)
        loop = asyncio.get_event_loop()
        page["content"] = await loop.run_in_executor(
            None, process_markdown, page.get("content", "")
        )
    except Exception as err:
        print(f"Error reading page {filename}")
        print(err)
        sys.exit(1)

    # Cache the result
    _page_cache[filename] = page
    return page


def read_page(filename: Path, config: dict[str, Any] | None = None) -> dict[str, Any]:
    """Read page using frontmatter library."""
    if config is None:
        config = {}

    with open(filename, "r") as f:
        content = f.read()

    return parse_frontmatter(content, filename, config)


async def read_page_async(filename: Path, config: dict[str, Any] | None = None) -> dict[str, Any]:
    """Async version: Read page using frontmatter library."""
    if config is None:
        config = {}

    async with aiofiles.open(filename, "r") as f:
        content = await f.read()

    return parse_frontmatter(content, filename, config)


def parse_frontmatter(content: str, filename: Path, config: dict[str, Any]) -> dict[str, Any]:
    """Parse frontmatter from content string."""
    fm = config.get("frontmatter", "yaml")
    match fm:
        case "toml":
            handler = frontmatter.default_handlers.TOMLHandler()
        case "json":
            handler = frontmatter.default_handlers.JSONHandler()
        case _:
            handler = frontmatter.default_handlers.YAMLHandler()

    page = frontmatter.loads(content, handler=handler).to_dict()
    page["filename"] = filename
    return page


def process_markdown(md: str) -> str:
    """Take markdown content and process to HTML."""
    exts = [
        "codehilite",
        "fenced_code",
        "tables",
        TocExtension(baselevel=2, toc_depth="2-3"),
    ]
    html = markdown(md, extensions=exts)
    return html


async def gather_pages_async(content_dir: Path, config: dict[str, Any]) -> list[dict[str, Any]]:
    """Async version: Build the list of pages from the file system."""
    baseurl = config["site"]["baseurl"]
    files = list(content_dir.glob("**/*.md"))

    # Filter out category index files first
    page_files = []
    for f in files:
        if f.name == "index.md":
            # We don't want category pages, they are special
            if f.parent.parent == content_dir:
                continue
        page_files.append(f)

    # Process all pages concurrently
    tasks = [get_page_async(f, config) for f in page_files]
    pages_raw = await asyncio.gather(*tasks)

    pages = []
    for f, page in zip(page_files, pages_raw):
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
            # Gather subpages (these will use cache)
            page["subpages"] = await gather_subpages_async(f, config)

        pages.append(page)

    return utils.human_sort(pages, "title")


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


async def gather_categories_async(
    content_dir: Path, config: dict[str, Any]
) -> list[dict[str, Any]]:
    """Async version: Build list of categories from the filesystem."""
    categories = []
    baseurl = config["site"]["baseurl"]

    ## do not recurse, categories are top level directories
    paths = list(content_dir.glob("*"))

    # Filter to valid category directories
    cat_data = []
    for p in paths:
        if not p.is_dir():
            continue
        index = Path(p, "index.md")
        if not index.is_file():
            continue
        cat_data.append((p, index))

    # Process all category pages concurrently
    tasks = [get_page_async(index, config) for _, index in cat_data]
    pages_raw = await asyncio.gather(*tasks)

    for (p, _), page in zip(cat_data, pages_raw):
        name = p.relative_to(content_dir).as_posix()
        page["category"] = name

        category = {
            "name": name,
            "page": page,
            "url": utils.urljoin([baseurl, name]),
            "parent": "",  # top-level categories have empty parent
        }
        categories.append(category)

    utils.human_sort(categories, "name")
    return categories


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
            "parent": "",  # top-level categories have empty parent
        }
        categories.append(category)

    utils.human_sort(categories, "name")

    return categories


async def gather_subpages_async(filepath: Path, config: dict[str, Any]) -> list[dict[str, Any]]:
    """Async version: Build the list of subpages from page system."""
    subpages = []
    content_dir = config["content_dir"]
    baseurl = config["site"]["baseurl"]

    dirs = [en for en in filepath.parent.iterdir() if en.is_dir()]

    # Collect all subpage files
    subpage_files = []
    for d in dirs:
        files = list(d.glob("*.md"))
        subpage_files.extend(files)

    if not subpage_files:
        return []

    # Process all subpages concurrently
    tasks = [get_page_async(f, config) for f in subpage_files]
    pages_raw = await asyncio.gather(*tasks)

    for f, page in zip(subpage_files, pages_raw):
        # determine name different for directory page
        if f.name == "index.md":
            page["name"] = f.parent.relative_to(content_dir).as_posix()
        else:
            page["name"] = (f.parent / f.stem).relative_to(content_dir).as_posix()

        page["url"] = utils.urljoin([baseurl, page["name"]])
        subpages.append(page)

    utils.human_sort(subpages, "title")
    return subpages


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


def build_category_pages_map(
    pages: list[dict[str, Any]]
) -> dict[str, list[dict[str, Any]]]:
    """Pre-compute mapping of category -> pages for fast lookup."""
    category_map: dict[str, list[dict[str, Any]]] = {}

    for page in pages:
        cat = page["category"]
        if cat not in category_map:
            category_map[cat] = []
        category_map[cat].append(page)

    return category_map


def get_category_pages_fast(
    category: str,
    category_map: dict[str, list[dict[str, Any]]],
    include_archived: bool = False
) -> list[dict[str, Any]]:
    """Fast category page lookup using pre-computed map."""
    pages = category_map.get(category, [])

    result = []
    for page in pages:
        if "draft" in page:
            continue
        if "archive" in page and not include_archived:
            continue
        result.append(page)

    return result
