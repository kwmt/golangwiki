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

Marshal <ha>    </ha>ndles an interface value bymarshalling the value it contains or,
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

<h2 id="MarshalIndent">func <a href="/src/pkg/encoding/xml/marshal.go?s=3243:3315#L71">MarshalIndent</a></h2>
<pre>func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)</pre>
<p>
MarshalIndent works like Marshal, but each XML element begins on a new
indented line that starts with prefix and is followed by one or more
copies of indent according to the nesting depth.
</p>

<p class="exampleHeading toggleButton">▹ <span class="text">Example</span></p>
<p>Code:</p>
<pre class="code">
type Address struct {
    City, State string
}
type Person struct {
    XMLName   xml.Name `xml:"person"`
    Id        int      `xml:"id,attr"`
    FirstName string   `xml:"name>first"`
    LastName  string   `xml:"name>last"`
    Age       int      `xml:"age"`
    Height    float32  `xml:"height,omitempty"`
    Married   bool
    Address
    Comment string `xml:",comment"`
}

v := &Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
v.Comment = " Need more details. "
v.Address = Address{"Hanga Roa", "Easter Island"}

output, err := xml.MarshalIndent(v, "  ", "    ")
if err != nil {
    fmt.Printf("error: %v\n", err)
}

os.Stdout.Write(output)
<span class="comment">
</pre>

<p>Output:</p>
<pre class="output">
<person id="13">
      <name>
          <first>John</first>
          <last>Doe</last>
      </name>
      <age>42</age>
      <Married>false</Married>
      <City>Hanga Roa</City>
      <State>Easter Island</State>
      <!-- Need more details. -->
</person>
</pre>

<h2 id="Unmarshal">func <a href="/src/pkg/encoding/xml/read.go?s=4824:4872#L103">Unmarshal</a></h2>
<pre>func Unmarshal(data []byte, v interface{}) error</pre>
<p>
Unmarshal parses the XML-encoded data and stores the result in
the value pointed to by v, which must be an arbitrary struct,
slice, or string.
Unmarshalは、XMLエンコードされたデータを解析し、vが指す値に結果を格納します。
Well-formed data that does not fit into v is discarded.
vに合わない整形式データは、破棄されます。
</p>
<p>
Because Unmarshal uses the reflect package, it can only assign
to exported (upper case) fields.
Unmarshalがreflectパッケージを使うので、エクスポートされたフィールド（大文字のフィールド）に割り当てるだけです。

Unmarshal uses a case-sensitive comparison to match XML element names
to tag values and struct field names.
Unmarshalは、XML要素名に一致するように大文字と小文字を区別する比較を使用しています
値をタグ付けし、フィールド名を構造体へ。
</p>
<p>
Unmarshal maps an XML element to a struct using the following rules.
Unmarshalは以下のルールを使って、XML要素を構造体へマップします。
In the rules, the tag of a field refers to the value associated with the
key 'xml'; in the struct field's tag (see the example above).
このルールでは、フィールドのタグは、構造体フィールドのタグにある'xml'キーに伴った値を参照します。(上の例を参照)
</p>
<pre>
* If the struct has a field of type []byte or string with tag
   ",innerxml", Unmarshal accumulates the raw XML nested inside the
   element in that field.  The rest of the rules still apply.
もし、構造体が".innerxml"タグの付いた[]byteやstring型のフィールドならば、
Unmarshal フィールド要素内にネストされたXMLを蓄えます。
残りのルールはまだ適用しています。

* If the struct has a field named XMLName of type xml.Name,
   Unmarshal records the element name in that field.
構造体が型xml.NameのXMLNameと名付けれたフィールドを持つなら、
Unmarshalはフィールドの要素名を記録します。

* If the XMLName field has an associated tag of the form
   "name" or "namespace-URL name", the XML element must have
   the given name (and, optionally, name space) or else Unmarshal
   returns an error.
XMLNameフィールドは、"name"や"namespace-URL name"の関連したタグを持つなら、
XML要素は、与えられた名前(optionally,名前空間)を持たなければなりません。
そうでなければ、Unmarshalはエラーを返します。

* If the XML element has an attribute whose name matches a
   struct field name with an associated tag containing ",attr" or
   the explicit name in a struct field tag of the form "name,attr",
   Unmarshal records the attribute value in that field.
XML要素に、構造体のフィールド名にマッチする属性　または　",attr"を含む関連したタグ　や
"name,attr"の構造体フィールドタグにある明示的な名前　を持つなら、
Unmarshalはそのフィールド内の属性の値を記録します

* If the XML element contains character data, that data is
   accumulated in the first struct field that has tag "chardata".
   The struct field may have type []byte or string.
   If there is no such field, the character data is discarded.
XML要素はキャラクターデータを含むなら、そのデータは"chardata"タグを持つ最初の構造体フィールドに蓄えられます。
構造体フィールドは型[]byteやstringを持つかもしれません。
そのようなフィールドがないなら、そのキャラクターデータは破棄されます。

* If the XML element contains comments, they are accumulated in
   the first struct field that has tag ",comments".  The struct
   field may have type []byte or string.  If there is no such
   field, the comments are discarded.
XML要素がコメントを含むなら、コメントは",comments"タグを持つ最初の構造体フィールドに蓄えられます。
構造体フィールドは型[]byteやstringを持つかもしれません。そのようなフィールドが無いなら、
そのコメントは破棄されます。

* If the XML element contains a sub-element whose name matches
   the prefix of a tag formatted as "a" or "a>b>c", unmarshal
   will descend into the XML structure looking for elements with the
   given names, and will map the innermost elements to that struct
   field. A tag starting with ">" is equivalent to one starting
   with the field name followed by ">".
XML要素に、"a"や"a>b>c"のようにフォーマットされたタグのプレフィックスにマッチしたサブ要素があるなら、
　


* If the XML element contains a sub-element whose name matches
   a struct field's XMLName tag and the struct field has no
   explicit name tag as per the previous rule, unmarshal maps
   the sub-element to that struct field.

* If the XML element contains a sub-element whose name matches a
   field without any mode flags (",attr", ",chardata", etc), Unmarshal
   maps the sub-element to that struct field.

* If the XML element contains a sub-element that hasn't matched any
   of the above rules and the struct has a field with tag ",any",
   unmarshal maps the sub-element to that struct field.

* A non-pointer anonymous struct field is handled as if the
   fields of its value were part of the outer struct.

* A struct field with tag "-" is never unmarshalled into.
</pre>
<p>
Unmarshal maps an XML element to a string or []byte by saving the
concatenation of that element&#39;s character data in the string or
[]byte. The saved []byte is never nil.
</p>
<p>
Unmarshal maps an attribute value to a string or []byte by saving
the value in the string or slice.
</p>
<p>
Unmarshal maps an XML element to a slice by extending the length of
the slice and mapping the element to the newly created value.
</p>
<p>
Unmarshal maps an XML element or attribute value to a bool by
setting it to the boolean value represented by the string.
</p>
<p>
Unmarshal maps an XML element or attribute value to an integer or
floating-point field by setting the field to the result of
interpreting the string value in decimal.  There is no check for
overflow.
</p>
<p>
Unmarshal maps an XML element to an xml.Name by recording the
element name.
</p>
<p>
Unmarshal maps an XML element to a pointer by setting the pointer
to a freshly allocated value and then mapping the element to that value.
</p>

