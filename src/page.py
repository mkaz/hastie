import frontmatter
from markdown import markdown
import os
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
        page = get_page(filename)
        page["filename"] = filename  # on disk file
        page["url"] = a["base_url"] + os.path.relpath(
            filename.with_suffix(".html"), start=a["content_dir"]
        )
        pages.append(page)
    return pages


def gather_categories(a: Dict) -> List:
    cats = []
    paths = a["content_dir"].glob("**")
    for p in paths:
        if p.is_dir():
            cat = os.path.relpath(p, start=a["content_dir"])
            if cat == ".":
                continue  # skip

            cat = {"name": cat, "url": a["base_url"] + cat + "/"}
            cats.append(cat)

    return cats
