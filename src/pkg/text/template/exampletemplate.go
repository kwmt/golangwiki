package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	// テンプレートを定義します。
	const letter = `
Dear {{.Name}},
{{if .Attended}}
It was a pleasure to see you at the wedding.{{else}}
It is a shame you couldn't make it to the wedding.{{end}}
{{with .Gift}}Thank you for the lovely {{.}}.
{{end}}
Best wishes,
Josie
`

	// テンプレートに埋め込むデータを準備します。
	type Recipient struct {
		Name, Gift string
		Attended   bool
	}
	var recipients = []Recipient{
		{"Aunt Mildred", "bone china tea set", true},
		{"Uncle John", "moleskin pants", false},
		{"Cousin Rodney", "", false},
	}

	// 新しいテンプレートを作成し、letterをテンプレートに入れてパースします。
	t := template.Must(template.New("letter").Parse(letter))

	// 変数recipientsの各要素に対してテンプレートを実行します。
	for _, r := range recipients {
		err := t.Execute(os.Stdout, r)
		if err != nil {
			log.Println("executing template:", err)
		}
	}

}
