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

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/russross/blackfriday"
)

// config file items
var config struct {
	SourceDir, LayoutDir, PublishDir, BaseUrl string
	CategoryMash                              map[string]string
	ProcessFilters                            map[string][]string
	UseMarkdown                               bool
}

var (
	verbose    = flag.Bool("v", false, "verbose output")
	help       = flag.Bool("h", false, "show this help")
	cfgfile    = flag.String("c", "hastie.json", "Config file")
	timing     = flag.Bool("t", false, "display timing")
	nomarkdown = flag.Bool("m", false, "do not use markdown conversion")
)

type Page struct {
	Content        string
	Title          string
	Category       string
	SimpleCategory string
	Layout         string
	OutFile        string
	Extension      string
	Url            string
	PrevUrl        string
	PrevTitle      string
	NextUrl        string
	NextTitle      string
	PrevCatUrl     string
	PrevCatTitle   string
	NextCatUrl     string
	NextCatTitle   string
	Params         map[string]string
	Recent         *PagesSlice
	Date           time.Time
	Categories     *CategoryList
	SourceFile     string
}

type PagesSlice []Page

func (p PagesSlice) Get(i int) Page         { return p[i] }
func (p PagesSlice) Len() int               { return len(p) }
func (p PagesSlice) Less(i, j int) bool     { return p[i].Date.Unix() < p[j].Date.Unix() }
func (p PagesSlice) Swap(i, j int)          { p[i], p[j] = p[j], p[i] }
func (p PagesSlice) Sort()                  { sort.Sort(p) }
func (p PagesSlice) Limit(n int) PagesSlice { return p[0:n] }
func (p PagesSlice) Reverse() PagesSlice {
	var rev PagesSlice
	for i := len(p) - 1; i >= 0; i-- {
		rev = append(rev, p[i])
	}
	return rev
}

type CategoryList map[string]PagesSlice

func (c CategoryList) Get(category string) PagesSlice { return c[category] }

var startTime time.Time
var lastTime time.Time

func init() {
	startTime = time.Now()
	lastTime = time.Now()
}

func elapsedTimer(str string) {
	if !*timing {
		return
	}
	fmt.Printf("Event: %-25s -- %9v  (%9v) \n", str, time.Since(lastTime), time.Since(startTime))
	lastTime = time.Now()
}

// holds lists of directories and files
var site = &SiteStruct{}

// Wrapper around Fprintf taking verbose flag in account.
func Printvf(format string, a ...interface{}) {
	if *verbose {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}

// Wrapper around Fprintln taking verbose flag in account.
func Printvln(a ...interface{}) {
	if *verbose {
		fmt.Fprintln(os.Stderr, a...)
	}
}

func PrintErr(str string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, str, a)
}

func usage() {
	PrintErr("usage: hastie [flags]", "")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {

	flag.Usage = usage
	flag.Parse()
	if *help {
		usage()
	}

	setupConfig()
	elapsedTimer("Config Setup")

	filepath.Walk(config.SourceDir, Walker)
	elapsedTimer("File Walker")

	/* ******************************************
	 * Loop through directories and build pages
	 * ****************************************** */
	var pages PagesSlice
	for _, dir := range site.Directories {
		pages = buildPagesSlice(dir, "/*.md", pages)
		pages = buildPagesSlice(dir, "/*.html", pages)
	}
	elapsedTimer("Loop and Parse")

	/* ******************************************
	 * Create any data needed from pages
	 * ****************************************** */

	// recent list if a sorted list of all pages
	recentList := getRecentList(pages)
	recentListPtr := &recentList

	// category list is sorted map of pages by category
	categoryList := getCategoryList(recentListPtr)
	categoryListPtr := &categoryList

	elapsedTimer("Recent and Category Lists")

	// read and parse all template files
	layoutsglob := config.LayoutDir + "/*.html"
	ts, err := template.ParseGlob(layoutsglob)
	if err != nil {
		PrintErr("Error Parsing Templates: ", err)
		os.Exit(1)
	}
	elapsedTimer("Parsed Templates")

	/* ******************************************
	 * Loop through pages and generate templates
	 * ****************************************** */
	for _, page := range pages {

		Printvf("  Generating Template: ", page.OutFile)

		// add recent pages lists to page object
		page.Recent = recentListPtr
		page.Categories = categoryListPtr

		// add prev-next links
		page.buildPrevNextLinks(recentListPtr)

		if config.BaseUrl != "" {
			page.Params["BaseUrl"] = config.BaseUrl
		}

		// Templating - writes page data to buffer
		buffer := new(bytes.Buffer)

		// pick layout based on specified in file
		templateFile := ""
		if page.Layout == "" {
			templateFile = "post.html"
		} else {
			templateFile = page.Layout + ".html"
		}

		if !exists(filepath.Join(config.LayoutDir, templateFile)) {
			PrintErr(" Missing template file:", templateFile)
			continue
		}
		ts.ExecuteTemplate(buffer, templateFile, page)

		// writing out file
		writedir := filepath.Join(config.PublishDir, page.Category)
		Printvln(" Write Directory:", writedir)
		os.MkdirAll(writedir, 0755) // does nothing if already exists

		outfile := filepath.Join(config.PublishDir, page.OutFile)
		Printvln(" Writing File:", outfile)
		ioutil.WriteFile(outfile, []byte(buffer.String()), 0644)
	}
	elapsedTimer("Generate Templates")

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
				var cmd *exec.Cmd
				// apply process filter command, capture output
				if len(filter) > 2 {
					opts := append(filter[2:], file)
					cmd = exec.Command(filter[0], opts...)
				} else {
					cmd = exec.Command(filter[0], file)
				}

				output, err := cmd.Output()
				if err != nil {
					PrintErr("Error Process Filter: "+file, err)
					continue
				}

				// determine output file path and extension
				outfile := file[strings.Index(file, string(os.PathSeparator))+1:]
				outfile = filepath.Join(config.PublishDir, outfile)
				outfile = strings.Replace(outfile, extStart, extEnd, 1)
				ioutil.WriteFile(outfile, output, 0644)
			}
		}
	}
	elapsedTimer("Process Filters")

	/* ******************************************
	 * Copy Theme Static Folder
	 * if a static directory exists in the theme, copy to publish/static
	 * TODO: process less files within theme
	 * ****************************************** */
	static_dir := config.LayoutDir + "/static"
	if exists(static_dir) {
		cmd := exec.Command("cp", "-rf", config.LayoutDir+"/static", config.PublishDir)
		cmd_err := cmd.Run()
		if cmd_err != nil {
			PrintErr("Error copying theme's static dir")
		}
	}

} // main

