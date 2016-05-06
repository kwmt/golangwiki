## Introduction
<p>
    このチュートリアルでカバーできる範囲
</p>
<ul>
    <li>loadやsaveメソッドを持つデータ構造を作成</li>
    <li><code>net/http</code>パッケージを使用してwebアプリケーションを作成</li>
    <li><code>html/template</code>パッケージを使用してHTMLテンプレートを処理</li>
    <li><code>regexp</code>パッケージを使用してユーザ入力をチェック</li>
    <li>クロージャを使用</li>
</ul>
<p>前提の知識</p>
<ul>
    <li>プログラミング経験</li>
    <li>基本的なWeb技術（HTTP,HTMLなど）の理解</li>
    <li>UNIX/DOSコマンドの知識</li>
</ul>
## Getting Started
<p>現時点では、Goを動かすにはFreeBSD,Linux,OS X,Windowsマシンが必要です。
    コマンドプロンプトを表す記号として、<code>$</code>を使用します。</p>
<p>
    Goをインストールします。(インストールは<a href="http://golang.org/doc/install">原文</a>を参照。
    <a href="http://androg.seesaa.net/article/270624478.html">私の記事</a>もよかったらどうぞ)
</p>
<pre class="go">
$ mkdir gowiki
$ cd gowiki
</pre>
<p>
    <code>wiki.go</code>という名前のファイルを作成し、エディタで開いてください。
    以下のラインを追加してください。
</p>
<pre class="go">
package main
import (
    "fmt"
    "io/ioutil"
)
</pre>
<p>
    私たちは<code>fmt</code>パッケージと<code>ioutil</code>パッケージをGoの標準ライブラリからインポートしました。
    Later, as we implement additional functionality, we will add more packages to this import declaration.
    後ほど、追加機能を実装していく度に、私たちは<code>import</code>宣言を追加していくことになります、
</p>
## データ構造
<p>
    データ構造を定義しましょう。このwikiは相互に関連した一連のページから構成され、
    各ページはタイトルとボディ(ページコンテンツ)を持っています。
    ここでは、titleとbodyを表す２つのフィールドを持つ構造体として
    <code>Page</code>を定義します。
</p>
<pre class="go">
type Page struct (
    Title string
    Body  []byte
)
</pre>
<p>
    <code>[]byte</code>型は、<code>byte</code>のスライスを意味しています。
    (スライスについては <a href="http://golang.org/doc/articles/slices_usage_and_internals.html">Slices:usage and internals</a>
    (<a href="https://github.com/kwmt/golangwiki/wiki/GoLang_Slices_usage_and_internals">日本語訳</a>)を参照してください)
    <code>Body</code>要素の型は<code>string</code>より、<code>[]byte</code>を使ったほうがいいでしょう。なぜなら、
    以下で見るように、<code>io</code>ライブラリを使うことを想定しているからです。
</p>
<p>
    <code>Page</code>構造体は、メモリ内でどのように保持するかわかります。
    しかし、永続的に保持したいものは何でしょうか。
    <code>Page</code>に対して<code>save</code>メソッドを作ることによって処理しましょう:
</p>
<pre class="go">
func (p *Page) save() error {
    filename := p.Title + ".txt"
    return ioutil.WriteFile(filname, p.Body, 0600)
}
</pre>
<p>
    このメソッドの特性は、次のとおりです:"これは<code>save</code>と名付けられたメソッドです。
    <code>Page</code>へのポインタ<code>p</code>を受け取ります
    これは引数をとらず、<code>error</code>を返します。"
</p>
<p>
    このメソッドは、<code>Page</code>の<code>Body</code>をテキストファイルにセーブします。
    単純にファイル名として、<code>Title</code>を使います。
</p>
<p>
    <code>save</code>メソッドが<code>error</code>値を返す理由は、
    <code>WriteFile</code>(ファイルにbyteのスライスを書き出す標準ライブラリ)の戻り値だからです。
    <code>save</code>メソッドは、ファイルに書き出し中に何か間違ってしまった場合、
    アプリケーションにハンドリングさせるために、エラー値を返します。
    もしすべて書き出しに成功したら、<code>Page.save()</code>は
    <code>nil</code>(ポインタやインターフェースや他の型にとってのゼロ値)を返します。
