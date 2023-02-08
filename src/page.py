import frontmatter
from markdown import markdown
import os
from typing import Dict


def get_page(filename: os.PathLike) -> Dict:
    page = {}
    exts = ["fenced_code"]
    with open(filename, "r") as f:
        page = frontmatter.load(f)

    page["filename"] = filename
    page.content = markdown(page.content, extensions=exts)
    return page
