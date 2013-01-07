package main

import (
	"fmt"
	"os"
	"text/template"
)

var test *template.Template
var tmpl *template.Template

func main() {

	exampleStr := []string{"{{\"output\"}}", //ここが合わない。。
		"{{`\"output\"`}}",
		"{{printf \"%q\" \"output\"}}",
		"{{\"output\" | printf \"%q\"}}",
		"{{\"put\" | printf \"%s%s\" \"out\" | printf \"%q\"}}",
		"{{\"output\" | printf \"%s\" | printf \"%q\"}}",
		"{{with \"output\"}}{{printf \"%q\" .}}{{end}}",
		"{{with $x := \"output\" | printf \"%q\"}}{{$x}}{{end}}",
		"{{with $x := \"output\"}}{{printf \"%q\" $x}}{{end}}",
		"{{with $x := \"output\"}}{{$x | printf \"%q\"}}{{end}}",
	}

	for i:= range exampleStr {
		test = template.New("test"+ string(i))
		tmpl,err := test.Parse(exampleStr[i])

		if err != nil {
			panic(err)
		}
		fmt.Printf(exampleStr[i] + ":")
		err = tmpl.Execute(os.Stdout, tmpl)
		if err != nil {
			panic(err)
		}
		fmt.Println()
	}


}
