"""Hastie file system module"""

from pathlib import Path
import hastie.hfs as hfs


def test_get_output_dir_index():
    content_dir = Path("/pages")
    output_dir = Path("/output")
    filepath = Path("/pages/index.md")
    p = hfs.get_output_file(filepath, content_dir, output_dir)
    assert p == Path("/output/index.html")


def test_get_output_dir_page():
    content_dir = Path("/pages")
    output_dir = Path("/output")
    filepath = Path("/pages/getting-started.md")
    p = hfs.get_output_file(filepath, content_dir, output_dir)
    assert p == Path("/output/getting-started.html")


def test_get_output_dir_category_index():
    content_dir = Path("/pages")
    output_dir = Path("/output")
    filepath = Path("/pages/pastas/index.md")
    p = hfs.get_output_file(filepath, content_dir, output_dir)
    assert p == Path("/output/pastas/index.html")


def test_get_output_dir_category_page():
    content_dir = Path("/pages")
    output_dir = Path("/output")
    filepath = Path("/pages/pastas/macaroni.md")
    p = hfs.get_output_file(filepath, content_dir, output_dir)
    assert p == Path("/output/pastas/macaroni.html")
