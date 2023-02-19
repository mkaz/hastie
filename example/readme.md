
# Example Site

This example is the documentation site for Hastie.

The content files are in: [pages/](https://github.com/mkaz/hastie/tree/trunk/example/pages) directory.

The template files are in: [templates](https://github.com/mkaz/hastie/tree/trunk/example/templates) directory.

The directories are specified in the [hastie.toml](https://github.com/mkaz/hastie/tree/trunk/example/hastie.toml) file.

To generate the site, `cd` to this directory and run `hastie` binary. The generated files will be put in `output/`

You can serve using `python3 -m http.server --directory output/`
