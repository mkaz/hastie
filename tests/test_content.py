""" Test module for Hastie resources."""

from pathlib import Path
import pytest

pytest.skip(allow_module_level=True)

import hastie


def test_get_page_basic():
    """Confirm title and filename reads in properly."""
    f = Path("./example/pages/index.md")
    page = hastie.content.get_page(f)
    assert page["filename"] == Path("example/pages/index.md")
    assert page["title"] == "Welcome to Hastie"


def test_gather_pages_basic():
    """Confirm pages gathered with correct, URL, and categories."""
    content_dir = Path("./example/pages")
    pages = hastie.content.gather_pages(content_dir)
    assert len(pages) == 3
    for p in pages:
        good = False
        if p["title"] == "Getting Started":
            assert p["category"] == ""
            assert p["parent"] == ""
            good = True
        elif p["title"] == "Markdown Content Page":
            assert p["category"] == ""
            assert p["parent"] == ""
            good = True
        elif p["title"] == "Templates":
            assert p["category"] == "templates"
            assert p["parent"] == ""
            good = True
    assert good


def test_gather_categories_basic():
    """Confirm categories gathered with correct, URL, and categories."""
    content_dir = Path("./example/pages")
    cats = hastie.content.gather_categories(content_dir)
    assert len(cats) == 2
    assert cats[0]["name"] == "templates"
    assert cats[0]["page"]["title"] == "Templates Index"
    assert cats[0]["parent"] == ""
    assert cats[1]["name"] == "sub"
    assert cats[1]["page"]["title"] == "Subcategory"
    assert cats[1]["parent"] == "templates"
