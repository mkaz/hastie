from pathlib import Path
from src.templates import what_template


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
