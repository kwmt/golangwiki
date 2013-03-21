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

// この例は、ドライバテンプレートの使い方のデモです。
// 明確なヘルパーテンプレートです。
func main() {
	// まずテンポラリディレクトリを作り、テンプレート定義ファイルを作成します。
	// 普通は、テンプレートファイルをあるディレクトリに作って置いていると思いますが、ここでは作成します。
	dir := createTestDir([]templateFile{
		// T0.tmplは普通のT1を呼ぶだけのテンプレートです。
		{"T0.tmpl", "T0 ({{.}} version) invokes T1: ({{template `T1`}})\n"},
		// T1.tmplはT2を呼び出すテンプレートT1を定義します。T2は定義されていないことに注意して下さい。
		{"T1.tmpl", `{{define "T1"}}T1 invokes T2: ({{template "T2"}}){{end}}`},
	})
	// テストが終わったら削除します。
	defer os.RemoveAll(dir)

	// パターンは、すべてのテンプレートファイルを検索したglobパターンです。
	pattern := filepath.Join(dir, "*.tmpl")

	// ここが本題です。
	// driversをロードします。
	drivers := template.Must(template.ParseGlob(pattern))

	// T2テンプレートを実装します。
	// まず、ドライバをクローンします。それから、T2定義をテンプレート名前空間へ追加します。

	// 1. Clone the helper set to create a new name space from which to run them.
	first, err := drivers.Clone()
	if err != nil {
		log.Fatal("cloning helpers: ", err)
	}
	// 2. Define T2, version A, and parse it.
	// 2. T2として"T2,varsiton A"を定義し、Parseします。
	_, err = first.Parse("{{define `T2`}}T2, version A{{end}}")
	if err != nil {
		log.Fatal("parsing T2: ", err)
	}

	// さて、T2の違うバージョンを使って、これまでのことを繰り返します。
	// 1. Clone the drivers.
	// 1. ドライバをクローンします。
	second, err := drivers.Clone()
	if err != nil {
		log.Fatal("cloning drivers: ", err)
	}
	// 2.T2として"T2,version B"を定義し、Parseします。
	_, err = second.Parse("{{define `T2`}}T2, version B{{end}}")
	if err != nil {
		log.Fatal("parsing T2: ", err)
	}

	// 逆の順番でテンプレートを実行するのは、
	// 1番目は2番めに影響を受けないこと確認するためです。
	err = second.ExecuteTemplate(os.Stdout, "T0.tmpl", "second")
	if err != nil {
		log.Fatalf("second execution: %s", err)
	}
	err = first.ExecuteTemplate(os.Stdout, "T0.tmpl", "first")
	if err != nil {
		log.Fatalf("first: execution: %s", err)
	}

	// Output:
	// T0 (second version) invokes T1: (T1 invokes T2: (T2, version B))
	// T0 (first version) invokes T1: (T1 invokes T2: (T2, version A))
}