https://code.google.com/p/go/source/browse/doc/go1.2.html 
23fc3139589c

https://code.google.com/p/go/source/browse/doc/go1.2.html?spec=svn5298d82f682fecf71a7570c6625b2d80880749ff&r=23fc3139589c0ba909c155153897159ca1972c20

<link rel="stylesheet" type="text/css" href="../css/main.css">

<h2 id="introduction">Introduction to Go 1.2</h2>

2013年4月に<a href="http://golang.org/doc/go1.1.html">Go version 1.1</a>をリリースして以来、リリーススケジュールが
リリースのプロセスをより効率的にするために、短くされていました。
このリリースGo バージョン1.2 （短縮表記すると、Go1.2）では、だいたい1.1以降6ヶ月になります。
1.1は1.0がリリースされてから1年以上かかりました。
タイムスケールが短くなった理由は、1.2は1.0から1.1の段階より差分が少ないからです。
しかし、よりよいスケジューラと1つの新しい言語の特徴を含む、重要な変更があります。

もちろんGo1.2は<a href="http://golang.org/doc/go1compat.html">promise
of compatibility</a>を守っています。

Go1.1で作ったほとんどのプログラムは、1.2へ移行したとしても変更なしで実行するでしょう。
コーナーケース（めったに発生しなやっかいなケース）に1つの制限の導入は、すでに誤ったコードを公開するかもしれませんが。
(<a href="#use_of_nil">use of nil</a>を参照)





<h2 id="language">Changes to the language</h2>
<p>
仕様を安定させるために、プログラムに対して重大な一つのコーナーケースが明らかになりました。
これは、一つの新しい言語の特徴です。
</p>

<h3 id="use_of_nil">Use of nil</h3>

<p>
現在、安全上の理由から、nilポインタの特定の用途がランタイムパニックを誘発することが保証されている、ことを規定しています。
例えば、Go1.0で、次のようなコードを与えられた場合
</p>
<pre>
type T struct {
    X [1&lt;&lt;24]byte
    Field int32
}

func main() {
    var x *T
    ...
}
</pre>

<p>
<code>nil</code>ポインタ<code>x</code>は間違ったメモリにアクセスできます:
<code>x.Field</code>はアドレス<code>1&lt;&lt;24</code>のメモリにアクセスできます。
このような安全でない振る舞いを防ぐために、Go1.2では、
配列へのnilポインタ、nilインターフェースバリュー、nilスライスなどのnilポインタを通しての関節参照がpanicになったり、
正しい安全な非nil値を返すことをコンパイラが保証します。
端的に言うと、明示的または暗黙的にnilをアドレスの評価を必要とする任意の式はエラーになります。
実装は、この仕様を施行するためにコンパイルされたプログラムの中に追加のテストを入れるかもしれません。
</p>

<p>
詳しい情報は<a href="http://golang.org/s/go12nil">design document</a>を参照してください。
</p>

<p>
<em>Updating</em>:
古い仕様に依存するほとんどのコードはエラーとなり、実行すると失敗します。
そのようなプログラムは、手動で更新する必要があります。
</p>

<h3 id="three_index">Three-index slices</h3>

<p>
Go1.2では、ある配列やスライスを扱うとき、容量だけななく長さも指定できるようにしました。
</p>

<pre>
var array [10]int
slice := array[2:4]
</pre>

<p>
スライスの容量は、スライスが保持している、再スライスした後の要素の最大数です。
もとの配列のサイズを反映します。
この例では、<code>slice</code>変数の容量は8です。
</p>


<p>
Go1.2では長さだけでなく容量を指定することができるスライス操作の新しい文法を追加しました。
2番目のコロンは容量の値を指定します。その値は、もともとのスライスや配列の容量と等しいかそれ以下であるはずです。例えば、
</p>

<pre>
slice = array[2:4:7]
</pre>


<p>
は、最初の例のように同じ長さを持つスライスですが、容量は5(7-2)だけになります。
元の配列の3つ目の要素にアクセスするために、新しいスライスの値を使うことは出来ません。(TODO:要見直し)

</p>


<p>
このthree-indexの記法では、最初のインデックスを書かなかれば(<code>[:i:j]</code>)０になりますが、他の２つのインデックスはいつもちゃんと書かなくてはいけません。
Goの将来のリリースでは、それらのインデックスにデフォルト値を導入することになるかもしれません。

<p>
詳しくはこちら
<a href="http://golang.org/s/go12slice">design document</a>.
</p>

<p>
<em>Updating</em>:
これは、既存のプログラムに影響しない下位互換性の変更です。
</p>

<h2 id="impl">Changes to the implementations and tools</h2>

<h3 id="preemption">Pre-emption in the scheduler</h3>


<p>
 以前のリリースでは、永遠にループしているゴルーチンが、同じスレッド上の他のゴルーチンを餓死させてしまうことがありました。
GOMAXPROCSが1つのユーザースレッドのみの場合は重要な問題です。
Go1.2では、部分的に処理されます：スケジューラが関数に入るとき関数に入る際に時々起動されます。
つまり、(組み込まれていない)関数を含むループが他のゴルーチンを同じスレッド上で実行することができるようになります。</span>
</p>
	

<h3 id="thread_limit">Limit on the number of threads</h3>

<p>
Go 1.2 introduces a configurable limit (default 10,000) to the total number of threads
a single program may have in its address space, to avoid resource starvation
issues in some environments.
Note that goroutines are multiplexed onto threads so this limit does not directly
limit the number of goroutines, only the number that may be simultaneously blocked
in a system call.
In practice, the limit is hard to reach.
</p>

<p>
Go1.2は、いくつかの環境でリソース枯渇する問題を避けるため、一つのプログラムがそのアドレス空間で持つスレッドの合計に上限(デフォルト 10,000)の設定を変更可能にしました。
ゴルーチンはそれぞれのスレッド上で送信されるので、この上限は、直接的なゴルーチンの数の上限ではないことに注意して下さい。
システムコールでブロックされるかもしれない数だけです。
実際は、この上限に到達することは難しいです。
</p>

<p>
The new <a href="http://golang.org/pkg/runtime/debug/#SetMaxThreads"><code>SetMaxThreads</code></a> function in the
<a href="http://golang.org/pkg/runtime/debug/"><code>runtime/debug</code></a> package controls the thread count limit.
</p>

<p>
<a href="http://golang.org/pkg/runtime/debug/"><code>runtime/debug</code></a>パッケージに新しい<a href="http://golang.org/pkg/runtime/debug/#SetMaxThreads"><code>SetMaxThreads</code></a>関数は、スレッドの上限値をコントロールします。
</p>



<p>
<em>Updating</em>:
Few functions will be affected by the limit, but if a program dies because it hits the
limit, it could be modified to call <code>SetMaxThreads</code> to set a higher count.
Even better would be to refactor the program to need fewer threads, reducing consumption
of kernel resources.
</p>

<p>
 いくつかの関数は制限の影響を受けますが、その上限に達するからプログラムが終了するような場合、より高い上限値を設定するために<code>SetMaxThreads</code>を呼び出す修正ができます。
</p>

<h3 id="stack_size">Stack size</h3>

<p>
In Go 1.2, the minimum size of the stack when a goroutine is created has been lifted from 4KB to 8KB.
Many programs were suffering performance problems with the old size, which had a tendency
to introduce expensive stack-segment switching in performance-critical sections.
The new number was determined by empirical testing.
</p>

<p>
Go1.2では、ゴルーチンが作られるときのスタックの最小サイズが、4KBから8KBに変更されました。
多くのプログラムは、旧サイズではパフォーマンス問題を抱えていました。
これは、パフォーマンスが重要なセクションで、高価なスタック·セグメントを変更することを導入する傾向がありました。
この数値は、テスト実験によって決定されました。
</p>

<p>
At the other end, the new function <a href="http://golang.org/pkg/runtime/debug/#SetMaxStack"><code>SetMaxStack</code></a>
in the <a href="http://golang.org/pkg/runtime/debug"><code>runtime/debug</code></a> package controls
the <em>maximum</em> size of a single goroutine's stack.
The default is 1GB on 64-bit systems and 250MB on 32-bit systems.
Before Go 1.2, it was too easy for a runaway recursion to consume all the memory on a machine.
</p>

<p>
他方で、<a href="http://golang.org/pkg/runtime/debug"><code>runtime/debug</code></a> パッケージにある新しい関数<a href="http://golang.org/pkg/runtime/debug/#SetMaxStack"><code>SetMaxStack</code></a>関数は、１つのゴルーチンスタックの<em>最大</em>サイズをコントロールします。
デフォルトは、64bitシステムでは1GB、32bitシステムでは250MBです。
Go1.2以前は、マシンのすべてのメモリを消費しやすかった。
</p>

<p>
<em>Updating</em>:
The increased minimum stack size may cause programs with many goroutines to use
more memory. There is no workaround, but plans for future releases
include new stack management technology that should address the problem better.
</p>

<p>
増加した最小スタックサイズは、多くのメモリを使用するため、多くのゴルーチンでは問題を引き起こすかもしれません。
回避策はありませんが、将来のリリースの計画には、問題に対処するべき新しいスタック管理技術が含まれます。
</p>

<h3 id="cgo_and_cpp">Cgo and C++</h3>

