package main

import (
	"log"
	"os"
	"strings"
	"text/template"
)

// この例は、テンプレートテキストを処理するための、カスタム関数のデモです。
// strings.Title関数を組み込みます。そしてそれをMake Title Text Look Goodのように出力するために使用します。
func main() {
	// まず、私たちはFuncMapを作ります。それは関数を登録します。
	funcMap := template.FuncMap{
		// "title"はテンプレートでコールされる関数名となります。
		"title": strings.Title,
	}

	// 先ほどの関数をテストするために、シンプルなテンプレートを定義します。
	// 私たちはいくつかの方法で入力されたテキストを表示します:
	// - オリジナルテキスト
	// - "title"をコール
	// - "title"をコールした後、%qを使って表示する
	// - %qを使って表示した後、"ttitle"をコール
	const templateText = `
Input: {{printf "%q" .}}
Output 0: {{title .}}
Output 1: {{title . | printf "%q"}}
Output 2: {{printf "%q" . | title}}
`

	// テンプレート作成し、関数マップを追加し、テキストをパースします。
	tmpl, err := template.New("titleTest").Funcs(funcMap).Parse(templateText)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	// 出力を確認するため、テンプレートを実行します。
	err = tmpl.Execute(os.Stdout, "the go programming language")
	if err != nil {
		log.Fatalf("execution: %s", err)
	}

	// Output:
	// Input: "the go programming language"
	// Output 0: The Go Programming Language
	// Output 1: "The Go Programming Language"
	// Output 2: "The Go Programming Language"
}
