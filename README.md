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

Untill Go v1.0 is released you need to install Go Weekly: <http://golang.org/doc/install.html#fetch>

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
      "TemplateDir": "posts",
      "LayoutDir": "layouts",
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
  * sub-directories are only one level deep

The usage of hastie is just as a template engine, it does not copy over any images, does not have a built-in web server or any of the other features that jekyll has.

I keep the `public` directory full with all of the assets for the site such as images, stylesheets, etc and hastie copies in the html files. So if you delete a template it won't be removed from `public`


Data available to templates:

    .Title -- Page Title
    .Date -- Page Date format using .Date.Format "Jan 2, 2006"
    .Content -- Converted HTML Content
    .Category -- Category (directory)
    .OutFile -- file path
    .Recent -- list of 3 most recent files
    .Pages -- list of all page obhects
    .Url   -- what is the url for this page

--------------------------------------------------------------------------------

### TODO

* Recent Files: Create List of Recent Files by Category
* Create Next-Prev Links by Category
* Create LESS converter for stylesheets
* Create syntax highlighting blocks

* Read .html files and apply template, no markdown
* Add ability to support rss.xml




--------------------------------------------------------------------------------

### CHANGE LOG



ver 0.3 

  * Merged Fredrik Steen changes in github.com/stone/hastie

  * Switched config to json format
    - removed dependency on old config

  * Moved to Go1 support 




ver 0.2 (unreleased)

  * In config, renamed `template_dir` to `source` This more accurately describes the directory, what I was thinking was templates to be expanded are really the source files for the site.

  * Added Url parameter to templates