/* ************************************************
 * Build Recent File List
 *    - return array sorted most recent first
 *    - array includes real link (no date)
 *    - does not include files without date
 * ************************************************ */
func getRecentList(pages PagesSlice) (list PagesSlice) {
	Printvln("Creating Recent File List")
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
* Build Category List
*    - return a map containing a list of pages for
       each category, the key being category name
* ************************************************ */
func getCategoryList(pages *PagesSlice) CategoryList {
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
		// this always the category to be referenced in template
		simpleCategory := strings.Replace(page.Category, string(os.PathSeparator), "_", -1)
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
func readParseFile(filename string) (page Page) {
	Printvln("Parsing File:", filename)
	epoch, _ := time.Parse("20060102", "19700101")

	// setup default page struct
	page = Page{
		Title: "", Category: "", SimpleCategory: "", Content: "", Layout: "", Date: epoch, OutFile: filename, Extension: ".html",
		Url: "", PrevUrl: "", PrevTitle: "", NextUrl: "", NextTitle: "",
		PrevCatUrl: "", PrevCatTitle: "", NextCatUrl: "", NextCatTitle: "",
		Params: make(map[string]string),
	}

	// read file
	var data, err = ioutil.ReadFile(filename)
	if err != nil {
		PrintErr("Error Reading: " + filename)
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
			found += 1
		}

	}

	// chop off first directory, since that is the template dir
	Printvln("Filename", filename)
	page.OutFile = filename[strings.Index(filename, string(os.PathSeparator))+1:]
	page.OutFile = strings.Replace(page.OutFile, ".md", page.Extension, 1)
	Printvln("page.Outfile", page.OutFile)

	// next directory(s) category, category includes sub-dir = solog/webdev
	if page.Category == "" {
		if strings.Contains(page.OutFile, string(os.PathSeparator)) {
			page.Category = page.OutFile[0:strings.LastIndex(page.OutFile, string(os.PathSeparator))]
			page.SimpleCategory = strings.Replace(page.Category, string(os.PathSeparator), "_", -1)
		}
	}
	Printvln("page.Category", page.Category)
	// parse date from filename
	base := filepath.Base(page.OutFile)
	if base[0:2] == "20" || base[0:2] == "19" { //HACK: if file starts with 20 or 19 assume date
		page.Date, _ = time.Parse("2006-01-02", base[0:10])
		page.OutFile = strings.Replace(page.OutFile, base[0:11], "", 1) // remove date from final filename
	}

	// add url of page, which includes initial slash
	// this is needed to get correct links for multi
	// level directories
	page.Url = "/" + page.OutFile

	// convert markdown content
	content := strings.Join(lines, "\n")
	if (config.UseMarkdown) && (page.Params["markdown"] != "no") {
		output := markdownRender([]byte(content))
		page.Content = string(output)
	} else {
		page.Content = content
	}

	return page
}

func markdownRender(content []byte) []byte {
	htmlFlags := 0
	//htmlFlags |= blackfriday.HTML_SKIP_SCRIPT
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS

	return blackfriday.Markdown(content, renderer, extensions)
}

func buildPagesSlice(dir string, globstr string, pages PagesSlice) PagesSlice {
	readglob := dir + globstr
	var dirfiles, _ = filepath.Glob(readglob)

	// loop through files in directory
	for _, file := range dirfiles {
		Printvln("  File:", file)
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

// Holds lists of Files, Directories and Categories
type SiteStruct struct {
	Files       []string
	Directories []string
	Categories  []string
}

// WalkFn that fills SiteStruct with data.
func Walker(fn string, fi os.FileInfo, err error) error {
	if err != nil {
		PrintErr("Walker: ", err)
		return nil
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

// Check if File / Directory Exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

// Read cfgfile or setup defaults.
func setupConfig() {
	file, err := ioutil.ReadFile(*cfgfile)
	if err != nil {
		// set defaults, no config file
		config.SourceDir = "_source"
		config.LayoutDir = "_layout"
		config.PublishDir = "public"
		config.UseMarkdown = true
	} else {
		// not required in config file, set default
		config.UseMarkdown = true
		if err := json.Unmarshal(file, &config); err != nil {
			fmt.Printf("Error parsing config: %s", err)
			os.Exit(1)
		}
	}

	// // command line overrides config file
	if *nomarkdown {
		config.UseMarkdown = false
	}

	Printvln("SourceDir", config.SourceDir)
	Printvln("LayoutDir", config.LayoutDir)
	Printvln("PublishDir", config.PublishDir)
}
