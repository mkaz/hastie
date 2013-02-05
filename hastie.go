/**
 *  _               _   _
 * | |             | | (_)
 * | |__   __ _ ___| |_ _  ___
 * | '_ \ / _` / __| __| |/ _ \
 * | | | | (_| \__ \ |_| |  __/
 * |_| |_|\__,_|___/\__|_|\___|
 *
 * Hastie - Static Site Generator
 * https://github.com/mkaz/hastie
 *
 */

package hastie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

type Config struct {
	SourceDir, LayoutDir, PublishDir, BaseUrl string
	CategoryMash                              map[string]string
	ProcessFilters                            map[string][]string
	NoMarkdown                                bool
}

var DefaultConfig = Config{SourceDir: "posts",
	LayoutDir:  "layouts",
	PublishDir: "public",
	NoMarkdown: false}

type Page struct {
	Content, Title, Category, SimpleCategory, Layout, OutFile, Extension, Url, PrevUrl, PrevTitle, NextUrl, NextTitle, PrevCatUrl, PrevCatTitle, NextCatUrl, NextCatTitle string
	Params                                                                                                                                                                map[string]string
	Recent                                                                                                                                                                *PagesSlice
	Date                                                                                                                                                                  time.Time
	Categories                                                                                                                                                            *CategoryList
	List                                                                                                                                                                  bool
}

type PagesSlice []Page

func (p PagesSlice) Get(i int) Page { return p[i] }
func (p PagesSlice) Len() int       { return len(p) }

// Less if used to sort PagesSlice by date descending then title ascending
func (p PagesSlice) Less(i, j int) bool {
	d := p[i].Date.Unix() - p[j].Date.Unix()
	if d < 0 {
		return false
	} else if d > 0 {
		return true
	}
	return p[i].Title < p[j].Title
}

func (p PagesSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PagesSlice) Sort()         { sort.Sort(p) }
func (p PagesSlice) Limit(n int) PagesSlice {
	if n < len(p) {
		return p[0:n]
	}
	return p
}
func (p PagesSlice) Listed() PagesSlice {
	listed := PagesSlice{}
	for _, page := range p {
		if page.List {
			listed = append(listed, page)
		}
	}
	return listed
}

type CategoryList map[string]PagesSlice

func (c CategoryList) Get(category string) PagesSlice { return c[category] }

