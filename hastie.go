/**
 *  _               _   _
 * | |             | | (_)
 * | |__   __ _ ___| |_ _  ___
 * | '_ \ / _` / __| __| |/ _ \
 * | | | | (_| \__ \ |_| |  __/
 * |_| |_|\__,_|___/\__|_|\___|
 *
 * Hastie - Static Site Generator
 *
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

const debug = false

type Page struct {
	Content  string
	Title    string
	Category string
	Layout   string
	Pages    PagesSlice
	Recent   PagesSlice
	Date     time.Time
	OutFile  string
}

var config map[string]string

type PagesSlice []Page

func (p PagesSlice) Len() int           { return len(p) }
func (p PagesSlice) Less(i, j int) bool { return p[i].Date.Unix() < p[j].Date.Unix() }
func (p PagesSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PagesSlice) Sort()              { sort.Sort(p) }

func main() {
	//TODO: add verbose / debug - command-line arguments

	setupConfig()

	// create list of directories and files
	site := &SiteStruct{}
	filepath.Walk(config["TemplateDir"], site, nil)

	/* ******************************************
	 * Loop through directories and build pages 
	 * ****************************************** */
	var pages PagesSlice
	for _, dir := range site.Directories {

		fmt.Println("Reading Directory: ", dir)
		directory_file_count := 0
		readglob := fmt.Sprintf("%s/*.md", dir)
		var dirfiles, _ = filepath.Glob(readglob)

		// loop through files in directory
		for _, file := range dirfiles {
			if debug {
				fmt.Println("  File: ", file)
			}
			outfile := filepath.Base(file)
			outfile = strings.Replace(outfile, ".md", ".html", 1)

			// read & parse file for parameters
			page := readParseFile(file)

			// skip file if no content
			if page.Content == "" {
				continue // skip to next file
			}
			pages = append(pages, page)
			directory_file_count += 1
		}
	}

	if debug { // spit out pages structure
		fmt.Println("################################################################################")
		fmt.Printf(" %-50s | %-10s | %s \n", "Title", "Category", "Outfile")
		fmt.Println("--------------------------------------------------------------------------------")
		for _, page := range pages {
			fmt.Printf(" %-50s | %-10s | %s \n", page.Title, page.Category, page.OutFile)
		}
		fmt.Println("################################################################################")
	}

	/* ******************************************
	 * Create any data needed from pages
	 * for example recent file list
	 * category list, etc...
	 * ****************************************** */

	// build recent file list, sorted
	recentList := getRecentList(pages)

	/* ******************************************
	 * Loop through pages and generate templates
	 * ****************************************** */
	for _, page := range pages {

		fmt.Println("  Generating Template: ", page.OutFile)

		/* Assign global data to page object
		 * Note: need better templating duplicating data
		         since no logic in templates to limit to 3 */
		page.Pages = recentList
		if len(recentList) > 3 {
			page.Recent = recentList[0:3]
		} else {
			page.Recent = recentList
		}

		/* Templating - writes page data to buffer 
		 * read and parse all template files          */
		buffer := new(bytes.Buffer)
		layoutsglob := fmt.Sprintf("%s/*.html", config["LayoutDir"])
		ts, err := template.ParseTemplateGlob(layoutsglob)
		if err != nil {
			fmt.Println("Error Parsing Templates: ", err)
			os.Exit(1)
		}
		// pick layout based on specified in file
		templateFile := ""
		if page.Layout == "" {
			templateFile = "post.html"
		} else {
			templateFile = fmt.Sprintf("%s.html", page.Layout)
		}
		ts.Execute(buffer, templateFile, page)

		// writing out file
		writedir := fmt.Sprintf("%s/%s", config["PublishDir"], page.Category)
		if debug {
			fmt.Println(" Write Directory: ", writedir)
		}
		os.MkdirAll(writedir, 0755) // does nothing if already exists

		outfile := fmt.Sprintf("%s/%s", config["PublishDir"], page.OutFile)
		if debug {
			fmt.Println(" Writing File: ", outfile)
		}
		ioutil.WriteFile(outfile, []byte(buffer.String()), 0644)
	}

}

