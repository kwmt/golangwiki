<h2 id="introduction">Introduction</h2>


<h2 id="formatting">Formatting</h2>

<h2 id="commentary">Commentary</h2>

<h2 id="name">Names</h2>

<h2 id="semicolons">Semicolons</h2>

<h2 id="control_structures">Control structures</h2>

<h2 id="functions">Functions</h2>

<h2 id="data">Data</h2>
<h3 id="allocation_new">Allocation with <code>new</code></h3>
<p>
    Goには２つの基本的なメモリ割り当て方法があります。
    ビルトイン関数の<code>new</code>と<code>make</code>です。
    これらは、異なった物で異なったタイプに割り当て、紛らわしいですが、ルールはシンプルです。
    まず、<code>new</code>について、説明します。
    これはメモリを確保するビルトイン関数ですが、いくつかの他の言語にある同じ名前とは違い、
    メモリを初期化しません。メモリをゼロにするだけです。
    <code>new(T)</code>は、型Tの新しい項目に対するゼロ化されたメモリを確保し、
    そのアドレス　型*T　を返します。
    Goでは、このことを、新しく確保したゼロ値の型Tのへのポインタを返す、と言います。
</p>
<p>
    <code>new</code>で返されたメモリがゼロになるので、
    さらに初期化することなしで、それぞれの型のゼロ値が使用されるデータ構造を設計するときに役に立ちます。
    つまり、ユーザーのデータ構造を、<code>new</code>を使って作成でき、すぐに作業に取り掛かれます。
    例えば、<code>bytes.Buffer</code>のドキュメントは、
    「<code>Buffer</code>のゼロ化は、空のバッファを準備します。」と述べてます。
    同様に、<code>sync.Mutex</code>は、明示的なコンストラクタうや
    <code>Init</code>メソッドを持っていません。
    その代わり、<code>sync.Mutex</code>に対してゼロ化が、アンロックなオブジェクトを定義します。
</p>
<p>
    ゼロ化が役に立つという特徴は、自動的に働きます。type宣言を考えてみてください。
</p>
<pre class="go">
type SyncedBuffer struct {
    lock    sync.Mutex
    buffer  bytes.Buffer
}
</pre>
<p>
    型<code>SyncedBuffer</code>の値は、割り当てられ、宣言した上にすぐに使う準備ができています。
    次では、<code>p</code>と<code>v</code>は、再度割り当てなくとも正しく動きます。
</p>
<pre class="go">
p := new(SyncedBuffer)  // type *SyncedBuffer
var v SyncedBuffer      // type  SyncedBuffer
</pre>

<h3 id="composite_literals">Constructors and composite literals</h3>

<p>
    ときどき、ゼロ初期化があまり良くない時があり、初期化するのに、コンストラクタが必要になる時があります。
    たとえば、<code>os</code>パッケージに見られます。
</p>
<pre class="go">
func NewFile(fd int, name string) *File {
    if fd &lt; 0 {
        return nil
    }
    f := new(File)
    f.fd = fd
    f.Name = name
    f.dirinfo = nil
    f.nepipe = 0
    return f
}
</pre>
<p>
    そこにはたくさんの例があります。私たちは"複合リテラル(composite literal)"を使って
    シンプルに書くことができます。これは評価されるたびに新しいインスタンスを作成する式です。
</p>
<pre class="go">
func NewFile(fd int, name string) *File {
    if fd &lt; 0 {
        return nil
    }
    f := File{fd, name, nil, 0}
    return &f
}
</pre>
<p>
    Note that, unlike in C, it's perfectly OK to return the address of a local variable;
    the storage associated with the variable survives after the function returns.
    Cとは異なり、

    In fact, taking the address of a composite literal allocates a fresh instance
    each time it is evaluated, so we can combine these last two lines.
    複合リテラルのアドレスを評価されるたびに更新するインスタンスを確保するので、
    最後の2行をまとめることができます。
</p>
<pre class="go">
return &amp;File{fd, name, nil, 0}
</pre>
<p>
    複合リテラルのフィールドを順番に配置し、すべてが存在しなければなりません。
    しかし、<i>field</i><code>:</code><i>value</i>のように、要素を明示的にラベルを付けることによって
    任意の順番で初期化され、記述してない残りの要素はゼロ初期化されます。
</p>
<pre class="go">
return &amp;File{fd: fd, name: name}
</pre>
<p>
    複合リテラルは、フィールドが全く含まれていない場合に制限するケースとして、ゼロの値の型を作成します。
    その時の式は、<code>new(File)</code>や<code>&amp;File{}</code>と等しいです。