// Compile uses the hastie config and generates the required static files updating monitor as is goes.
// If monitor is nil all output will be discarded by delegating to DiscardMonitor.
func (config Config) Compile(monitor Monitor) error {
	if monitor == nil {
		monitor = DiscardMonitor
	}

	monitor.Start()
	defer monitor.End()

	site := SiteStruct{}

	// Need source dirs relative to their parent so can correctly extract categories
	filepath.Walk(config.SourceDir, site.Walker())
	monitor.Walked()

	/* ******************************************
	 * Loop through directories and build pages
	 * ****************************************** */
	var pages PagesSlice
	for _, dir := range site.Directories {

		readglob := dir + "/*.md"
		var dirfiles, _ = filepath.Glob(readglob)

		// loop through files in directory
		for _, file := range dirfiles {
			monitor.ParsingSource(file)

			// Make outfile relative to source dir
			outfile, err := filepath.Rel(config.SourceDir, file)
			if err != nil {
				return err
			}
			// let parsing below work out the the extension
			// read & parse file for parameters
			page, err := config.readParseFile(file, outfile, config.NoMarkdown)
			if err != nil {
				return err
			}

			// create array of parsed pages
			pages = append(pages, page)
		}
	}
	monitor.ParsedSources()

	/* ******************************************
	 * Create any data needed from pages
	 * ****************************************** */

	// Sort page into order that we want to display them
	pages.Sort()

	// Filter out those pages that are not listed
	recentList := pages.Listed()
	recentListPtr := &recentList

	// category list is sorted map of pages by category
	categoryList := config.getCategoryList(recentListPtr)
	categoryListPtr := &categoryList

	monitor.Listed()

	// read and parse all template files
	layoutsglob := config.LayoutDir + "/*.html"
	ts, err := template.ParseGlob(layoutsglob)
	if err != nil {
		return fmt.Errorf("Error Parsing Templates: %s", err)
	}
	monitor.ParsedTemplates()

	/* ******************************************
	 * Loop through pages and generate templates
	 * ****************************************** */
	for _, page := range pages {
		monitor.ParsingTemplate(page.OutFile)

		// add recent pages lists to page object
		page.Recent = recentListPtr
		page.Categories = categoryListPtr

		// add prev-next links
		page.buildPrevNextLinks(recentListPtr)

		// Templating - writes page data to buffer
		buffer := new(bytes.Buffer)

		// pick layout based on specified in file
		templateFile := ""
		if page.Layout == "" {
			templateFile = "post.html"
		} else {
			templateFile = page.Layout + ".html"
		}
		ts.ExecuteTemplate(buffer, templateFile, page)

		// writing out file
		writedir := config.PublishDir + "/" + page.Category
		os.MkdirAll(writedir, 0755) // does nothing if already exists

		outfile := config.PublishDir + "/" + page.OutFile
		monitor.WritingTemplate(outfile)
		ioutil.WriteFile(outfile, []byte(buffer.String()), 0644)
	}
	monitor.GeneratedTemplates()

	/* ******************************************
		 * Process Filters
	   * proces filters are a mapping of file extensions to commands
	   * and an output extensions. find files with extension, run
	   * command which should spit out text and create new file.extension
	   * For example: Less CSS or CoffeeSript
		 * ****************************************** */
	for ext, filter := range config.ProcessFilters {
		extStart := "." + ext
		extEnd := "." + filter[1]

		for _, dir := range site.Directories {
			readglob := dir + "/*" + extStart
			var dirfiles, _ = filepath.Glob(readglob)
			for _, file := range dirfiles {
				// TODO: check for filter exists
				//apply process filter command, capture output
				cmd := exec.Command(filter[0], file)
				output, err := cmd.Output()
				if err != nil {
					fmt.Errorf("Error Process Filter: %s caused by: %s", file, err)
					return err
				}

				// determine output file path and extension
				// Make outfile relative to source dir
				outfile, err := filepath.Rel(config.SourceDir, file)
				if err != nil {
					return err
				}
				outfile = config.PublishDir + "/" + outfile
				outfile = strings.Replace(outfile, extStart, extEnd, 1)
				ioutil.WriteFile(outfile, output, 0644)
			}
		}
	}
	monitor.Filtered()

	return nil
}

/* ************************************************
* Build Category List
*    - return a map containing a list of pages for
       each category, the key being category name
* ************************************************ */
func (config Config) getCategoryList(pages *PagesSlice) CategoryList {
	mapList := make(CategoryList)
	// recentList is passed in which is already sorted
	// just need to map the pages to category

	// read category mash config, which allows to create
	// a new category based on combining multiple categories
	// this is used on my site when I want to display a list
	// of recent items from similar categories together
	reverseMap := make(map[string]string)

	// config consists of a hash with new category being the
	// key and a comma separated list of existing categories
	// being the value, create a reverse map
	for k, v := range config.CategoryMash {
		cats := strings.Split(string(v), ",")
		//loop through split and add to reverse map
		for _, cat := range cats {
			reverseMap[cat] = string(k)
		}
	}

	for _, page := range *pages {

		// create new category from category mash map
		if reverseMap[page.Category] != page.Category {
			thisCategory := reverseMap[page.Category]
			mapList[thisCategory] = append(mapList[thisCategory], page)
		}

		// still want a list of regular categories
		// simpleCategory replaces / in sub-dir categories to _
		// this always the categorty to be referenced in template
		simpleCategory := strings.Replace(page.Category, "/", "_", -1)
		mapList[simpleCategory] = append(mapList[simpleCategory], page)
	}
	return mapList
}

/* ************************************************
 * Add Prev Next Links to Page Object
 * ************************************************ */
func (page *Page) buildPrevNextLinks(recentList *PagesSlice) {
	foundPage := false

	nextPage := Page{}
	prevPage := Page{}
	nextPageCat := Page{}
	prevPageCat := Page{}
	lastPageCat := Page{}

	for i, rp := range *recentList {
		if rp.Category == page.Category {
			if foundPage {
				prevPageCat = rp
				break
			}
		}

		if rp.Title == page.Title {
			foundPage = true
			nextPageCat = lastPageCat
			if i != 0 {
				nextPage = recentList.Get(i - 1)
			}
			if i+1 < recentList.Len() {
				prevPage = recentList.Get(i + 1)
			}
		}

		if rp.Category == page.Category {
			lastPageCat = rp // previous page
		}
	}

	page.NextUrl = nextPage.Url
	page.NextTitle = nextPage.Title
	page.PrevUrl = prevPage.Url
	page.PrevTitle = prevPage.Title

	page.NextCatUrl = nextPageCat.Url
	page.NextCatTitle = nextPageCat.Title
	page.PrevCatUrl = prevPageCat.Url
	page.PrevCatTitle = prevPageCat.Title
}

