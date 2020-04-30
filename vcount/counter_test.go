package vcount

import (
	"testing"
)

func TestVisitCounter(t *testing.T) {
	type pageUser struct {
		page page
		user user
	}
	type incCalls []pageUser
	type pageUserVisits map[pageUser]int
	type pageTotalVisits map[page]int
	type userTotalVisits map[user]int

	page1 := "test_page_1"
	page2 := "test_page_2"
	user1 := "test_user_1"
	user2 := "test_user_2"

	tests := []struct {
		name                string
		incCalls            incCalls
		wantPageUserVisits  pageUserVisits
		wantPageTotalVisits pageTotalVisits
		wantUserTotalVisits userTotalVisits
		wantTotalVisits     int
	}{
		{	
			"Empty counter", 
			incCalls{}, 
			pageUserVisits{{page1, user1}: 0},
			pageTotalVisits{page1: 0},
			userTotalVisits{user1: 0}, 
			0,
		},
		{
			"One page, one user: visits once",
			incCalls{{page1, user1}},
			pageUserVisits{{page1, user1}: 1},
			pageTotalVisits{page1: 1},
			userTotalVisits{user1: 1},
			1,
		},
		{
			"One page, one user: visits twice",
			incCalls{{page1, user1}, {page1, user1}},
			pageUserVisits{{page1, user1}: 2},
			pageTotalVisits{page1: 2},
			userTotalVisits{user1: 2},
			2,
		},
		{
			"Two pages, one user: visits each",
			incCalls{{page1, user1}, {page2, user1}},
			pageUserVisits{{page1, user1}: 1, {page2, user1}: 1},
			pageTotalVisits{page1: 1, page2: 1},
			userTotalVisits{user1: 2},
			2,
		},
		{
			"Two pages, two users: one unique user visit per page",
			incCalls{{page1, user1}, {page2, user2}},
			pageUserVisits{{page1, user1}: 1, {page1, user2}: 0, {page2, user1}: 0, {page2, user2}: 1},
			pageTotalVisits{page1: 1, page2: 1},
			userTotalVisits{user1: 1, user2: 1},
			2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			counter := &VisitCounter{}

			for _, args := range test.incCalls {
				counter.Inc(args.page, args.user)
			}

			for args, want := range test.wantPageUserVisits {
				got := counter.PageUserVisits(args.page, args.user)
				if got != want {
					t.Errorf("PageUserVisits: got %d, want %d", got, want)
				}
			}

			for page, want := range test.wantPageTotalVisits {
				got := counter.PageTotalVisits(page)
				if got != want {
					t.Errorf("PageTotalVisits: got %d, want %d", got, want)
				}
			}

			for user, want := range test.wantUserTotalVisits {
				got := counter.UserTotalVisits(user)
				if got != want {
					t.Errorf("UserTotalVisits: got %d, want %d", got, want)
				}
			}

			gotTotalVisits := counter.TotalVisits()
			if gotTotalVisits != test.wantTotalVisits {
				t.Errorf("TotalVisits: got %d, want %d", gotTotalVisits, test.wantTotalVisits)
			}
		})
	}
}

// Also tried a non table approach. Left it here for the history.

// func TestVisitCounter(t *testing.T) {
// 	page1 := "test_page_1"
// 	page2 := "test_page_2"
// 	user1 := "test_user_1"
// 	user2 := "test_user_2"

// 	t.Run("Empty counter", func(t *testing.T) {
// 		counter := &VisitCounter{}

// 		gotUserVisits := counter.PageUserVisits(page1, user1)
// 		wantUserVisits := 0
// 		if gotUserVisits != wantUserVisits {
// 			t.Errorf("PageUserVisits: got: %d want: %d", gotUserVisits, wantUserVisits)
// 		}

// 		gotPageVisits := counter.PageTotalVisits(page1)
// 		wantPageVisits := 0
// 		if gotPageVisits != wantPageVisits {
// 			t.Errorf("PageTotalVisits: got: %d want: %d", gotUserVisits, wantUserVisits)
// 		}
// 	})

// 	t.Run("One page: one user visits once", func(t *testing.T) {
// 		counter := &VisitCounter{}

// 		counter.Inc(page1, user1)

// 		gotUserVisits := counter.PageUserVisits(page1, user1)
// 		wantUserVisits := 1
// 		if gotUserVisits != wantUserVisits {
// 			t.Errorf("PageUserVisits: got: %d want: %d", gotUserVisits, wantUserVisits)
// 		}

// 		gotPageVisits := counter.PageTotalVisits(page1)
// 		wantPageVisits := 1
// 		if gotPageVisits != wantPageVisits {
// 			t.Errorf("PageTotalVisits: got: %d want: %d", gotUserVisits, wantUserVisits)
// 		}
// 	})

// 	t.Run("One page: one user visits twice", func(t *testing.T) {
// 		counter := &VisitCounter{}

// 		counter.Inc(page1, user1)
// 		counter.Inc(page1, user1)

// 		gotUserVisits := counter.PageUserVisits(page1, user1)
// 		wantUserVisits := 2
// 		if gotUserVisits != wantUserVisits {
// 			t.Errorf("PageUserVisits: got: %d want: %d", gotUserVisits, wantUserVisits)
// 		}

// 		gotPageVisits := counter.PageTotalVisits(page1)
// 		wantPageVisits := 2
// 		if gotPageVisits != wantPageVisits {
// 			t.Errorf("PageTotalVisits: got: %d want: %d", gotUserVisits, wantUserVisits)
// 		}
// 	})

// 	t.Run("Two pages: one user visits each", func(t *testing.T) {
// 		counter := &VisitCounter{}

// 		counter.Inc(page1, user1)
// 		counter.Inc(page2, user1)

// 		gotUserVisits := counter.PageUserVisits(page1, user1)
// 		wantUserVisits := 2
// 		if gotUserVisits != wantUserVisits {
// 			t.Errorf("PageUserVisits: got: %d want: %d", gotUserVisits, wantUserVisits)
// 		}

// 		gotPageVisits := counter.PageTotalVisits(page1)
// 		wantPageVisits := 2
// 		if gotPageVisits != wantPageVisits {
// 			t.Errorf("PageTotalVisits: got: %d want: %d", gotUserVisits, wantUserVisits)
// 		}
// 	})

// 	t.Run("Two pages: two users visit each once", func(t *testing.T) {
// 		counter := &VisitCounter{}

// 		counter.Inc(page1, user1)
// 		counter.Inc(page2, user2)

// 		gotUserVisits := counter.PageUserVisits(page1, user1)
// 		wantUserVisits := 2
// 		if gotUserVisits != wantUserVisits {
// 			t.Errorf("PageUserVisits: got: %d want: %d", gotUserVisits, wantUserVisits)
// 		}

// 		gotPageVisits := counter.PageTotalVisits(page1)
// 		wantPageVisits := 2
// 		if gotPageVisits != wantPageVisits {
// 			t.Errorf("PageTotalVisits: got: %d want: %d", gotUserVisits, wantUserVisits)
// 		}
// 	})
// }
