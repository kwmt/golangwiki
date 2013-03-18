<h2>Overview</h2>
<p>
templateパッケージは、テキスト出力するデータ駆動型のテンプレートを実装しています。
</p>
<p>
HTMLを生成するには、html/templateパッケージを参照ください。
それは、このパッケージと同じインターフェースを持ちますが、ある攻撃に対してセキュアなHTMLを自動的に生成します。
</p>
<p>
templateはデータ構造に適用することで実行されます。
このtemplateの記法は、実行をコントールし値を表示するには、データ構造の要素（一般には構造体のフィールドかマップのキー）を参照します。

Execution of the template walks the structure and sets the cursor, represented by a period '.' and called "dot",
to the value at the current location in the structure as execution proceeds.
templateの実行で、その構造を解析しピリオドやドットとよばれる'.'で表されるカーソルをセットし、実行の手続きとして構造にある現在の場所へ
</p>
<p>
templateに対する入力テキストは、どんなフォーマットにでもあうUTF-8エンコードされたテキストです。
データを評価したり、構造を制御したりする"アクション(Actions)"は、"{{"と"}}"で区切られます。
アクション外のすべてのテキストは、出力は変更されず、コピーされます。
アクションなコメントできますが、改行をまたぐことができません。
</p>
<p>
一度構造化されたテンプレートは並列に安全に実行されるかもしれません。
</p>
<p>
簡単な例ですが、"17 items are made of wool"を表示します。
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
    If the value is a map and the
    keys are of basic type with a defined order ("comparable"), the
    elements will be visited in sorted key order.

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
- A boolean, string, character, integer, floating-point, imaginary
  or complex constant in Go syntax. These behave like Go's untyped
  constants, although raw strings may not span newlines.
  Goの文法での真偽値、文字列、文字、整数、浮動小数点、虚数、複素数。
  raw文字列は改行をまたがることはできませんが、
  これらは、Goの型を持たない定数のように動作します。
- 文字は '.' (ピリオド):
    .
  The result is the value of dot.
  この結果はドットの値です。
- 変数名は
    $piOver2
  あるいは
    $
  のような$記号を前に置いた英数字の文字列です。
  The result is the value of the variable.
  Variables are described below.
- The name of a field of the data, which must be a struct, preceded
  by a period, such as
  構造体であるはずのデータのフィールド名は、
    .Field
  のように、ドットが前に置かれます。
  The result is the value of the field. Field invocations may be
  chained:
  フィールドをつなげることで、フィールドのフィールドをよびだせます:
    .Field1.Field2
  Fields can also be evaluated on variables, including chaining:
  変数からフィールドを評価することもできます:
    $x.Field1.Field2
- The name of a key of the data, which must be a map, preceded
  by a period, such as
    .Key
  マップデータのキーは
    .Key
  のようにピリオドを前に置きます。
  The result is the map element value indexed by the key.
  結果は、そのキーに対応した要素の値となります。
  Key invocations may be chained and combined with fields to any
  depth:
    .Field1.Key1.Field2.Key2
  このように連結することもできます:
    .Field1.Key1.Field2.Key2
  Although the key must be an alphanumeric identifier, unlike with
  field names they do not need to start with an upper case letter.
  キーは一意でなければなりませんが、フィールドとは違って大文字から始める必要はありません。
  Keys can also be evaluated on variables, including chaining:
    $x.key1.key2
- The name of a niladic method of the data, preceded by a period,
  such as
    .Method
  引数を持たないメソッド
    .Method
  The result is the value of invoking the method with dot as the
  receiver, dot.Method().
  レシーバーのようにドットを使ってメソッドを呼び出した結果です。
   Such a method must have one return value (of
  any type) or two return values, the second of which is an error.
  そのようなメソッドは、1つの戻り値あるいは2つの戻り値(2つ目はerror)をもつ必要があります。
  If it has two and the returned error is non-nil, execution terminates
  and an error is returned to the caller as the value of Execute.
  もし2つの戻り値を持って、errorがnilでなかったら、
  実行は中止し、実行の値として呼び出し元へエラーが返ってきます。
  Method invocations may be chained and combined with fields and keys
  to any depth:
  メソッドはフィールドとキーを組み合わせて連鎖して呼び出すことができます。
    .Field1.Key1.Method1.Field2.Key2.Method2
  Methods can also be evaluated on variables, including chaining:
    $x.Method1.Field