</p>
<p>
    <code>WriteFile</code>の3番目の引数の8進数の値<code>0600</code>は、
    カレントユーザーに対して読み書き権限を付与した状態でファイルを作成しますということ意味しています。
</p>
<p>
   ページもロードしたいですよね:
</p>
<pre class="go">
func loadPage(title string) *Page {
    filename := title + ".txt"
    body, _ := ioutil.ReadFile(filename)
    return &amp;Page{Title: title, Body: body}
}
</pre>
<p>
    関数<code>loadPage</code>は、<code>Title</code>カラファイル名を作り、
    新しい<code>Page</code>にファイルの内容を読み込んで、その新しい<code>page</code>へのポインタを返します。
</p>
<p>
    関数は複数の戻り値を返すことができます。
    標準ライブラリ関数の<code>io.ReadFile</code>は<code>[]byte</code>と<code>error</code>を返します。
    <code>loadPage</code>では、エラーは扱いませんので、
    エラーの戻り値(実際には値無し)を無視するために、アンダースコア(_)によって表現されるブランク識別子を使っています。
</p>
<p>
    しかし、<code>ReadFile</code>でエラーになったらどうなりますか？
    たとえば、ファイルが存在しなかったり。
    このようなエラーを無視すべきではありません。
    この関数を<code>*Page</code>と<code>error</code>を返すように修正しましょう。
</p>
<pre class="go">
func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil,err
    }
    return &amp;Page{Title: title, Body: body}, nil
}
</pre>
<p>
    この関数の呼び出しは、2番目の引数をチェックします。
    もし2番目の引数が<code>nil</code>だったら、めでたくPageをロードします。
    もし2番目の引数が<code>nil</code>でなかったら、呼び出し側で扱える<code>error</code>になります。
    (詳しくは、<a href="http://golang.org/ref/spec#Errors">language specification</a>を参照ください)
</p>
<p>
    この時点で、簡単なデータ構造と、ファイルにセーブとファイルからロードする機能を持っています。
    いままで書いたのをテストするために、<code>main</code>関数を書きましょう。
</p>
<pre class="go">
func main() {
    p1 := &amp;Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
    p1.save()
    p2, _ := loadPage("TestPage")
    fmt.Println(string(p2.Body))
}
</pre>
<p>
    このコードをコンパイルして実行すると、<code>p1</code>の内容が書かれた
    <code>TestPage.txt</code>というファイルが作成されます。
    このファイルを構造体<code>p2</code>に読み込んで、<code>p2</code>の<code>Body</code>要素を
    出力します。
</p>
<p>
    下記のようにして、プログラムをコンパイル・実行できます:
</p>
<pre class="go">
$ go build wiki.go
$ ./wiki
This is a sample page.
</pre>
<p>
    (Windowsを使っている場合、実行させるには"./"を入力しないで"wiki"だけ入力する必要があるかもしれません。)
</p>
<p>
    これまでに書いてきた全体のコードを表示するには、
    <a href="http://golang.org/doc/articles/wiki/part1.go">こちら</a>
    をクリックしてください。
</p>
## <code>net/http</code>パッケージの紹介(an interlude)
<p>
    簡単なwebサーバーの動く例を記述します:
</p>
<pre class="go">
package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
</pre>
<p>
    <code>main</code>関数は、<code>http.HandleFunc</code>のコールで始まります。
    <code>handler</code>を使用したWebルート（"/"）へのすべての要求を処理するように、
    <code>http</code>パッケージに問い合わせます。
</p>
<p>
    <code>http.ListenAndServe</code>をコールして、
    ポート8080でリッスン（待ち受け）すべきということを指定しています。
    (2番目のパラメータについては、今は<code>nil</code>で問題ありません。)
    この関数は、プログラムを止めるまで待ち受け状態になっています。
</p>
<p>
    The function handler is of the type http.HandlerFunc.
    関数<code>handler</code>は、<code>http.HandlerFunc</code>の引数で、
    <code>http.ResponseWriter</code>と<code>http.Request</code>
    の引数をとります。
