import os


def what_template(filename: os.PathLike, content_dir: os.PathLike) -> str:

    if filename.name == "index.md":
        if filename.parent == content_dir:
            return "index.html"
        return "category.html"

    return "page.html"
