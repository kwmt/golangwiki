<h2>Overview</h2>
<p>
templateパッケージは、テキスト出力するデータ駆動型のテンプレートを実装しています。
</p>
<p>
HTMLを生成するには、html/templateパッケージを参照ください。
それは、このパッケージと同じインターフェースを持ちますが、ある攻撃に対してセキュアなHTMLを自動的に生成します。
</p>
<p>
テンプレートはそれらをデータ構造へ適用することで実行されます。
このtemplateの記法は、実行を制御し表示される値を得るには、データ構造の要素（一般には構造体のフィールドかマップのキー）を参照します。
テンプレートの実行は、構造を解析し、カーソルをセットします。
カーソルは、ピリオドやドットとよばれる'.'によって表現されます。
</p>
<p>
テンプレートの入力テキストは、どんなフォーマットにでもあうUTF-8エンコードされたテキストです。
データを評価したり、構造を制御したりする"アクション(Actions)"は、"{{"と"}}"で区切られます。
アクション外のすべてのテキストは、出力は変更されず、コピーされます。
アクションはコメントできますが、改行をまたぐことができません。
</p>
<p>
一度構造化されたテンプレートは並列に安全に実行されるはずです。
</p>
<p>
ここで"17 items are made of wool"を表示させる簡単な例を見てみましょう。
</p>
<pre>
type Inventory struct {
    Material    string
    Count       uint
}
sweaters := Inventory{"wool", 17}
tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
if err != nil { panic(err) }
err = tmpl.Execute(os.Stdout, sweaters)
if err != nil { panic(err) }
</pre>
<p>
もっと複雑な例は以下に示します。
</p>
<h3 id="hdr-Actions">Actions</h3>
<p>
これはアクションのリストです。
</p>
<pre>
{{/* コメント */}}
    コメントです。無視されます。改行を含めます。
    複数のコメントのネストはできません。

{{パイプライン}}
    パイプラインの文字通りの表現を出力先へコピーします。

{{if パイプライン}} T1 {{end}}
    もしパプラインの値が空だったら、出力はありません。空でなければ、T1が実行されます。
    空とは、false,0, nilポインタやインターフェースの値,そして配列・スライス・マップの長さがゼロのことを言います。
    ドットは影響を受けません。

{{if パイプライン}} T1 {{else}} T0 {{end}}
    もしパイプラインの値が空だったら、T0を実行します。
    空でなければ、T1を実行します。
    ドットは影響を受けません。

{{range パイプライン}} T1 {{end}}
    パイプラインの値は、配列・スライス・マップのいずれかでなければなりません。
    もしその値の長さがゼロなら、何も出力しません。
    ゼロでなければ、ドットは配列・スライス・マップの連続した要素をセットし、T1が実行されます。
    もしその値がマップで、keys are of basic type with a defined order ("comparable") ならば、
    その要素はソートされたキーの順番でアクセスされるでしょう。

{{range パイプライン}} T1 {{else}} T0 {{end}}
    パイプラインの値は、配列・スライス・マップのいずれかでなければなりません。
    もしその値の長さがゼロなら、ドットは影響を受けず、T0が実行されます。
    ゼロでなければ、ドットは配列・スライス・マップの連続した要素をセットし、T1が実行されます。

{{template "name"}}
    指定したnameのテンプレートが、nilデータと共に実行されます。

{{template "name" パイプライン}}
    指定したnameのテンプレートが、ドットがパイプラインの値にセットされて実行されます。

{{with パイプライン}} T1 {{end}}
    パイプラインの値が空だったら、出力されません。
    空でなければ、ドットはパイプラインの値にセットされ、T1が実行されます。

{{with パイプライン}} T1 {{else}} T0 {{end}}
    パイプラインの値が空だったら、ドットは影響を受けずT0が実行されます。
    空でなければ、ドットはパイプラインの値にセットされ、T1が実行されます。
</pre>
<h2 id="Arguments">Arguments</h2>
<p>
引数は、以下のいずれかで表現されるシンプルな値です。
</p>
<pre>
- Goの文法での真偽値、文字列、文字、整数、浮動小数点、虚数、複素数。
  これらは、Goの型無し定数のように動作します。生の文字列は改行をまたがることはできません。
  
- 文字 '.' (ピリオド):
    .
  この結果はドットの値です。
