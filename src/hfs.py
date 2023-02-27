import os
from pathlib import Path


def get_output_file(f: os.PathLike, c: os.PathLike, o: os.PathLike) -> os.PathLike:
    """Takes filename, content directory, output directory and returns output file"""
    # get file relative to content directory
    jf = Path(os.path.relpath(f, start=c))

    if jf.name == "index.md":
        return Path(o, jf.parent, "index.html")

    # filename without extension is .stem
    # create directory from stem - write as index.html
    return Path(o, jf.parent, jf.stem, "index.html")