- The name of a niladic function, such as
    fun
  funのように引数をもたない関数です。
  The result is the value of invoking the function, fun().
  結果は関数fun()を呼び出した値です。
  The return
  types and values behave as in methods.
  返ってくる型と値は、メソッドのように振る舞います。
  Functions and function
  names are described below.
  関数と関数名については以下に記述します。
</pre>
<p>
    引数はどんな型も評価するかもしれません。
    もし引数がポインタなら、実装は自動的に必要なときに元になる肩を指します。
    If an evaluation yields a function value, such as a function-valued field of a struct,
    the function is not invoked automatically,
    but it can be used as a truth value for an if action and the like.
    To invoke it, use the call function, defined below.
</p>
<p>
A pipeline is a possibly chained sequence of "commands".
A command is a simple value (argument) or a function or method call, possibly with multiple arguments:
</p>
<pre>
Argument
    The result is the value of evaluating the argument.
.Method [Argument...]
    The method can be alone or the last element of a chain but,
    unlike methods in the middle of a chain, it can take arguments.
    The result is the value of calling the method with the
    arguments:
        dot.Method(Argument1, etc.)
functionName [Argument...]
    The result is the value of calling the function associated
    with the name:
        function(Argument1, etc.)
    Functions and function names are described below.
</pre>

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

実行関数が２つの関数マップで見つかっている間、このテンプレートの最初の関数はグローバル関数マップです。
デフォルトでは、テンプレートでは関数は定義されていませんが、Funcsメソッドを追加することができます。

</p>
<p>
事前に定義されたグローバル関数は以下のとおりに名前付けされています。
</p>
<pre>and
    Returns the boolean AND of its arguments by returning the
    first empty argument or the last argument, that is,
    "and x y" behaves as "if x then y else x". All the
    arguments are evaluated.
    "and x y"は"if x then y else x"のように振舞います。
call
    Returns the result of calling the first argument, which
    must be a function, with the remaining arguments as parameters.
    最初の関数であるはずの引数をコールした結果を返します。
    Thus 'call .X.Y 1 2' is, in Go notation, dot.X.Y(1, 2) where
    Y is a func-valued field, map entry, or the like.
    "call .X.Y 1 2"はGo記法では dot.X.Y(1, 2)と書き、
    Yは関数フィールドかマップエントリーです。

    The first argument must be the result of an evaluation
    that yields a value of function type (as distinct from
    a predefined function such as print). The function must
    return either one or two result values, the second of which
    is of type error. If the arguments don't match the function
    or the returned error value is non-nil, execution stops.
html
    Returns the escaped HTML equivalent of the textual
    representation of its arguments.
    テキストで表された引数をエスケープしたHTMLを返します。
index
    Returns the result of indexing its first argument by the
    following arguments. Thus 'index x 1 2 3' is, in Go syntax,
    x[1][2][3]. Each indexed item must be a map, slice, or array.
    一番目の引数の次の引数を添え字とした結果を消します。
    "index x 1 2 3"はGoでは x[1][2][3]となります。
    indexを使うアイテムは、マップかスライスか配列出なければなりません。
js
    Returns the escaped JavaScript equivalent of the textual
    representation of its arguments.
    テキストで表された引数をエスケープしたJavaScriptを返します。
len
    Returns the integer length of its argument.
    引数の長さ（整数）を返します。
not
    Returns the boolean negation of its single argument.
    一つの引数の否定の真偽値を返します。
or
    Returns the boolean OR of its arguments by returning the
    first non-empty argument or the last argument, that is,
    'or x y' behaves as 'if x then x else y'. All the
    arguments are evaluated.
    "or x y"は"if x then x else y"のように振舞います。
print
    An alias for fmt.Sprint
    fmt.Sprintのエイリアス
printf
    An alias for fmt.Sprintf
    fmt.Sprintfのエイリアス
println
    An alias for fmt.Sprintln
    fmt.Sprintlnのエイリアス
urlquery
    Returns the escaped value of the textual representation of
    its arguments in a form suitable for embedding in a URL query.
</pre>


<span class="text"><a id="example_Template_glob" href="../../../src/pkg/text/template/exampletemplate_glob.go">Example (Glob)</a></span>
<p> ディレクトリにあるテンプレート郡をロードするデモです。</p>


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