- 変数名は
    $piOver2
  あるいは
    $
  のような$記号を前に置いた英数字の文字列（空文字も可能）です。
  結果はその変数の値です。
  変数については以下で説明します。
- 構造体データのフィールド名は、
    .Field
  のように、ドットが前に置かれます。
  結果はフィールドの値です。
  フィールドをつなげることで、フィールドのフィールドをよびだせます:
    .Field1.Field2
  フィールドは変数に対して評価することもできます:
    $x.Field1.Field2
- マップデータのキー名は
    .Key
  のようにピリオドを前に置きます。
  結果は、そのキーに対応した要素の値となります。
  キーの呼び出しはこのように連結することもできます:
    .Field1.Key1.Field2.Key2
  Although the key must be an alphanumeric identifier, unlike with
  field names they do not need to start with an upper case letter.
  キーは英数字で一意でなければなりませんが、フィールドとは違って大文字から始める必要はありません。
  キーは次のように変数に対して評価することもできます。
    $x.key1.key2
- The name of a niladic method of the data, preceded by a period,
  such as
    .Method
  データの引数を持たないメソッド名は
    .Method
  のようにピリオドを前に置きます。
  The result is the value of invoking the method with dot as the
  receiver, dot.Method().
  結果は、レシーバーのように(dot.Method())ドットを使ってメソッドを呼び出した値です。
  そのようなメソッドは、1つの戻り値あるいは2つの戻り値(2つ目はerror)をもつ必要があります。
  もし2つの戻り値を持って、errorがnilでなかったら、
  実行は中止し、実行の値として呼び出し元へエラーが返ってきます。
  メソッド呼び出しはフィールドとキーを組み合わせて連鎖して呼び出すことができます。
    .Field1.Key1.Method1.Field2.Key2.Method2
  メソッドは次のように変数に対して評価することもできます。
    $x.Method1.Field
- 次のように引数をもたない関数名です。
    fun
  結果は関数fun()を呼び出した値です。
  返ってくる型と値は、メソッドのように振る舞います。
  関数と関数名については以下に記述します。
</pre>
<p>
    引数はどんな型も評価するかもしれません。
    もし引数がポインタなら、実装は自動的に必要なときに元になる型を指します。
    If an evaluation yields a function value, such as a function-valued field of a struct,
    the function is not invoked automatically,
    その関数は自動的に呼ばれませんが、
    but it can be used as a truth value for an if action and the like.
    それを呼び出すには、以下で定義されるcall関数を使用します。
</p>
<p>
パイプラインは、コマンドを連結することが可能です。
コマンドは、簡単な値や関数やメソッド呼び出し、複数の引数をもつことが可能です。
</p>
<pre>
Argument
    結果は、この引数の評価の値になります。
.Method [Argument...]
    このメソッドは、単独でもいいし連結の最後の要素で構いませんが、
    連消した途中のメソッドはダメで、複数の引数を取ることができます。
    結果は、その引数を使ってコールしたメソッドの値です:
        dot.Method(Argument1, etc.)
functionName [Argument...]
    結果は、名前に関連した関数をコールした値です:
        function(Argument1, etc.)
    関数と関数名は以下で説明します。
</pre>

<h3 id="hdr-Pipelines">Pipelines</h3>
<p>
パイプラインは、パイプラインキャラクタ'|'があるコマンドの列を分けることで"連結"されます。
連結されたパイプラインでは、各コマンドの結果は、次のコマンドの引数として渡されます。
パイプラインでの最後のコマンドの出力はそのパイプラインの値です。  
</p>
<p>
コマンドの出力は、1つの値か2つの値のどちらかになり、2番めの出力は、型errorを持ちます。
もし2番めの値が存在し、nilでなかった場合は実行が終了し、エラーはExecuteの呼び出し元に返されます。
</p>