</p>
<p>
    <code>http.ResponseWriter</code>は、HTTPサーバーの応答を組み立て、
    それに書き込むことによって、HTTPクライアントにデータを送信します。
</p>
<p>
    <code>http.Request</code>は、クライアントのHTTP要求を表すデータ構造です。
    文字列<code>r.URL.Path</code>は、リクエストされたURLのパスコンポーネントです。
    末尾に<code>[1:]</code>としているのは、
    1文字目から最後までの<code>Path</code>のサブスライスを作成することを意味しています。
    これは、パス名から先頭の"/"を削除するということです。
</p>
<p>
    もし、以下のプログラムを作り、そのURLにアクセスした場合:
</p>
<pre class="go">
http://localhost:8080/monkeys
</pre>
<p>
    出力は以下のように表示されるでしょう:
</p>
<pre class="go">
Hi there, I love monkeys!
</pre>


## wikiページに<code>net/http</code>を使用します
<p>
    <code>net/http</code>パッケージを使用するには、インポートしなければなりません:
</p>
<pre class="go">
import (
    "fmt"
    "<strong>net/http</strong>"
    "io/ioutil"
)
</pre>
<p>
    wikiページを見るためのハンドラを作成しましょう:
</p>
<pre class="go">
const lenPath = len("/view/")
</pre>
<pre class="go">
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[lenPath:]
    p, _ := loadPage(title)
    fmt.Fprintf(w, "&lt;h1&gt;%s&lt;/h1&gt;&lt;div&gt;%s&lt;/div&gt;"), p.Title, p.Body)
}
</pre>
<p>
    まず、この関数は<code>r.URL.Path</code>からページタイトルを取り出します。
    グローバル定数<code>lenPath</code>はリクエストパスの先頭の"/view/"の文字列の長さです。
    この<code>Path</code>は<code>[lenPath:]</code>で最スライスされ、
    文字列の最初の６文字を捨てます。
    ページは必ず"/view/"で始まり、この"/view/"はページタイトルの部分ではないからです。
</p>
<p>
    それから、この関数はページデータをロードします。
    シンプルなHTMLの文字列でページをフォーマットし、<code>http.ResponseWriter</code>の<code>w</code>
    に、それを書き込みます。
</p>
<p>
    再度になりますが、<code>loadPage</code>からの戻り値<code>error</code>は、
    <code>_</code>を使って無視していることに注意してください。
    これで、シンプルな全体的な悪い習慣を考察することは終わりです。
    後ほど修正します。
</p>
<p>
    ハンドラを使うには、<code>main</code>関数を作成します。
    <code>viewHandler</code>を使って<code>http</code>を初期化し、
    <code>/view/</code>パス配下のリクエストをハンドリングします。
</p>
<pre class="go">
func main() {
    http.HandleFunc("/view/", viewHandler)
    http.ListenAndServe(":8080", nil)
}
</pre>
<p>
    これまでに書いてきた全体のコードを表示するには、
    <a href="http://golang.org/doc/articles/wiki/part2.go">こちら</a>
    をクリックしてください。
</p>
<p>
    <code>test.txt</code>のようなページを作成して、コードをコンパイルし、
    wikiページ表示させてみましょう。
</p>
<p>
    <code>test.txt</code>ファイルをエディタで開いて、
    "Hello world"(引用符は外してね)を書いて、セーブします。
</p>
<pre class="go">
$ go build wiki.go
$ ./wiki
</pre>
<p>
    webサーバを実行させて、
    <a href="http://localhost:8080/view/test"><code>http://localhost:8080/view/test</code></a>
    にアクセスすると、ページタイトルに"test"、bodyに"Hello world"と表示されるはずです。
</p>

## ページを編集する
<p>
    wikiは、ページを任意に編集できなければ、wikiではありません。
    ２つのハンドラを作成しましょう:１つは'編集ページ'フォームを表示するための<code>editHandler</code>を、
    もう１つは、フォームに入力したデータをセーブするための<code>saveHandler</code>を作成します。
