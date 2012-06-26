package main

import (
	"html/template"
	"net/http"
	"io/ioutil"
)

type Page struct {
	Title string
        Body []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

const lenPath = len("/view/")

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[lenPath:]
	p, _ := loadPage(title)
	t, _ := template.ParseFiles("view.html")
	t.Execute(w, p)
}


/*
func main(){
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}
*/

func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[lenPath:]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
/*
    fmt.Fprintf(w, "<h1>Editing %s</h1>" +
        "<form action=\"/save/%s\" method=\"POST¥\">" +
        "<textarea name=\"body\">%s</textarea><br>" +
        "</form>",
        p.Title, p.Title, p.Body)
*/
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, p)
}



func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.ListenAndServe(":8080", nil)
}