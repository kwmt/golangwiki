package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)


// templateFile構造体は、テスト用に、ファイルに保存するテンプレートの内容を定義します。
type templateFile struct {
	name     string
	contents string
}

func createTestDir(files []templateFile) string {
	dir, err := ioutil.TempDir("", "template") // $TMPDIR/tempalteにサフィックスが付いたディレクトリを作成
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		f, err := os.Create(filepath.Join(dir, file.name)) // $TMPDIR/tempalte + suffix/file.name を作成
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		_, err = io.WriteString(f, file.contents) // file.contentsをfに記入
		if err != nil {
			log.Fatal(err)
		}
	}
	return dir
}

// ディレクトリにあるテンプレート郡をロードするデモです。
func main() {
	dir := createTestDir([]templateFile{
		// T0.tmpl はT1をコールします。
		{"T0.tmpl", `T0 invokes T1: ({{template "T1"}})`},
		// T1.tmpl T2をコールするテンプレートを"T1"として定義しています。
		{"T1.tmpl", `{{define "T1"}}T1 invokes T2: ({{template "T2"}}){{end}}`},
		// T2.tmpl は"T2"としてテンプレートを定義しています。
		{"T2.tmpl", `{{define "T2"}}This is T2{{end}}`},
	})
	// main関数終了後、作成したディレクトリを削除します。
	defer os.RemoveAll(dir)

	// dirにあるすべてのテンプレートファイルを検索するようにパターンを作成しています。
	pattern := filepath.Join(dir, "*.tmpl")

	// ここから本題です。
	// T0.tmpl は初めにマッチする名前ですので、始めのテンプレート(ParseGlobの戻り値の値)になります。
	tmpl := template.Must(template.ParseGlob(pattern))
	log.Println(tmpl.Name())

	err := tmpl.Execute(os.Stdout, nil)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	// Output:
	// T0 invokes T1: (T1 invokes T2: (This is T2))
}

