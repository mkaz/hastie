---
title: Welcome to Hastie
---

Hastie is a static site generator similar to [Jekyll](https://github.com/mojombo/jekyll), [Pelican](https://getpelican.com/), or the many others. It takes content written in markdown and mashes them up with templates written in Jinja to create HTML to be published anywhere.

## Structure

The file system of folders and files in the content directory define the structure for the web site.

The top level directories define the categories. These category directories contain pages defined by a single markdown file, or a directory with an index.md file. Directories are useful when including additional assets, such as images, for the page.

```
content/

	index.md	# home page
	about.md	# about page
	now.md		# now page

	category-one/
		index.md			# index file for category
		cat-page.md			# page inside category

		another-cat-page/
			index.md		# another cat page
			image.jpg		# image for another
```

## Background

I created Hastie as a project to play with and learn Go, I'm a firm believer the best way to learn something is to build something. I've tweaked Hastie over the years, so now it is less blogging focus and more page generation.

After almost a decade working at Automattic and contributing to WordPress, I've returned to Hastie, rewritten it in Python, my preferred language, and using it again to publish my site: [mkaz.blog](https://mkaz.blog/).

The tool is a personal hobby to play with programming. If you are looking for a feature-filled easy way to get a site up and running, go with [WordPress.com](https://wordpress.com/), or one of the more mature site generators.

If you are looking for something to tinker with, probably requires tweaking,  and learning along the way, then Hastie might be it. It is open sourced to hopefully help others learn or play with.

Note: The name Hastie is from a character in the novel Dr. Jekyll and Mr. Hyde
