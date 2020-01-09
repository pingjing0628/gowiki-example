package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"text/template"
)

// A wiki consists of a series of interconnected pages,
// each of which has a title and a body (the page content)

// type []byte means "a byte slice"
// the body element is a []byte rather than string because that is the type expected by the io libraries we will use
// Page struct describes how page data will be stored in memory
type Page struct {
	Title string
	Body  []byte
}

// this is method named save that takes as its receiver p , a pointer to Page.
// It takes no parameters, and returns a value of type error
// the method will save the Page's Body to a text file.
// The save method returns an error value because that is the return type of WriteFile (a standard library function that writes a byte slice to a file).
// The save method returns the error value, to let the application handle it should anything go wrong while writing the file.
// If all goes well, Page.save() will return nil (the zero-value for pointers, interfaces, and some other types).
// The octal integer literal 0600, passed as the third parameter to WriteFile,
// which indicates that the file should be created with read-write permissions for the current user only.

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// The function loadPage constructs the file name from the title parameter, reads the file's contents into a new variable body,
// and returns a pointer to a Page literal constructed with the proper title and body values.

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
// 	m := validPath.FindStringSubmatch(r.URL.Path)
// 	if m == nil {
// 		http.NotFound(w, r)
// 		return "", errors.New("Invalid Page Title")
// 	}
// 	return m[2], nil // The title is the second subexpression.
// }

// viewHandler extracts the page title from r.URL.Path,
// re-sliced with [len("/view/"):] to drop the leading "/view/" component of the request path.
// because the path will invariably begin with "/view/", which is not part of the page's title.
// then loads the page data, formats the page with a string of simple HTML, and writes it to w.

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	// title, err := getTitle(w, r)
	// if err != nil {
	// 	return
	// }
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	// title, err := getTitle(w, r)
	// if err != nil {
	// 	return
	// }
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	// title, err := getTitle(w, r)
	// if err != nil {
	// 	return
	// }
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	// p1 := &Page{Title: "TestPage", Body: []byte("This is a test Page.")}
	// p1.save()
	// p2, _ := loadPage("TestPage")
	// fmt.Println(string(p2.Body))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