<p>
The <a href="http://golang.org/cmd/cgo/"><code>cgo</code></a> command will now invoke the C++
compiler to build any pieces of the linked-to library that are written in C++;
<a href="http://golang.org/cmd/cgo/">the documentation</a> has more detail.
</p>

<p>
<a href="http://golang.org/cmd/cgo/"><code>cgo</code></a>のコマンドは、C++で書かれているライブラリにリンクされたいずれかの部分をビルドするためにC++コンパイラを起動します。
<a href="http://golang.org/cmd/cgo/">the documentation</a> に詳細がありあます。
</p>

<h3 id="go_tools_godoc">Godoc and vet moved to the go.tools subrepository Godocとvetコマンドがサブリポジトリgo.toolsに移動しました。</h3>

<p>
Both binaries are still included with the distribution, but the source code for the
godoc and vet commands has moved to the
<a href="http://code.google.com/p/go.tools">go.tools</a> subrepository.
</p>

<p>
両方のコマンドはまだディストリビューションにありますが、godocとvetコマンドのソースコードは、サブリポジトリ<a href="http://code.google.com/p/go.tools">go.tools</a>に移動しました。
</p>

<p>
Also, the core of the godoc program has been split into a
<a href="https://code.google.com/p/go/source/browse/?repo=tools#hg%2Fgodoc">library</a>,
while the command itself is in a separate
<a href="https://code.google.com/p/go/source/browse/?repo=tools#hg%2Fcmd%2Fgodoc">directory</a>.
The move allows the code to be updated easily and the separation into a library and command
makes it easier to construct custom binaries for local sites and different deployment methods.
</p>

<p>
また、godocプログラムのコアは、<a href="https://code.google.com/p/go/source/browse/?repo=tools#hg%2Fgodoc">ライブラリ</a>にあり、コマンド自体は、<a href="https://code.google.com/p/go/source/browse/?repo=tools#hg%2Fcmd%2Fgodoc">ディレクトリ</a>にわかれています。
移動したのは、コードを簡単に更新できるようにし、ライブラリとコマンドに分けたのは、ローカルサイトや異なる開発方法に対して、カスタムコマンドを構築しやすくするためです。
</p>

<p>
<em>Updating</em>:
Since godoc and vet are not part of the library,
no client Go code depends on the their source and no updating is required.
</p>

<p>
godocとvetはライブラリの一部ではありませんので、Goコードのクライアントはこれらのソースに依存しません、更新も必要ありません。
</p>

<p>
The binary distributions available from <a href="http://golang.org">golang.org</a>
include these binaries, so users of these distributions are unaffected.
</p>

<p>
 <a href="http://golang.org">golang.org</a>から使用できるバイナリディストリビューションには、これらのバイナリがありますので、ユーザーは影響は受けません。
</p>

<p>
When building from source, users must use "go get" to install godoc and vet.
(The binaries will continue to be installed in their usual locations, not
<code>$GOPATH/bin</code>.)
</p>

<p>
ソースからビルドする場合、ユーザーは"go get"コマンドを使ってgodocとvetをインストールする必要があります。
(これらのバイナリは、いつもの場所にインストールされます。<code>$GOPATH/bin</code>ではありません。)
</p>

<pre>
$ go get code.google.com/p/go.tools/cmd/godoc
$ go get code.google.com/p/go.tools/cmd/vet
</pre>

<h3 id="gccgo">Status of gccgo</h3>

<p>
We expect the future GCC 4.9 release to include gccgo with full
support for Go 1.2.
In the current (4.8.2) release of GCC, gccgo implements Go 1.1.2.
</p>

<p>
私達は、将来のGCC4.9のリリースに、Go1.2のフルサポートがついたgccgoが入ることを期待しています。
GCCの現在のリリース(4.8.2)では、gccgoはGo1.1.2を実装しています。
</p>

<h3 id="gc_changes">Changes to the gc compiler and linker</h3>

<p>
Go 1.2 has several semantic changes to the workings of the gc compiler suite.
Most users will be unaffected by them.
</p>

<p>
 Go1.2は、gcコンパイラの動きに関していくつかセマンティックな変更をしています。
 ほとんどのユーザーは、それらに影響を受けないでしょう。
</p>

<p>
The <a href="http://golang.org/cmd/cgo/"><code>cgo</code></a> command now
works when C++ is included in the library being linked against.
See the <a href="http://golang.org/cmd/cgo/"><code>cgo</code></a> documentation
for details.
</p>

<p>
C++がリンクされているライブラリに含まれているとき、<a href="http://golang.org/cmd/cgo/"><code>cgo</code></a>コマンドが動きます。
詳細は<a href="http://golang.org/cmd/cgo/"><code>cgo</code></a>を参照下さい。
</p>

<p>
The gc compiler displayed a vestigial detail of its origins when
a program had no <code>package</code> clause: it assumed
the file was in package <code>main</code>.
The past has been erased, and a missing <code>package</code> clause
is now an error.
</p>

<p>
プログラムに<code>package</code>句がなければ、gcコンパイラが元の詳細な痕跡を表示していました：ファイルが<code>main</code>パッケージにあるとした場合。
昔にそれは消され、<code>package</code>句が無い場合は、いまはエラーになります。
</p>


<p>
On the ARM, the toolchain supports "external linking", which
is a step towards being able to build shared libraries with the gc
tool chain and to provide dynamic linking support for environments
in which that is necessary.
</p>

<p>
ARMでは、Goツールは"外部リンク"をサポートしています。
GCツールを使った共有ライブラリをビルドすることができ、
動的リンクのサポートを提供できるようにする第一歩です。
</p>

<p>
In the runtime for the ARM, with <code>5a</code>, it used to be possible to refer
to the runtime-internal <code>m</code> (machine) and <code>g</code>
(goroutine) variables using <code>R9</code> and <code>R10</code> directly.
It is now necessary to refer to them by their proper names.
</p>

<p>
ARMのランタイムにある、<code>5a</code>を使って、
直接<code>R9</code>と<code>R10</code>を使用して
ランタイム内部の<code>m</code>（マシン）と<code>g</code>（ゴルーチン）変数を参照可能にするために使いました。
適切な名前でそれらを参照することが必要となります。
</p>

<p>
Also on the ARM, the <code>5l</code> linker (sic) now defines the
<code>MOVBS</code> and <code>MOVHS</code> instructions
as synonyms of <code>MOVB</code> and <code>MOVH</code>,
to make clearer the separation between signed and unsigned
sub-word moves; the unsigned versions already existed with a
<code>U</code> suffix.
</p>

<p>
<code>5l</code>リンカ(sic)は、<code>MOVBS</code>と<code>MOVHS</code>命令を
<code>MOVB</code>と<code>MOVH</code>の同義語として定義します。
siegnedとunsignedのサブワードmoveの区別を明確にするためです。
unsignedについては、すでに<code>U</code>のサフィックスがついていました。

</p>

<h3 id="cover">Test coverage</h3>

<p>
One major new feature of <a href="http://golang.org/pkg/go/"><code>go test</code></a> is
that it can now compute and, with help from a new, separately installed
"go tool cover" program, display test coverage results.
</p>

<p>
<a href="http://golang.org/pkg/go/"><code>go test</code></a>の主な新しい特徴の一つとして、
テストのカバレッジを計算でき、結果を表示することができるようになったということです。
それは、"go tool cover"プログラムで、別にインストールする必要があります。
</p>

<p>
The cover tool is part of the
<a href="https://code.google.com/p/go/source/checkout?repo=tools"><code>go.tools</code></a>
subrepository.
It can be installed by running
</p>

<pre>
$ go get code.google.com/p/go.tools/cmd/cover
</pre>

<p>
coverツールはサブリポジトリ<a href="https://code.google.com/p/go/source/checkout?repo=tools"><code>go.tools</code></a>にあります。
それは、下記のコマンドを実行することでインストールすることが出来ます。
</p>
<pre>
$ go get code.google.com/p/go.tools/cmd/cover
</pre>

<p>
The cover tool does two things.
First, when "go test" is given the <code>-cover</code> flag, it is run automatically 
to rewrite the source for the package and insert instrumentation statements.
The test is then compiled and run as usual, and basic coverage statistics are reported:
</p>

<p>
coverツールは2つの事ができます。
1つは、"go test"に <code>-cover</code>フラグを付けた時、
パッケージのソースを書き換え、計測ステートメントを挿入することを自動で行います。
そのテストはコンパイルされ、いつものとおりに実行し、基本的なカバレッジの測定結果がレポートされます:
</p>

<pre>
$ go test -cover fmt
ok      fmt     0.060s  coverage: 91.4% of statements
$
</pre>

<p>
Second, for more detailed reports, different flags to "go test" can create a coverage profile file,
which the cover program, invoked with "go tool cover", can then analyze.
</p>

<p>
2つ目は、もっと詳細にレポートしてくれます。"go test"に違うフラグを付けると、カバレッジプロファイルのフィアルを作成することができます。
"go tool cover"でcoverプログラムが実行し解析します。
</p>

<p>
Details on how to generate and analyze coverage statistics can be found by running the commands
</p>

<p>
カバレッジの計測結果を生成の仕方や解析の仕方についての詳しいことは、次のコマンドで確認できます。
</p>

