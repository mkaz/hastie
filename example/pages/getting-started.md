---
title: Getting Started
---

## Binary Install

Download binaries for Linux, Mac, and Windows from [Github releases tab](https://github.com/mkaz/hastie/releases).

TBD

## Usage

```bash
$ hastie [-flags]

Flags:
  -version
        Display version and quit
```

Configuration file hastie.conf is in TOML format

```toml
base_url = "/"
content = "./pages"
output = "./output"
static = "./static"
templates = "../themes/docs/"

[site]
title = "Example Site"
author = "Marcus Kazmierczak"
```

Hastie walks through the `content` directory and finds all `.md` files. It applies the templates from `templates` directory and generates HTML copying to `output`. 

The file front matter can specify the template, using `template: filename` (without `.html` extension).

If not specified, the default templates are:
- `home` for the top-level `index.md`
- `category` for `index.md` files in direcotries
- `page` for all other pages

For this example content structure:

```
content/
    about.md
    index.md
    pastas/
        index.md
        macaroni.md
        spaghetti.md
```

- `about.md`, `macaroni.md`, and `spaghetti.md` would use `page` template
- `pastas/index.md` would use `category` template
- `index.md` would use `home` template

