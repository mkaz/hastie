#!/usr/bin/env python3

from config import init_args
from jinja2 import Environment, FileSystemLoader, select_autoescape
import shutil
import time

# internal imports
from page import get_page
from templates import *

VERSION = "1.0.0"


def main():
    start_time = time.time()

    args = init_args(VERSION)
    cdir = args["content_dir"]
    odir = args["output_dir"]
    tdir = args["templates_dir"]

    print(f"Hastie v{VERSION}")

    page_files = args["content_dir"].glob("**/*.md")

    # load in jinja templates
    jinja = Environment(loader=FileSystemLoader(tdir), autoescape=select_autoescape())

    # copy templates static dir to output
    tpl_static = Path(tdir, "static")
    out_static = Path(odir)
    shutil.copytree(tpl_static, out_static, dirs_exist_ok=True)

    # gather pages info
    pages = []
    for filename in page_files:
        page = get_page(filename)
        page["filename"] = filename
        page["url"] = os.path.relpath(page["filename"], start=cdir)
        pages.append(page)

    # generate pages
    for page in pages:
        tpl_name = what_template(page["filename"], cdir)
        if "template" in page:
            tpl = page["template"]

        tpl = jinja.get_template(tpl_name)
        html = tpl.render(page=page, pages=pages)
        outfile = get_output_file(page["filename"], cdir, odir)

        # create directories if they don't exist
        outfile.parent.mkdir(exist_ok=True, parents=True)
        outfile.write_text(html)

    elapsed = time.time() - start_time
    print(f"Generated {len(pages)} files in {elapsed:.3f} sec")


if __name__ == "__main__":
    main()
