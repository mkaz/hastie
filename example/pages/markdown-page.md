---
title: Markdown Page
layout: page
order: 2
---

Hastie uses front matter to specify parameters for each page. The parameters are specified at the top of the markdown document in a section delimited with `---`.

Example top of markdown document:

```markdown
---
title: Blog carefully my friend
layout: post
date: 2012-02-14
---

Page content here...
```

The list of standard parameters are:

title
: Page Title

Category
: Page Category

layout
: Layout template file, without .html extension, default post

extension
: Page extension to save, default .html

date
: Page date, YYYY-mm-dd format

order
: Page order in AllPages list

unlisted
: Do not include in AllPages list


## User-defined Parameters

Hastie supports user-defined parameters. These are made available to the templates using `.Params.YOURPARAM`

### Example parameter for an individual page

Setting parameter in front matter:

```markdown
---
title: Blog carefully my friend
layout: post
date: 2012-02-14
guest: Hemingway
---
```

Using parameter in a template:

```html
{{ if .Params.guest }}
    Guest Author: {{ .Params.guest }}
{{ end }}
```

### Example parameter in global config

In `hastie.json`

```json
{
  "Params": {
    "SiteName": "Hastie"
  }
}
```

Usage in template is the same:

```html
<title> {{ .Title }} - {{ .Params.SiteName }} </title>
```

## Disable Markdown

`UseMarkdown` is an optional parameter in the config. By default, Hastie will convert documents from markdown to HTML. If you don't want documents to be converted globally, you can specify it false on the command-line or `UseMarkdown: false` in the JSON config.

If you want to disable markdown on a per document basis, you can put `markdown: no` in the front matter of the document.
