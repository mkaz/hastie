package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildOutFile(t *testing.T) {
	var testData = []struct {
		sourceDir string
		filename  string
		ext       string
		expected  string
	}{
		{"pages/", "pages/index.md", ".html", "index.html"},
		{"pages/", "pages/subdir/index.md", ".html", "subdir/index.html"},
		{"pages/", "pages/subdir/2006-01-04-test.md", ".html", "subdir/test.html"},
		{"/tmp/posts", "/tmp/posts/index.md", ".html", "index.html"},
		{"/tmp/posts", "/tmp/posts/subdir/index.md", ".html", "subdir/index.html"},
		{"/tmp/posts", "/tmp/posts/subdir/2006-01-04-test.md", ".html", "subdir/test.html"},
	}
	for _, td := range testData {
		config.SourceDir = td.sourceDir // global used in function
		result := buildOutFile(td.filename, td.ext)
		assert.Equal(t, td.expected, result, "Filename should match")
	}
}

func TestGetDateFromFilename(t *testing.T) {
	var testData = []struct {
		filename string
		expected string
	}{
		{"pages/subdir/2006-01-04-test.md", "2006-01-04"},
		{"pages/subdir/test.md", "0001-01-01"},
	}

	for _, td := range testData {
		result := getDateFromFilename(td.filename)
		assert.Equal(t, td.expected, result.Format("2006-01-02"), "Date should match")
	}
}

func TestGetCategoryFromFilename(t *testing.T) {
	var testData = []struct {
		filename string
		expected string
	}{
		{"subdir/2006-01-04-test.md", "subdir"},
		{"subdir/test.md", "subdir"},
		{"subdir/math/2006-01-04-test.md", "subdir/math"},
		{"subdir/math/test.md", "subdir/math"},
	}

	for _, td := range testData {
		result := getCategoryFromFilename(td.filename)
		assert.Equal(t, td.expected, result, "Category should match")
	}
}
