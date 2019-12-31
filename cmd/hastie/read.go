package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/mkaz/hastie/pkg/utils"
)

/* ************************************************
 * Read and Parse File
 * ************************************************ */
func readParseFile(filename string) (page Page) {
	log.Debug("Parsing File:", filename)
	epoch, _ := time.Parse("20060102", "19700101")

	// setup default page struct
	page = Page{
		Date:      epoch,
		OutFile:   filename,
		Extension: ".html",
		Params:    make(map[string]string),
	}

	// read file
	var data, err = ioutil.ReadFile(filename)
	if err != nil {
		log.Warn("Error Reading: " + filename)
		return
	}

	// go through content parse from --- to ---
	var lines = strings.Split(string(data), "\n")
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
			// params over
			lines = lines[i:]
			break
		}

		if line == "---" {
			found++
		}

	}

	// chop off first directory, since that is the template dir
	log.Debug("Filename", filename)
	page.OutFile = filename[strings.Index(filename, string(os.PathSeparator))+1:]
	page.OutFile = strings.Replace(page.OutFile, ".md", page.Extension, 1)
	log.Debug("page.Outfile", page.OutFile)

	// next directory(s) category, category includes sub-dir = solog/webdev
	if page.Category == "" {
		if strings.Contains(page.OutFile, string(os.PathSeparator)) {
			page.Category = page.OutFile[0:strings.LastIndex(page.OutFile, string(os.PathSeparator))]
			page.SimpleCategory = strings.Replace(page.Category, string(os.PathSeparator), "_", -1)
		}
	}
	log.Debug("page.Category", page.Category)
	// parse date from filename
	base := filepath.Base(page.OutFile)
	if base[0:2] == "20" || base[0:2] == "19" { //HACK: if file starts with 20 or 19 assume date
		page.Date, _ = time.Parse("2006-01-02", base[0:10])
		page.OutFile = strings.Replace(page.OutFile, base[0:11], "", 1) // remove date from final filename
	}

	// add url of page, which includes initial slash
	// this is needed to get correct links for multi
	// level directories

	page.Url = "/" + utils.RemoveIndexHTML(page.OutFile)

	// convert markdown content
	content := strings.Join(lines, "\n")
	if (config.UseMarkdown) && (page.Params["markdown"] != "no") {
		output := markdown.ToHTML([]byte(content), nil, nil)
		page.Content = string(output)
	} else {
		page.Content = content
	}

	return page
}
