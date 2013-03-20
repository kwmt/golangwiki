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

// この例は、いくつかのテンプレートを共有し、それらを違う文脈・コンテキスト使う方法の1つをデモします。
// 変わった形では、存在するテンプレート郡に渡すことで私たちは複数のドライバテンプレートを追加します。
func main() {
	// まず、テンポラリディレクトリを作成し、テンプレート定義ファイルを準備します。
	// 普通は、テンプレートファイルをあるディレクトリに作って置いていると思いますが、ここでは作成します。
	dir := createTestDir([]templateFile{
		// T1.tmpl defines a template, T1 that invokes T2.
		{"T1.tmpl", `{{define "T1"}}T1 invokes T2: ({{template "T2"}}){{end}}`},
		// T2.tmpl defines a template T2.
		{"T2.tmpl", `{{define "T2"}}This is T2{{end}}`},
	})
	// テストが終わったら削除します。
	defer os.RemoveAll(dir)

	// パターンは、すべてのテンプレートファイルを検索したglobパターンです。
	pattern := filepath.Join(dir, "*.tmpl")

	// ここが本題です。
	// helperをロードします。
	templates := template.Must(template.ParseGlob(pattern))
	// Add one driver template to the bunch; we do this with an explicit template definition.
	_, err := templates.Parse("{{define `driver1`}}Driver 1 calls T1: ({{template `T1`}})\n{{end}}")
	if err != nil {
		log.Fatal("parsing driver1: ", err)
	}
	// 別のドライバテンプレートを追加します。
	_, err = templates.Parse("{{define `driver2`}}Driver 2 calls T2: ({{template `T2`}})\n{{end}}")
	if err != nil {
		log.Fatal("parsing driver2: ", err)
	}
	// 実行前にすべてのテンプレートをロードします。この操作は必要としませんが、
	// html/templateはエスケープする為、良い習慣です。
	err = templates.ExecuteTemplate(os.Stdout, "driver1", nil)
	if err != nil {
		log.Fatalf("driver1 execution: %s", err)
	}
	err = templates.ExecuteTemplate(os.Stdout, "driver2", nil)
	if err != nil {
		log.Fatalf("driver2 execution: %s", err)
	}
	// Output:
	// Driver 1 calls T1: (T1 invokes T2: (This is T2))
	// Driver 2 calls T2: (This is T2)
}
