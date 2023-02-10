import frontmatter
from markdown import markdown
import os
from pathlib import Path
from typing import Dict, List


def get_page(filename: os.PathLike) -> Dict:
    page = {}
    exts = ["fenced_code", "tables"]
    with open(filename, "r") as f:
        page = frontmatter.load(f)

    page["filename"] = filename
    page.content = markdown(page.content, extensions=exts)
    return page


def gather_pages(a: Dict) -> List:
    pages = []
    files = a["content_dir"].glob("**/*.md")
    for filename in files:
        if filename.name == "index.md":
            continue
        page = get_page(filename)
        page["category"] = os.path.relpath(filename.parent, a["content_dir"])
        page["url"] = a["base_url"] + os.path.relpath(
            filename.with_suffix(".html"), start=a["content_dir"]
        )
        pages.append(page)
    return pages


## TODO - pages need to know own category


def gather_categories(a: Dict) -> List:
    categories = []

    paths = a["content_dir"].glob("**")
    for p in paths:
        if not p.is_dir():
            continue  # skip

        name = os.path.relpath(p, start=a["content_dir"])
        if name == ".":
            continue  # skip

        index = Path(p, "index.md")
        if not index.is_file():
            continue  # skip - it's not a category without a index.md

        page = get_page(index)
        category = {
            "name": name,
            "page": page,
            "url": a["base_url"] + name + "/",
        }
        categories.append(category)

    return categories