<pre>
$ go help testflag
$ go tool cover -help
</pre>

<h3 id="go_doc">The go doc command is deleted</h3>

<p>
The "go doc" command is deleted.
Note that the <a href="http://golang.org/cmd/godoc/"><code>godoc</code></a> tool itself is not deleted,
just the wrapping of it by the <a href="http://golang.org/cmd/go/"><code>go</code></a> command.
All it did was show the documents for a package by package path,
which godoc itself already does with more flexibility.
It has therefore been deleted to reduce the number of documentation tools and,
as part of the restructuring of godoc, encourage better options in future.
</p>

<p>
"go doc"コマンドは削除されました。
<a href="http://golang.org/cmd/go/"><code>go</code></a>コマンドでツールをラッピングしているだけで、
<a href="http://golang.org/cmd/godoc/"><code>godoc</code></a>ツール自体が削除されたわけではないことに注意してください。
All it did was show the documents for a package by package path,
which godoc itself already does with more flexibility.
したがって、文書化ツールの数を減らし、godocの再構築の一環として、
将来的により良い選択を奨励するために削除されました。
</p>

<p>
<em>Updating</em>: For those who still need the precise functionality of running
</p>

<pre>
$ go doc
</pre>

<p>
in a directory, the behavior is identical to running
</p>

<p>
の動きは、
</p>

<pre>
$ godoc .
</pre>

<p>
と同じです。
</p>

<h3 id="gocmd">Changes to the go command</h3>

<p>
The <a href="http://golang.org/cmd/go/"><code>go get</code></a> command
now has a <code>-t</code> flag that causes it to download the dependencies
of the tests run by the package, not just those of the package itself.
By default, as before, dependencies of the tests are not downloaded.
</p>

<p>
<a href="http://golang.org/cmd/go/"><code>go get</code></a>コマンドに、<code>-t</code>フラグが追加されました。
それは、パッケージにより実行するテストの依存されるものダウンロードします。パッケージ自身の依存関係ではありません。
デフォルトでは、以前のように、テストの依存されるものはダウンロードされません。
</p>

<h2 id="performance">Performance</h2>

<p>
There are a number of significant performance improvements in the standard library; here are a few of them.
</p>

<p>
標準ライブラリのパフォーマンスの大幅な改善が多くあります。ここではそのうち少し紹介します。
</p>

<ul> 

<li class="en">
The <a href="http://golang.org/pkg/compress/bzip2/"><code>compress/bzip2</code></a>
decompresses about 30% faster.
</li>

<li>
<a href="http://golang.org/pkg/compress/bzip2/"><code>compress/bzip2</code></a>は約30%早く復元します。
</li>

<li>
The <a href="http://golang.org/pkg/crypto/des/"><code>crypto/des</code></a> package
is about five times faster.
</li>

<li>
<a href="http://golang.org/pkg/crypto/des/"><code>crypto/des</code></a>パッケージは、約5倍早くなりました。
</li>

<li>
The <a href="http://golang.org/pkg/encoding/json/"><code>encoding/json</code></a> package
encodes about 30% faster.
</li>

<li>
<a href="http://golang.org/pkg/encoding/json/"><code>encoding/json</code></a>パッケージは、約30%早くエンコードなりました。
</li>

<li>
Networking performance on Windows and BSD systems is about 30% faster through the use
of an integrated network poller in the runtime, similar to what was done for Linux and OS X
in Go 1.1.
</li>

<li>
Go1.1でLinuxとOS Xに対して行われたものと同様に、ランタイムにあるintegrated network pollerを使用して、WindowsとBSDのネットワークパフォーマンスは約30%早くなりました。
</li>

</ul>

<h2 id="library">Changes to the standard library 標準ライブラリの変更点</h2>


<h3 id="archive_tar_zip">The archive/tar and archive/zip packages</h3>

<p>
The
<a href="http://golang.org/pkg/archive/tar/"><code>archive/tar</code></a>
and
<a href="http://golang.org/pkg/archive/zip/"><code>archive/zip</code></a>
packages have had a change to their semantics that may break existing programs.
The issue is that they both provided an implementation of the
<a href="http://golang.org/pkg/os/#FileInfo"><code>os.FileInfo</code></a>
interface that was not compliant with the specification for that interface.
In particular, their <code>Name</code> method returned the full
path name of the entry, but the interface specification requires that
the method return only the base name (final path element).
</p>

<p>
<a href="http://golang.org/pkg/archive/tar/"><code>archive/tar</code></a>と<a href="http://golang.org/pkg/archive/zip/"><code>archive/zip</code></a>パッケージは、既存のプログラムを破壊することがあり、セマンティックへの変更がありました。
問題は、両方共、<a href="http://golang.org/pkg/os/#FileInfo"><code>os.FileInfo</code></a>インターフェースの実装を提供したことです。
それは、それらのインターフェースの仕様に準拠していませんでした。
特に、<code>Name</code>メソッドがエントリのフルパス名を返していましたが、インターフェースの仕様は、そのメソッドはパスのベースの名前（最後のパス要素）だけを返すことが必要です。
</p>	

<p>
<em>Updating</em>: Since this behavior was newly implemented and
a bit obscure, it is possible that no code depends on the broken behavior.
If there are programs that do depend on it, they will need to be identified
and fixed manually.
</p>

<p>
<em>Updating</em>:この振る舞いは、新しく実装され、少し隠蔽されましたので、コードが壊れた振る舞いに依存しない可能性があります。
もしそれに依存する問題があれば、確認し手動で修正する必要があります。

</p>	

<h3 id="encoding">The new encoding package</h3>

<p>
There is a new package, <a href="http://golang.org/pkg/encoding/"><code>encoding</code></a>,
that defines a set of standard encoding interfaces that may be used to
build custom marshalers and unmarshalers for packages such as
<a href="http://golang.org/pkg/encoding/xml/"><code>encoding/xml</code></a>,
<a href="http://golang.org/pkg/encoding/json/"><code>encoding/json</code></a>,
and
<a href="http://golang.org/pkg/encoding/binary/"><code>encoding/binary</code></a>.
These new interfaces have been used to tidy up some implementations in
the standard library.
</p>

<p>
新パッケージ<a href="http://golang.org/pkg/encoding/"><code>encoding</code></a>は、<a href="http://golang.org/pkg/encoding/xml/"><code>encoding/xml</code></a>,<a href="http://golang.org/pkg/encoding/json/"><code>encoding/json</code></a>,<a href="http://golang.org/pkg/encoding/binary/"><code>encoding/binary</code></a>のようなパッケージに対して、カスタムmarshalers and unmarshalersを構築するかもしれない標準encodingインターフェースのセットを定義しています。
</p>

<p>
The new interfaces are called
<a href="http://golang.org/pkg/encoding/#BinaryMarshaler"><code>BinaryMarshaler</code></a>,
<a href="http://golang.org/pkg/encoding/#BinaryUnmarshaler"><code>BinaryUnmarshaler</code></a>,
<a href="http://golang.org/pkg/encoding/#TextMarshaler"><code>TextMarshaler</code></a>,
and
<a href="http://golang.org/pkg/encoding/#TextUnmarshaler"><code>TextUnmarshaler</code></a>.
Full details are in the <a href="http://golang.org/pkg/encoding/">documentation</a> for the package
and a separate <a href="http://golang.org/s/go12encoding">design document</a>.
</p>

<p>
新しいインターフェースです。
<a href="http://golang.org/pkg/encoding/#BinaryMarshaler"><code>BinaryMarshaler</code></a>,<a href="http://golang.org/pkg/encoding/#BinaryUnmarshaler"><code>BinaryUnmarshaler</code></a>,<a href="http://golang.org/pkg/encoding/#TextMarshaler"><code>TextMarshaler</code></a>,<a href="http://golang.org/pkg/encoding/#TextUnmarshaler"><code>TextUnmarshaler</code></a>
詳細は<a href="http://golang.org/pkg/encoding/">documentation</a>と<a href="http://golang.org/s/go12encoding">design document</a>にあります。
</p>

<h3 id="fmt_indexed_arguments">The fmt package</h3>

<p>
The <a href="http://golang.org/pkg/fmt/"><code>fmt</code></a> package's formatted print
routines such as <a href="http://golang.org/pkg/fmt/#Printf"><code>Printf</code></a>
now allow the data items to be printed to be accessed in arbitrary order
by using an indexing operation in the formatting specifications.
Wherever an argument is to be fetched from the argument list for formatting,
either as the value to be formatted or as a width or specification integer,
a new optional indexing notation <code>[</code><em>n</em><code>]</code>
fetches argument <em>n</em> instead.
The value of <em>n</em> is 1-indexed.
After such an indexing operating, the next argument to be fetched by normal
processing will be <em>n</em>+1.
</p>

<p>
<a href="http://golang.org/pkg/fmt/#Printf"><code>Printf</code></a>のような<a href="http://golang.org/pkg/fmt/"><code>fmt</code></a>パッケージのフォーマットを出力するものは、書式仕様におけるインデックス操作を使うことによって、データを任意の順番で出力できるようになりました。
Wherever an argument is to be fetched from the argument list for formatting,
either as the value to be formatted or as a width or specification integer,
新しいオプショナルなインデックス記法<code>[</code><em>n</em><code>]</code>は、<em>n</em> 番目の引数をかわりに取ります。
<em>n</em>の値は1始まりです。
インデックス操作したあとの次の引数は<em>n</em>+1として取得されます。
</p>

