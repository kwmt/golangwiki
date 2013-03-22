package main

import (
	"fmt"
	"os"
	"text/template"
	"text/template/parse"
)

func main() {
	root := template.Must(template.New("root").Parse(`{{define "a"}} {{.}} {{template "b" }} {{.}} ">{{.}} </a>{{end}}`))
	tree, err := parse.Parse("t", `{{define "b"}}<a href="{{end}}`, "", "", nil, nil)
	if err != nil {
		fmt.Println(err)
	}
	added := template.Must(root.AddParseTree("b", tree["b"]))

	err = added.ExecuteTemplate(os.Stdout, "a", "1>0")
	if err != nil {
		fmt.Println(err)
	}

}
