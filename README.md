
# Hastie - Static Site Generator

Hastie is a static site generator, it processes a folder of markdown text files, applies a template, and generates an HTML site.

## Example Sites

* [Hastie Documentation](https://mkaz.github.io/hastie)
* [Working with Go](https://mkaz.github.io/working-with-go/)
* [Working with Vim](https://mkaz.github.io/working-with-vim/)
* [My Recipes site](http://kazmierczaks.com/)

Do you use Hastie? Submit your site via PR.

## Using Hastie

To use Hastie, first create a directory of source files in markdown, and specify templates to use, either pre-built or create your own. You then run hastie to smash the two together producing a site of HTML files. Upload and serve.

For the [Hastie documentation site](https://mkaz.github.io/hastie), the template files are available at [themes/docs](https://github.com/mkaz/hastie/tree/master/themes/docs) and the pages and config are in the [example directory](https://github.com/mkaz/hastie/tree/master/example).

To generate the documentation, after downloading a binary:

1. Clone the repository.
2. Change to `hastie/example` directory.
3. Run `hastie` to generate.
4. Files output to `docs/` directory per `hastie.json` config.

[Read the docs](https://mkaz.github.io/hastie) for customization and usage.

## Contribute

All contributions are welcome. Please use Github to submit feature requests, bug reports, documentation, or pull requests.

Thanks to all [the contributors](https://github.com/mkaz/hastie/graphs/contributors)!

## License

The project is licensed under the [MIT LICENSE](https://github.com/mkaz/hastie/blob/master/LICENSE).
