## Hastie - Static Site Generator in Go

Created by [Marcus Kazmierczak](http://mkaz.com/) and contributions by [Fredrik Steen](http://github.com/stone)

Started: Feb 13, 2012

License: MIT, see [LICENSE](https://github.com/mkaz/hastie/blob/master/LICENSE)

**About** 

Hastie is a simple static site generator in Go. I use it as a replacement of jekyll on [mkaz.com](http://mkaz.com/). I wanted a project to play with and learn Go and jekyll was starting to slow and ruby dependencies give me a headache. I think I switched systems and everything broke and I couldn't publish. The Go binary is completely portable and all includes all dependencies.

If you are looking for a tool to tweek and play with the Go language, then this might be fun. Most customizations will probably require code changes.  The reason I created the tool was to learn Go, I'm open sourcing to hopefully help others play with Go.

If you just want simple blogging and no headaches, setup a hosted blog on [WordPress.com](http://wordpress.com) easiest platform and you'll be up in minutes.

Note: The name Hastie is from a character in the novel Dr. Jekyll and Mr. Hyde

--------------------------------------------------------------------------------

## Install Notes

Install Go: <http://golang.org/doc/install.html#fetch>

Get Hastie: `go get github.com/mkaz/hastie`

If you have your Go environment setup, `go get` will automatically create the binary in $GOPATH/bin.

My setup is

```bash
mkdir -p $HOME/gocode/src
mkdir -p $HOME/gocode/bin

export GOPATH="$HOME/gocode"
export PATH="$GOPATH/bin:$PATH"
```

#### Libraries

Uses **blackfriday** for markdown conversion. `go get github.com/russross/blackfriday`


--------------------------------------------------------------------------------

## Usage

    usage: hastie [flags]
      -c="hastie.json": Config file
      -h=false: show this help
      -http="": HTTP service address (e.g., ':8080')
      -m=false: do not use markdown conversion
      -t=false: display timing
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
    .Params.BaseUrl -- BaseUrl as defined in hastie.json

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
* Create syntax highlighting blocks
* Add ability to support rss.xml

* Read .html files and apply template, no markdown
* Add Less and Filter processing of static files

#### Bugs
* Add nicer error message/detection when no config found
* Add nicer error detection when error with template
* Should template variables work in source files ??

--------------------------------------------------------------------------------

### CHANGE LOG

ver 0.5.1 - Feb 2014

 * Rename ThemeDir back to LayoutDir
    - wrong direction, confusing name

 * Alter markdown rendering so allows `<script>` tags
    - hastie is used on personal site, so does not need to be so strict


ver 0.5.0 - June 2013
 
 * Change LayoutDir to ThemeDir parameter

 * Add copying of {ThemeDir}/static directory to {PublishDir}/static
   This allows you to create a theme with static assets, such as css

 * Use category if specified in file



ver 0.?.? - June 2013

  * Reverted branch changes which monitor and served the site
    This was adding too much complexity and not core to the tool

  * If you want a web server to test, do the following from public directory
		$ python -m SimpleHTTPServer
		$ open http://localhost:8000/


ver 0.4.4 - March 23, 2012

  * Add BaseURL config parameter
  * Update Test Site 

  
ver 0.4.3 - March 22, 2012
  
  * Speed optimizations improved rendering of mkaz.com time by 75% from 2sec down to 500ms
  * Pass more objects by reference to reduce copying of large arrays
  * Moved template parsing out of loop of pages, only needed to do once
  * Added Timing flag to see what areas take most time ( -t )


ver 0.4.2 - March 20, 2012

  * Run file through "go fmt" for proper formating
  * Convert fmt.Sprintf to simple "+" for string concat


ver 0.4.1 - March 15, 2012

  * Categories under subdirected changed to _ instead of / in index
  * Cleaned up code comments
  * Updated documentation to include all new features
  * Add Extension to Template Parameter to Support RSS

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

