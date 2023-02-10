""" Test module for Hastie resources."""

from pathlib import Path
import src.hres as hres


def test_get_page_basic():
    """Confirm title and filename reads in properly."""
    f = Path("./example/pages/index.md")
    page = hres.get_page(f)
    assert page["filename"] == Path("example/pages/index.md")
    assert page["title"] == "Welcome to Hastie"


def test_gather_pages_basic():
    """Confirm pages gathered with correct, URL, and categories."""
    content_dir = Path("./example/pages")
    pages = hres.gather_pages(content_dir)
    assert len(pages) == 3
    assert pages[0]["title"] == "Getting Started"
    assert pages[0]["category"] == ""
    assert pages[0]["parent"] == ""
    assert pages[1]["title"] == "Markdown Content Page"
    assert pages[1]["category"] == ""
    assert pages[1]["parent"] == ""
    assert pages[2]["title"] == "Templates"
    assert pages[2]["category"] == "templates"


def test_gather_categories_basic():
    """Confirm categories gathered with correct, URL, and categories."""
    content_dir = Path("./example/pages")
    cats = hres.gather_categories(content_dir)
    assert len(cats) == 2
    assert cats[0]["name"] == "templates"
    assert cats[0]["page"]["title"] == "Templates Index"
    assert cats[0]["parent"] == ""
    assert cats[1]["name"] == "sub"
    assert cats[1]["page"]["title"] == "Subcategory"
    assert cats[1]["parent"] == "templates"
