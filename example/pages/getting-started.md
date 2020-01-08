---
title: Getting Started
layout: page
order: 1
---

## Binary Install

Download binaries for Linux, Mac, and Windows from [Github releases tab](https://github.com/mkaz/hastie/releases).


## Usage

```bash
$ hastie [-flags]

Flags:

  -config string
        Config file (default "hastie.json")
  -debug
        Debug output (verbose)
  -help
        show this help
  -nomarkdown
        do not use markdown conversion
  -verbose
        Show info level
  -version
        Display version and quit
```

Configuration file format (default ./hastie.json)

```json
{
  "SourceDir" : "posts",
  "LayoutDir" : "layouts",
  "PublishDir": "public"
}
```

Hastie walks through the `SourceDir` and finds all `.md` and `.html` files. It applies the templates from `LayoutDir` and generates HTML copying to `PublishDir`. It uses Go's template language for templates and markdown for content.

The [Hastie documentation site](https://mkaz.github.io/hastie/) is generated using Hastie. See the [github example directory](https://github.com/mkaz/hastie/tree/master/example) for configuration and pages and the [themes/docs](https://github.com/mkaz/hastie/tree/master/themes/docs) for templates files.

