from pathlib import Path
from src.templates import *


def test_what_template_index():
    content_dir = Path("/pages")
    filepath = Path("/pages/index.md")
    tpl = what_template(filepath, content_dir)
    assert tpl == "index.html"


def test_what_template_page():
    content_dir = Path("/pages")
    filepath = Path("/pages/getting-started.md")
    tpl = what_template(filepath, content_dir)
    assert tpl == "page.html"


def test_what_template_category_index():
    content_dir = Path("/pages")
    filepath = Path("/pages/pastas/index.md")
    tpl = what_template(filepath, content_dir)
    assert tpl == "category.html"


def test_what_template_category_page():
    content_dir = Path("/pages")
    filepath = Path("/pages/pastas/macaroni.md")
    tpl = what_template(filepath, content_dir)
    assert tpl == "page.html"


def test_get_output_dir_index():
    content_dir = Path("/pages")
    output_dir = Path("/output")
    filepath = Path("/pages/index.md")
    p = get_output_file(filepath, content_dir, output_dir)
    assert p == Path("/output/index.html")


def test_get_output_dir_page():
    content_dir = Path("/pages")
    output_dir = Path("/output")
    filepath = Path("/pages/getting-started.md")
    p = get_output_file(filepath, content_dir, output_dir)
    assert p == Path("/output/getting-started.html")


def test_get_output_dir_category_index():
    content_dir = Path("/pages")
    output_dir = Path("/output")
    filepath = Path("/pages/pastas/index.md")
    p = get_output_file(filepath, content_dir, output_dir)
    assert p == Path("/output/pastas/index.html")


def test_get_output_dir_category_page():
    content_dir = Path("/pages")
    output_dir = Path("/output")
    filepath = Path("/pages/pastas/macaroni.md")
    p = get_output_file(filepath, content_dir, output_dir)
    assert p == Path("/output/pastas/macaroni.html")
