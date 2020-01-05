package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupConfigDefaults(t *testing.T) {
	config := parseConfig("")
	assert.Equal(t, config.SourceDir, "_source", "Is source directory set?")
	assert.Equal(t, config.LayoutDir, "_layout", "Is layout directory set?")
	assert.Equal(t, config.PublishDir, "public", "Is publish directory set?")
	assert.Equal(t, config.UseMarkdown, true, "Is default use markdown?")
}

func TestSetupConfigJson(t *testing.T) {
	jsonConfig := `{"SourceDir":"pages","LayoutDir":"../themes/docs/","PublishDir":"../docs/"}`
	config := parseConfig(jsonConfig)
	assert.Equal(t, config.SourceDir, "pages", "Is source directory set?")
	assert.Equal(t, config.LayoutDir, "../themes/docs/", "Is layout directory set?")
	assert.Equal(t, config.UseMarkdown, true, "Is default use markdown?")
}

func TestSetupConfigParams(t *testing.T) {
	jsonConfig := `{"SourceDir":"pages","Params":{"SiteName":"Hastie"}}`
	config := parseConfig(jsonConfig)
	assert.Equal(t, config.SourceDir, "pages", "Is source directory set?")
	assert.Equal(t, config.Params["SiteName"], "Hastie", "Is custom parameter set?")
}
