---
title: Markdown Content Page
---

## Frontmatter

Hastie uses front matter to specify parameters for each page. The parameters are specified at the top of the markdown document. Hastie can support both YAML and TOML frontmatter, personal preference for what you want to use.

If converting from another static site generator Jekyll sites use YAML and Zola sites use TOML.

### YAML Frontmatter

If using YAML frontmatter (default) then the section is delimited by `---` and
variables are defined using `:`

Markdown frontmatter example using YAML:

```markdown
---
title: Blog carefully my friend
template: page
date: 2012-02-14
---

Page content here...
```

### TOML Frontmatter

The TOML frontmatter uses `+++` to delimite the section, and `=` for defining
variables.

Markdown frontmatter example using TOML:

```markdown
+++
title = Blog carefully my friend
template = page
date = 2012-02-14
+++

Page content here...
```


## Document States

There are two properties you can use to set the document state, **draft** and
**archive**.

A draft document will not show in page lists and will not be
generated. A draft is used so you can check in a document as you are working on
it, but not publish. (Unfortunately, no dev mode yet to view drafts)

An archive document will not show in page lists but is generated amd written
out. Archive documents are used to avoid removing content, it just becomes
unlisted.

Archive example:

```markdown
+++
title = Some old document
date = 2004-05-15

archive = true
+++

...
```

Draft example:

```markdown
+++
title = Fresh new document
date = 2023-02-27

draft = true
+++

...
```

## Code Highlighting

Use fenced code blocks to define a code section with triple ticks and the name
of the language. See the [Python Markdown
docs](https://python-markdown.github.io/extensions/fenced_code_blocks/) for additional details.

~~~markdown
Here is my fenced code block.

```python
for i in range(10):
    print(i)
```
~~~

Hastie uses the codehilite extension to provide syntax highlighting. This requires including the necessary CSS to colorize the content. See the [templates page](../templates/templates/) for details.
