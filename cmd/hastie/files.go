package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"
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

// buildOutFile for where to save the page
// generate the public file path based on the source file path
func buildOutFile(filename, ext string) (outfile string) {
	sourceDir := strings.TrimRight(config.SourceDir, "/")
	outfile = filename[strings.Index(filename, sourceDir)+len(sourceDir)+1:]
	outfile = strings.Replace(outfile, ".md", ext, 1)

	base := filepath.Base(outfile)
	// HACK: if file starts with 20 or 19 assume date
	// TODO: deprecate setting date via filename
	if base[0:2] == "20" || base[0:2] == "19" {
		// remove date from filename
		outfile = strings.Replace(outfile, base[0:11], "", 1)
	}
	return outfile
}

// getDateFromFilename parses a filename of 2010-03-25-file.md
func getDateFromFilename(filename string) (dt time.Time) {
	base := filepath.Base(filename)
	// HACK: if file starts with 20 or 19 assume date
	// TODO: deprecate setting date via filename
	if base[0:2] == "20" || base[0:2] == "19" {
		dt, _ = time.Parse("2006-01-02", base[0:10])
	}
	return dt
}

// getCategoryFromFilename parses a filename of foo/bar/file.md into a category of foo_bar
func getCategoryFromFilename(filename string) (category string) {
	if strings.Contains(filename, string(os.PathSeparator)) {
		category = filename[0:strings.LastIndex(filename, string(os.PathSeparator))]
	}
	return category
}
