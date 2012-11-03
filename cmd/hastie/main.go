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
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mkaz/hastie"
	"io/ioutil"
	"os"
	"time"
)

const (
	cfgFiledefault = "hastie.json"
)

var (
	verbose    = flag.Bool("v", false, "verbose output")
	help       = flag.Bool("h", false, "show this help")
	cfgfile    = flag.String("c", cfgFiledefault, "Config file")
	timing     = flag.Bool("t", false, "display timing")
	nomarkdown = flag.Bool("m", false, "do not use markdown conversion")
	config     hastie.Config
)

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

type monitor struct{}

func (monitor) Walked() {
	elapsedTimer("File Walker")
}

func (monitor) ParsingSource(file string) {
	Printvln("  File:", file)
}

func (monitor) ParsedSources() {
	elapsedTimer("Loop and Parse")
}

func (monitor) Listed() {
	elapsedTimer("Recent and Category Lists")
}

func (monitor) ParsedTemplates() {
	elapsedTimer("Parsed Templates")
}

func (monitor) ParsingTemplate(file string) {
	Printvln("  Generating Template:", file)
}

func (monitor) WritingTemplate(file string) {
	Printvln("  Writing File:", file)
}

func (monitor) GeneratedTemplates() {
	elapsedTimer("Generate Templates")
}

func (monitor) Filtered() {
	elapsedTimer("Process Filters")
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

	if err := hastie.Compile(config, monitor{}); err != nil {
		PrintErr(err.Error())
	}
}

// Read cfgfile or setup defaults.
func setupConfig() {
	file, err := ioutil.ReadFile(*cfgfile)
	if err != nil {
		// set defaults
		config.SourceDir = "posts"
		config.LayoutDir = "layouts"
		config.PublishDir = "public"
		config.NoMarkdown = false
	} else {
		if err := json.Unmarshal(file, &config); err != nil {
			fmt.Printf("Error parsing config: %s", err)
			os.Exit(1)
		}
	}
}
