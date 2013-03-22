<h2 id="Parse">func <a href="http://golang.org/src/pkg/text/template/parse/parse.go?s=1186:1309#L25">Parse</a></h2>
<pre>func Parse(name, text, leftDelim, rightDelim string, funcs ...map[string]interface{}) (treeSet map[string]*Tree, err error)</pre>
<p>
Parseはテンプレートnameからparse.Treeへのマップを返します。
それは、引数の文字列で表されるテンプレートを解析することで作成されます。
トップレベルのテンプレートは、指定された名前を与えられます。
もしエラーが起こったら、解析をストップし、空のマップと一緒にerrorを返します。
</p>
