
## Hastie - Static Site Generator in Go

Started: Feb 13, 2012

Hastie is my replacement of Jekyll, the name derived from the novel.

The primary reason I created hastie is it seems like a good project
to learn the Go language. There was no necessary improvement, outside
of possible speed, that I was looking to develop over Jekyll.

--------------------------------------------------------------------------------

## Install Notes

Uses blackfriday for markdown conversion. 
Install:
  $ goinstall github.com/russross/blackfriday

Users goconf for reading configuration file
Install:
  $ goinstall goconf.googlecode.com/hg

  If the above does not work, try
  $ cd $GOROOT/src/pkg/goconf/googlecode.com/hg
  $ gomake install


--------------------------------------------------------------------------------

## Usage

Hastie takes a templates directory and generates HTML files to a publish 
directory. It uses Go's template language, similar to mustache for templates
and markdown for files.

Planned to also support less stylesheets 

Hastie is intended as replacement of jekyll (for myself), but jekyll has a 
robust plugin, extensibility and community that I do not expect to even attempt.
If you are looking for a flexible tool to publish your site use jekyll.

If you are looking for a tool to tweek and play with the Go language, then this
might be your choice. Most customizations will probably require code changes.


--------------------------------------------------------------------------------


### TODO
* Recent Files: Recent By Category

* create next-prev link in category
* create less converter
* read .html files and apply template, no markdown
* rss.xml

* Command Line Arguments
* --verbose : --help

### DONE
* convert markdown
* read posts
  * parameters, title, date
  * apply layout

* write out file

* read files by category
* create category pages
* create index page
* Recent Files: Filter out non-date from Recent
* Recent Files: Include Full Link
* Recent Files: Need Titles
* Recent Files: Need Dates

* Add debug parameter to spit out info every step


### DEPLOY
* create github
* write docs






