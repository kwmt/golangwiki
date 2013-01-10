<h2>Overview</h2>
<p>
xmlパッケージは、XML名前空間をし解釈できるXML1.0のパーサを実装しています。
</p>
<h2 id="pkg-variables">Variables</h2>

<pre>var HTMLAutoClose = htmlAutoClose</pre>
<p>
HTMLAutoClassは、機会的に閉じることが考慮できるHTML要素の集まりです。
</p>

<pre>var HTMLEntity = htmlEntity</pre>
<p>
HTMLEntityは、標準のHTMLエンティティに変換するエンティティマップです。
</p>

<h2 id="Escape">func <a href="http://golang.org/src/pkg/encoding/xml/xml.go?s=37968:38002#L1641">Escape</a></h2>
<pre>func Escape(w io.Writer, s []byte)</pre>
<p>
Escapeは、平文テキストデータsと同じ、適切にエスケープされたXMLをwに出力します。
</p>

<h2 id="Marshal">func <a href="/src/pkg/encoding/xml/marshal.go?s=2882:2925#L60">Marshal</a></h2>
<pre>func Marshal(v interface{}) ([]byte, error)</pre>
<p>
Marshal returns the XML encoding of v.
MarshalはvのXMLエンコードを返します。
</p>
<p>
Marshal handles an array or slice by marshalling each of the elements.
Marshalは、それぞれの要素を適正に並べることで、配列またはスライスを処理します。

Marshal handles a pointer by marshalling the value it points at or, if the
pointer is nil, by writing nothing.

Marshal handles an interface value bymarshalling the value it contains or,
if the interface value is nil, by　writing nothing.

Marshal handles all other data by writing one or more XML elements containing the data.

</p>
<p>
The name for the XML elements is taken from, in order of preference:
XML要素の名前は優先順に付けられる。
</p>
<pre>
- the tag on the XMLName field, if the data is a struct
  データが構造体なら、XMLNameフィールド上のタグ
- the value of the XMLName field of type xml.Name
  型xml.NameのXMLNameフィールドの値
- the tag of the struct field used to obtain the data
  データを取得するために使用される構造体フィールドのタグ
- the name of the struct field used to obtain the data
  データを取得するために使用される構造体フィールドの名前
- the name of the marshalled type
  マーシャルされた型の名前
</pre>
<p>
The XML element for a struct contains marshalled elements for each of the
exported fields of the struct, with these exceptions:
構造体のXML要素は、構造体の各エクスポートフィールドに対して、
例外があります：
</p>
<pre>
- the XMLName field, described above, is omitted.
上記で説明されるXMLNameフィールドは除外されます。
- a field with tag &#34;-&#34; is omitted.
タグ "-"のフィールドは除外されます。
- a field with tag &#34;name,attr&#34; becomes an attribute with
  the given name in the XML element.
  タグ"name,attr"のフィールドは、XML要素の中で与えられた名前の属性になります。
- a field with tag &#34;,attr&#34; becomes an attribute with the
  field name in the in the XML element.
  タグ",attr"のフィールドは、XML要素の中のフィールド名の属性となります。
- a field with tag &#34;,chardata&#34; is written as character data,
  not as an XML element.
  タグ",chardata"を持つフィールドは文字データとして書かれ、XML要素としてではありません。
- a field with tag &#34;,innerxml&#34; is written verbatim, not subject
  to the usual marshalling procedure.
  タグ",innerxml"を持つフィールドは、文字通りに書かれます。通常の整列化手続きはとりません。
- a field with tag &#34;,comment&#34; is written as an XML comment, not
  subject to the usual marshalling procedure. It must not contain
  the &#34;--&#34; string within it.
  タグ",comment"を持つフィールドは、XMLコメントして書かれます。
  通常の整列化手続きはとりません。
  その中に文字列"--"を含める必要はありません。
- a field with a tag including the &#34;omitempty&#34; option is omitted
  if the field value is empty. The empty values are false, 0, any
  nil pointer or interface value, and any array, slice, map, or
  string of length zero.
  フィールドの値が空ならば、"omitempty"オプションを含むタグをもつフィールドは、除外されます。
  空の値は、false,0,nilポインタ,インターフェース値と配列,スライス,マップの長さが０のことです。

- a non-pointer anonymous struct field is handled as if the
  fields of its value were part of the outer struct.
  ポインタではない無名構造体フィールドは、そのフィールドの値が外側の構造体の一部であれば、使用できます。
</pre>
<p>
If a field uses a tag &#34;a&gt;b&gt;c&#34;, then the element c will be nested inside
parent elements a and b.
もしフィールドがタグ"a>b>c"を使っていたら、その要素cは親要素aとbの内側にネストされます。

 Fields that appear next to each other that name
the same parent will be enclosed in one XML element.
同じ名前の親要素が現れたら、そのひとつの親要素のなかに閉じられます。
</p>
<p>
See MarshalIndent for an example.
サンプルはMarshalIndentを参照。
</p>
<p>
Marshal will return an error if asked to marshal a channel, function, or map.
Marshalは、チャネル、関数、またはマップをmarshalするとき、エラーを返します。
</p>
