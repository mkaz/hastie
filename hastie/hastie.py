#!/usr/bin/env python3

from config import init_args
from jinja2 import Environment, FileSystemLoader, select_autoescape
from pathlib import Path
import shutil
import time

# internal imports
import hres
import hfs

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

    # copy templates static dir to output - content ends up top level
    tpl_static = Path(tdir, "static")
    out_tpl_static = Path(odir)
    shutil.copytree(tpl_static, out_tpl_static, dirs_exist_ok=True)

    # copy site static dir to output - content ends in static/ dir
    site_static = Path("./", "static")
    out_static = Path(odir, "static")
    if site_static.is_dir():
        shutil.copytree(site_static, out_static, dirs_exist_ok=True)

    # gather site info - all pages, categories
    pages = hres.gather_pages(cdir, base_url=args["base_url"])
    categories = hres.gather_categories(cdir, base_url=args["base_url"])
    site = []
    if "site" in args:
        site = args["site"]

    # generate pages
    for page in pages:
        tpl_name = "page.html"
        if "template" in page:
            tpl_name = page["template"]

        # filter categories to the page
        page["categories"] = filter(
            lambda c: page["category"] == c["parent"], categories
        )

        tpl = jinja.get_template(tpl_name)
        html = tpl.render(page=page, pages=pages, categories=categories, site=site)
        outfile = hfs.get_output_file(page["filename"], cdir, odir)

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

        cat["page"]["categories"] = filter(
            lambda c: (c["parent"] == cat["parent"] or c["parent"] == cat["name"])
            and (c["name"] != cat["name"]),
            categories,
        )

        tpl = jinja.get_template(tpl_name)
        html = tpl.render(
            page=cat["page"], pages=category_pages, categories=categories, site=site
        )
        outfile = hfs.get_output_file(cat["page"]["filename"], cdir, odir)

        # create directories if they don't exist
        outfile.parent.mkdir(exist_ok=True, parents=True)
        outfile.write_text(html)
        count += 1

    # generate home page
    home = hres.get_page(Path(cdir, "index.md"))
    home["url"] = args["base_url"]
    tpl_name = "index.html"
    if "template" in home:
        tpl_name = home["template"]

    home["categories"] = filter(lambda c: c["parent"] == "", categories)
    tpl = jinja.get_template(tpl_name)
    html = tpl.render(page=home, pages=pages, categories=categories, site=site)
    outfile = hfs.get_output_file(home["filename"], cdir, odir)
    outfile.write_text(html)
    count += 1

    elapsed = time.time() - start_time
    if not args["quiet"]:
        print(f"Generated {count} files in {elapsed:.3f} sec")


if __name__ == "__main__":
    main()