<h3 id="hdr-Variables">Variables</h3>
<p>
アクション内のパイプラインでは、結果を取得する変数を初期化することができます。
初期化は次の文法となります。
</p>
<pre>
$variable := pipeline
</pre>
<p>
ただし、$variableは変数の名前です。変数を宣言するアクションは何も出力しません。
</p>
<p>
"range"アクションで変数を初期化する場合、変数は繰り返しの連続した要素にセットされます。
また、"range"はカンマで区切って2つの変数宣言します:
</p>
<pre>
range $index, $element := pipeline
</pre>
<p>
この場合、$indexと$elementは、それぞれ、配列/スライスのインデックスまたはマップのキーと要素の
連続した値がセットされています。
もし1つの変数だけなら、要素が割り当てられます。これは、Goのrangeとは反対になります。
</p>
<p>
変数のスコープは、制御構造("if","with",あるいは"range")で宣言されてから"end"アクションまで、または、もしこのような制御構造がない場合はテンプレートの終わりまで有効です。
テンプレートの呼び出しでは、その呼び出しのポイントから変数を受け取りません。
</p>
<p>
When execution begins, $ is set to the data argument passed to Execute, that is,
to the starting value of dot.
</p>


<h2 id="Examples">Examples</h2>
<p>
    ここでは、1行テンプレートのいくつかの例を見てみます。パイプラインと変数のデモです。
    結果はすべて、ダブルクオーテーションがついた"output"です。
</p>
<pre>
{{"\"output\""}}
    文字列定数
{{`"output"`}}
    A raw string constant 生文字列定数
{{printf "%q" "output"}}
    関数を呼びます
{{"output" | printf "%q"}}
    最後の引数は前のコマンドからきて関数をコールします
{{"put" | printf "%s%s" "out" | printf "%q"}}
    もっと凝った関数コール
{{"output" | printf "%s" | printf "%q"}}
    長い連鎖
{{with "output"}}{{printf "%q" .}}{{end}}
    ドットを使ったwithアクション
{{with $x := "output" | printf "%q"}}{{$x}}{{end}}
    変数を定義して使ったwithアクション
{{with $x := "output"}}{{printf "%q" $x}}{{end}}
    他のアクションの中で変数を使ったwithアクション
{{with $x := "output"}}{{$x | printf "%q"}}{{end}}
    上と同じだが、パイプラインを使用
</pre>

<h3 id="hdr-Functions">Functions</h3>
<p>
During execution functions are found in two function maps: first in the
template, then in the global function map. By default, no functions are defined
in the template but the Funcs method can be used to add them.
</p>
<p>
事前に定義されたグローバル関数は以下のとおりに名前付けされています。
</p>
<pre>
and
    Returns the boolean AND of its arguments by returning the
    first empty argument or the last argument, that is,
    "and x y" behaves as "if x then y else x". All the
    arguments are evaluated.
    "and x y"は"if x then y else x"のように振舞います。
call
    最初の引数をコールした結果を返します。この時引数は関数です。残り引数はパラメータとして使用します。
    Goでは"call .X.Y 1 2"はdot.X.Y(1, 2)と書き、
    Yは関数フィールドかマップエントリーです。

    The first argument must be the result of an evaluation
    that yields a value of function type (as distinct from
    a predefined function such as print). The function must
    return either one or two result values, the second of which
    is of type error. If the arguments don't match the function
    or the returned error value is non-nil, execution stops.
html
    テキストで表された引数をエスケープしたHTMLを返します。
index
    一番目の引数の次の引数を添え字とした結果を消します。
    "index x 1 2 3"はGoでは x[1][2][3]となります。
    indexを使うアイテムは、マップかスライスか配列出なければなりません。
js
    テキストで表された引数をエスケープしたJavaScriptを返します。
len
    引数の長さ（整数）を返します。
not
    一つの引数の否定の真偽値を返します。
or
    Returns the boolean OR of its arguments by returning the
    first non-empty argument or the last argument, 
    "or x y"は"if x then x else y"のように振舞います。
    All the arguments are evaluated.
print
    fmt.Sprintのエイリアス
printf
    fmt.Sprintfのエイリアス
println
    fmt.Sprintlnのエイリアス
urlquery
    Returns the escaped value of the textual representation of
    its arguments in a form suitable for embedding in a URL query.
</pre>

<h2 id="FuncMap">type <a href="http://golang.org/src/pkg/text/template/funcs.go?s=612:647#L13">FuncMap</a></h2>
<pre>type FuncMap map[string]interface{}</pre>
<p>
FuncMapは型mapで、名前から関数へのマッピングを定義します。
関数には、1つの値を返すか、2つの値を返す必要があります。2つを返す方の2つ目は型errorとなります。
もし2つ目のerrorが実行中nilでなかったら、実行を中断し、Executeはerrorを返します。
</p>

