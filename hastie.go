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
	"flag"
	"fmt"
	/* switched to dhconnelly fork which works with Go1
   "github.com/russross/blackfriday" */
  "github.com/dhconnelly/blackfriday"
	"io/ioutil"
	"os"
  "os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

const (
	cfgFiledefault = "hastie.json"
)

// config file items
var config struct {
  SourceDir, LayoutDir, PublishDir string
  CategoryMash map[string]string
  ProcessFilters map[string][]string
}

var (
	verbose = flag.Bool("v", false, "verbose output")
	help    = flag.Bool("h", false, "show this help")
	cfgfile = flag.String("c", cfgFiledefault, "Config file")
)

type Page struct {
	Content, Title, Category, Layout, OutFile, Url, PrevUrl, PrevTitle, NextUrl, NextTitle, PrevCatUrl, PrevCatTitle, NextCatUrl, NextCatTitle  string
  Params        map[string]string
	Recent        PagesSlice
	Date          time.Time
  Categories    CategoryList
}

type PagesSlice []Page
func (p PagesSlice) Len() int               { return len(p) }
func (p PagesSlice) Less(i, j int) bool     { return p[i].Date.Unix() < p[j].Date.Unix() }
func (p PagesSlice) Swap(i, j int)          { p[i], p[j] = p[j], p[i] }
func (p PagesSlice) Sort()                  { sort.Sort(p) }
func (p PagesSlice) Limit(n int) PagesSlice { return p[0:n] }

type CategoryList map[string]PagesSlice
func (c CategoryList) Get(category string) PagesSlice   { return c[category] }


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

func usage() {
	fmt.Fprintln(os.Stderr, "usage: hastie [flags]")
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

	filepath.Walk(config.SourceDir, Walker)

	/* ******************************************
	 * Loop through directories and build pages 
	 * ****************************************** */
	var pages PagesSlice
	for _, dir := range site.Directories {

		readglob := fmt.Sprintf("%s/*.md", dir)
		var dirfiles, _ = filepath.Glob(readglob)

		// loop through files in directory
		for _, file := range dirfiles {
			Printvln("  File:", file)
			outfile := filepath.Base(file)
			outfile = strings.Replace(outfile, ".md", ".html", 1)

			// read & parse file for parameters
			page := readParseFile(file)

			pages = append(pages, page)
		}
	}

	if *verbose { // spit out pages structure
		Printvln("################################################################################")
		Printvf(" %-50s | %-10s | %s \n", "Title", "Category", "Outfile")
		Printvln("--------------------------------------------------------------------------------")
		for _, page := range pages {
			Printvf(" %-50s | %-10s | %s \n", page.Title, page.Category, page.OutFile)
		}
		Printvln("################################################################################")
	}

	/* ******************************************
	 * Create any data needed from pages
	 * for example recent file list
	 * category list, etc...
	 * ****************************************** */

	// build recent file list, sorted
	recentList := getRecentList(pages)
  categoryList := getCategoryList(recentList)


	/* ******************************************
	 * Loop through pages and generate templates
	 * ****************************************** */
	for _, page := range pages {

		Printvf("  Generating Template: ", page.OutFile)

    // added recent pages lists to each page object
		page.Recent = recentList
    page.Categories = categoryList

    // add prev-next links
    page.buildPrevNextLinks(recentList)
    page.buildCatPrevNextLinks(recentList)

		/* Templating - writes page data to buffer 
		 * read and parse all template files          */
		buffer := new(bytes.Buffer)
		layoutsglob := fmt.Sprintf("%s/*.html", config.LayoutDir)
		ts, err := template.ParseGlob(layoutsglob)
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
		ts.ExecuteTemplate(buffer, templateFile, page)

		// writing out file
		writedir := fmt.Sprintf("%s/%s", config.PublishDir, page.Category)
		Printvln(" Write Directory:", writedir)
		os.MkdirAll(writedir, 0755) // does nothing if already exists

		outfile := fmt.Sprintf("%s/%s", config.PublishDir, page.OutFile)
		Printvln(" Writing File:", outfile)
		ioutil.WriteFile(outfile, []byte(buffer.String()), 0644)
	}


	/* ******************************************
	 * Loop through processFilters and process
	 * ****************************************** */
  for ext,filter := range config.ProcessFilters {
    extStart := fmt.Sprintf(".%s", ext)
    extEnd := fmt.Sprintf(".%s", filter[1])

    for _, dir := range site.Directories {
      readglob := fmt.Sprintf("%s/*%s", dir, extStart)
      var dirfiles, _ = filepath.Glob(readglob)
      for _, file := range dirfiles {
        Printvf("Found file to Process: %s in Directory: %s \n", file, dir)
        // check for filter
        //apply process filter command, capture output
        cmd := exec.Command(filter[0], file)
        output, _ := cmd.Output()

        // determine output file path and extension
        outfile := file[strings.Index(file, "/")+1:]
        outfile = fmt.Sprintf("%s/%s", config.PublishDir, outfile)
        outfile = strings.Replace(outfile, extStart, extEnd, 1)
        Printvf("Writing File: %s ", outfile)

        ioutil.WriteFile(outfile, output, 0644)

      }
    }
  }


} // main


