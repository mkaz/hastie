#!/usr/bin/env python3

from jinja2 import Environment, FileSystemLoader, select_autoescape
from pathlib import Path

import sys
import time

# internal imports
from hastie.config import config, __version__
import hastie.content as content
import hastie.hfs as hfs
from hastie.rss import generate_rss
from hastie.utils import human_sort, date_sort


def generate_site():
    """Generate the static site."""
    start_time = time.time()
    count = 0

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

    # copy all the static assets
    hfs.copy_static_assets(cdir, odir, static_dir)

    # load in jinja templates
    jinja = Environment(loader=FileSystemLoader(tdir), autoescape=select_autoescape())

    # gather site info - all pages, categories
    pages = content.gather_pages(cdir, config)
    categories = content.gather_categories(cdir, config)
    site = []
    if "site" in config:
        site = config["site"]

    # sort by date, most recent
    recent_pages = date_sort(pages)

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
        # If this page is archived or has include_archived: true, include archived pages in the listing
        include_archived = "archive" in page or page.get("include_archived", False)
        category_pages = content.filter_category_pages(
            page["category"], pages, include_archived
        )

        # sort pages
        # human_sort(category_pages, "title")

        try:
            html = tpl.render(
                page=page,
                pages=category_pages,
                categories=categories,
                recent_pages=recent_pages,
                site=site,
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

        # remove archived from category pages (unless this category page itself is archived or has include_archived: true)
        if "archive" not in cat["page"] and not cat["page"].get(
            "include_archived", False
        ):
            category_pages = list(filter(lambda p: "archive" not in p, category_pages))

        # sort by title
        human_sort(category_pages, "name")

        cat["page"]["categories"] = filter(
            lambda c: (c["parent"] == cat["parent"] or c["parent"] == cat["name"])
            and (c["name"] != cat["name"]),
            categories,
        )

        tpl = jinja.get_template(tpl_name)
        try:
            html = tpl.render(
                page=cat["page"],
                pages=category_pages,
                categories=categories,
                recent_pages=recent_pages,
                site=site,
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
    html = tpl.render(
        page=home,
        pages=pages,
        categories=categories,
        recent_pages=recent_pages,
        site=site,
    )
    outfile = hfs.get_output_file(home["filename"], cdir, odir)
    outfile.write_text(html)
    count += 1

    # generate RSS
    if "rss" in config["site"]:
        rss = generate_rss(config, pages)
        outfile = Path(odir, "rss.xml")
        outfile.write_text(rss)
        count += 1

    elapsed = time.time() - start_time
    if not config["quiet"]:
        print(f"Generated {count} files in {elapsed:.3f} sec")


def main():
    if not config["quiet"]:
        print(f"Hastie v{__version__}")

    # Run initial generation
    generate_site()

    # If watch mode is enabled, start monitoring for changes
    if config.get("watch", False):
        from watchdog.observers import Observer
        from watchdog.events import FileSystemEventHandler

        class RegenerateHandler(FileSystemEventHandler):
            def __init__(self):
                self.last_regeneration = time.time()
                self.debounce_seconds = 1.0

            def should_regenerate(self, event):
                # Skip directory events
                if event.is_directory:
                    return False

                # Skip temporary and backup files
                path = event.src_path
                if path.endswith(('~', '.swp', '.swx', '.tmp', '.bak')):
                    return False

                # Skip hidden files
                if Path(path).name.startswith('.'):
                    return False

                # Only regenerate on actual file modifications or creations
                # Ignore delete and move events
                if event.event_type not in ('modified', 'created'):
                    return False

                return True

            def on_modified(self, event):
                if not self.should_regenerate(event):
                    return

                # Debounce multiple rapid changes
                now = time.time()
                if now - self.last_regeneration < self.debounce_seconds:
                    return

                self.last_regeneration = now

                if not config["quiet"]:
                    print(f"\nChange detected: {event.src_path}")
                    print("Regenerating site...")

                try:
                    generate_site()
                except Exception as err:
                    print(f"Error during regeneration: {err}")

            def on_created(self, event):
                # Use the same logic as on_modified
                self.on_modified(event)

        handler = RegenerateHandler()
        observer = Observer()

        # Watch both content and templates directories
        cdir = config["content_dir"]
        tdir = config["templates_dir"]

        observer.schedule(handler, str(cdir), recursive=True)
        observer.schedule(handler, str(tdir), recursive=True)
        observer.start()

        if not config["quiet"]:
            print(f"\nWatching for changes in {cdir} and {tdir}...")
            print("Press Ctrl+C to stop")

        try:
            while True:
                time.sleep(1)
        except KeyboardInterrupt:
            observer.stop()
            if not config["quiet"]:
                print("\nStopping watch mode...")

        observer.join()


if __name__ == "__main__":
    main()