<h2 id="Template">type <a href="http://golang.org/src/pkg/text/template/template.go?s=740:837#L16">Template</a></h2>
<pre>
type Template struct {
*parse.Tree
<span class="comment">// フィルタあるいはエクスポートされないフィールドを含んでいます。</span>
}
</pre>
<p>
Templateはパースされたテンプレートを表します。*parse.Treeフィールドは、
html/templateで使うためだけにエクスポートされています。
他の全てのクライアントでエクスポートされないように取り扱われます。
</p>

<span class="text"><a id="example_Template_glob" href="../../../src/pkg/text/template/exampletemplate.go">Example</a></span>

<span class="text"><a id="example_Template_func" href="../../../src/pkg/text/template/exampletemplate_func.go">Example (Func)</a></span>
<p> 
この例は、テンプレートテキストを処理するための、カスタム関数のデモです。
strings.Title関数を組み込みます。そしてそれをMake Title Text Look Goodのように出力するために使用します。
</p>

<span class="text"><a id="example_Template_glob" href="../../../src/pkg/text/template/exampletemplate_glob.go">Example (Glob)</a></span>
<p> ディレクトリにあるテンプレート郡をロードするデモです。</p>

<span class="text"><a id="example_Template_helpers" href="../../../src/pkg/text/template/exampletemplate_helpers.go">Example (Helpers)</a></span>
<p> 
この例は、いくつかのテンプレートを共有し、それらを違う文脈・コンテキスト使う方法の1つをデモします。
変わった形では、存在するテンプレート郡に渡すことで私たちは複数のドライバテンプレートを追加します。
</p>

<span class="text"><a id="example_Template_share" href="../../../src/pkg/text/template/exampletemplate_share.go">Example (Share)</a></span>
<p>
この例は、ヘルパーテンプレートを使ってドライバテンプレートグループの使い方をデモします。
</p>


<h3 id="Must">func <a href="../../../src/pkg/text/template/helper.go?s=576:619#L11">Must</a></h3>
<pre>func Must(t *Template, err error) *Template</pre>
<p>
Mustは、 (*Template, error)を返す関数をラップし、errorがnilだったらpanicするヘルパー関数です。
下記のように、変数を初期化するのに使用します。
</p>
<pre>
var t = template.Must(template.New(&#34;name&#34;).Parse(&#34;text&#34;))
</pre>

<h3 id="New">func <a href="../../../src/pkg/text/template/template.go?s=892:923#L25">New</a></h3>
<pre>func New(name string) *Template</pre>
<p>
Newは、与えられた名前で、新しく1つのテンプレートを割り当てます
</p>

<h3 id="ParseFiles">func <a href="http://golang.org/src/pkg/text/template/helper.go?s=970:1025#L22">ParseFiles</a></h3>
<pre class="go">func ParseFiles(filenames ...string) (*Template, error)</pre>
<p>
ParseFilesは新しいテンプレートを作成し、指定した名前のファイルからテンプレート定義を解析します。
戻り値のtemplateは、最初のファイルの(ベース)名前と(パースした)内容を格納します。
それらは、少なくとも1つファイルを指定する必要があります。
エラーが発生したら、解析を止め、nilの*Templateを返します。
</p>

<h3 id="ParsrGlob">func <a href="http://golang.org/src/pkg/text/template/helper.go?s=2861:2910#L75">ParseGlob</a></h3>
<pre class="go">
func ParseGlob(patter string) (*Template, error)
</pre>
<p>
ParseGlobはパターンで識別されたファイルから、新しいTemplateを作成し、テンプレート定義を解析します。
パターンは、少なくとも１つのファイルとマッチングしている必要があります。
戻り値のtemplateは、パターンにマッチした最初のファイルの(ベース)名前と(パースした)内容を格納します。
ParseGlobは、パターンにマッチしたファイルリストをもったParseFilesと同じです。
</p>


<h3 id="Template.ExecuteTemplate">func (*Template) <a href="http://golang.org/src/pkg/text/template/exec.go?s=2276:2361#L85">ExecuteTemplate</a></h3>
        <pre>func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error</pre>
        <p>
  ExecuteTemplate applies the template associated with t that has the given name
  to the specified data object and writes the output to wr.
</p>
