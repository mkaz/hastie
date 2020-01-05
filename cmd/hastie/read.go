package main

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/mkaz/hastie/pkg/utils"
)

/* ************************************************
 * Read and Parse File
 * ************************************************ */
func readParseFile(filename string) Page {
	data, _ := ioutil.ReadFile(filename)

	// setup default page struct
	// inherit paramters from config
	page := Page{
		OutFile:   filename,
		Extension: ".html",
		Params:    config.Params,
	}

	page = parseContent(string(data), page)

	// only parse date from filename if not already set
	if page.Date.Format("2006") == "1970" {
		page.Date = getDateFromFilename(filename)
	}

	page.OutFile = buildOutFile(filename, page.Extension)

	// next directory(s) category, category includes sub-dir = solog/webdev
	if page.Category == "" {
		page.Category = getCategoryFromFilename(filename)
	}

	// add url of page, which includes initial slash
	// this is needed to get correct links for multi
	// level directories
	page.Url = "/" + utils.RemoveIndexHTML(page.OutFile)

	// convert content to markdown
	if (config.UseMarkdown) && (page.Params["markdown"] != "no") {
		output := markdown.ToHTML([]byte(page.Content), nil, nil)
		page.Content = string(output)
	}

	return page
}

// parseContent front matter
func parseContent(content string, page Page) Page {

	var lines = strings.Split(content, "\n")
	var found = 0
	for i, line := range lines {
		line = strings.TrimSpace(line)

		if found == 1 {
			// parse line for param
			colonIndex := strings.Index(line, ":")
			if colonIndex > 0 {
				key := strings.ToLower(strings.TrimSpace(line[:colonIndex]))
				value := strings.TrimSpace(line[colonIndex+1:])
				value = strings.Trim(value, "\"") //remove quotes
				switch key {
				case "title":
					page.Title = value
				case "category":
					page.Category = value
				case "layout":
					page.Layout = value
				case "extension":
					page.Extension = "." + value
				case "date":
					page.Date, _ = time.Parse("2006-01-02", value[0:10])
				default:
					page.Params[key] = value
				}
			}

		} else if found >= 2 {
			// frontmatter over
			// slurp up the rest
			lines = lines[i:]
			break
		}

		if line == "---" {
			found++
		}

	}
	page.Content = strings.Join(lines, "\n")
	return page
}
