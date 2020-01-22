package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/mkaz/hastie/pkg/utils"
)

// global config
type Config struct {
	SourceDir, LayoutDir, PublishDir, BaseURL string
	CategoryMash                              map[string]string
	ProcessFilters                            map[string][]string
	UseMarkdown                               bool
	UsePrism                                  bool
	UseAsciinema                              bool
	Params                                    map[string]string
}

// Read cfgfile or setup defaults.
func setupConfig(filename string) Config {
	var jsonConfig string

	// check if config file exists
	if utils.FileExists(filename) {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Warn("Error reading config file, using defaults", err)
		}
		jsonConfig = string(data)
	} else {
		log.Warn("No config file, using defaults")
	}

	return parseConfig(jsonConfig)
}

func parseConfig(jsonConfig string) (config Config) {
	// set default values
	config.SourceDir = "_source"
	config.LayoutDir = "_layout"
	config.PublishDir = "public"
	config.UseMarkdown = true
	config.UsePrism = true
	config.UseAsciinema = false
	config.Params = make(map[string]string)

	// if no config file, than data is empty
	if jsonConfig == "" {
		return
	}

	err := json.Unmarshal([]byte(jsonConfig), &config)
	if err != nil {
		log.Fatal("Error parsing config: %s", err)
	}
	return config
}
