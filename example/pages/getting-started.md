---
title: Getting Started
layout: page
order: 1
---

## Binary Install

Download binaries for Linux, Mac, and Windows from [Github releases tab](https://github.com/mkaz/hastie/releases).


## Usage

```
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


### Static  Directory

If a directory named `static` exists in the `LayoutDir`, Hastie will copy it as-is to the root of the `PublishDir`.


### Markdown Front Matter

Hastie uses front matter to specify custom parameters. The parameters are specified at the top of the markdown document in a section delimited with `---`.

Example top of markdown document:

```
---
title: Blog carefully my friend
layout: post
date: 2012-02-14
---

This is my content...
```

### User-defined Parameters

Hastie supports user-defined parameters and makes them available to the templates using `.Params.YOURPARAM`

#### Example setting and using a parameter.

Setting parameter:
```
---
title: Blog carefully my friend
layout: post
date: 2012-02-14
guest: Hemingway
---
```

Using parameter in a template:

```
{{ if .Params.guest }}
    Guest Author: {{ .Params.guest }}
{{ end }}
```

### Disable Markdown

`UseMarkdown` is optional parameter in the config. By default, Hastie will convert documents from markdown to HTML. If you don't want documents to be converted globally, you can specify it false on the command-line or `UseMarkdown: false` in the JSON config.

If you want to disable markdown on a per document basis, you can put `markdown: no` in the front matter of the document.