</p>
<p>
    複合リテラルは、マップのキーやインデックスのようらフィールドラベルを持つ
    配列やスライス、マップに対しても作成することができます。
     In these examples, the initializations work regardless of the values of
     Enone, Eio, and Einval, as long as they are distinct.
</p>
<pre class="go">
a := [...]string    {Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
s := []string       {Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
m := map[int]string {Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
</pre>
<h3 id="allocation_make">Allocation with <code>make</code></h3>
<p>
    メモリの割り当てに戻ります。
    ビルトイン関数<code>make(T,<i>args</i>)</code>は、
    <code>new(T)</code>とは異なる目的を提供します。
    それは、スライス、マップ、チャネルのみを作成し、
    初期化された(ゼロ初期化ではなく)型T(*Tではない)を返します。
    その特徴の理由は、それら3つの型を使用する前には、
    必ず初期化しなければならないデータ構造だからです。
    たとえば、スライスは、3つの情報(データ(内部では配列)へのポインタ、長さ、キャパシティ)から構成されており、
    初期化されるまでは、スライスは<code>nil</code>となります。
    スライス、マップ、チャネルに対して、<code>make</code>は内部データ構造を初期化し、使う準備をします。
    例えば、
</p>
<pre class="go">
make([]int, 10, 100)
</pre>
<p>
    は、100個のint配列を確保し、長さ10で、配列の最初の10個の要素を指す100個のキャパシティを持つ
    スライス構造を作成します。
    (スライスを作るとき、キャパシティを省略することができます。詳細はスライスセクションにて。)
    対照的に、<code>new([]int)</code>は、新しく確保され、ゼロ化されたスライス構造へのポインタを返し、
    それはスライスの値<code>nil</code>へのポインタとなります。
</p>
<p>
    <code>new</code>と<code>make</code>の違いを例に説明します。
</p>
<pre class="go">
var p *[]int = new([]int)       // スライス構造 *p == nil を割り当てますが、めったに使用しません。
var v  []int = make([]int, 100) // スライスvは100個のint配列の一つを参照しています

// 不必要に複雑な例です
var p *[]int = new([]int)
*p = make([]int, 100, 100)

// 慣習的な方法
v := make([]int, 100)
</pre>
<p>
    <code>make</code>はマップ、スライス、チャネルだけに適用し、ポインタを返さないことを覚えてください。
    明示的にポインタを得るためには、<code>new</code>を使ってください。
</p>

<h2 id="initialization">Initialization</h2>

<h2 id="methods">Methods</h2>
<h3 id="pointer_vs_Values">Pointers vs. Values</h3>
<p>
    メソッドは、任意の名前付けされた型に対して定義することができます。
    その型は、ポインタでないか、インターフェースのどちらかで、レシーバーは構造体である必要ありません。
</p>
<p>
    上で行ったスライスの議論では、私たちは、<code>Append</code>関数を書きました。
    私たちは、代わりにスライス上のメソッドとして定義することができます。
    これを実現するには、まず、メソッドをバインドさせたい名前付けされた型を宣言します。
    それから、宣言した型のレシーバーに対してメソッドを作成します。
</p>
<pre class="go">
type ByteSlice []byte
func (slice ByteSlice) Append(data []byte) []byte {
    // 中身は上述の<code>Append</code>の中身と同じ
}
</pre>
<p>
    これではまだ、更新されたスライスを返す記述が必要になります。
    レシーバーとして<code>ByteSlice</code>へのポインタをとるメソッドに再定義することで、
    その不恰好な記述を取り去ることができます。
    ですので、呼び出すスライスをオーバーライドできます。
</p>
<pre class="go">
func (p *ByteSlice) Append(data []byte) {
    slice := *p
    // 中身は上述と同じ（ただし<code>return</code>文はありません)
    *p = slice
}
</pre>
<p>
we can do even better.
If we modify our function
so it looks like a standard Write method, like this,
</p>
<pre class="go">
func (p *ByteSlice) Write(data []byte) (n int, err error) {
    slice := p
    // 上部と同じです。
    *p = slice
    return len(data), nil
}
</pre>
<p>
    型<code>*ByteSlice</code>は標準インターフェース<code>io.Writer</code>を満たします。
    例えば、we can print into one.
</p>
<pre class="go">
    var b ByteSlice
    fmt.Fprintf(&b, "This hour has %d days¥n", 7)
</pre>
<p>
     <codeByteSlice</code>が<code>io.Writer</code>を満たすため、
     <code>ByteSlice</code>のアドレスを渡します。
     レシーバーのポインタvs値についてのルールは、
     値のメソッドは、ポインタと値上で呼び出すことができますが、
     ポインタのメソッドは、ポインタ上でしか呼び出すことができません。
     これは、ポインタメソッドがレシーバーを修正することができるからです。
     つまり、値のコピー上でそれらを呼び出すということは、
     それらの変更が破棄されることを引き起こします。
</p>
<p>
    ところで、バイトのスライス上で<code>Write</code>を使うアイデアは、
    <code>bytes.Buffer</code>で実装されています。
</p>

<h2 id="interfaces_and_other_types">Interfaces and other types</h2>
<h3 id="interfaces">Interfaces</h3>
<p>
    Goでのインターフェースは、オブジェクトの振る舞いを特定するための方法を提供します。
    if something can do this, then it can be used here.
    私たちはすでに、1組の簡単な例を見ています。
    <code>Fprintf</code>が<code>Write</code>メソッドを使って出力出来る間、
    カスタムプリントは、<code>String</code>メソッドによって実装されています。
    Goコードでは、1つや2つのメソッドだけを使ったインターフェースは共通です。
    and are usually given a name derived from the method,
    such as io.Writer for something that implements Write.
</p>
<p>
    1つの型は、複数のインターフェースを実装することができます。
    例えば、以下です。
    <code>Len()</code>、<code>Less(i, j int)</code>、<code>Swap(i, j int)</code>
    および、カスタムして出力することもできる<code>sort.Interface</code>を実装しているならば、
    <code>sort</code>パッケージで繰り返すことにより、コレクションを
    並び替えをすることができます。
    継承した例<code>Sequence</code>では、両方を満たしています。
</p>
<pre class="go">
type Sequence []int

// sort.Interfaceで必要なメソッド
func (s Sequence) Len() int {
    return len(s)
}
func (s Sequence) Less(i, j int) bool {
    return s[i] < s[j]
}
func (s Sequence) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

// 出力するためのメソッドです。 - 出力前に各要素をソートしています。
func (s Sequence) String() string {
    sort.Sort(s)
    str := "["
    fot i, elem := range s {
        if i > 0 {
            str += " "
        }
        str += fmt.Sprint(elem)
    }
    return str + "]"
}
</pre>
<h3 id="conversions">Conversions</h3>
<p>
    <code>Sequence</code>の<code>String</code>メソッドは、
    <code>Sprint</code>がスライスに対して実施する作業を、再作成しています。
    私たちが、<code>Sprint</code>をコールする前に、
    <code>Sequence</code>を普通の<code>[]int</code>に変えれば
    その努力を共有することができます。
</p>
<pre class="go">
func (s Seqeunce) String() string {
    sort.Sort(s)
    return fmt.Sprint([]int(s))
}
</pre>
<p>
    この変更は、<code>s</code>が通常のスライスとして扱われるため、デフォルトの書式を受け取ります。
    この変更をしなければ、<code>Sprint</code>は<code>Sequence</code>の<code>String</code>
    メソッドを見つけて、無限に繰り返されます。
    理由は、2つの型（<code>Sequence</code>と<code>[]int</code>)は、
    型名を無視すれば、同じだからです。
    それは、それら2つの間で変換するのに正当な方法です。
    この変更は、新しい値を作成せず、既存の値が新しい型を持っているかのように、一時的に動きます。
    （整数から浮動小数点に変換するような、新しく値を作る正当な方法はあります。）
</p>
<p>
    Goプログラミングでは、異なったメソッドにアクセスするために、型を変換することは、慣習的です。
    例として、既存の型である
    <a href="http://golang.org/pkg/sort/#IntSlice">
    <code>sort.IntSlice</code></a>を使うことができます:
</p>
<pre class="go">
type Sequence []int
func (s Sequence) String() string {
    sort.IntSlice(s).Sort()
    return fmt.Sprint([]int(s))
}
</pre>
<p>
    今や、<code>Sequence</code>を持つ複数のインターフェース(並び替えと出力)を実装する代わりに、
    任意のデータを複数の型(<code>Sequence,sort.IntSlice,[]int</code>)に
    変換する
</p>
<h2 id="embedding">Embedding</h2>

<h2 id="concurrency">Concurrency</h2>

<h2 id="errors">Errors</h2>
<p>
ライブラリールーチン（ライブラリーで提供されているルーチン）は、
多くの場合、呼び出し側にエラー表示のいくつかの並べ替えを返す必要があります。
前述したように、Goの戻り値を複数返せる特徴は、簡単に、
正常な戻り値と一緒にエラーの詳細な説明を返せます。
慣例では、errorsは組み込まれたシンプルなインターフェース<code>error</code>を持っています。
</p>
<pre class="go">
type error interface {
    Error() string
}
</pre>
<p>

</p>

<h2 id="a_web_server">A web server</h2>

