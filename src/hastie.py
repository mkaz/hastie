#!/usr/bin/env python3

from config import init_args
from jinja2 import Environment, FileSystemLoader, select_autoescape
from pathlib import Path
import shutil
import time

# internal imports
from page import gather_pages, gather_categories, get_page
from templates import get_output_file

VERSION = "1.0.0"


def main():
    start_time = time.time()
    count = 0

    args = init_args(VERSION)
    cdir = args["content_dir"]
    odir = args["output_dir"]
    tdir = args["templates_dir"]

    if not args["quiet"]:
        print(f"Hastie v{VERSION}")

    # load in jinja templates
    jinja = Environment(loader=FileSystemLoader(tdir), autoescape=select_autoescape())

    # copy templates static dir to output
    tpl_static = Path(tdir, "static")
    out_static = Path(odir)
    shutil.copytree(tpl_static, out_static, dirs_exist_ok=True)

    # gather site info
    pages = gather_pages(args)
    categories = gather_categories(args)
    site = []
    if "site" in args:
        site = args["site"]

    # generate pages
    for page in pages:
        tpl_name = "page.html"
        if "template" in page:
            tpl_name = page["template"]

        tpl = jinja.get_template(tpl_name)
        html = tpl.render(page=page, pages=pages, categories=categories, site=site)
        outfile = get_output_file(page["filename"], cdir, odir)

        # create directories if they don't exist
        outfile.parent.mkdir(exist_ok=True, parents=True)
        outfile.write_text(html)
        count += 1

    # generate category pages
    for cat in categories:
        tpl_name = "category.html"
        if "template" in cat["page"]:
            tpl_name = cat["page"]["template"]

        ## filter pages to those within category
        category_pages = filter(lambda p: p["category"] == cat["name"], pages)

        tpl = jinja.get_template(tpl_name)
        html = tpl.render(
            page=cat["page"], pages=category_pages, categories=categories, site=site
        )
        outfile = get_output_file(cat["page"]["filename"], cdir, odir)

        # create directories if they don't exist
        outfile.parent.mkdir(exist_ok=True, parents=True)
        outfile.write_text(html)
        count += 1

    # generate home page
    home = get_page(Path(cdir, "index.md"))
    home["url"] = args["base_url"]
    tpl_name = "index.html"
    if "template" in home:
        tpl_name = home["template"]

    tpl = jinja.get_template(tpl_name)
    html = tpl.render(page=home, pages=pages, categories=categories, site=site)
    outfile = get_output_file(home["filename"], cdir, odir)
    outfile.write_text(html)
    count += 1

    elapsed = time.time() - start_time
    if not args["quiet"]:
        print(f"Generated {count} files in {elapsed:.3f} sec")


if __name__ == "__main__":
    main()
