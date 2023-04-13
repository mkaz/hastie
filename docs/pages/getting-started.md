---
title: Getting Started
date: 2023-01-03
---

## Install

```
pip install git+https://github.com/mkaz/hastie
```

## Usage

```bash
usage: hastie [-h] [-q] [-v] [-c CONF] [--baseurl BASEURL]

options:
  -h, --help            show this help message and exit
  -q, --quiet
  -v, --version
  -c CONF, --conf CONF  Config file
  --baseurl BASEURL     Override base url in config
```
## Config

Create a configuration file hastie.toml

```toml
content = "./pages"
output = "./output"
static = "./static"
templates = "./templates/"

[site]
title = "Example Site"
description = "Just another example site"
author = "Marcus Kazmierczak"

# set if site hosted at a subdirectory
# for example: https://mkaz.github.io/hastie
baseurl = "/hastie"
```

## Run

With the above config, just run `hastie`

```bash
$ hastie
Hastie v0.9.3
Generated 6 files in 0.086 sec
```

Hastie walks through the `content` directory and finds all `.md` files. It applies the templates from `templates` directory and generates HTML copying to `output` using the same directory structure.

Use directories to create categories for your ontent.

The default templates are:

- `index.html` for the top-level `index.md`
- `category.html` for `index.md` files in directories
- `page.html` for all other pages

The frontmatter in the markdown document can specify a different template, use `template: filename` (without `.html` extension).