</p>
<p>
    まず、<code>main()</code>関数にそれぞれ追加します:
</p>
<pre class="go">
func main() {
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/save/", saveHandler)
    http.ListenAndServe(":8080", nil)
}
</pre>
<p>
    関数<code>editHandler</code>は、ページをロードしてから
    （あるいは、ページが存在しなかったら、空のページ構造体を作成してから）、HTMLフォームを表示します。
</p>
<pre class="go">
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[lenPath:]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    fmt.Fprintf(w, "&lt;h1&gt;Editing %s&lt;/h1&gt;" +
        "&lt;form action=\"/save/%s\" method=\"POST¥\"&gt;" +
        "&lt;textarea name=\"body\"&gt;%s&lt;/textarea&gt;&lt;br&gt;" +
        "&lt;/form&gt;",
        p.Title, p.Title, p.Body)
}
</pre>
<p>
    この関数は正常に動作しますが、ハードコーディングされたHTMLはかっこ悪いです。
    いい方法がもちろんあります。
</p>

## <code>html/template</code>パッケージ
<p>
    <code>html/template</code>パッケージGo標準ライブラリの一つです。
    <code>html/template</code>は、HTMLファイルを分けて記述することができ、
    Goコード中で編集することなく編集ページのレイアウトを変更することができるようになります。
</p>
<p>
    まず、<code>html/template</code>をimport文に追加する必要があります:
</p>
<pre class="go">
import (
    <strong>"html/template"</strong>
    "fmt"
    "net/http"
    "io/ioutil"
)
</pre>
<p>
    HTMLフォームのテンプレートファイルを作成しましょう。<code>edti.html</code>というファイルを
    新しく作成し、下記を追加してください。
</p>
<pre class="go">
&lt;h1&gt;Editing {{.Title}}&lt;/h1&gt;

&lt;form action="/save/{{.Title}}" method="POST"&gt;
&lt;div&gt;&lt;textarea name="body" rows="20" cols="80"&gt;{{printf "%s" .Body}}&lt;/textarea&gt;&lt;/div&gt;
&lt;div&gt;&lt;input type="submit" value="Save"&gt;&lt;/div&gt;
&lt;/form&gt;
</pre>
<p>
    ハードコーディングされたHTMLの代わりにテンプレートを使うために、<<code>editHandler</code>を修正します:
</p>
<pre class="go">
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[lenPath:]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    t, _ := template.ParseFiles("edit.html")
    t.Execute(w, p)
}
</pre>
<p>
    <code>template.ParseFiles</code>関数は、<code>edit.html</code>を読み込んで、
    <code>*template.Template</code>を返します。
</p>
<p>
    <code>t.Execute</code>メソッドは、
    生成されたHMTLを<code>http.ResponseWriter</code>に書き込んで、
    テンプレートを実行します。
    <code>.Title</code>と<code>.Body</code>ｍのドット識別子は、それぞれ、
    <code>p.Title</code>,<code>p.Body</code>を参照します。
</p>
<p>
    テンプレートのディレクティブは、２重の中括弧で囲まれます。
    <code>printf "%s" .Body</code>命令は、1つの関数コールです。
    バイトストリームの代わりの文字列として<code>.Body</code>をアウトプットします。
    <code>fmt.Printf</code>と同じです。
    <code>html/template</code>パッケージは、安全で正しいHTMLを作成しやすくなります。
    たとえば、<code>></code>のような記号を自動的にエスケープし、<code>&gt;</code>に置き換え、
    ユーザーデータはHTMLフォームを壊しません。
</p>
<p>
    さて、<code>fmt.Fprintf</code>を削除しました。
    <code>import</code>リストから"fmt"を削除できます。
</p>
<p>
    テンプレートを使っていますので、<code>view.html</code>をコールする<code>viewHandler</code>
    用のテンプレートを作成しましょう:
</p>
<pre class="go">
&lt;h1&gt;{{.Title}}&lt;/h1&gt;
&lt;p&gt;[&lt;a href="/edit/{{.Title}}"&gt;edit&lt;/a&gt;]&lt;/p&gt;
&lt;div&gt;{{printf "%s" .Body}}&lt;/div&gt;
</pre>
<p>
    <code>viewHandler</code>関数の中身も修正しましょう:
