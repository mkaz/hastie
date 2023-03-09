""" Test module for Hastie resources."""

from pathlib import Path

# import pytest
# pytest.skip(allow_module_level=True)
import hastie.content


def test_read_page_basic():
    """Read docs page and confirm frontmatter parse."""
    f = Path("./docs/pages/index.md")
    page = hastie.content.read_page(f)
    assert page["filename"] == Path("docs/pages/index.md")
    assert page["title"] == "Welcome to Hastie"


def test_read_page_custom_var():
    """Read docs page and check for custom variable in frontmatter."""
    f = Path("./docs/pages/templates/sub/example.md")
    page = hastie.content.read_page(f)
    assert page["topic"] == "moon"


def test_determine_categories_from_path_single_category():
    """Test determining categories from path with a single categoey"""
    content_dir = Path("./pages")
    file_parent = Path("./pages/templates")
    d = hastie.content.determine_categories_from_path(file_parent, content_dir)
    assert d["parent"] == ""
    assert d["category"] == "templates"


def test_determine_categories_from_path_no_category():
    """Test determining categories from file in root, no category"""
    content_dir = Path("./pages")
    file_parent = Path("./pages")
    d = hastie.content.determine_categories_from_path(file_parent, content_dir)
    assert d["parent"] == ""
    assert d["category"] == ""


def test_determine_categories_from_path_sub_category():
    """Test determining categories from file in sub category"""
    content_dir = Path("./pages")
    file_parent = Path("./pages/templates/sub")
    d = hastie.content.determine_categories_from_path(file_parent, content_dir)
    assert d["parent"] == "templates"
    assert d["category"] == "sub"