/* ************************************************
 * Read and Parse File
 * @param filename
 * @return Page object
 * ************************************************ */
func readParseFile(filename string) (page Page) {
	if debug {
		fmt.Println("Parsing File: ", filename)
	}
	epoch, _ := time.Parse("20060102", "19700101")

	// setup default page struct
	page = Page{
		Title:    "",
		Category: "",
		Content:  "",
		Layout:   "",
		Date:     epoch,
		OutFile:  filename}

	// read file
	var data, err = ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error Reading: ", filename)
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
				key := strings.TrimSpace(line[:colonIndex])
				value := strings.TrimSpace(line[colonIndex+1:])
				switch key {
				case "title":
					page.Title = value
				case "category":
					page.Category = value
				case "layout":
					page.Layout = value
				}
			}

		} else if found >= 2 {
			// params over
			lines = lines[i:]
			break
		}

		if line == "---" {
			found += 1
		}

	}

	// switch directory name to just category

	// chop off first directory, since that is the template dir
	page.OutFile = filename[strings.Index(filename, "/")+1:]
	page.OutFile = strings.Replace(page.OutFile, ".md", ".html", 1)

	// next first directory is category
	if strings.Contains(page.OutFile, "/") {
		page.Category = page.OutFile[0:strings.Index(page.OutFile, "/")]
	}

	// parse date from filename
	base := filepath.Base(page.OutFile)
	if base[0:2] == "20" || base[0:2] == "19" { //HACK: if file starts with 20 or 19 assume date
		page.Date, _ = time.Parse("2006-01-02", base[0:10])
		page.OutFile = strings.Replace(page.OutFile, base[0:11], "", 1) // remove date from final filename
	}

	// convert markdown content
	content := strings.Join(lines, "\n")
	output := blackfriday.MarkdownCommon([]byte(content))
	page.Content = string(output)

	return page
}

/* ************************************************
 * Build Recent File List
 *    - return array sorted most recent first
 *    - array includes real link (no date)
 *    - does not include files without date
 * ************************************************ */
func getRecentList(pages PagesSlice) (list PagesSlice) {
	fmt.Println("Creating Recent File List")
	for _, page := range pages {
		// pages without dates are set to epoch
		if page.Date.Format("2006") != "1970" {
			list = append(list, page)
		}
	}
	list.Sort()

	// reverse
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}

	return list
}

/* ************************************************
 *  File / Directory Walker
 * ************************************************ */
type SiteStruct struct {
	Files       []string
	Directories []string
	Categories  []string
}

func (v *SiteStruct) VisitDir(path string, f os.FileInfo) bool {
	v.Categories = append(v.Categories, f.Name())
	v.Directories = append(v.Directories, path)
	/* fmt.Println("Directory: ", f.Name)*/
	return true
}

func (v *SiteStruct) VisitFile(path string, f os.FileInfo) {
	v.Files = append(v.Files, path)
	/* fmt.Printf("File: %s  -- Path: %s \n", f.Name, path)*/
}

/* ************************************************
 * Check if File / Directory Exists
 * ************************************************ */
func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if e, ok := err.(*os.PathError); ok && (e.Error == os.ENOENT || e.Error == os.ENOTDIR) {
			return false
		}
		return false
	}
	return true
}

/* ************************************************
 * Read Config file or set defaults 
 * ************************************************ */
func setupConfig() {
	file, err := ioutil.Readfile("hastie.json")
	if err != nil {
		// set defaults
		config["TemplateDir"] = "posts"
		config["LayoutDir"] = "layouts"
		config["PublishDir"] = "public"
	} else {
		if err := json.Unmarshal(file, &config); err != nil {
			fmt.Printf("Error parsing config: %s", err)
			os.Exit(1)
		}
	}
}
