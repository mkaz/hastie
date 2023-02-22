"""Hastie file system module"""

from pathlib import Path
import pytest

pytest.skip(allow_module_level=True)

import hastie


def test_get_output_dir_index():
    content_dir = Path("/pages")
    output_dir = Path("/output")
    filepath = Path("/pages/index.md")
    p = hastie.hfs.get_output_file(filepath, content_dir, output_dir)
    assert p == Path("/output/index.html")


def test_get_output_dir_page():
    content_dir = Path("/pages")
    output_dir = Path("/output")
    filepath = Path("/pages/getting-started.md")
    p = hastie.hfs.get_output_file(filepath, content_dir, output_dir)
    assert p == Path("/output/getting-started.html")


def test_get_output_dir_category_index():
    content_dir = Path("/pages")
    output_dir = Path("/output")
    filepath = Path("/pages/pastas/index.md")
    p = hastie.hfs.get_output_file(filepath, content_dir, output_dir)
    assert p == Path("/output/pastas/index.html")


def test_get_output_dir_category_page():
    content_dir = Path("/pages")
    output_dir = Path("/output")
    filepath = Path("/pages/pastas/macaroni.md")
    p = hastie.hfs.get_output_file(filepath, content_dir, output_dir)
    assert p == Path("/output/pastas/macaroni.html")
