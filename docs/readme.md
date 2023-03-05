
# Hastie Documentation Site

An example Hastie site to documentat Hastie. View published at: <https://mkaz.github.io/hastie>

The content files are in: [pages/](https://github.com/mkaz/hastie/tree/trunk/example/pages)

The template files are in: [templates](https://github.com/mkaz/hastie/tree/trunk/example/templates)

See the config: [hastie.toml](https://github.com/mkaz/hastie/tree/trunk/example/hastie.toml)

To generate the site:

```bash
# install hastie
pip install git+https://github.com/mkaz/hastie

# clone site
git clone https://github.com/mkaz/hastie

# change to docs and run
cd hastie/docs/
hastie
```

To generate the site from development (without install):

```bash
# clone site
git clone https://github.com/mkaz/hastie
cd hastie

poetry install
cd docs/
poetry run python ../hastie/main.py

```

For both methods the generated files are put in `output/`

Serve using `python3 -m http.server --directory output/`

---

💡 If you followed the above locally, you should notice the CSS doesn't load, for that you need to specify a `baseurl` on the command-line.

The default config is setup to run from: `mkaz.github.io/hastie` with the `baseurl` set to `/hastie` in the `hastie.toml` config file.

For local, it runs top level so uses command-line parameter to override:

```
hastie --baseurl "/"
```
