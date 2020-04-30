package vcount

import "sync"

type user = string
type page = string

type pageVisits map[user]int

// VisitCounter counts each user visit per page.
type VisitCounter struct {
	sync.RWMutex
	visits map[page]pageVisits
}

// New returns a new instance of VisitCounter.
func New() *VisitCounter {
	return &VisitCounter{}
}

// Inc registers a page visit by the user.
func (c *VisitCounter) Inc(p page, u user) {
	c.Lock()
	defer c.Unlock()

	if c.visits == nil {
		c.visits = make(map[page]pageVisits)
	}
	if c.visits[p] == nil {
		c.visits[p] = make(pageVisits)
	}
	c.visits[p][u]++
}

// PageUserVisits returns the number of page visits by the user.
func (c *VisitCounter) PageUserVisits(p page, u user) int {
	c.RLock()
	defer c.RUnlock()

	return c.visits[p][u]
}

// PageTotalVisits returns the total number of page visits.
func (c *VisitCounter) PageTotalVisits(p page) int {
	c.RLock()
	defer c.RUnlock()

	total := 0
	for _, v := range c.visits[p] {
		total += v
	}
	return total
}

// UserTotalVisits returns the total number of visits by the user.
func (c *VisitCounter) UserTotalVisits(u user) int {
	c.RLock()
	defer c.RUnlock()

	total := 0
	for _, pv := range c.visits {
		total += pv[u]
	}
	return total
}

// TotalVisits returns the total number of visits.
func (c *VisitCounter) TotalVisits() int {
	c.RLock()
	defer c.RUnlock()

	total := 0
	for _, p := range c.visits {
		for _, v := range p {
			total += v
		}
	}
	return total
}

// This goes out of the scope of a counter but is interesting to implement.
// Decided to stop for now.

// UserPages returns a slice of pages visited by user.
// func (c *VisitCounter) UserPages(u user) []page {
// 	pages := []page{}
// 	for p, pv := range c.visits {
// 		if _, ok := pv[u]; ok {
// 			pages = append(pages, p)
// 		}
// 	}
// 	return pages
// }
