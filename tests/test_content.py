""" Test module for Hastie resources."""

from pathlib import Path
from typing import Dict

# import pytest
# pytest.skip(allow_module_level=True)
import hastie.content

# Fixtures
config: Dict = {}
config["site"] = {}
config["site"]["baseurl"] = "/"

config_base: Dict = {}
config_base["site"] = {}
config_base["site"]["baseurl"] = "/hastie"


def test_read_page_basic():
    """Read docs page and confirm frontmatter parse."""
    f = Path("./docs/pages/index.md")
    page = hastie.content.read_page(f)
    assert page["filename"] == Path("docs/pages/index.md")
    assert page["title"] == "Welcome to Hastie"


def test_read_page_custom_var():
    """Read docs page and check for custom variable in frontmatter."""
    f = Path("./docs/pages/templates/index.md")
    page = hastie.content.read_page(f)
    assert page["topic"] == "templates"


def test_gather_categories():
    """Test determining categories from path with a single categoey"""
    global config

    content_dir = Path("./docs/pages")
    categories = hastie.content.gather_categories(content_dir, config)
    assert len(categories) == 1
    category = categories.pop()
    assert category["name"] == "templates"
    assert category["url"] == "/templates/"


def test_gather_categories_with_base():
    """Test determining categories from path with a single categoey"""
    global config_base

    content_dir = Path("./docs/pages")
    categories = hastie.content.gather_categories(content_dir, config_base)
    assert len(categories) == 1
    category = categories.pop()
    assert category["name"] == "templates"
    assert category["url"] == "/hastie/templates/"


# Add test for shortcodes
