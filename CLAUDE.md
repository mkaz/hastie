# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Hastie is a static site generator written in Python 3.13+ that processes markdown files and applies Jinja2 templates to generate HTML sites. It takes a convention-over-configuration approach with minimal setup requirements and leverages modern Python features.

## Development Commands

**Install dependencies:**
```bash
uv sync
```

**Run the static site generator:**
```bash
uv run python hastie/main.py
# Or from within a site directory:
hastie
```

**Run tests:**
```bash
uv run pytest
```

**Code formatting and linting:**
```bash
uv run ruff format .
uv run ruff check .
uv run mypy .
```

**Pre-commit hooks (automatically runs ruff, mypy):**
```bash
pre-commit install
pre-commit run --all-files
```

## Architecture

**Core modules:**
- `hastie/main.py` - Entry point, orchestrates the site generation process
- `hastie/config.py` - Configuration management with CLI args and TOML config merging
- `hastie/content.py` - Page parsing, frontmatter processing, and markdown rendering
- `hastie/hfs.py` - File system operations for copying assets and organizing output
- `hastie/rss.py` - RSS feed generation
- `hastie/utils.py` - Sorting utilities (human_sort, date_sort)

**Site generation flow:**
1. Parse configuration from `hastie.toml` and CLI arguments
2. Copy static assets from content and template directories
3. Gather all pages and categories from content directory
4. Process each page: parse frontmatter → render markdown → apply Jinja2 template
5. Generate category index pages
6. Generate home page (index.md)
7. Optionally generate RSS feed

**Key concepts:**
- Pages are markdown files with YAML/TOML/JSON frontmatter
- Categories are inferred from directory structure
- Templates use Jinja2 with access to page, pages, categories, recent_pages, and site variables
- Output preserves directory structure from content directory
- Drafts (pages with `draft: true`) are excluded from output
- Configuration merges CLI args over TOML file settings
- Static assets are copied from configurable `static_dir` (defaults to `templates/static`)

**Template variables available in Jinja2:**
- `page` - Current page data
- `pages` - All pages (filtered by category for category templates)
- `categories` - All categories
- `recent_pages` - Pages sorted by date (most recent first)
- `site` - Site-wide configuration from hastie.toml

**Testing:**
Uses pytest with tests in the `tests/` directory. Test configuration is in `pyproject.toml` under `[tool.pytest.ini_options]`.