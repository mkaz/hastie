package main

import (
	"os"
	"path/filepath"
	"strings"
)

// getSiteFiles walks down the path given and builds the
// arrays of pages and directories for the site. Ignores
// dirs starting with _ underscopre
func getSiteFiles(path string) (pages PageList, dirs []string) {

	filepath.Walk(path, func(fn string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Fatal("Error building site directories", err)
			return nil
		}

		if fi.IsDir() {
			// ignore directories starting with _
			if strings.HasPrefix(fi.Name(), "_") {
				return filepath.SkipDir
			}
			dirs = append(dirs, fn)
			return nil
		}

		// Not directory so a file
		if filepath.Ext(fn) == ".md" || filepath.Ext(fn) == ".html" {
			page := readParseFile(fn)
			page.SourceFile = fn
			pages = append(pages, page)
		}
		return nil

	}) // end walk

	return pages, dirs
}
