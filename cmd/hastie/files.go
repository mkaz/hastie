package main

import (
	"os"
	"path/filepath"
	"strings"
)

// buildSiteDirecrtories walks down the path and builds
// an array of directories for the site, ignoring dirs
// starting with _ underscopre
func buildSiteDirectories(path string) (dirs []string) {
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
		return nil
	})
	return dirs
}

func buildPagesSlice(dir string, globstr string, pages PagesSlice) PagesSlice {
	readglob := dir + globstr
	var dirfiles, _ = filepath.Glob(readglob)

	// loop through files in directory
	for _, file := range dirfiles {
		log.Debug("  File:", file)
		outfile := filepath.Base(file)
		outfile = strings.Replace(outfile, ".md", ".html", 1)

		// read & parse file for parameters
		page := readParseFile(file)
		page.SourceFile = file

		// create array of parsed pages
		pages = append(pages, page)
	}
	return pages
}
