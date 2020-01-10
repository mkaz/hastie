/**
 * Hastie - Static Site Generator
 * https://github.com/mkaz/hastie
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/mkaz/hastie/pkg/logger"
	"github.com/mkaz/hastie/pkg/utils"
)

var log logger.Logger
var config Config

// Page main page object
type Page struct {
	Title          string
	Content        string
	Date           time.Time
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
	Order          int
	Params         map[string]string
	AllPages       *PageList
	Recent         *PageList
	Categories     *CategoryList
	Unlisted       bool
	SourceFile     string
}

func main() {
	var helpFlag = flag.Bool("help", false, "show this help")
	var versionFlag = flag.Bool("version", false, "Display version and quit")
	var noMarkdown = flag.Bool("nomarkdown", false, "do not use markdown conversion")
	var configFile = flag.String("config", "hastie.json", "Config file")
	flag.BoolVar(&log.DebugLevel, "debug", false, "Debug output (verbose)")
	flag.BoolVar(&log.Verbose, "verbose", false, "Show info level")
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Println("hastie v0.9.1")
		os.Exit(0)
	}

	config = setupConfig(*configFile)
	if *noMarkdown {
		config.UseMarkdown = false
	}

	// get pages and directories
	pages, dirs := getSiteFiles(config.SourceDir)

	allPages := filterUnlisted(pages)
	By(orderOrder).Sort(allPages)
	allPagesPointer := &allPages

	// recent list is a sorted list of all pages with dates
	By(dateOrder).Sort(pages)
	recentList := pages.Reverse()
	recentListPointer := &recentList

	// category list is sorted map of pages by category
	categoryList := getCategoryList(recentListPointer)
	categoryListPointer := &categoryList

	// functions made available to templates
	funcMap := template.FuncMap{
		"trim":    utils.TrimSlash,
		"Title":   strings.Title,
		"ToLower": strings.ToLower,
		"ToUpper": strings.ToUpper,
	}

	// read and parse all template files
	layoutsglob := config.LayoutDir + "/*.html"
	ts, err := template.New("master").Funcs(funcMap).ParseGlob(layoutsglob)
	if err != nil {
		log.Fatal("Error Parsing Templates: ", err)
	}

	// loop through each page
	// add extra data to page to be available to template
	// apply templates and write out generated files
	for _, page := range pages {

		// add recent pages lists to page object
		page.AllPages = allPagesPointer
		page.Recent = recentListPointer
		page.Categories = categoryListPointer

		// add prev-next links
		page.buildPrevNextLinks(recentListPointer)

		page.Params["BaseURL"] = config.BaseURL

		// applyTemplate to page
		buffer, err := applyTemplate(ts, page)
		if err != nil {
			log.Warn("Error applying template", err)
			continue
		}

		// confirm directory exists
		writedir := filepath.Join(config.PublishDir, page.Category)
		log.Debug(" Writing Directory:", writedir)
		os.MkdirAll(writedir, 0755) // does nothing if already exists

		// write out file
		outfile := filepath.Join(config.PublishDir, page.OutFile)
		log.Debug(" Writing File:", outfile)
		ioutil.WriteFile(outfile, []byte(buffer.String()), 0644)
	}

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

		for _, dir := range dirs {
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
					log.Warn("Error Process Filter: "+file, err)
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

	/* ******************************************
	 * Copy Theme Static Folder
	 * if a static directory exists in the theme, copy to publish/static
	 * TODO: process less files within theme
	 * ****************************************** */
	staticDir := config.LayoutDir + "/static"
	if utils.FileExists(staticDir) {
		cmd := exec.Command("cp", "-rf", config.LayoutDir+"/static", config.PublishDir)
		cmdErr := cmd.Run()
		if cmdErr != nil {
			log.Warn("Error copying theme's static dir")
		}
	}

} // main

/* ************************************************
* Build Category List
*    - return a map containing a list of pages for
       each category, the key being category name
* ************************************************ */
func getCategoryList(pages *PageList) CategoryList {
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
func (page *Page) buildPrevNextLinks(recentList *PageList) {
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

func applyTemplate(ts *template.Template, page Page) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)

	// pick layout based on specified in file
	templateFile := ""
	if page.Layout == "" {
		templateFile = "post.html"
	} else {
		templateFile = page.Layout + ".html"
	}

	if !utils.FileExists(filepath.Join(config.LayoutDir, templateFile)) {
		return nil, fmt.Errorf("Missing template file %s", templateFile)
	}
	ts.ExecuteTemplate(buffer, templateFile, page)
	return buffer, nil
}

func filterUnlisted(pages PageList) (filtered PageList) {
	for _, page := range pages {
		if !page.Unlisted {
			filtered = append(filtered, page)
		}
	}
	return filtered
}
