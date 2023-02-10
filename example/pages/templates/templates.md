---
title: Templates
---

Hastie uses Jinja2 templates, for template syntax see their [Template Designer Documentation](https://jinja.palletsprojects.com/en/3.1.x/templates/).


## Page Variables

The following variables are made available to each page template:

Variable      | Description            
------------- | ----------------------
page          | Current page object
page.title    | Title for page
page.date     | Date for page
page.content  | HTML content for page
page.category | Category for page
page.url      | Permalink for page
pages         | List of all pages 
site          | site data object


## Category Variables

The following variables are made available to each category page (index.md):

Variable      | Description            
------------- | ----------------------
page          | Current category page
page.title    | Title for page
page.date     | Date for page
page.content  | HTML content for page
page.category | Category for page
page.url      | Permalink for page
pages         | List of all pages in category
site          | site data object



## Static  Directory

If a directory named `static` exists in the `templates` direcotry, Hastie will copy its contents as-is to the root of the `output_dir`.

So the following would be copied:

```
templates/static/css/style.css  -> output/css/style.css
templates/static/favicon.ico    -> output/favicon.ico
```