/* ************************************************
 * Read and Parse File
 * @param filename
 * @return Page object
 * ************************************************ */
func readParseFile(filename string) (page Page) {
	Printvln("Parsing File:", filename)
	epoch, _ := time.Parse("20060102", "19700101")

	// setup default page struct
	page = Page{
		Title:     "",
		Category:  "",
		Content:   "",
		Layout:    "",
		Date:      epoch,
		OutFile:   filename,
    Url:       "",
    Params:   make(map[string]string),
    PrevUrl:   "",
    PrevTitle: "",
    NextUrl:   "",
    NextTitle: "",
    PrevCatUrl:   "",
    PrevCatTitle: "",
    NextCatUrl:   "",
    NextCatTitle: ""}

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
        value = strings.Trim(value, "\"")  //remove quotes
				switch key {
				case "title":
					page.Title = value
				case "category":
					page.Category = value
				case "layout":
					page.Layout = value
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

	// switch directory name to just category

	// chop off first directory, since that is the template dir
	page.OutFile = filename[strings.Index(filename, "/")+1:]
	page.OutFile = strings.Replace(page.OutFile, ".md", ".html", 1)

	// next directory(s) category, category includes sub-dir = solog/webdev
  // TODO: allow category parameter
	if strings.Contains(page.OutFile, "/") {
		page.Category = page.OutFile[0:strings.LastIndex(page.OutFile, "/")]
	}

	// parse date from filename
	base := filepath.Base(page.OutFile)
	if base[0:2] == "20" || base[0:2] == "19" { //HACK: if file starts with 20 or 19 assume date
		page.Date, _ = time.Parse("2006-01-02", base[0:10])
		page.OutFile = strings.Replace(page.OutFile, base[0:11], "", 1) // remove date from final filename
	}

  // add url of page, which includes initial slash
  // this is needed to get correct links for multi 
  // level directories 
  page.Url = fmt.Sprintf("/%s", page.OutFile)

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
	Printvf("Creating Recent File List")
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
func getCategoryList(pages PagesSlice) CategoryList {
  mapList := make(CategoryList)
  // recentList is passed in which is already sorted
  // just need to map the pages to category

  // read category mash config, which allows to create
  // a new category based on combining multiple categories
  // this is used on my site when I want to display a list 
  // of items from similar categories together
  reverseMap := make(map[string]string)

  for k,v := range config.CategoryMash {
    //split v on comma
    cats := strings.Split(string(v), ",")
    //loop through split and add to reverse map
    for _,cat := range cats {
      reverseMap[cat] = string(k)
    }
  }

  for _, page := range pages {

    // create new category from category mash map
    if reverseMap[page.Category] != page.Category {
      thisCategory := reverseMap[page.Category]
      mapList[thisCategory] = append(mapList[thisCategory], page) 
    }
    // still want a list of regular categories
    mapList[page.Category] = append(mapList[page.Category],page)
  }
	return mapList
}



/* ************************************************
 * Add Prev Next Links to Page Object 
 * ************************************************ */
func (page *Page) buildPrevNextLinks(recentList PagesSlice) {
    foundIt := false
    nextPage := Page{}
    prevPage := Page{}
    pp := Page{}
    for _, rp := range recentList {

      if foundIt {
        prevPage = rp
        break
      }

      if (rp.Title == page.Title) {
        nextPage = pp
        foundIt = true
      }
      pp = rp   // previous page
    }
    page.NextUrl = nextPage.Url
    page.NextTitle = nextPage.Title
    page.PrevUrl = prevPage.Url
    page.PrevTitle = prevPage.Title
}


/* ************************************************
 * Add Prev Next Links by Category to Page Object
 * ************************************************ */
func (page *Page) buildCatPrevNextLinks(recentList PagesSlice) {
   foundIt := false
   nextPage := Page{}
   prevPage := Page{}
   pp := Page{}

   for _, rp := range recentList {
     if rp.Category == page.Category {
       if foundIt {
         prevPage = rp
         break
       }

       if (rp.Title == page.Title) {
         nextPage = pp
         foundIt = true
       }
       pp = rp   // previous page
     }
   }
   page.NextCatUrl = nextPage.Url
   page.NextCatTitle = nextPage.Title
   page.PrevCatUrl = prevPage.Url
   page.PrevCatTitle = prevPage.Title
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
		fmt.Println("Walker:", err)
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

/* ************************************************
 * Check if File / Directory Exists
 * ************************************************ */
func exists(path string) bool {
	// TODO: Check if regular file
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
		// set defaults
		config.SourceDir = "posts"
		config.LayoutDir = "layouts"
		config.PublishDir = "public"
	} else {
		if err := json.Unmarshal(file, &config); err != nil {
			fmt.Printf("Error parsing config: %s", err)
			os.Exit(1)
		}
	}
}
