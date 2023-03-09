---
title: Markdown Content Page
---

[TOC]

## Frontmatter

Hastie uses frontmatter to specify parameters for each page. The parameters are specified at the top of the markdown document. Hastie can support both YAML and TOML frontmatter, personal preference for what you want to use.

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
title = 'Blog carefully my friend'
template = 'page'
date = 2012-02-14
+++

Page content here...
```


### Document States

There are two properties in the frontmatter that can specify the document state, **draft** and **archive**.

A draft document will not show in page lists and will not be generated. A draft is used so you can check in a document as you are working on it, but not publish. (There is no dev mode yet to view drafts)

An archive document will not show in page lists, but is still generated amd written out. Archive documents are used to avoid removing content and creating dead links on the internet, they just become unlisted.

Archive example:

```markdown
+++
title = 'Some old document'
date = 2004-05-15

archive = true
+++

...
```

Draft example:

```markdown
+++
title = 'Fresh new document'
date = 2023-02-27

draft = true
+++

...
```

### Custom Parameters

You can create your own custom parameters in the frontmatter that can be used in your templates however you wish. Anything specified is attached to the page object. For example with this frontmatter:

```markdown
---
title = 'My Document'

next_url = 'https://mkaz.blog'
next_name = 'mkaz.blog'
---
```

You could add a section in a template that uses it:

```jinja
{% if page.next_url and page.next_name %}
  <a href="{{ page.next_url}}">{{ page.next_name }}</a>
{% endif %}
```

## Page Content

After the frontmatter, add your page content in markdown format. This content will be translated to HTML and then made available to the template in the `page.content` variable.

### Table of Contents

Use the shortcode `[TOC]` in your markdown to place a table of contents based on the headings in your document. The markup inserted will be a nested list wrapped in a div with `toc` class. For example:

```html
<div class="toc">
  <ul>
    <li><a href="#header-1">Header 1</a></li>
      <ul>
        <li><a href="#header-2">Header 2</a></li>
      </ul>
  </ul>
</div>
```

### Code Highlighting

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
