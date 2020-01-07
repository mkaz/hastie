package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var pages = PageList{
	{
		Title: "Apple",
		Date:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		Order: 1,
	},
	{
		Title: "Banana",
		Date:  time.Date(2009, time.November, 10, 01, 0, 0, 0, time.UTC),
		Order: 3,
	},
	{
		Title: "Carrot",
		Date:  time.Date(2009, time.November, 10, 31, 0, 0, 0, time.UTC),
		Order: 2,
	},
}

func TestPagelist(t *testing.T) {
	By(orderOrder).Sort(pages)
	assert.Equal(t, pages[0].Order, 1, "In order order")
	assert.Equal(t, pages[1].Order, 2, "In order order")
	assert.Equal(t, pages[2].Order, 3, "In order order")

	By(dateOrder).Sort(pages)
	assert.Equal(t, pages[0].Title, "Banana", "In date order")
	assert.Equal(t, pages[1].Title, "Apple", "In date order")
	assert.Equal(t, pages[2].Title, "Carrot", "In date order")

	pages = pages.Reverse()
	assert.Equal(t, pages[0].Title, "Carrot", "In reverse date order")
	assert.Equal(t, pages[1].Title, "Apple", "In reverse date order")
	assert.Equal(t, pages[2].Title, "Banana", "In reverse date order")
}
