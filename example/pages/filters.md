---
title: Using Filters
layout: page
order: 4
---

Hastie allows for the use of command-line processing of files, provided the process takes the filename as input and spits out the results. It does so using `processFilters` configuration. You set a file extension mapped to the utility to process and the final extension.

Configuration in hastie.json

```json
"processFilters": {
    "less": ["/usr/local/bin/lessc", "css"]
}
```

The above example processes files with the extension `.less` and converts to `.css` using the lessc binary. The generated file is copied to the public directory at the same spot in the directory tree as the original less file.
