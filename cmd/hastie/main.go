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
	config     hastie.Config
)

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

	monitor := hastie.NewLogMonitor(Printvln, *timing)
	if err := config.Compile(monitor); err != nil {
		PrintErr("Error compiling config: ", err.Error())
	}
}
