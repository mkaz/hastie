import os
from pathlib import Path


def get_output_file(f: os.PathLike, c: os.PathLike, o: os.PathLike) -> os.PathLike:
    """Takes filename, content directory, output directory and returns output file"""
    # get file relative to content directory
    jf = Path(os.path.relpath(f, start=c))

    # switch markdown to html
    jf = jf.with_suffix(".html")

    # return relative file joined with output directory
    return Path(o, jf)
