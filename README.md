
## Hastie - Static Site Generator in Go

Hastie is a simple static site generator in Go. I created it as a replacement of Jekyll. I wanted a project to play with and learn Go and Jekyll was slow and Ruby dependencies give me a headache. The Go binary is completely portable and includes all dependencies, so once built makes it easy in the future.

If you are looking for a tool to tweak and play with the Go language, then this might be fun. Most customizations will probably require code changes.  The reason I created the tool was to learn Go, I'm open sourcing to hopefully help others play with Go.

If you just want simple blogging and no headaches, setup a hosted blog on [WordPress.com](http://wordpress.com) it's an easy platform and you'll be up in minutes.

Note: The name Hastie is from a character in the novel Dr. Jekyll and Mr. Hyde

--------------------------------------------------------------------------------

## Install Notes

Install Go: <http://golang.org/doc/install.html#fetch>

Get Hastie: `go get github.com/mkaz/hastie`

If you have your Go environment setup, `go get` will automatically create the binary in `$GOPATH/bin`.

I set GOPATH to my home directory, so binaries go in `bin`

### Binaries

Download binaries for Linux, Mac, and Windows from [Github releases tab](https://github.com/mkaz/hastie/releases).


--------------------------------------------------------------------------------

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

    {
      "SourceDir" : "posts",
      "LayoutDir" : "layouts",
      "PublishDir": "public",
      "UseMarkdown": true
    }

UseMarkdown is optional parameter, by default it will convert documents to markdown. If you don't want documents to be converted globally, you can specify it false on the command-line or here in the config. If you want to disable markdown on a per document basis, you can put "markdown: no" in the front matter of the document.

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

  * all files must be have .md or .html extension

The usage of hastie is mainly a template engine, it does not have a built-in web server or watch, or many of the features that jekyll has.

If a directory named `static` exists in the Layout Directory, Hastie will copy it over to the root of the Publis Directory as `static`


### Template Variables

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

	.Reverse				  -- will reverse sort order of list
    .Recent.Limit n           -- will limit recent list to n items
    .Categories.Get CATEGORY  -- will fetch category list CATEGORY, useful for dynamic categories


#### Examples:

Show 3 most recent titles:

        {{ range .Recent.Limit 3 }}
          {{ .Title }}
        {{ end }}

Show 3 most recent from math category:

        {{ range .CategoryList.math }}
          {{ .Title }}
        {{ end }}

Show oldest items first:

		{{ range .Recent.Reverse }}
			{{.Title }}
		{{ end }}

### Using Filters (Example: Less CSS, CoffeeScript)

Hastie allows for the use of any command-line processing of files, provided the process takes the filename as input and spits out the results. It does so using `processFilters` configuration. You set a file extension mapped to the utility to process and the final extension.

Add follow configuration to hastie.json

      "processFilters": {
        "less": ["/usr/local/bin/lessc", "css"]
      }

So the above example any files with the extension `.less` will be converted to `.css` using lessc binary and copied to the public directory at the same spot in the directory tree as the original less file.

## Contributions

Thanks to all [the contributors](https://github.com/mkaz/hastie/graphs/contributors)!

Any contributions are welcome, be it in feature requests, bug reports, documentation, or pull requests.


## License

Licensed under MIT see [LICENSE](https://github.com/mkaz/hastie/blob/master/LICENSE) file.

