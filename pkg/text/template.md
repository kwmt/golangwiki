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
    A with action using dot.
    ドットを使ってwithアクション
{{with $x := "output" | printf "%q"}}{{$x}}{{end}}
    変数を定義して使ったwithアクション
{{with $x := "output"}}{{printf "%q" $x}}{{end}}
    A with action that uses the variable in another action.
    他のアクションの中で変数を使うwithアクション
{{with $x := "output"}}{{$x | printf "%q"}}{{end}}
    上と同じだが、パイプラインを使用
</pre>
