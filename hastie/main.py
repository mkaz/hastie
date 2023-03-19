#!/usr/bin/env python3

from jinja2 import Environment, FileSystemLoader, select_autoescape
import operator
from pathlib import Path
import shutil
import sys
import time

# internal imports
from hastie.config import config, __version__
import hastie.content as content
import hastie.hfs as hfs
from hastie.rss import generate_rss


def main():
    start_time = time.time()
    count = 0

    cdir = config["content_dir"]
    odir = config["output_dir"]
    tdir = config["templates_dir"]

    if not config["quiet"]:
        print(f"Hastie v{__version__}")

    # Confirm content and template directories exists
    if not cdir.is_dir():
        print(f"Content directory {cdir} not found")
        sys.exit()

    if not tdir.is_dir():
        print(f"Templates directory {tdir} not found")
        sys.exit()

    # copy templates static dir to output
    # - content at top level: for example /favicon.ico
    tpl_static = Path(tdir, "static")
    out_tpl_static = Path(odir)
    shutil.copytree(tpl_static, out_tpl_static, dirs_exist_ok=True)

    # copy site static dir to output
    # - content within under /static dir
    site_static = Path("./", "static")
    out_static = Path(odir, "static")
    if site_static.is_dir():
        shutil.copytree(site_static, out_static, dirs_exist_ok=True)

    # load in jinja templates
    jinja = Environment(loader=FileSystemLoader(tdir), autoescape=select_autoescape())

    # gather site info - all pages, categories
    pages = content.gather_pages(cdir, config)
    categories = content.gather_categories(cdir, config)
    site = []
    if "site" in config:
        site = config["site"]

    # generate pages
    for page in pages:
        tpl_name = "page.html"
        if "template" in page:
            tpl_name = page["template"]

        # filter categories to the page
        page["categories"] = filter(
            lambda c: page["category"] == c["parent"], categories
        )

        # check for template
        try:
            tpl = jinja.get_template(tpl_name)
        except Exception as err:
            print("Error getting template")
            print(f"    Template: {tpl_name}")
            print(err)
            sys.exit()

        ## filter pages to those within category
        category_pages = list(
            filter(lambda p: p["category"] == page["category"], pages)
        )

        # remove drafts from category pages
        category_pages = list(filter(lambda p: "draft" not in p, category_pages))

        # remove archived from category pages
        category_pages = list(filter(lambda p: "archive" not in p, category_pages))

        # sort pages
        category_pages.sort(key=operator.itemgetter("title"))

        try:
            html = tpl.render(
                page=page, pages=category_pages, categories=categories, site=site
            )
        except Exception as err:
            print("Error rendering page with template")
            print(f"    Page    : {page['filename']}")
            print(f"    Template: {tpl_name}")
            print(err)
            sys.exit()

        outfile = hfs.get_output_file(page["filename"], cdir, odir)

        # do not write out drafts
        if "draft" in page:
            continue

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
        category_pages = list(filter(lambda p: p["category"] == cat["name"], pages))

        # remove drafts from category pages
        category_pages = list(filter(lambda p: "draft" not in p, category_pages))

        # remove archived from category pages
        category_pages = list(filter(lambda p: "archive" not in p, category_pages))

        # sort by title
        category_pages.sort(key=operator.itemgetter("title"))

        cat["page"]["categories"] = filter(
            lambda c: (c["parent"] == cat["parent"] or c["parent"] == cat["name"])
            and (c["name"] != cat["name"]),
            categories,
        )

        tpl = jinja.get_template(tpl_name)
        try:
            html = tpl.render(
                page=cat["page"], pages=category_pages, categories=categories, site=site
            )
        except Exception as err:
            print("Error rendering category with template")
            print(f"    Page    : {cat.page.filename}")
            print(f"    Template: {tpl_name}")
            print(err)
            sys.exit()

        outfile = hfs.get_output_file(cat["page"]["filename"], cdir, odir)

        # create directories if they don't exist
        outfile.parent.mkdir(exist_ok=True, parents=True)
        outfile.write_text(html)
        count += 1

    # generate home page
    home = content.get_page(Path(cdir, "index.md"), config)
    home["url"] = config["site"]["baseurl"]
    tpl_name = "index.html"
    if "template" in home:
        tpl_name = home["template"]

    home["categories"] = filter(lambda c: c["parent"] == "", categories)
    tpl = jinja.get_template(tpl_name)
    html = tpl.render(page=home, pages=pages, categories=categories, site=site)
    outfile = hfs.get_output_file(home["filename"], cdir, odir)
    outfile.write_text(html)
    count += 1

    # generate RSS ?
    if "rss" in config["site"]:
        rss = generate_rss(config, pages)
        outfile = Path(odir, "rss.xml")
        outfile.write_text(rss)
        count += 1

    elapsed = time.time() - start_time
    if not config["quiet"]:
        print(f"Generated {count} files in {elapsed:.3f} sec")


if __name__ == "__main__":
    main()
