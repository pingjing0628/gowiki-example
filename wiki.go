package main

import (
	"fmt"
	"io/ioutil"
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

func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a test Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}
