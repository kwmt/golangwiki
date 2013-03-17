##                             Slices:usage and internals
<p>
    Goのスライスタイプは型付けされたデータの配列での作業の便利で効率的な手段を提供します。
    スライスは他の言語での配列に似ていますが、珍しい（変わった）特性を持っています。
    この記事では、スライスとは何か、どうやってスライスを使うかを見ていきます。
</p>
## Array
<p>
    スライスタイプはGoの配列上に抽象的に構築されています。
    それで、スライスを理解するには、まず配列を理解する必要があります。
</p>
<p>
    配列の定義は長さと要素の型を指定します。
    たとえば、型<code>[4]int</code>は、4つのintegerの配列を表します。
    配列のサイズは固定されます。長さはその型の一部です
    (<code>[4]int</code>　や <code>[5]int</code>はまったく別物で、互換性はありません）。
    通常の方法では配列はインデックス化されるので、式<code>s[n]</code>は<code>n</code>番目の要素にアクセスできます。
</p>
<pre class="go">
var a [4]int
a[0] = 1
i := a[0]
//i == 1
</pre>
<p>
    配列を明示的に初期化する必要はありません。
    the zero value of an array is a ready-to-use array whose elements are themselves zeroed:
    配列は０で初期化されます。
</p>
<pre class="go">
// a[2] == 0, the zero value of the int type
</pre>
<p>
    <code>[4]int</code>メモリ内での表現は、ちょうど４つのintegerの値が順番にメモリに割り当てられます。
</p>
<p>
    <img src="http://golang.org/doc/articles/slice-array.png" alt="" />
</p>
<p>
    Goの配列は値です。配列変数は配列全体を示しています。
    先頭の配列要素へのポインタではありません（Cではそのようなケースになります）。
    This means that when you assign or pass around an array value you will make a copy of its contents.
    つまり、配列の値を割り当て渡すときに、その内容のコピーを作成します。
    （コピーではなくて、配列への<i>ポインタ</i>を渡すことができますが、配列へのポインタであって配列ではありません。）
    配列について考えるための１つの方法は、構造体のようなものとしてではなく、名前付けされたフィールドよりインデックスを
    使用することです。
</p>
<p>
    配列リテラルは次のように指定できます：
</p>
<pre class="go">
b := [2]string{"Penn","Teller"}
</pre>
<p>
    あるいは、配列要素数をコンパイラに数えさせることもできます:
</p>
<pre class="go">
b := [...]string{"Penn","Teller"}
</pre>
<p>
    上の両方のケースで、<code>b</code>の型は、<code>[2]string</code>となります。
</p>
##                             Slices
<p>
    配列はその場所を持っていますが、すこし融通がききませんので、
    Goコード内ではあまり頻繁には見られません。
    しかし、スライスはあらゆるところで見うけられます。
    スライスは、配列上に構築され、すばらしい力と利便性を提供します。
</p>
<p>
    スライスの型の記述は、<code>[]T</code>とかきます。<code>T</code>はスライスの要素の型です。
    配列とは違い、スライスは配列で言うところの長さ（サイズ）を指定しません。
</p>
<p>
    スライスリテラルは、ちょうど配列リテラルのように記述され、要素数の指定以外は同じになります。
</p>
<pre class="go">
letters := []string{"a", "b", "c", "d"}
</pre>
<p>
    スライスは、ビルトイン関数の<code>make</code>を使って作ることができます。
</p>
<pre class="go">
func make([]T, len, cap) []T
</pre>
<p>
    ただし、<code>T</code>は、作られるスライスの要素の型を表します。
    <code>make</code>関数は、type(型)とlenght(長さ)とオプションのcapacity(容量)を引数にとります。
    この<code>make</code>関数がコールされると、配列を割り当て、その配列を参照するスライスを返します。
</p>
<pre class="go">
var s []byte
s = make([]btye, 5, 5)
// s == []byte{0, 0, 0, 0, 0}
</pre>
<p>
    capacityの引数が省略されると、capacityの値は、指定されたlengthと同じになります。
    上記と同じコードの簡易版が次のようになります：
</p>
<pre class="go">
s := make([]byte, 5)
</pre>
<p>
    スライスのlengthとcapacityはビルトイン関数の
    <code>len</code>と<code>cap</code>関数を使って調べることが出来ます。
</p>
<pre class="go">
len(s) == 5
cap(s) == 5
</pre>
<p>
    次の2つのセクションで、lengthとcapacityの関係性を説明します。
</p>
<p>
    スライスのゼロを表す値は、<code>nil</code>です。スライスがnilだったら、
    <code>len</code>と<code>cap</code>関数は両方とも0を返します。
</p>
<p>
    スライスは、存在するスライスや配列を"slicing"することで、整形することができます。
    slicingは、2つの数字をコロンで分けます。
    たとえば、<code>b[1:4]</code>は、<code>b</code>の要素1から3を含むスライスを作成します。
    (作成されたスライスの結果の添字は、0から2となります。）
</p>
<pre class="go">
b := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
// b[1:4] == []byte{'o', 'l', 'a'}, bと同じ記憶領域を共有します。
</pre>
<p>
    The start and end indices of a slice expression are optional;
    they default to zero and the slice's length respectively:
    スライスの最初と最後の数字は、オプションです。下の例を見た方が分かりやすい。（訳注：訳が分からなかったわけでは・・・）
</p>
<pre class="go">
// b[:2] == []byte{'g', 'o'}
// b[2:] == []byte{'l', 'a', 'n', 'g'}
// b[:] == b
</pre>
<p>
    下記の構文も、与えられた配列からスライスを作る構文です。
</p>
<pre class="go">
x := [3]string{"Лайка", "Белка", "Стрелка"}
s := x[:] // スライスsはxの記憶領域を参照します。
</pre>
##                             Slice Internals
<p>
    スライスは配列内の連続した領域への参照であり、
    配列へのポインタ、セグメントの長さ、その容量（セグメントの最大の長さ）から構成されます。
</p>
<p>
    <img src="http://golang.org/doc/articles/slice-struct.png" alt="" />
</p>
<p>
    <code>make([]byte, 5)</code>から作られた変数<code>s</code>は、このように構造化されています：
</p>
<p>
    <img src="http://golang.org/doc/articles/slice-1.png" alt="" />
</p>
<p>
    lengthはスライスに参照される要素数です。
    capacityは配列の基盤（スライスのポインタによって参照される要素の最初）となる要素の数です。
    lengthとcapacityとの区別は、次の少しの例で明確になるはずです。
</p>
<p>
    As we slice s,
    observe the changes in the slice data structure and their relation to the underlying array:
    先程スライス<code>s</code>を作ったように、スライスデータ構造の変化と基になる配列の関係を観察します：
</p>
<pre class="go">
s = s[2:4]
</pre>
<p>
    <img src="http://golang.org/doc/articles/slice-2.png" alt="" />
</p>
<p>
    slicingはスライスのデータをコピーするわけではなく、元の配列を指す新しいスライスを作成します。
    これは、配列のインデックスを操作するように効率的にスライス操作を行います。
    したがって、再スライスされた要素を変更すると、元のスライスの要素は変更されます：
</p>
<pre class="go">
d := []byte{'r', 'o', 'a', 'd'}
e := d[2:]
// e == []byte{'a', 'd'}
e[1] == 'm'
// e == []byte{'a', 'm'}
// d := []byte{'r', 'o', 'a', 'm'}
</pre>
<p>
    さきほどのスライスした<code>s</code>は、
    capacityよりも短いlengthにスライスしました。
    それを再度slicingすることで、
    lengthをcapacityまで<code>s</code>を拡張することができます:
</p>
<pre class="go">
s = s[:cap(s)]
</pre>
<p>
    <img src="http://golang.org/doc/articles/slice-3.png" alt="" />
</p>
<p>
    スライスはそのcapacityを超えて伸ばすことはできません。
    スライスや配列の上限を超えてインデックス作成するのと同じように、
    capacityを超えて伸ばそうとすると、ランタイムパニックを起こします。
    同様に、スライスは最初の配列内の要素にアクセスするために０未満で再スライスすることはできません。
</p>
##                             スライスを拡張する(the copy and append functions)
<p>
    スライスのcapacityを増やすには、
    新しく、より大きなスライスを作成し、その作成したスライスに元のスライスの内容をコピーしなければなりません。
    このテクニックは他の言語から動的配列の実装では、舞台裏で動く方法です。
    次の例は、新しくスライス<code>t</code>を作ることで、<code>s</code>のcapacityを2倍にします。
    <code>s</code>の内容を<code>t</code>にコピーし、スライスの値<code>t</code>を<code>s</code>
    に割り当てます:
</p>
<pre class="go">
t := make([]byte, len(s), (cap(s)+1)*2) // cap(s)==0 の場合は +1
for i := range s {
    t[i] = s[i]
}
s = t
</pre>
<p>
    訳注：ここから<br />
    これを実際にやってみたところ次のようになりました。
</p>
<pre class="go">
s := make([]byte, 5)
fmt.Println("s=", s, "len(s)=", len(s), "cap(s)=", cap(s))

t := make([]byte, len(s), (cap(s)+1)*2) // cap(s)==0 の場合は +1
for i := range s {
    t[i] = s[i]
}
s = t
fmt.Println("t=", t, "len(t)=", len(t), "cap(t)=", cap(t))
fmt.Println("s=", s, "len(s)=", len(s), "cap(s)=", cap(s))

// 出力
// s= [0 0 0 0 0] len(s)= 5 cap(s)= 5
// t= [0 0 0 0 0] len(t)= 5 cap(t)= 12
// s= [0 0 0 0 0] len(s)= 5 cap(s)= 12
</pre>
<p>
    訳注：ここまで
</p>
<p>
    この共通操作のループ部分は、ビルトインの<code>copy</code>関数で簡単につくれます。
    名前が示すように、<code>copy</code>は元のスライスから先のスライスにデータをコピーします。
    そして、コピーされた要素の数を返します。
</p>
<pre class="go">
func copy(dst, src []T) int
</pre>
<p>
    <code>copy</code>関数は、異なるlengthのスライス間のコピーもサポートしています。（小さい方の要素の数分コピーします）
     In addition, copy can handle source and destination slices
     that share the same underlying array, handling overlapping slices correctly.
</p>
<p>
    <code>copy</code>を使って、以下のようにシンプルになります：
</p>
<pre class="go">
t := make([]byte, len(s), (cap(s)+1)*2)
copy(t, s)
s = t
</pre>
<p>
    一般的な操作は、スライスの最後にデータを追加することです。
    下記の関数は、byte要素をbyteのスライスに加えます。必要なスライスを拡張させ、
    更新されたスライス値を返します：
</p>
<pre class="go">
func AppendByte(slice []byte, data ...byte) []byte {
    m := len(slice)
    n := m + len(data)
    if n &gt; cap(slice) { // 必要なら、割り当てし直します
        // 今後のために2倍に拡張します。
        newSlice := make([]byte, (n+1)*2)
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[0:n]
    copy(slice[m:n], data)
    return slice
}
</pre>
<p>
    <code>AppendByte</code>関数をこの様に使うこと事ができます：
</p>
<pre class="go">
p := []byte{2, 3, 5}
p = AppendByte(p, 7, 11, 13)
// p == []byte{2, 3, 5, 7, 11, 13}
</pre>
<p>
    Functions like AppendByte are useful
    because they offer complete control over the way the slice is grown.
    Depending on the characteristics of the program,
    it may be desirable to allocate in smaller or larger chunks,
    or to put a ceiling on the size of a reallocation. <br />

    スライスを拡張する以上の制御をしばしばするので、<code>AppendByte</code>のような関数は役に立ちます。
    プログラムの特性に応じて、小さいかたまりか大きい固まりに配分するのが良いかもしれないし、
    再配分したサイズの上限に置くのが良いかもしれない
</p>
<p>
    多くのプログラミングは、完璧な制御は必要としませんので、
    Goはビルトイン<code>append</code>関数を提供しています。
</p>
<pre class="go">
func append(s []T, x ...T) []T
</pre>
<p>
    <code>append</code>関数は要素<code>x</code>をスライス<code>s</code>の末尾に追加し、
    より大きなcapacityが必要ならスライスを拡張します。
</p>
<pre class="go">
a := make([]int, 1)
// a == []int{0}
a = append(a, 1, 2, 3)
// a == []int{0, 1, 2, 3}
</pre>
<p>
    あるスライスに別のスライスを追加するには、<code>...</code>を使い、
    2つ目の引数をリストに拡張します。
</p>
<pre class="go">
a := []string{"John", "Paul"}
b := []string{"George", "Ringo", "Pete"}
a = append(a, b...) // "append(a, b[0], b[1], b[2])"と同じことです。
// a == []string{"John", "Paul", "George", "Ringo", "Pete"}
</pre>
<p>
    スライスのゼロの値(nil)は、ゼロの長さのスライスのように動作するので、
    スライス変数を宣言してからループ内でそれに追加することができます。
</p>
<pre class="go">
// Filterは関数f()を満たすsの要素だけを
//　保持する新しいスライスを返します。
func Filter(s []int, fn func(int) bool) []int {
    var p []int // == nil
    for _, i := range s {
        if fn(i) {
            p = append(p, i)
        }
    }
    return p
}
</pre>
##                             A possible "gotcha"
<p>
    前述したように、スライスを再度スライスすると、
    元になる配列をコピーしません。
    配列は、それが参照されなくなるまで、メモリ内に保持し続けます。
    Occasionally this can cause the program to hold all the data
    in memory when only a small piece of it is needed.
</p>
<p>
    たとえば、<code>FindDigits</code>関数は、ファイルをメモリにロードして、
    連続した数値の最初のグループを検索し、新しいスライスとしてそれらを返します。
</p>
<pre class="go">
var digitRegexp = regexp.MustCompile("[0-9]+")

func FindDigits(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    return digitRegexp.Find(b)
}
</pre>
<p>
    このコードは、記述通りに動きますが、
    ファイル全体のスライス<code>[]byte</code>を返しています。
    スライスは元の配列を参照しますので、ガベッジコレクター上でスライスが保持される限り、
    配列(メモリにファイル全体の内容を保持するbyte)を開放することができません。
</p>
<p>
    この問題を修正するための１つは、スライスを返す前に、
    新しいスライスにデータをコピーすることです：
</p>
<pre class="go">
var digitRegexp = regexp.MustCompile("[0-9]+")

func CopyDigits(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    b = digitRegexp.Find(b)
    c := make([]byte, len(b))
    copy(c, b)
    return c
}
</pre>
<p>
    この関数のもっと簡潔なバージョンを、<code>append</code>を使って構成することができます。
    これは読者の演習として残しておきます。
</p>
##                             Futher Reading
<p>
    <a href="http://golang.org/doc/effective_go.html">Effective Go</a>には、
    <a href="http://golang.org/doc/effective_go.html#slices">スライス</a>と
    <a href="http://golang.org/doc/effective_go.html#arrays">配列</a>の
    詳細な内容が書かれており、
    Go <a href="http://golang.org/doc/go_spec.html">言語の詳細</a>は、
    <a href="http://golang.org/doc/go_spec.html#Slice_types">スライス</a>を定義し、
    スライスに<a href="http://golang.org/doc/go_spec.html#Length_and_capacity">関連した</a>
    <a href="http://golang.org/doc/go_spec.html#Making_slices_maps_and_channels">ヘルパー</a>
    <a href="http://golang.org/doc/go_spec.html#Appending_and_copying_slices">機能</a>を提供します。
</p>
