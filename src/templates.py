import os
from pathlib import Path


def what_template(filename: os.PathLike, content_dir: os.PathLike) -> str:

    if filename.name == "index.md":
        if filename.parent == content_dir:
            return "index.html"
        return "category.html"

    return "page.html"


def get_output_file(f: os.PathLike, c: os.PathLike, o: os.PathLike) -> os.PathLike:
    """Takes filename, content directory, output directory and returns output file"""
    # get file relative to content directory
    jf = Path(os.path.relpath(f, start=c))

    # switch markdown to html
    jf = jf.with_suffix(".html")

    # return relative file joined with output directory
    return Path(o, jf)