<p>
For example, the normal <code>Printf</code> call
</p>

<p>
例えば、<code>Printf</code>は
</p>

<pre>
fmt.Sprintf("%c %c %c\n", 'a', 'b', 'c')
</pre>

<p>
would create the string <code>"a b c"</code>, but with indexing operations like this,
</p>

<p>
文字列<code>"a b c"</code>を作成しますが、次のようにインデックス操作をすると、
</p>

<pre>
fmt.Sprintf("%[3]c %[1]c %c\n", 'a', 'b', 'c')
</pre>

<p>
the result is "<code>"c a b"</code>. The <code>[3]</code> index accesses the third formatting
argument, which is <code>'c'</code>, <code>[1]</code> accesses the first, <code>'a'</code>,
and then the next fetch accesses the argument following that one, <code>'b'</code>.
</p>

<p>
結果は、<code>"c a b"</code>となります。<code>[3]</code>は3番目の引数にアクセスし、'c'を出力しています。<code>[1]</code>は最初の<code>'a'</code>、それから＋1した次の引数'b'にアクセスしています。(訳者注：<a hfre="http://play.golang.org/p/T2DaHNcOzq">http://play.golang.org/p/T2DaHNcOzq</a>)
</p>

<p>
The motivation for this feature is programmable format statements to access
the arguments in different order for localization, but it has other uses:
</p>

<p>
この特徴における興味は、局所的に異なる順番で引数にアクセスするために、プログラム制御できるフォーマット文であることですが、次のように他の用途でも使えます。
</p>

<pre>
log.Printf("trace: value %v of type %[1]T\n", expensiveFunction(a.b[c]))
</pre>

<p>
<em>Updating</em>: The change to the syntax of format specifications
is strictly backwards compatible, so it affects no working programs.
</p>

<p>
<em>Updating</em>:フォーマット仕様シンタックスの変更は、厳密に後方互換ですので、作業プログラムには影響ありません。
</p>

<h3 id="text_template">The text/template and html/template packages</h3>

<p>
The
<a href="http://golang.org/pkg/text/template/"><code>text/template</code></a> package
has a couple of changes in Go 1.2, both of which are also mirrored in the
<a href="http://golang.org/pkg/html/template/"><code>html/template</code></a> package.
</p>

<p>
<a href="http://golang.org/pkg/text/template/"><code>text/template</code></a>パッケージはGo1.2で2つの変更があります。
<a href="http://golang.org/pkg/html/template/"><code>html/template</code></a>も同様です。
</p>

<p>
First, there are new default functions for comparing basic types.
The functions are listed in this table, which shows their names and
the associated familiar comparison operator.
</p>

<p>
1つ目の変更点は、基本タイプを比較するために、新しくデフォルトの関数を追加しました。それらの関数は、下記テーブルに、名前と関連する比較演算子を示しています。
</p>

<table cellpadding="0" summary="Template comparison functions">
<tr>
<th width="50"></th><th width="100">Name</th> <th width="50">Operator</th>
</tr>
<tr>
<td></td><td><code>eq</code></td> <td><code>==</code></td>
</tr>
<tr>
<td></td><td><code>ne</code></td> <td><code>!=</code></td>
</tr>
<tr>
<td></td><td><code>lt</code></td> <td><code>&lt;</code></td>
</tr>
<tr>
<td></td><td><code>le</code></td> <td><code>&lt;=</code></td>
</tr>
<tr>
<td></td><td><code>gt</code></td> <td><code>&gt;</code></td>
</tr>
<tr>
<td></td><td><code>ge</code></td> <td><code>&gt;=</code></td>
</tr>
</table>

<p>
These functions behave slightly differently from the corresponding Go operators.
First, they operate only on basic types (<code>bool</code>, <code>int</code>,
<code>float64</code>, <code>string</code>, etc.).
(Go allows comparison of arrays and structs as well, under some circumstances.)
Second, values can be compared as long as they are the same sort of value:
any signed integer value can be compared to any other signed integer value for example. (Go
does not permit comparing an <code>int8</code> and an <code>int16</code>).
Finally, the <code>eq</code> function (only) allows comparison of the first
argument with one or more following arguments. The template in this example,
</p>

<p>
これらの関数は、対応するGoの演算子とはわずかに動きが異なります。
1つは、basic types (<code>bool</code>, <code>int</code>,
<code>float64</code>, <code>string</code>, など)だけ扱います。
(Goプログラムは基本タイプの他に、配列や構造体も比較可能です。)

2つ目に、値同士が同じ種類であれば比較できます。例えば、sigined integerの値と他のsigine integerの値を比較できます。(Goプログラムは<code>int8</code>と<code>int16</code>を比較することができません。)

最後に、<code>eq</code>関数は、最初の引数と、それに続く1つか複数の引数と比較します。下記の例では、
</p>

<pre>
{{"{{"}}if eq .A 1 2 3 {{"}}"}} equal {{"{{"}}else{{"}}"}} not equal {{"{{"}}end{{"}}"}}
</pre>

<p>
reports "equal" if <code>.A</code> is equal to <em>any</em> of 1, 2, or 3.
</p>

<p>
<code>.A</code> が1,2,3の<em>いずれか</em>と等しい場合、"equal"を出力します。
</p>

<p>
The second change is that a small addition to the grammar makes "if else if" chains easier to write.
Instead of writing,
</p>

<p>
2つ目の変更点は、"if else if"の追加です。以前までは下記のように書く必要がありましたが、
</p>

<pre>
{{"{{"}}if eq .A 1{{"}}"}} X {{"{{"}}else{{"}}"}} {{"{{"}}if eq .A 2{{"}}"}} Y {{"{{"}}end{{"}}"}} {{"{{"}}end{{"}}"}} 
</pre>

<p>
one can fold the second "if" into the "else" and have only one "end", like this:
</p>

<p>
下記のように、"else"に"if"を入れることができるようになり、"end"が1つだけでよくなりました。
</p>

<pre>
{{"{{"}}if eq .A 1{{"}}"}} X {{"{{"}}else if eq .A 2{{"}}"}} Y {{"{{"}}end{{"}}"}}
</pre>

<p>
The two forms are identical in effect; the difference is just in the syntax.
</p>

<p>
2つの形式は結果的には同じで、書き方だけが違います。
</p>

<p>
<em>Updating</em>: Neither the "else if" change nor the comparison functions
affect existing programs. Those that
already define functions called <code>eq</code> and so on through a function
map are unaffected because the associated function map will override the new
default function definitions.
</p>

<p>
<em>Updating</em>: "else if"の変更も比較関数の変更もどちらも既存プログラムに影響を与えません。
関数マップを使って<code>eq</code>と名づけた関数などを既に定義しているものは、関連する関数マップが新しいデフォルト関数定義を上書きするので、影響を受けません。
</p>

<h3 id="new_packages">New packages</h3>

<p>
There are two new packages.
</p>

<ul>
<li>
The <a href="http://golang.org/pkg/encoding/"><code>encoding</code></a> package is
<a href="#encoding">described above</a>.
</li>
<li>
The <a href="http://golang.org/pkg/image/color/palette/"><code>image/color/palette</code></a> package
provides standard color palettes.
</li>
</ul>

<p>
2つ新しいパッケージを追加しました。
</p>

<ul>
<li>
<a href="http://golang.org/pkg/encoding/"><code>encoding</code></a> パッケージは<a href="#encoding">上述しました</a>。
</li>
<li>
<a href="http://golang.org/pkg/image/color/palette/"><code>image/color/palette</code></a> パッケージは、標準カラーパレットを提供します。
</li>
</ul>


<h3 id="minor_library_changes">Minor changes to the library</h3>

<p>
The following list summarizes a number of minor changes to the library, mostly additions.
See the relevant package documentation for more information about each change.
</p>

<p>
以下のリストは、ライブラリのマイナー変更の要約で、ほとんどが追加になります。
各変更のより詳細な情報は、関連ドキュメントをご参照下さい。
</p>

<ul>

<li>
The <a href="http://golang.org/pkg/archive/zip/"><code>archive/zip</code></a> package
adds the
<a href="http://golang.org/pkg/archive/zip/#File.DataOffset"><code>DataOffset</code></a> accessor
to return the offset of a file's (possibly compressed) data within the archive.
</li>

<li>
<a href="http://golang.org/pkg/archive/zip/"><code>archive/zip</code></a>パッケージに、<a href="http://golang.org/pkg/archive/zip/#File.DataOffset"><code>DataOffset</code></a>関数が追加されました。
これはそのアーカイブ内のファイルの（圧縮されているかもしれない）データのオフセットを返します。
</li>

<li>
The <a href="http://golang.org/pkg/bufio/"><code>bufio</code></a> package
adds <a href="http://golang.org/pkg/bufio/#Reader.Reset"><code>Reset</code></a>
methods to <a href="http://golang.org/pkg/bufio/#Reader"><code>Reader</code></a> and
<a href="http://golang.org/pkg/bufio/#Writer"><code>Writer</code></a>.
These methods allow the <a href="http://golang.org/pkg/io/#Reader"><code>Readers</code></a>
and <a href="http://golang.org/pkg/io/#Writer"><code>Writers</code></a>
to be re-used on new input and output readers and writers, saving
allocation overhead. 
</li>

<li>
<a href="http://golang.org/pkg/bufio/"><code>bufio</code></a>パッケージには、<a href="http://golang.org/pkg/bufio/#Reader"><code>Reader</code></a>と<a href="http://golang.org/pkg/bufio/#Writer"><code>Writer</code></a>に<a href="http://golang.org/pkg/bufio/#Reader.Reset"><code>Reset</code></a>メソッドが追加されました。
このメソッドは、<a href="http://golang.org/pkg/io/#Reader"><code>Readers</code></a>
と <a href="http://golang.org/pkg/io/#Writer"><code>Writers</code></a>新しいインプットやアウトプットのreadersやwritersを再利用でき、メモリ割付のオーバーヘッドをセーブすることができます。
</li>

<li>
The <a href="http://golang.org/pkg/compress/bzip2/"><code>compress/bzip2</code></a>
can now decompress concatenated archives.
</li>

<li>
<a href="http://golang.org/pkg/compress/bzip2/"><code>compress/bzip2</code></a>は連結しているアーカイブを復元することができるようになりました。
</li>

<li>
The <a href="http://golang.org/pkg/compress/flate/"><code>compress/flate</code></a>
package adds a <a href="http://golang.org/pkg/compress/flate/#Writer.Reset"><code>Reset</code></a> 
method on the <a href="http://golang.org/pkg/compress/flate/#Writer"><code>Writer</code></a>,
to make it possible to reduce allocation when, for instance, constructing an
archive to hold multiple compressed files.
</li>

<li>
<a href="http://golang.org/pkg/compress/flate/"><code>compress/flate</code></a>パッケージには、<a href="http://golang.org/pkg/compress/flate/#Writer"><code>Writer</code></a>に<a href="http://golang.org/pkg/compress/flate/#Writer.Reset"><code>Reset</code></a>メソッドが追加されました。
たとえば、複数の圧縮ファイルをホールドしておくためにアーカイブを作成するとき、アロケーションを減らすことが可能になります。
</li>

<li>
The <a href="http://golang.org/pkg/compress/gzip/"><code>compress/gzip</code></a> package's
<a href="http://golang.org/pkg/compress/gzip/#Writer"><code>Writer</code></a> type adds a
<a href="http://golang.org/pkg/compress/gzip/#Writer.Reset"><code>Reset</code></a>
so it may be reused.
</li>

<li>
<a href="http://golang.org/pkg/compress/gzip/"><code>compress/gzip</code></a>パッケージの<a href="http://golang.org/pkg/compress/gzip/#Writer"><code>Writer</code></a>型に<a href="http://golang.org/pkg/compress/gzip/#Writer.Reset"><code>Reset</code></a>メソッドが追加されましたので、再利用されるかもしれません。
</li>

<li>
The <a href="http://golang.org/pkg/compress/zlib/"><code>compress/zlib</code></a> package's
<a href="http://golang.org/pkg/compress/zlib/#Writer"><code>Writer</code></a> type adds a
<a href="http://golang.org/pkg/compress/zlib/#Writer.Reset"><code>Reset</code></a>
so it may be reused.
</li>

<li>
<a href="http://golang.org/pkg/compress/zlib/"><code>compress/zlib</code></a>パッケージの<a href="http://golang.org/pkg/compress/zlib/#Writer"><code>Writer</code></a>型に<a href="http://golang.org/pkg/compress/zlib/#Writer.Reset"><code>Reset</code></a>メソッドが追加されましたので、再利用されるかもしれません。
</li>

<li>
The <a href="http://golang.org/pkg/container/heap/"><code>container/heap</code></a> package
adds a <a href="http://golang.org/pkg/container/heap/#Fix"><code>Fix</code></a>
method to provide a more efficient way to update an item's position in the heap.
</li>

<li>
<a href="http://golang.org/pkg/container/heap/"><code>container/heap</code></a>パッケージに<a href="http://golang.org/pkg/container/heap/#Fix"><code>Fix</code></a>メソッドが追加されました。heapのアイテム位置を更新するもっと効率的な手段を提供します。
</li>

<li>
The <a href="http://golang.org/pkg/container/list/"><code>container/list</code></a> package
adds the <a href="http://golang.org/pkg/container/list/#List.MoveBefore"><code>MoveBefore</code></a>
and
<a href="http://golang.org/pkg/container/list/#List.MoveAfter"><code>MoveAfter</code></a>
methods, which implement the obvious rearrangement.
</li>

<li>
<a href="http://golang.org/pkg/container/list/"><code>container/list</code></a>パッケーには、<a href="http://golang.org/pkg/container/list/#List.MoveBefore"><code>MoveBefore</code></a>メソッドと<a href="http://golang.org/pkg/container/list/#List.MoveAfter"><code>MoveAfter</code></a>メソッドが追加されました。これは明らかに再整理を実装しています。
</li>

<li>
The <a href="http://golang.org/pkg/crypto/cipher/"><code>crypto/cipher</code></a> package
adds the a new GCM mode (Galois Counter Mode), which is almost always
used with AES encryption.
</li>

<li>
<a href="http://golang.org/pkg/crypto/cipher/"><code>crypto/cipher</code></a>パッケージには、新しくGCM(Galois Counter Mode)を追加しました。それはほとんどの場合、AES暗号が使用されます。
</li>

<li>
The 
<a href="http://golang.org/pkg/crypto/md5/"><code>crypto/md5</code></a> package
adds a new <a href="http://golang.org/pkg/crypto/md5/#Sum"><code>Sum</code></a> function
to simplify hashing without sacrificing performance.
</li>

<li>
<a href="http://golang.org/pkg/crypto/md5/"><code>crypto/md5</code></a>パッケージには、新しく<a href="http://golang.org/pkg/crypto/md5/#Sum"><code>Sum</code></a>関数が追加されました。パフォーマンスが犠牲にならずに簡単にハッシュ化します。
</li>

<li>
Similarly, the 
<a href="http://golang.org/pkg/crypto/md5/"><code>crypto/sha1</code></a> package
adds a new <a href="http://golang.org/pkg/crypto/sha1/#Sum"><code>Sum</code></a> function.
</li>

<li>
同様に、<a href="http://golang.org/pkg/crypto/md5/"><code>crypto/sha1</code></a>パッケージにも<a href="http://golang.org/pkg/crypto/sha1/#Sum"><code>Sum</code></a>関数が追加されています。

</li>

<li>
Also, the
<a href="http://golang.org/pkg/crypto/sha256/"><code>crypto/sha256</code></a> package
adds <a href="http://golang.org/pkg/crypto/sha256/#Sum256"><code>Sum256</code></a>
and <a href="http://golang.org/pkg/crypto/sha256/#Sum224"><code>Sum224</code></a> functions.
</li>

<li>
また、<a href="http://golang.org/pkg/crypto/sha256/"><code>crypto/sha256</code></a>パッケージには、<a href="http://golang.org/pkg/crypto/sha256/#Sum256"><code>Sum256</code></a>関数と<a href="http://golang.org/pkg/crypto/sha256/#Sum224"><code>Sum224</code></a>関数が追加されました。
</li>

<li>
Finally, the <a href="http://golang.org/pkg/crypto/sha512/"><code>crypto/sha512</code></a> package
adds <a href="http://golang.org/pkg/crypto/sha512/#Sum512"><code>Sum512</code></a> and
<a href="http://golang.org/pkg/crypto/sha512/#Sum384"><code>Sum384</code></a> functions.
</li>

<li>
最後に、<a href="http://golang.org/pkg/crypto/sha512/"><code>crypto/sha512</code></a>パッケージに、<a href="http://golang.org/pkg/crypto/sha512/#Sum512"><code>Sum512</code></a>関数と <a href="http://golang.org/pkg/crypto/sha512/#Sum384"><code>Sum384</code></a>関数が追加されました。
</li>

<li>
The <a href="http://golang.org/pkg/crypto/x509/"><code>crypto/x509</code></a> package
adds support for reading and writing arbitrary extensions.
</li>

<li>
<a href="http://golang.org/pkg/crypto/x509/"><code>crypto/x509</code></a>パッケージは、任意の拡張を読み書きのサポートを追加しました。
</li>

<li>
The <a href="http://golang.org/pkg/crypto/tls/"><code>crypto/tls</code></a> package adds
support for TLS 1.1, 1.2 and AES-GCM.
</li>

<li>
<a href="http://golang.org/pkg/crypto/tls/"><code>crypto/tls</code></a>パッケージには、TLS 1.1, 1.2 と AES-GCMのサポートを追加しました。
</li>

<li>
The <a href="http://golang.org/pkg/database/sql/"><code>database/sql</code></a> package adds a
<a href="http://golang.org/pkg/database/sql/#DB.SetMaxOpenConns"><code>SetMaxOpenConns</code></a>
method on <a href="http://golang.org/pkg/database/sql/#DB"><code>DB</code></a> to limit the
number of open connections to the database.
</li>

<li>
<a href="http://golang.org/pkg/database/sql/"><code>database/sql</code></a>パッケージは、たくさんのデータベースとのオープン接続を制限するため、<a href="http://golang.org/pkg/database/sql/#DB"><code>DB</code></a>に<a href="http://golang.org/pkg/database/sql/#DB.SetMaxOpenConns"><code>SetMaxOpenConns</code></a>メソッドを追加しました。
</li>

<li>
The <a href="http://golang.org/pkg/encoding/csv/"><code>encoding/csv</code></a> package
now always allows trailing commas on fields.
</li>

<li>
<a href="http://golang.org/pkg/encoding/csv/"><code>encoding/csv</code></a>パッケージは、フィールドの最後がカンマの場合エラーになっていたのが、最後はカンマでもよくなりました。
</li>

<li>
The <a href="http://golang.org/pkg/encoding/gob/"><code>encoding/gob</code></a> package
now treats channel and function fields of structures as if they were unexported,
even if they are not. That is, it ignores them completely. Previously they would
trigger an error, which could cause unexpected compatibility problems if an
embedded structure added such a field.
The package also now supports the generic <code>BinaryMarshaler</code> and
<code>BinaryUnmarshaler</code> interfaces of the
<a href="http://golang.org/pkg/encoding/"><code>encoding</code></a> package
described above.
</li>

<li>
<a href="http://golang.org/pkg/encoding/gob/"><code>encoding/gob</code></a>パッケージは、チャネルと構造体の関数フィールドを、たとえそれらが無くても、アンエクスポートされているかのように扱うようになりました。これは、完全にそれらを無視します。
埋め込まれた構造体にそのようなフィールドを追加していた場合、予期せぬ互換性の問題を引き起こすことがあります。
また、上述された<a href="http://golang.org/pkg/encoding/"><code>encoding</code></a>パッケージの<code>BinaryMarshaler</code> と<code>BinaryUnmarshaler</code>インターフェースをサポートします。
</li>

<li>
The <a href="http://golang.org/pkg/encoding/json/"><code>encoding/json</code></a> package
now will always escape ampersands as "\u0026" when printing strings.
It will now accept but correct invalid UTF-8 in
<a href="http://golang.org/pkg/encoding/json/#Marshal"><code>Marshal</code></a>
(such input was previously rejected).
Finally, it now supports the generic encoding interfaces of the
<a href="http://golang.org/pkg/encoding/"><code>encoding</code></a> package
described above.
</li>

<li>
<a href="http://golang.org/pkg/encoding/json/"><code>encoding/json</code></a>パッケージは、文字列を出力するとき常にアンパサンドを"\u0026"としてエスケープします。
それは受け入れますが、(インプットが以前に拒否された)<a href="http://golang.org/pkg/encoding/json/#Marshal"><code>Marshal</code></a>では無効なUTF-8を訂正します。結局、上述した<a href="http://golang.org/pkg/encoding/"><code>encoding</code></a>パッケージの一般的なエンコードインターフェースをサポートします。
</li>

<li>
The <a href="http://golang.org/pkg/encoding/xml/"><code>encoding/xml</code></a> package
now allows attributes stored in pointers to be marshaled.
It also supports the generic encoding interfaces of the
<a href="http://golang.org/pkg/encoding/"><code>encoding</code></a> package
described above through the new
<a href="http://golang.org/pkg/encoding/xml/#Marshaler"><code>Marshaler</code></a>,
<a href="http://golang.org/pkg/encoding/xml/#Unmarshaler"><code>Unmarshaler</code></a>,
and related
<a href="http://golang.org/pkg/encoding/xml/#MarshalerAttr"><code>MarshalerAttr</code></a> and
<a href="http://golang.org/pkg/encoding/xml/#UnmarshalerAttr"><code>UnmarshalerAttr</code></a>
interfaces.
The package also adds a
<a href="http://golang.org/pkg/encoding/xml/#Encoder.Flush"><code>Flush</code></a> method
to the
<a href="http://golang.org/pkg/encoding/xml/#Encoder"><code>Encoder</code></a>
type for use by custom encoders. See the documentation for
<a href="http://golang.org/pkg/encoding/xml/#Encoder.EncodeToken"><code>EncodeToken</code></a>
to see how to use it.
</li>

<li>
<a href="http://golang.org/pkg/encoding/xml/"><code>encoding/xml</code></a>パッケージは属性をmarshalされたポインタに格納することができます。
また、新しい<a href="http://golang.org/pkg/encoding/xml/#Marshaler"><code>Marshaler</code></a>,
<a href="http://golang.org/pkg/encoding/xml/#Unmarshaler"><code>Unmarshaler</code></a>,<a href="http://golang.org/pkg/encoding/xml/#MarshalerAttr"><code>MarshalerAttr</code></a> and
<a href="http://golang.org/pkg/encoding/xml/#UnmarshalerAttr"><code>UnmarshalerAttr</code></a>を通して、上述した<a href="http://golang.org/pkg/encoding/"><code>encoding</code></a>パッケージの一般的なencodingインターフェースをサポートします。
このパッケージは、カスタムエンコーダーによる使用に対して、<a href="http://golang.org/pkg/encoding/xml/#Encoder"><code>Encoder</code></a>型に<a href="http://golang.org/pkg/encoding/xml/#Encoder.Flush"><code>Flush</code></a>メソッドも追加しました。使い方は<a href="http://golang.org/pkg/encoding/xml/#Encoder.EncodeToken"><code>EncodeToken</code></a>のドキュメントを見て下さい。
</li>

<li>
The <a href="http://golang.org/pkg/flag/"><code>flag</code></a> package now
has a <a href="http://golang.org/pkg/flag/#Getter"><code>Getter</code></a> interface
to allow the value of a flag to be retrieved. Due to the
Go 1 compatibility guidelines, this method cannot be added to the existing
<a href="http://golang.org/pkg/flag/#Value"><code>Value</code></a>
interface, but all the existing standard flag types implement it.
The package also now exports the <a href="http://golang.org/pkg/flag/#CommandLine"><code>CommandLine</code></a>
flag set, which holds the flags from the command line.
</li>

<li>
<a href="http://golang.org/pkg/flag/"><code>flag</code></a>パッケージは<a href="http://golang.org/pkg/flag/#Getter"><code>Getter</code></a>インターフェースを持ちます。これはflagの値を取得することができます。Go1の互換性ガイドラインにより、このメソッドを既存の<a href="http://golang.org/pkg/flag/#Value"><code>Value</code></a>に追加できませんが、すべての既存の標準flag型はそれを実装します。
また、<a href="http://golang.org/pkg/flag/#CommandLine"><code>CommandLine</code></a>変数を公開しました。この変数はFlagSetをコマンドラインから保持します。
</li>

<li>
The <a href="http://golang.org/pkg/go/ast/"><code>go/ast</code></a> package's
<a href="http://golang.org/pkg/go/ast/#SliceExpr"><code>SliceExpr</code></a> struct
has a new boolean field, <code>Slice3</code>, which is set to true
when representing a slice expression with three indices (two colons).
The default is false, representing the usual two-index form.
</li>

<li>
<a href="http://golang.org/pkg/go/ast/"><code>go/ast</code></a>パッケージの<a href="http://golang.org/pkg/go/ast/#SliceExpr"><code>SliceExpr</code></a>構造体は、新しいbooleanフィールド<code>Slice3</code>を持ちます。これは、スリーインデックススライス(2つコロンがあるスライス)を表すとき、trueがセットされます。
デフォルではfalseで、通常のツーインデックス形式を表します。
</li>

<li>
The <a href="http://golang.org/pkg/go/build/"><code>go/build</code></a> package adds
the <code>AllTags</code> field
to the <a href="http://golang.org/pkg/go/build/#Package"><code>Package</code></a> type,
to make it easier to process build tags.
</li>

<li>
<a href="http://golang.org/pkg/go/build/"><code>go/build</code></a> パッケージは、<a href="http://golang.org/pkg/go/build/#Package"><code>Package</code></a>型に<code>AllTags</code>フィールドを追加しました。これはタグをビルドするプロセスを簡単にします。
</li>

<li>
The <a href="http://golang.org/pkg/image/draw/"><code>image/draw</code></a> package now
exports an interface, <a href="http://golang.org/pkg/image/draw/#Drawer"><code>Drawer</code></a>,
that wraps the standard <a href="http://golang.org/pkg/image/draw/#Draw"><code>Draw</code></a> method.
The Porter-Duff operators now implement this interface, in effect binding an operation to
the draw operator rather than providing it explicitly.
Given a paletted image as its destination, the new
<a href="http://golang.org/pkg/image/draw/#FloydSteinberg"><code>FloydSteinberg</code></a>
implementation of the
<a href="http://golang.org/pkg/image/draw/#Drawer"><code>Drawer</code></a>
interface will use the Floyd-Steinberg error diffusion algorithm to draw the image.
To create palettes suitable for such processing, the new
<a href="http://golang.org/pkg/image/draw/#Quantizer"><code>Quantizer</code></a> interface
represents implementations of quantization algorithms that choose a palette
given a full-color image.
There are no implementations of this interface in the library.
</li>

<li>
<a href="http://golang.org/pkg/image/draw/"><code>image/draw</code></a>パッケージは、<a href="http://golang.org/pkg/image/draw/#Drawer"><code>Drawer</code></a>インターフェースを公開しました。これは標準の<a href="http://golang.org/pkg/image/draw/#Draw"><code>Draw</code></a>メソッドをラップしています。Porter-Duff操作がこのインターフェースで実装されました。
明示的にそれを提供するというより、drawオペレータに動作をバインドしています。
宛先としてパレットイメージが与えられ、<a href="http://golang.org/pkg/image/draw/#Drawer"><code>Drawer</code></a>インターフェースの新しい<a href="http://golang.org/pkg/image/draw/#FloydSteinberg"><code>FloydSteinberg</code></a>実装が、フロイド-スタインバーグ誤差拡散法をイメージを描くために使用します。
そのような処理に対して安定にパレットを作るには、新しい<a href="http://golang.org/pkg/image/draw/#Quantizer"><code>Quantizer</code></a>インターフェースが、フルカラーイメージを与えたパレットを選ぶ量子アルゴリズムの実装を表します。ライブラリには、このインターフェースの実装はありません。
</li>


<li>
The <a href="http://golang.org/pkg/image/gif/"><code>image/gif</code></a> package
can now create GIF files using the new
<a href="http://golang.org/pkg/image/gif/#Encode"><code>Encode</code></a>
and <a href="http://golang.org/pkg/image/gif/#EncodeAll"><code>EncodeAll</code></a>
functions.
Their options argument allows specification of an image
<a href="http://golang.org/pkg/image/draw/#Quantizer"><code>Quantizer</code></a> to use;
if it is <code>nil</code>, the generated GIF will use the 
<a href="http://golang.org/pkg/image/color/palette/#Plan9"><code>Plan9</code></a>
color map (palette) defined in the new
<a href="http://golang.org/pkg/image/color/palette/"><code>image/color/palette</code></a> package.
The options also specify a
<a href="http://golang.org/pkg/image/draw/#Drawer"><code>Drawer</code></a>
to use to create the output image;
if it is <code>nil</code>, Floyd-Steinberg error diffusion is used.
</li>

<li>
<a href="http://golang.org/pkg/image/gif/"><code>image/gif</code></a>パッケージは、新しい<a href="http://golang.org/pkg/image/gif/#Encode"><code>Encode</code></a>関数と
<a href="http://golang.org/pkg/image/gif/#EncodeAll"><code>EncodeAll</code></a>関数を使って、GIFファイルを作ることができるようになりました。それらの引数のOptionは、イメージのスペックを指定することができます。<a href="http://golang.org/pkg/image/draw/#Quantizer"><code>Quantizer</code></a>を使った場合、もし<code>nil</code>なら、生成されたGIFは、新しく<a href="http://golang.org/pkg/image/color/palette/"><code>image/color/palette</code></a> パッケージで定義された<a href="http://golang.org/pkg/image/color/palette/#Plan9"><code>Plan9</code></a>のカラーマップ(パレット)を使用します。アウトプットイメージを作成するのに<a href="http://golang.org/pkg/image/draw/#Drawer"><code>Drawer</code></a>を使った場合、もし<code>nil</code>なら、フロイド-スタインバーグ誤差拡散法が使われます。

</li>

<li>
The <a href="http://golang.org/pkg/io/#Copy"><code>Copy</code></a> method of the
<a href="http://golang.org/pkg/io/"><code>io</code></a> package now prioritizes its
arguments differently.
If one argument implements <a href="http://golang.org/pkg/io/#WriterTo"><code>WriterTo</code></a>
and the other implements <a href="http://golang.org/pkg/io/#ReaderFrom"><code>ReaderFrom</code></a>,
<a href="http://golang.org/pkg/io/#Copy"><code>Copy</code></a> will now invoke
<a href="http://golang.org/pkg/io/#WriterTo"><code>WriterTo</code></a> to do the work,
so that less intermediate buffering is required in general.
</li>

<li>
<a href="http://golang.org/pkg/io/"><code>io</code></a>パッケージの<a href="http://golang.org/pkg/io/#Copy"><code>Copy</code></a>メソッドは、2つの引数の優先順位を変更しました。1つの引数が<a href="http://golang.org/pkg/io/#WriterTo"><code>WriterTo</code></a>を実装し、もう1つの引数が<a href="http://golang.org/pkg/io/#ReaderFrom"><code>ReaderFrom</code></a>を実装している場合、<a href="http://golang.org/pkg/io/#Copy"><code>Copy</code></a>は<a href="http://golang.org/pkg/io/#WriterTo"><code>WriterTo</code></a>を呼びします。so that less intermediate buffering is required in general.
</li>

<li>
The <a href="http://golang.org/pkg/net/"><code>net</code></a> package requires cgo by default
because the host operating system must in general mediate network call setup.
On some systems, though, it is possible to use the network without cgo, and useful
to do so, for instance to avoid dynamic linking.
The new build tag <code>netgo</code> (off by default) allows the construction of a
<code>net</code> package in pure Go on those systems where it is possible.
</li>

<li>
<a href="http://golang.org/pkg/net/"><code>net</code></a>パッケージは、デフォルトではcgoを必要とします。なぜなら、ホストOSが一般にはネットワークコールセットアップを仲介するからです。いくつかのシステムでは、cgo無しでネットワークを使うことが可能で、動的リンクを避けるインスタンスに対して、そのようにすることが有用です。新しいビルドタグ<code>netgo</code>(デフォルトではオフ)は、可能なシステム上で純粋なGoで<code>net</code>パッケージの構築をすることができます。
</li>

<li>
The <a href="http://golang.org/pkg/net/"><code>net</code></a> package adds a new field
<code>DualStack</code> to the <a href="http://golang.org/pkg/net/#Dialer"><code>Dialer</code></a>
struct for TCP connection setup using a dual IP stack as described in
<a href="http://tools.ietf.org/html/rfc6555">RFC 6555</a>.
</li>

<li>
<a href="http://golang.org/pkg/net/"><code>net</code></a>パッケージは、<a href="http://tools.ietf.org/html/rfc6555">RFC 6555</a>で説明されるように、デュアルIPスタックを使っているTCP接続セットアップに対して、<a href="http://golang.org/pkg/net/#Dialer"><code>Dialer</code></a>構造体に、新しいフィールド<code>DualStack</code>を追加しました。
</li>

<li>
The <a href="http://golang.org/pkg/net/http/"><code>net/http</code></a> package will no longer
transmit cookies that are incorrect according to
<a href="http://tools.ietf.org/html/rfc6265">RFC 6265</a>.
It just logs an error and sends nothing.
Also,
the <a href="http://golang.org/pkg/net/http/"><code>net/http</code></a> package's
<a href="http://golang.org/pkg/net/http/#ReadResponse"><code>ReadResponse</code></a>
function now permits the <code>*Request</code> parameter to be <code>nil</code>,
whereupon it assumes a GET request.
Finally, an HTTP server will now serve HEAD
requests transparently, without the need for special casing in handler code.
While serving a HEAD request, writes to a 
<a href="http://golang.org/pkg/net/http/#Handler"><code>Handler</code></a>'s
<a href="http://golang.org/pkg/net/http/#ResponseWriter"><code>ResponseWriter</code></a>
are absorbed by the
<a href="http://golang.org/pkg/net/http/#Server"><code>Server</code></a>
and the client receives an empty body as required by the HTTP specification.
</li>

<li>
<a href="http://golang.org/pkg/net/http/"><code>net/http</code></a>パッケージは、<a href="http://tools.ietf.org/html/rfc6265">RFC 6265</a>により、有効でないクッキーを送信しません。
それは、エラーを記録するだけで、なにも送信しません。(訳者補足：<a href="http://bit.ly/1cJcMcY">メインの差分箇所</a> )
<a href="http://golang.org/pkg/net/http/"><code>net/http</code></a>パッケージの<a href="http://golang.org/pkg/net/http/#ReadResponse"><code>ReadResponse</code></a>関数の<code>*Request</code> 引数を<code>nil</code>としてもよくなりました。(訳者補足：<a href="http://bit.ly/1cJdmHz">メインの差分箇所</a> ) その場合、GETリクエストを指定したことになります。
最後に、HTTPサーバーは、ハンドラコードで特殊ケースを必要とせずに、HEADリクエストを明示的に取り扱います。(訳者補足：<a href="http://bit.ly/1cJf9wg">コミットログメッセージ</a>)
HEADリクエストのサーブ中、<a href="http://golang.org/pkg/net/http/#Handler"><code>Handler</code></a>の<a href="http://golang.org/pkg/net/http/#ResponseWriter"><code>ResponseWriter</code></a>への書き込みは、<a href="http://golang.org/pkg/net/http/#Server"><code>Server</code></a>によって吸収され、クライアントは、HTTP仕様で要求されているように、空のメッセージボディを受け取ります。
</li>

<li>
The <a href="http://golang.org/pkg/os/exec/"><code>os/exec</code></a> package's 
<a href="http://golang.org/pkg/os/exec/#Cmd.StdinPipe"><code>Cmd.StdinPipe</code></a> method 
returns an <code>io.WriteCloser</code>, but has changed its concrete
implementation from <code>*os.File</code> to an unexported type that embeds
<code>*os.File</code>, and it is now safe to close the returned value.
Before Go 1.2, there was an unavoidable race that this change fixes.
Code that needs access to the methods of <code>*os.File</code> can use an
interface type assertion, such as <code>wc.(interface{ Sync() error })</code>.
</li>

<li>
<a href="http://golang.org/pkg/os/exec/"><code>os/exec</code></a>パッケージの<a href="http://golang.org/pkg/os/exec/#Cmd.StdinPipe"><code>Cmd.StdinPipe</code></a>メソッドは、<code>io.WriteCloser</code>を返しますが、<code>*os.File</code>から、<code>*os.File</code>を埋め込んだ非公開型に変更し、安全にcloseできるようになりました。(訳者補足：<a href="http://bit.ly/1cJgVxt">メインの差分箇所</a> )  Go1.2以前は、避けることの出来ない現象がありました。
<code>*os.File</code>のメソッド郡にアクセスを必要とするコードは、<code>wc.(interface{ Sync() error })</code>のように、interfaceの型アサーションを使うことができます。
</li>

<li>
The <a href="http://golang.org/pkg/runtime/"><code>runtime</code></a> package relaxes
the constraints on finalizer functions in
<a href="http://golang.org/pkg/runtime/#SetFinalizer"><code>SetFinalizer</code></a>: the
actual argument can now be any type that is assignable to the formal type of
the function, as is the case for any normal function call in Go.
</li>

<li>
<a href="http://golang.org/pkg/runtime/"><code>runtime</code></a> パッケージは、<a href="http://golang.org/pkg/runtime/#SetFinalizer"><code>SetFinalizer</code></a>で、finalizer関数の制約をゆるめました。実引数は、Goでの任意の通常の関数呼び出しの場合のように、関数の型に割り当て可能な任意の型になりました。
</li>

<li>
The <a href="http://golang.org/pkg/sort/"><code>sort</code></a> package has a new
<a href="http://golang.org/pkg/sort/#Stable"><code>Stable</code></a> function that implements
stable sorting. It is less efficient than the normal sort algorithm, however.
</li>

<li>
<a href="http://golang.org/pkg/sort/"><code>sort</code></a>パッケージは新しく<a href="http://golang.org/pkg/sort/#Stable"><code>Stable</code></a>関数を追加しました。これは安定ソートを実装しています。しかし、これは通常のソートアルゴリズムよりもあまり効率的ではありません。
</li>

<li>
The <a href="http://golang.org/pkg/strings/"><code>strings</code></a> package adds
an <a href="http://golang.org/pkg/strings/#IndexByte"><code>IndexByte</code></a>
function for consistency with the <a href="http://golang.org/pkg/bytes/"><code>bytes</code></a> package.
</li>

<li>
<a href="http://golang.org/pkg/strings/"><code>strings</code></a>パッケージは、<a href="http://golang.org/pkg/bytes/"><code>bytes</code></a>との整合性の為、<a href="http://golang.org/pkg/strings/#IndexByte"><code>IndexByte</code></a>関数を追加しました。
</li>

<li>
The <a href="http://golang.org/pkg/sync/atomic/"><code>sync/atomic</code></a> package
adds a new set of swap functions that atomically exchange the argument with the
value stored in the pointer, returning the old value.
The functions are
<a href="http://golang.org/pkg/sync/atomic/#SwapInt32"><code>SwapInt32</code></a>,
<a href="http://golang.org/pkg/sync/atomic/#SwapInt64"><code>SwapInt64</code></a>,
<a href="http://golang.org/pkg/sync/atomic/#SwapUint32"><code>SwapUint32</code></a>,
<a href="http://golang.org/pkg/sync/atomic/#SwapUint64"><code>SwapUint64</code></a>,
<a href="http://golang.org/pkg/sync/atomic/#SwapUintptr"><code>SwapUintptr</code></a>,
and
<a href="http://golang.org/pkg/sync/atomic/#SwapPointer"><code>SwapPointer</code></a>,
which swaps an <code>unsafe.Pointer</code>.
</li>

<li>
<a href="http://golang.org/pkg/sync/atomic/"><code>sync/atomic</code></a>パッケージは、新しいswap関数一式を追加しました。
<a href="http://golang.org/pkg/sync/atomic/#SwapInt32"><code>SwapInt32</code></a>,
<a href="http://golang.org/pkg/sync/atomic/#SwapInt64"><code>SwapInt64</code></a>,
<a href="http://golang.org/pkg/sync/atomic/#SwapUint32"><code>SwapUint32</code></a>,
<a href="http://golang.org/pkg/sync/atomic/#SwapUint64"><code>SwapUint64</code></a>,
<a href="http://golang.org/pkg/sync/atomic/#SwapUintptr"><code>SwapUintptr</code></a>,
<a href="http://golang.org/pkg/sync/atomic/#SwapPointer"><code>SwapPointer</code></a>(<code>unsafe.Pointer</code>をswapします)
これらは、ポインタに格納された値と引数を交換し、古い値を返します。
</li>

<li>
The <a href="http://golang.org/pkg/syscall/"><code>syscall</code></a> package now implements
<a href="http://golang.org/pkg/syscall/#Sendfile"><code>Sendfile</code></a> for Darwin.
</li>

<li>
<a href="http://golang.org/pkg/syscall/"><code>syscall</code></a>パッケージは、Darwinに<a href="http://golang.org/pkg/syscall/#Sendfile"><code>Sendfile</code></a>を実装しました。
(訳者補足：<a href="http://bit.ly/1k2Hils">Revision</a>)
</li>

<li>
The <a href="http://golang.org/pkg/testing/"><code>testing</code></a> package
now exports the <a href="http://golang.org/pkg/testing/#TB"><code>TB</code></a> interface.
It records the methods in common with the
<a href="http://golang.org/pkg/testing/#T"><code>T</code></a>
and
<a href="http://golang.org/pkg/testing/#B"><code>B</code></a> types,
to make it easier to share code between tests and benchmarks.
Also, the
<a href="http://golang.org/pkg/testing/#AllocsPerRun"><code>AllocsPerRun</code></a>
function now quantizes the return value to an integer (although it
still has type <code>float64</code>), to round off any error caused by
initialization and make the result more repeatable. 
</li>

<li>
<a href="http://golang.org/pkg/testing/"><code>testing</code></a>は、<a href="http://golang.org/pkg/testing/#TB"><code>TB</code></a>インターフェースを公開しました。それは、テストとベンチマーク間のコードを共有しやすくするために、<a href="http://golang.org/pkg/testing/#T"><code>T</code></a>型と<a href="http://golang.org/pkg/testing/#B"><code>B</code></a>型が持つ共通のメソッドを登録しています。
<a href="http://golang.org/pkg/testing/#AllocsPerRun"><code>AllocsPerRun</code></a>関数は、戻り値を(まだ<code>float64</code>ですが)整数にしました。初期化により引き起こされるエラーを丸め、結果をより再利用可能させます。
</li>

<li>
The <a href="http://golang.org/pkg/text/template/"><code>text/template</code></a> package
now automatically dereferences pointer values when evaluating the arguments
to "escape" functions such as "html", to bring the behavior of such functions
in agreement with that of other printing functions such as "printf".
</li>

<li>
<a href="http://golang.org/pkg/text/template/"><code>text/template</code></a>パッケージは、ポインターから参照先の値を取得します。when evaluating the arguments
to "escape" functions such as "html", to bring the behavior of such functions
in agreement with that of other printing functions such as "printf".
</li>

<li>
In the <a href="http://golang.org/pkg/time/"><code>time</code></a> package, the
<a href="http://golang.org/pkg/time/#Parse"><code>Parse</code></a> function
and
<a href="http://golang.org/pkg/time/#Time.Format"><code>Format</code></a>
method
now handle time zone offsets with seconds, such as in the historical
date "1871-01-01T05:33:02+00:34:08".
Also, pattern matching in the formats for those routines is stricter: a non-lowercase letter
must now follow the standard words such as "Jan" and "Mon".
</li>

<li>
<a href="http://golang.org/pkg/time/"><code>time</code></a>パッケージでは、<a href="http://golang.org/pkg/time/#Parse"><code>Parse</code></a>関数と<a href="http://golang.org/pkg/time/#Time.Format"><code>Format</code></a>メソッドは、歴史的な日"1871-01-01T05:33:02+00:34:08"のような秒が入ったタイムゾーンのオフセットを扱えるようになりました。(訳者補足：<a href="http://bit.ly/1fv2bTe">Revision</a>) また、フォーマットのパターンマッチがすこし厳しくなりました。最初が小文字でない"Jan"や "Mon"のような標準的な単語を使わなければなりません。
</li>

<li>
The <a href="http://golang.org/pkg/unicode/"><code>unicode</code></a> package
adds <a href="http://golang.org/pkg/unicode/#In"><code>In</code></a>,
a nicer-to-use but equivalent version of the original
<a href="http://golang.org/pkg/unicode/#IsOneOf"><code>IsOneOf</code></a>,
to see whether a character is a member of a Unicode category.
</li>

<li>
<a href="http://golang.org/pkg/unicode/"><code>unicode</code></a>パッケージは、<a href="http://golang.org/pkg/unicode/#In"><code>In</code></a>を追加しました。使いやすくなりましたが、最初からの<a href="http://golang.org/pkg/unicode/#IsOneOf"><code>IsOneOf</code></a>と等しいのですが、キャラクタがUnicodeカテゴリのメンバかどうか確認することができます。
</li>

</ul>