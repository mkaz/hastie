
# Hastie - static site generator

Hastie is a static site generator, it processes a folder of markdown text files, applies a template, and generates an HTML site.


## Using Hastie

To use Hastie:

1. Create your content in markdown organized by directory
2. Specify templates to use in hastie.conf
3. Run Hastie to smash the two together, output to `output/`
4. Upload and serve.


## Install

For development, I symlink ~/bin/hastie to ./src/hastie.py

Eventually I'll make it pippable.


## History

Hastie started as a Go project, but rewritten to Python. See [golang branch](https://github.com/mkaz/hastie/tree/golang) for archived code.

**Why?** I switched my [mkaz.blog] site over to static for easier maintenance and looked briefly at Pelican but it does too much. Domain name aside, my site isn't really a blog in the reverse-chronological order of posts sense; it is a collection of pages in categories.

Pelican has a lot of requirements around dates, and template types. It generates numerous additional pages for author, tags, and other things I don't really want.

So I dusted off Hastie which does exactly what I want, but since I haven't coded in Go in awhile, I switched it over to Python.


## License

The project is licensed under the [MIT LICENSE](https://github.com/mkaz/hastie/blob/master/LICENSE).
