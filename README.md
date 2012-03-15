## Hastie - Static Site Generator in Go

Author :

 - Marcus Kazmierczak, http://mkaz.com/
 - Fredrik Steen, http://github.com/stone

Started: Feb 13, 2012
Project: https://github.com/mkaz/hastie

Hastie is intended as replacement of jekyll (for myself), but jekyll has a robust plugin, extensibility and community that I do not expect to even attempt.  If you are looking for a flexible tool to publish your site use jekyll.

If you are looking for a tool to tweek and play with the Go language, then this might be your choice. Most customizations will probably require code changes.  The reason I created the tool was to learn Go, I'm publishing to hopefully help others with playing with the language.

Note: The name Hastie also comes from the novel Dr. Jekyll and Mr. Hyde

--------------------------------------------------------------------------------

## Install Notes

Until Go v1.0 is released you need to install Go Weekly: <http://golang.org/doc/install.html#fetch>

    $ cd $HOME
    $ hg clone -u weekly https://go.googlecode.com/hg/ go
    $ cd $HOME/go/src; ./all.bash
    $ echo "GOROOT=$HOME/go" >> $HOME/.bashrc


#### Libraries

Uses **blackfriday** for markdown conversion. `go get github.com/russross/blackfriday`


--------------------------------------------------------------------------------

## Usage

    usage: hastie [flags]
      -c="hastie.json": Config file
      -h=false: show this help
      -v=false: verbose output

Configuration file format (default ./hastie.json)

    {
      "SourceDir" : "posts",
      "LayoutDir" : "layouts",
      "PublishDir": "public"

    }


Hastie walks through a templates directory and generates HTML files to a publish directory. It uses Go's template language for templates and markdown for content.

Here is sample site layout: (see test directory)

    layouts/footer.html
    layouts/header.html
    layouts/indexpage.html
    layouts/post.html
    posts/2011-03-02-angelica.html
    posts/index.md
    posts/zebra/2009-12-12-sample-post.md
    posts/zebra/2012-02-14-hastie-intro.md


This will generate:

    public/angelica.html
    public/index.html
    public/zebra/sample-post.html
    public/zebra/hastie-intro.html


A few current limitations:

  * all files must be have .md extension

The usage of hastie is just as a template engine, it does not copy over any images, does not have a built-in web server or any of the other features that jekyll has.

I keep the `public` directory full with all of the assets for the site such as images, stylesheets, etc and hastie copies in the html files. So if you delete a template it won't be removed from `public`


Data available to templates:

    .Title        -- Page Title
    .Date         -- Page Date format using .Date.Format "Jan 2, 2006"
    .Content      -- Converted HTML Content
    .Category     -- Category (directory)
    .OutFile      -- file path
    .Recent       -- list most recent files, latest first
    .Url          -- Url for this page
    .PrevUrl      -- Previous Page Url
    .PrevTitle    -- Previous Page Title
    .NextUrl      -- Next Page Url
    .NextTitle    -- Next Page Title
    .PrevCatUrl   -- Previous Page Url by Category
    .PrevCatTitle -- Previous Page Title by Category
    .NextCatUrl   -- Next Page Url by Category
    .NextCatTitle -- Next Page Title by Category
    .Params       -- Map of User Parameters, set in page head

    .Categories.CATEGORY -- list of most recent files for CATEGORY


Functions Available:
    
    .Recent.Limit n           -- will limit recent list to n items
    .Categories.Get CATEGORY  -- will fetch category list CATEGORY, useful for dynamic categories


Examples:

    Show 3 most recent titles: 
        {{ range .Recent.Limit 3 }}
          {{ .Title }}
        {{ end }}
    
    Show 3 most recent from math category:
        {{ range .CategoryList.math }}
          {{ .Title }}
        {{ end }}


### Using Filters (Example: Less CSS, CoffeeScript)

Hastie allows for the use of any command-line processing of files, provided the process takes the filename as input and spits out the results. It does so using `processFilters` configuration. You set a file extension mapped to the utility to process and the final extension.

Add follow configuration to hastie.json

      "processFilters": {
        "less": ["/usr/local/bin/lessc", "css"]
      }

So the above example any files with the extension `.less` will be converted to `.css` using lessc binary and copied to the public directory at the same spot in the directory tree as the original less file.


--------------------------------------------------------------------------------

### TODO
* Create LESS converter for stylesheets
* Create syntax highlighting blocks

* Add ability to support rss.xml

* Read .html files and apply template, no markdown
* Expand examples to use categories, limit and new feature sets


#### Bugs
* Add nicer error message/detection when no config found
* Add nicer error detection when error with template
* Shouldn't templates work in source files ??

--------------------------------------------------------------------------------

### CHANGE LOG

ver 0.4.0 - March 14, 2012
  * Add Process Filters allows processing using any third party such as Less CSS or Coffee Script


ver 0.3.5 - March 12, 2012

  * Add Prev-Next Links by Category
  * Add Parameters to Header, allows user created parameters, stored in .Params
  * Removed skip if empty content, can build page based on parameters


ver 0.3.4 - March 10, 2012

  * Add Recent List by Category
  * Switched Config from string map to struct
  * Created new config element called CategoryMash which allows
    the combination of multiple categories into a single category.
    This allows for displaying a list of combined categories


ver 0.3.3 - March 9, 2012

  * Add Limit function to PagesSlice, availabe in templates
  * Removed Pages data, since duplicated


ver 0.3.2 - March 8, 2012

  * Add Prev-Next Links to Page Object


ver 0.3.1 - March 7, 2012

  * Change category to include full directory path
  * Trimmed begin-end quotes from passed in parameters
    no need to quote parameters in template files


ver 0.3.0 - March 2, 2012

  * Merged Fredrik Steen changes in github.com/stone/hastie
  * Switched config to json format
    - removed dependency on old config
  * Moved to Go1 support 


ver 0.2 (unreleased)
  * In config, renamed `template_dir` to `source` This more accurately describes the directory, what I was thinking was templates to be expanded are really the source files for the site.

  * Added Url parameter to templates

