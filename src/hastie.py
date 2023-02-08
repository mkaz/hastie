#!/usr/bin/env python3

from config import init_args
from jinja2 import Environment, FileSystemLoader, select_autoescape
from page import get_page
from templates import what_template


def main():
    args = init_args()
    print("Welcome to Hastie")

    page_files = args["content_dir"].glob("**/*.md")

    # load in jinja templates
    jinja = Environment(
        loader=FileSystemLoader(args["templates_dir"]), autoescape=select_autoescape()
    )

    for filename in page_files:
        tpl_name = what_template(filename, args["content_dir"])
        page = get_page(filename)
        if "template" in page:
            tpl = page["template"]

        tpl = jinja.get_template(tpl_name)
        html = tpl.render(page=page)
        print(html)
        print("--------------------------")


if __name__ == "__main__":
    main()