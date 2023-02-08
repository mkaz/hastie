#!/usr/bin/env python3

from config import init_args
import frontmatter
from jinja2 import Environment, FileSystemLoader, select_autoescape
from markdown import markdown
from typing import Dict, List


def main():
    args = init_args()
    print("Welcome to Hastie")

    page_files = args["content_dir"].glob("**/*.md")

    jinja = Environment(
        loader=FileSystemLoader(args["templates_dir"]), autoescape=select_autoescape()
    )

    # instantiate jinja templating
    page_template = jinja.get_template("page.html")

    for filename in page_files:
        if filename.name != "index.md":
            # read in page data from f
            page = get_page(filename)
            html = page_template.render(page=page)
            print(html)
            print("--------------------------")


def get_page(filename) -> Dict:
    page = {}
    with open(filename, "r") as f:
        page = frontmatter.load(f)

    page["filename"] = filename
    page["content"] = markdown(page["content"])
    return page


if __name__ == "__main__":
    main()