</p>
<pre class="go">
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[lenPath:]
    p, _ := loadPage(title)
    t, _ := template.ParseFiles("view.html")
    t.Execute(w, p)
}
</pre>
<p>
    2つのハンドラでほぼ同じテンプレートコードとなっていることに注目してください。
    重複部分を関数にまとめることで、重複を削除しましょう:
</p>
<pre class="go">
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[lenPath:]
    p, _ := loadPage(title)
    renderTemplate(w, "view", p)
}
</pre>
<pre class="go">
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[lenPath:]
    p, err != loadPage(title)
    renderTemplate(w, "edit", p)
}
</pre>
<pre class="go">
func renderTemplate(w http.ResponseWrite, tmpl string, p *Page) {
    t, _ := template.ParseFiles(tmpl + ".html")
    t.Execute(w, p)
}
</pre>
<p>
    ハンドラは、短くシンプルになりました。
</p>

## 存在しないページの取り扱い(Handling non-existent pages)
<p>
    もし、<code>/view/APageThatDoesntExist</code>を訪れたらどうなりますか？
    クラッシュするでしょう。これは、<code>loadPage</code>の戻り値である<code>error</code>を
    無視しているからです。
    代わりに、ページが存在しないなら、クライアントを編集ページへリダイレクトさせるべきです。
</p>
<pre class="go">
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/" + title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}
</pre>
<p>
    <code>http.Redirect</code>関数は、
    <code>http.StatusFound</code>(302)のHTTPステータスコードを追加し、
    <code>Location</code>ヘッダーをHTTPレスポンスに追加します。
</p>

## ページをセーブする
<p>
    関数<code>saveHandler</code>は送信フォームを扱います。
</p>
<pre class="go">
func saveHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[lenPath:]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    p.save()
    http.Redirect(w, r, "/view/" + title, http.StatusFound)
}
</pre>
<p>
    URLから与えられるページのタイトルと<code>Body</code>のみのフォームのフィールドは、
    新しいページに保存されます。
    データをファイルに保存するために<code>save()</code>メソッドをコールし、
    クライアントは<code>/view/</code>ページにリダイレクトされます。
</p>
<p>
    <code>FormValue</code>の戻り値は、<code>string</code>です。
    <code>Page</code>構造体にフィットさせるためには、
    <code>FormValue</code>の戻り値を<code>[]byte</code>に変換する必要があり、
    <code>[]byte(body)</code>で変換します。
</p>

## エラーの取り扱い
<p>
    これまでのプログラムには、エラーが無視される箇所がいくつかあります。
    エラーが発生した時にクラッシュしますので、悪い習慣ではなく最低です。
    よい解決策は、エラーを扱い、ユーザーにエラーメッセージを表示させることです。
    何かがうまくいかない場合に、サーバーはクラッシュさせないし、ユーザーに通知します。
</p>
<p>
    まず、<code>renderTemplate</code>メソッドでエラーを扱いましょう:
</p>
<pre class="go">
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, err := template.ParseFiles(tmpl + "html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = t.Execute(w, p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
</pre>
<p>
    <code>http.Error</code>関数は、特定のHTTPレスポンスコード(ここでは、"Internal Server Error")
    と、エラーメッセージを送ります。
    Already the decision to put this in a separate function is paying off.
</p>
<p>
    さて、<code>saveHanler</code>を修正しましょう:
</p>
<pre class="go">
func saveHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err = p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
</pre>
<p>
    <code>p.save()</code>中に起こりうるどんなエラーも、ユーザーに通知します。
</p>

## テンプレートのキャッシュ
<p>
    このコードは非効率な箇所があります:
    ページがレンダリングされるたびに毎回<code>renderTemplete</code>が<code>ParseFiles</code>をコールしています。
    よいアプローチは、プログラムの初期化の時点で一度だけ<code>ParseFiles</code>をコールすることです。
    一つの<code>*Tmplate</code>にすべてのテンプレートをパースします。
    それから、指定したテンプレートをレンダリングするために、<code>ExecuteTemplate</code>メソッドを使います。
</p>
<p>
    <code>templates</code>と名付けたグローバル変数を作り、
    <code>ParseFiles</code>を使って初期化します。
</p>
<pre class="go">
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
</pre>
<p>
    <code>template.Must</code>関数は、便利にラッパーです。
    <code>nil</code>でない<code>error</code>が起こった時panicになりますが、
    <code>*Template</code>を返します。
    A panic is appropriate here; if the templates can't be loaded
    the only sensible thing to do is exit the program.
</p>
<p>
    パースしたいテンプレート名をイテレートするのに、
    <code>range</code>ステートメントと一緒に<code>for</code>ループが使われています。
    私達がもっと多くのテンプレートを追加する際、その配列に名前を追加します。
</p>
<p>
    対応したテンプレート名で、<code>templates.ExecuteTemplate</code>をコールするために
    <code>renderTemplate</code>関数を修正します。
</p>
<pre class="go">
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := template.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
</pre>
<p>
    テンプレート名が、テンプレートのファイル名であることに注意してください。
    ですので、<code>tmpl</code>の引数に<code>".html"</code>を追加する必要があります。
</p>

## Validation
<p>
    気づいているかもしれませんが、このプログラムには、重大なセキュリティ欠陥が存在します。
    ユーザーはサーバー上で任意のパスで読み書きできることです。
    この欠陥を防ぐために、正規表現を使った妥当性確認をする関数を書くことができます。
</p>
<p>
    それにはまず、<code>import</code>リストに<code>"regexp"</code>を追加してください。
    そして、グローバル変数を作り、妥当性確認用の正規表現を変数に格納してください:
</p>
<pre class="go">
var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")
</pre>
<p>
    <code>regexp.MustCompile</code>関数は、正規表現をパースしコンパイルし,
    regexp.Regexpを返します。
    <code>Compile</code>の2番目のパラメータとして<code>error</code>を返す一方で、
    <code>MustCompile</code>は、式のコンパイルに失敗した場合パニックになるという点で
    <code>Compile</code>とは異なります。
</p>
<p>
    さて、リクエストURLからタイトル文字を取得する関数を書きましょう。
    そして、テストしましょう:
</p>
<pre class="go">
func getTitle(w http.ResponseWriter, r *http.Request) (title string, err error) {
    title = r.URL.Path[lenPath:]
    if !titleValidator.MatchString(title) {
        http.NotFound(w, r)
        err = errors.New("Invalid Page Title")
    }
    return
}
</pre>
<p>
    （訳者注：<code>errors</code>を使っていますが、これまでに「<code>errors</code>
    をimportしてください」のような説明がなかったのですが、<code>"errors"</code>を
    importに追加してください）
</p>
<p>
    もしタイトルが有効な場合、<code>nil</code>のエラー値と一緒に<code>title</code>を返します。
    もしタイトルが無効な場合、この関数は"404 Not Found"エラーをHTTPコネクションに書き込み、
    エラーをハンドラに返します。
</p>
<p>
    各ハンドラで<code>getTitle</code>をコールするようにしてみましょう:
</p>
<pre class="go">
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}
</pre>
<pre class="go">
func editHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}
</pre>
<pre class="go">
func saveHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err = p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
</pre>

## Introducing Function Literals and Closures
<p>
    各ハンドラでエラー条件をキャッチすることは、繰り返したくさんのコードを引き起こします。
    ハンドラのそれぞれの妥当性確認とエラーチェックを、１つの関数でラップできたらどうなりますか？
    Go言語の<a href="http://golang.org/ref/spec#Function_declarations">関数リテラル</a>は、
    ここで私たちを助けることができる抽象的な機能の強力な手段を提供します。
</p>
<p>
    まず、それぞれのハンドラの関数定義を、<code>title</code>を引数として追加します。
</p>
<pre class="go">
func viewHandler(w http.ResponseWriter, r *http.Request, title string)
func editHandler(w http.ResponseWriter, r *http.Request, title string)
func saveHandler(w http.ResponseWriter, r *http.Request, title string)
</pre>
<p>
    ラッパー関数を定義し、<code>http.HandlerFunc</code>の関数を戻り値とします:
