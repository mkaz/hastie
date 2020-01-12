---
title: Templates
layout: page
order: 3
---

Hastie uses Go's standard template package to provide templating. See the [two themes in the repository](https://github.com/mkaz/hastie/tree/master/themes) for an example of template files.

Hastie looks for the template files in the `LayoutDir` defined in hastie.json config. When parsing a markdown page, the default template file is `post.html`. You can specify the `layout` parameter in the front-matter to use a different template. For example, the following front matter uses the `page.html` template:

```markdown
---
title: Templates
layout: page
---
```

## Template Variables

The variables available to each template:

.Title
: Page Title

.Date
: Page Date format using .Date.Format "Jan 2, 2006"

.Content
: Converted HTML Content

.Category
: Category (directory)

.OutFile
: file path

.Recent
: list most recent files, latest first

.Url
: Url for this page

.PrevUrl
: Previous Page Url

.PrevTitle
: Previous Page Title

.NextUrl
: Next Page Url

.NextTitle
: Next Page Title

.PrevCatUrl
: Previous Page Url by Category

.PrevCatTitle
: Previous Page Title by Category

.NextCatUrl
: Next Page Url by Category

.NextCatTitle
: Next Page Title by Category

.Params
: Map of User Parameters set in front matter

.Params.BaseURL
: BaseURL as defined in hastie.json

.Categories.CATEGORY
: list of most recent files for CATEGORY


## Template functions

.Reverse
: reverse sort order of list

.Recent.Limit n
: limit recent list to n items

.Trim
: trim leading/trailing slashes (relative links)

.Title
: convert string to title case

.ToLower
: convert string to lower case

.ToUpper
: convert string to upper case


### Examples:

Show 3 most recent titles:

```html
    {{ range .Recent.Limit 3 }}
        {{ .Title }}
    {{ end }}
```

Show 3 most recent from math category:

```html
    {{ range .CategoryList.math }}
        {{ .Title }}
    {{ end }}
```

Show oldest items first:

```html
    {{ range .Recent.Reverse }}
        {{.Title }}
    {{ end }}
```

Trim leading slash to make links relative:

```html
    <a href="{{ .Url | trim }}"> Relative link </a>
```

## Static  Directory

If a directory named `static` exists in the `LayoutDir`, Hastie will copy it as-is to the root of the `PublishDir`.
