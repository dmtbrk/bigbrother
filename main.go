package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/ortymid/bigbrother/vcount"
)

var templates = template.Must(template.ParseFiles("./templates/index.html"))

func logVisits(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("visit: %s - %s", r.RemoteAddr, r.URL.Path)

		h(w, r)
	}
}

var counter = vcount.New()

func countVisits(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := getPageID(r)
		// filtering out favicon requests. may grow in the future as more 
		// additional requests will come in, but not a problem for this demo
		if page != "/favicon.ico" { 
			user := getUserID(r)
			counter.Inc(page, user)
		}

		h(w, r)
	}
}

// getPageID returns the unique identifier of the viewed page.
// It is a path to the page in the current implementation.
func getPageID(r *http.Request) string {
	// removing the trailing slash to count 'page/' and 'page' as the same page
	return strings.TrimSuffix(r.URL.Path, "/")
}

// getUserID returns the unique identifier of the user by the given request.
// It is a remote addrres in the current implementation.
func getUserID(r *http.Request) string {
	return r.RemoteAddr
}

func allHandler(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		Home            bool
		Title           string
		PageUserVisits  int
		PageTotalVisits int
		TotalVisits     int
	}

	page := getPageID(r)
	user := getUserID(r)

	title := getPageTitle(r)
	pd := &PageData{
		Title:           title,
		PageUserVisits:  counter.PageUserVisits(page, user),
		PageTotalVisits: counter.PageTotalVisits(page),
		TotalVisits:     counter.TotalVisits(),
	}

	err := templates.ExecuteTemplate(w, "index.html", pd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getTitle returns the title for the given page.
func getPageTitle(r *http.Request) string {
	// have to use RawPath to distinguish '%2F' from '/' in further parsing,
	// but RawPath is not set most of the time
	path := r.URL.Path
	if len(r.URL.RawPath) != 0 {
		path = r.URL.RawPath
	}
	// removing the trailing slash that doesn't needed in title
	// also makes the path to root an empty string
	path = strings.TrimSuffix(path, "/")

	if len(path) == 0 {
		return "Home"
	}

	sp := strings.Split(path, "/")
	title, err := url.PathUnescape(sp[len(sp)-1]) // path may be raw sometimes
	if err != nil {
		log.Println(err)
	}

	return strings.Title(title)
}

func main() {
	// register handlers
	http.HandleFunc("/", countVisits(logVisits(allHandler)))

	// serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Starting the server...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