/* ************************************************
 * Read and Parse File
 * ************************************************ */
func (config Config) readParseFile(filename string, outfile string, nomarkdown bool) (Page, error) {
	epoch, _ := time.Parse("20060102", "19700101")

	// setup default page struct
	page := Page{
		Title: "", Category: "", SimpleCategory: "", Content: "", Layout: "", Date: epoch, OutFile: outfile, Extension: ".html",
		Url: "", PrevUrl: "", PrevTitle: "", NextUrl: "", NextTitle: "",
		PrevCatUrl: "", PrevCatTitle: "", NextCatUrl: "", NextCatTitle: "",
		Params: make(map[string]string),
		List:   false,
	}

	// read file
	var data, err = ioutil.ReadFile(filename)
	if err != nil {
		return page, fmt.Errorf("Error Reading: %s", filename)
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
				value = strings.Trim(value, "\"") //remove quotes
				switch key {
				case "title":
					page.Title = value
				case "category":
					page.Category = value
				case "layout":
					page.Layout = value
				case "list":
					page.List = strings.ToLower(value) == "true"
				case "extension":
					page.Extension = "." + value
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
			found += 1
		}

	}

	// Change extension
	if filepath.Ext(page.OutFile) == ".md" {
		page.OutFile = page.OutFile[:len(page.OutFile)-3] + page.Extension
	}

	// next directory(s) category, category includes sub-dir = solog/webdev
	// TODO: allow category parameter
	if strings.Contains(page.OutFile, "/") {
		page.Category = page.OutFile[0:strings.LastIndex(page.OutFile, "/")]
		page.SimpleCategory = strings.Replace(page.Category, "/", "_", -1)
	}

	// parse date from filename
	base := filepath.Base(page.OutFile)
	if base[0:2] == "20" || base[0:2] == "19" { //HACK: if file starts with 20 or 19 assume date
		page.Date, _ = time.Parse("2006-01-02", base[0:10])
		page.OutFile = strings.Replace(page.OutFile, base[0:11], "", 1) // remove date from final filename
		page.List = true                                                // If we have a date prefix then this file is always listed in recents & categories
	}

	// add url of page, which includes initial slash
	// this is needed to get correct links for multi
	// level directories
	page.Url = "/" + page.OutFile

	// convert markdown content
	content := strings.Join(lines, "\n")
	if !nomarkdown {
		output := blackfriday.MarkdownCommon([]byte(content))
		page.Content = string(output)
	} else {
		page.Content = content
	}

	// add in default BaseUrl to params if not set in page
	if _, ok := page.Params["BaseUrl"]; !ok {
		page.Params["BaseUrl"] = config.BaseUrl
	}

	return page, nil
}

// Holds lists of Files, Directories and Categories
type SiteStruct struct {
	Files       []string
	Directories []string
	Categories  []string
}

// WalkFn that fills SiteStruct with data.
func (site *SiteStruct) Walker() filepath.WalkFunc {
	f := func(fn string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			site.Categories = append(site.Categories, fi.Name())
			site.Directories = append(site.Directories, fn)
			return nil
		} else {
			site.Files = append(site.Files, fn)
			return nil
		}
		return nil
	}
	return f
}

// Check if File / Directory Exists
func exists(path string) bool {
	// TODO: Check if regular file
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

// Read cfgfile or setup defaults.
func ReadConfig(basedir string, cfgfile string) (Config, error) {
	var config Config

	if file, err := ioutil.ReadFile(cfgfile); err != nil {
		return config, err
	} else {
		if err := json.Unmarshal(file, &config); err != nil {
			return config, err
		}
	}

	if basedir != "" {
		config.SourceDir = path.Join(basedir, config.SourceDir)
		config.LayoutDir = path.Join(basedir, config.LayoutDir)
		config.PublishDir = path.Join(basedir, config.PublishDir)
	}

	return config, nil
}
