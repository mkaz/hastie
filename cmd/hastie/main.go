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
	"flag"
	"fmt"
	"github.com/mkaz/hastie"
	hhttp "github.com/mkaz/hastie/http"
	"net/http"
	"os"
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
	httpAddr   = flag.String("http", "", "HTTP service address (e.g., ':8080')")
	config     hastie.Config
)

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

	var config hastie.Config
	var err error

	if _, err = os.Stat(*cfgfile); err != nil {
		config = hastie.DefaultConfig
	} else {
		if config, err = hastie.ReadConfig("", *cfgfile); err != nil {
			PrintErr("Error parsing config: ", err.Error())
			os.Exit(1)
		}
	}

	monitor := hastie.NewLogMonitor(func(msg string) {
		fmt.Fprintln(os.Stderr, msg)
	}, *timing, *verbose)

	if *httpAddr == "" {
		// Perform simple compile
		if err := config.Compile(monitor); err != nil {
			PrintErr("Error compiling config: ", err.Error())
		}
	} else {
		// Start http handler
		handler := hhttp.Handle(config, monitor)

		// Watch for changes in SourceDir & LayoutDir
		go func() {
			if err := handler.Watch(); err != nil {
				PrintErr("Error watching for changes: ", err.Error())
			}
		}()

		// Start http server
		if err := http.ListenAndServe(*httpAddr, handler); err != nil {
			PrintErr("Error serving: ", err.Error())
		}
	}
}