</p>
<pre class="go">
func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // ここでリクエストからページタイトルを取り出し、
        // 提供されたハンドラ'fn'をコールします。
    }
}
</pre>
<p>
    外部で定義された値を覆うために返された関数は、クロージャと呼ばれます。
    この場合、変数<code>fn</code>(<code>makeHandler</code>の一つの引数)は、クロージャによって覆われます。
    この変数<code>fn</code>は、私たちのsave,edit,viewハンドラの１つとなります。
</p>
<p>
    私たちは<code>getTitle</code>のコードを使います（いくつかマイナー修正しています):
</p>
<pre class="go">
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        title := r.URL.Path[lenPath:]
        if !titleValidator.MatchString(title) {
            http.NotFound(w, r)
            return
        }
        fn(w, r, title)
    }
}
</pre>
<p>
    <code>makeHandler</code>で返されたクロージャは、<code>http.ResponseWriter</code>と
    <code>http.Request</code>(他の言葉では、<code>http.HandlerFunc</code>)を取る関数です。
    クロージャはリクエストパスから<code>title</code>を取り出し、
    <code>TitleValidator</code>の正規表現で、妥当性を確認します。
    <code>title</code>が無効の場合、エラーが<code>http.NotFound</code>関数を使って
    <code>ResponseWriter</code>に書き込まれます。
    <code>title</code>が有効の場合、<code>ResponseWriter</code>,<code>Request</code>,
    <code>title</code>を引数として、内包されたハンドラ関数<code>fn</code>をコールします。
</p>
<p>
    <code>http</code>パッケージで登録される前に、<code>main</code>に<code>makeHandler</code>で
    ハンドラ関数をラップします:
</p>
<pre class="go">
func main() {
    http.HandlerFunc("/view/", makeHandler(viewHandler))
    http.HandlerFunc("/edit/", makeHandler(editHandler))
    http.HandlerFunc("/save/", makeHandler(saveHandler))
    http.ListenAndServe(":8080", nil)
}
</pre>
<p>
    最終的にハンドラ関数から<code>getTitle</code>をコールする部分を削除します。
    シンプルになります:
</p>
<pre class="go">
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}
</pre>
<pre class="go">
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err != loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}
</pre>
<pre class="go">
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []body(body)}
    err := p.save()
    if err != nil {
        http.Error(w, Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
</pre>

## 動かしてみよう！
<p>
    <a href="http://golang.org/doc/articles/wiki/final.go">
    これまでのコードの完成版をみるには、こちらをクリックしてください。
    </a>
</p>
<p>
    コードをコンパイルして、アプリケーションを実行します:
</p>
<pre class="go">
$ go build wiki.go
$ ./wiki
</pre>
<p>
    <a href="http://localhost:8080/view/ANewPage">http://localhost:8080/view/ANewPage</a>
    ページを訪れると、編集フォームのページが現れるはずです。
    テキストを入力することができ、'Save'ボタンをクリックすると、
    新しく生成されたページにリダイレクトされるはずです。
</p>

## 他の課題
<p>
    ここでは、いくつかの課題を与えますので、あなた自身で取り組んでください。
</p>
<p>
<ul>
    <li>
        テンプレートを<code>tmpl/</code>に、
        ページデータを<code>data/</code>に格納してください。
    </li>
    <li>
        webルートを<code>/view/FrontPage</code>にリダイレクトさせるハンドラを追加してください。
    </li>
    <li>
        適切なHTMLといくつかのCSSを追加して、テンプレートページを綺麗にしてください。
    </li>
    <li>
        <code>[PageName]</code>を<code><a href="/view/PageName">PageName</a></code>に
        変換することで、内部ページ間のリンクができるような実装をしてください。
        (ヒント：<a href="http://golang.org/src/pkg/regexp/regexp.go?h=ReplaceAllFunc#L541">
        <code>regexp.ReplaceAllFunc</code></a>を使うと良いかも）
    </li>
</ul>
</p>
