#!/usr/bin/env python3

import asyncio
from concurrent.futures import ThreadPoolExecutor
from jinja2 import Environment, FileSystemLoader, select_autoescape
from pathlib import Path

import aiofiles
import sys
import time

# internal imports
from hastie.config import config, __version__
import hastie.content as content
import hastie.hfs as hfs
from hastie.rss import generate_rss
from hastie.utils import human_sort, date_sort


async def generate_site_async():
    """Generate the static site using async I/O."""
    start_time = time.time()

    cdir = config["content_dir"]
    odir = config["output_dir"]
    tdir = config["templates_dir"]
    static_dir = config["static_dir"]

    # Confirm content and template directories exists
    if not cdir.is_dir():
        print(f"Content directory {cdir} not found")
        sys.exit()

    if not tdir.is_dir():
        print(f"Templates directory {tdir} not found")
        sys.exit()

    # Clear the page cache and reset semaphore for fresh generation
    content.clear_cache()
    content.reset_semaphore()

    # copy all the static assets (synchronous, uses shutil.copytree)
    hfs.copy_static_assets(cdir, odir, static_dir)

    # load in jinja templates
    jinja = Environment(loader=FileSystemLoader(tdir), autoescape=select_autoescape())

    # Template cache for fast lookups
    template_cache: dict[str, any] = {}

    def get_template_cached(name: str):
        if name not in template_cache:
            template_cache[name] = jinja.get_template(name)
        return template_cache[name]

    # gather site info - all pages, categories (async)
    pages, categories = await asyncio.gather(
        content.gather_pages_async(cdir, config),
        content.gather_categories_async(cdir, config)
    )

    site = config.get("site", [])

    # sort by date, most recent
    recent_pages = date_sort(pages)

    # Pre-compute category -> pages mapping for O(1) lookups
    category_map = content.build_category_pages_map(pages)

    # Prepare all page render tasks
    write_tasks = []

    # Use thread pool for CPU-bound template rendering
    loop = asyncio.get_event_loop()
    executor = ThreadPoolExecutor(max_workers=8)

    # Generate pages
    for page in pages:
        # do not write out drafts
        if "draft" in page:
            continue

        tpl_name = page.get("template", "page.html")

        # filter categories to the page
        page["categories"] = [
            c for c in categories if page["category"] == c["parent"]
        ]

        # check for template
        try:
            tpl = get_template_cached(tpl_name)
        except Exception as err:
            print("Error getting template")
            print(f"    Template: {tpl_name}")
            print(err)
            sys.exit()

        # Use fast category lookup
        include_archived = "archive" in page or page.get("include_archived", False)
        category_pages = content.get_category_pages_fast(
            page["category"], category_map, include_archived
        )

        outfile = hfs.get_output_file(page["filename"], cdir, odir)

        # Queue render and write task
        write_tasks.append(
            render_and_write(
                loop, executor, tpl, outfile,
                page=page,
                pages=category_pages,
                categories=categories,
                recent_pages=recent_pages,
                site=site,
            )
        )

    # Generate category pages
    for cat in categories:
        tpl_name = cat["page"].get("template", "category.html")

        # Use fast category lookup
        include_archived = "archive" in cat["page"] or cat["page"].get("include_archived", False)
        category_pages = content.get_category_pages_fast(
            cat["name"], category_map, include_archived
        )

        # sort by title
        human_sort(category_pages, "name")

        cat["page"]["categories"] = [
            c for c in categories
            if (c["parent"] == cat["parent"] or c["parent"] == cat["name"])
            and (c["name"] != cat["name"])
        ]

        try:
            tpl = get_template_cached(tpl_name)
        except Exception as err:
            print("Error getting template")
            print(f"    Template: {tpl_name}")
            print(err)
            sys.exit()

        outfile = hfs.get_output_file(cat["page"]["filename"], cdir, odir)

        write_tasks.append(
            render_and_write(
                loop, executor, tpl, outfile,
                page=cat["page"],
                pages=category_pages,
                categories=categories,
                recent_pages=recent_pages,
                site=site,
            )
        )

    # Generate home page
    home = content.get_page(Path(cdir, "index.md"), config)
    home["url"] = config["site"]["baseurl"]
    tpl_name = home.get("template", "index.html")

    home["categories"] = [c for c in categories if c["parent"] == ""]
    tpl = get_template_cached(tpl_name)
    outfile = hfs.get_output_file(home["filename"], cdir, odir)

    write_tasks.append(
        render_and_write(
            loop, executor, tpl, outfile,
            page=home,
            pages=pages,
            categories=categories,
            recent_pages=recent_pages,
            site=site,
        )
    )

    # Execute all render and write tasks in parallel
    results = await asyncio.gather(*write_tasks, return_exceptions=True)

    # Check for errors
    count = 0
    for result in results:
        if isinstance(result, Exception):
            print(f"Error: {result}")
        else:
            count += 1

    # generate RSS (small file, do synchronously)
    if "rss" in config["site"]:
        rss = generate_rss(config, pages)
        outfile = Path(odir, "rss.xml")
        outfile.write_text(rss)
        count += 1

    executor.shutdown(wait=False)

    elapsed = time.time() - start_time
    if not config["quiet"]:
        print(f"Generated {count} files in {elapsed:.3f} sec")


async def render_and_write(loop, executor, tpl, outfile, **context):
    """Render template and write to file asynchronously."""
    # Render in thread pool (CPU-bound)
    try:
        html = await loop.run_in_executor(
            executor,
            lambda: tpl.render(**context)
        )
    except Exception as err:
        page = context.get("page", {})
        print("Error rendering page with template")
        print(f"    Page    : {page.get('filename', 'unknown')}")
        print(f"    Template: {tpl.name}")
        print(err)
        raise

    # create directories if they don't exist
    outfile.parent.mkdir(exist_ok=True, parents=True)

    # Write asynchronously (use semaphore to limit concurrent open files)
    sem = content._get_file_semaphore()
    async with sem:
        async with aiofiles.open(outfile, "w") as f:
            await f.write(html)


def generate_site():
    """Generate the static site (wrapper for async version)."""
    asyncio.run(generate_site_async())


def main():
    if not config["quiet"]:
        print(f"Hastie v{__version__}")

    generate_site()


if __name__ == "__main__":
    main()
