package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseContent(t *testing.T) {

	content := `---
title: Hola Mundo
date: 2012-05-15
category: foo
dummy: bar
---
Here is my content`

	// setup default page struct
	page := Page{
		Params: make(map[string]string),
	}
	page = parseContent(content, page)

	assert.Equal(t, page.Title, "Hola Mundo", "Title should match")
	assert.Equal(t, page.Date.Format("2006-01-02"), "2012-05-15", "Date should match")
	assert.Equal(t, page.Category, "foo", "Category should match")
	assert.Equal(t, page.Params["dummy"], "bar", "User param should match")

}
