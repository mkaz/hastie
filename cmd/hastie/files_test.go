package main

import "testing"

func TestBuildOutFile(t *testing.T) {
	var testData = []struct {
		sourceDir string
		filename  string
		ext       string
		expected  string
	}{
		{"pages", "pages/index.md", ".html", "index.html"},
		{"pages", "pages/subdir/index.md", ".html", "subdir/index.html"},
		{"pages", "pages/subdir/2006-01-04-test.md", ".html", "subdir/test.html"},
		{"/tmp/posts", "/tmp/posts/index.md", ".html", "index.html"},
		{"/tmp/posts", "/tmp/posts/subdir/index.md", ".html", "subdir/index.html"},
		{"/tmp/posts", "/tmp/posts/subdir/2006-01-04-test.md", ".html", "subdir/test.html"},
	}
	for _, td := range testData {
		config.SourceDir = td.sourceDir // global used in function
		result := buildOutFile(td.filename, td.ext)
		if td.expected != result {
			t.Errorf("Expected %s but received %s", td.expected, result)
		}
	}
}
