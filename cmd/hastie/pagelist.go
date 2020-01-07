package main

import "sort"

// PageList is an array of Page objects
type PageList []Page

// By is the type of less function to use for sorting
type By func(p1, p2 *Page) bool

func (by By) Sort(pages []Page) {
	ps := &pageListSorter{
		pages: pages,
		by:    by,
	}
	sort.Sort(ps)
}

// pageListSorter joins a By function and PageList to be sorted.
// See Example (SortKeys): https://golang.org/pkg/sort/#Sort
type pageListSorter struct {
	pages []Page
	by    func(p1, p2 *Page) bool
}

// Len is part of sort.Interface.
func (s *pageListSorter) Len() int {
	return len(s.pages)
}

// Swap is part of sort.Interface.
func (s *pageListSorter) Swap(i, j int) {
	s.pages[i], s.pages[j] = s.pages[j], s.pages[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *pageListSorter) Less(i, j int) bool {
	return s.by(&s.pages[i], &s.pages[j])
}

func (p PageList) Limit(n int) PageList {
	if len(p) > n {
		return p[0:n]
	}
	return p
}

func (p PageList) Get(i int) Page {
	return p[i]
}

func (p PageList) Len() int {
	return len(p)
}

func (p PageList) Reverse() PageList {
	var rev PageList
	for i := len(p) - 1; i >= 0; i-- {
		rev = append(rev, p[i])
	}
	return rev
}

func orderOrder(p1, p2 *Page) bool {
	return p1.Order < p2.Order
}

func dateOrder(p1, p2 *Page) bool {
	return p1.Date.Before(p2.Date)
}

// CategoryList is a map of category to PageList
type CategoryList map[string]PageList

// Get returns the PageList for given category
func (c CategoryList) Get(category string) PageList { return c[category] }
