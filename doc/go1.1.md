<h2 id="introduction">Introduction to Go 1.1</h2>

<p>
2012年３月の<a href="/doc/go1.html">Go version 1</a>のリリースは、
Go言語とライブラリの安定版の新しい区切りとなりました。
その安定性は、世界中のGoユーザーのコミュニティとシステムを育ててきました。
それ以来、1.0.1,1.0.2,1.0.3をリリースしました。
これらのリリースは既知のバグの修正でしたが、
非クリティカルな実装はされていません。	
</p>

<p>
この新しいリリースGo 1.1は<a href="http://golang.org/doc/go1compat.html">互換性の約束</a>をキープしていまが、
数個の重要な言語変化を追加して(もちろん、後方互換)、ライブラリの変更(これも互換性あり)、
コンパイラ、ライブラリ、ランタイムの実装上の主要な仕事が含まれています。
焦点は、性能にあります。
ベンチマークは最善の方法ではありませんが、
多くのテストプログラムが劇的にスピードアップするのを見ることは重要です。
私たちは、これらのGoの組み込みや再コンパイルをアップデートすることで、
多くのユーザーのプログラムが改善されると信じています。
</p>

<p>
このドキュメントは、Go1とGo1.1の変更点を要約しています。
Very little if any code will need modification to run with Go 1.1,
although a couple of rare error cases surface with this release
and need to be addressed if they arise.
以下に詳細を示します。
特に<a href="#int">64-bit 整数</a> と 
<a href="#unicode_literals">Unicodeリテラル</a>
の議論があります。
</p>

<h2 id="language">Changes to the language</h2>

<p>
<a href="http://golang.org/doc/go1compat.html">Go 互換性のドキュメント </a>
では、Go1で書かれたプログラムは、継続して動作することを約束しています。
In the interest of firming up the specification, though, 
明確にされたエラーケースについて詳細があります。
新しい言語の特徴もあります。
</p>

<h3 id="divzero">ゼロ除算</h3>

<p>
Go1では、ゼロでの割り算はランタイムパニックになっていました。
</p>

<pre>
func f(x int) int {
	return x/0
}
</pre>

<p>
Go1.1では、ゼロでの割り算は正当なプログラムではありませんので、
コンパイルエラーとなります。
</p>

<h3 id="unicode_literals">Unicodeリテラルの代用</h3>

<p>
stringとruneリテラルの定義は、代用の半分をUnicodeの集合を除いて複製されてきました。
詳細は、<a href="#unicode">Unicode</a>セクションを見て下さい。
</p>

<h3 id="method_values">Method values</h3>
<p>
Go1.1は<a href="http://tip.golang.org/ref/spec#Method_values">method values</a>を実装しています。
それらは特定のレシーバの値に結び付けられる関数です。
例えば、<a href="/pkg/bufio/#Writer"><code>Writer</code></a>の値<code>w</code>が与えられ、
式<code>w.Write</code>のmethod value は、常に<code>w</code>に対して書き込む関数です。
それは<code>w</code>に関して閉じている関数リテラルと同じことです。
</p>
<pre9>
func (p []byte) (n int, err error) {
	return w.Write(p)
}
</pre>

<p>
Method values are distinct from method expressions, which generate functions
from methods of a given type; the method expression <code>(*bufio.Writer).Write</code>
is equivalent to a function with an extra first argument, a receiver of type
<code>(*bufio.Writer)</code>:
Method valuesはmethod expressionsとは違います。
与えられた型のメソッドから関数を作ります。
method expression <code>(*bufio.Writer).Write</code>は
第一引数、型<code>(*bufio.Writer)</code>のレシーバを持つ関数と同じです:
</p>

<pre>
func (w *bufio.Writer, p []byte) (n int, err error) {
	return w.Write(p)
}
</pre>

<p>
<em>Updating</em>: No existing code is affected; the change is strictly backward-compatible.
既存のコードには影響ありません。この変更は下位互換性があります。
</p>

<h3 id="return">Return requirements</h3>

<p>
Go1.1以前は、値を返す関数は、明示的な"return"か、
関数の最後に<code>panic</code>へのコールが必要でした。
これは、プログラマが関数の意味について明確にする簡単な方法でした。
しかし、無限の"for"ループだけの関数のように、最後の"return"が明らかに不要な多くのケースがあります。
</p>

<p>
In Go 1.1, the rule about final "return" statements is more permissive.
It introduces the concept of a
<a href="/ref/spec/#Terminating_statements"><em>terminating statement</em></a>,
a statement that is guaranteed to be the last one a function executes.
Examples include 
"for" loops with no condition and "if-else"
statements in which each half ends in a "return".
If the final statement of a function can be shown <em>syntactically</em> to
be a terminating statement, no final "return" statement is needed.

Go1.1では、最後の"return"ステートメントについてのルールは、より寛容になります。
ステートメントは、関数を実行する最後であることが保証される
<a href="http://tip.golang.org/ref/spec/#Terminating_statements"><em>terminating statement</em></a>
のコンセプトを取り入れています。

無条件の"for"ループや"if-else"ステートメントを含む半分の例は、"return"で終わります。
もし、関数の最後のステートメントが、構文的に(syntactically)、ステートメントの終わりであることが示されれば、
最後の"return"文は必要ありません。
</p>

<p>
Note that the rule is purely syntactic: it pays no attention to the values in the
code and therefore requires no complex analysis.

このルールは、純粋に構文規則であることに注意です。
それは、コード内の値に注意を払わなくていいので、
複雑な分析は必要ありません。

</p>

<p>
<em>Updating</em>: The change is backward-compatible, but existing code
with superfluous "return" statements and calls to <code>panic</code> may
be simplified manually.
Such code can be identified by <code>go vet</code>.


この変更は下位互換性がありますが、
なくてもよい"return"文と<code>panic</code>のコールを持つ既存のコードを、
手作業でシンプルにするかもしません。

そのようなコードは<code>go vet</code>で確認することができます。

</p>

<h2 id="impl">Changes to the implementations and tools
実装とツールの変更点</h2>


<h3 id="gccgo">Status of gccgo</h3>

<p>
The GCC release schedule does not coincide with the Go release schedule, so some skew is inevitable in
<code>gccgo</code>'s releases.
The 4.8.0 version of GCC shipped in March, 2013 and includes a nearly-Go 1.1 version of <code>gccgo</code>.
Its library is a little behind the release, but the biggest difference is that method values are not implemented.
Sometime around May 2013, we expect 4.8.1 of GCC to ship with a <code>gccgo</code>
providing a complete Go 1.1 implementaiton.
</p>
<p>
GCCのリリーススケジュールは、Goのリリーススケジュールと一致するとは限りませんので、
いくつかの歪みは避けれません。
GCCのバージョン4.8.0は2013年3月に送り出しました。
<code>gccgo</code>のGo1.1のバージョンに近いです。
そのライブラリはリリースがすこし遅れますが、
大きな変更は、method valuesが実行されないことです。
2013年5月頃には、Go1.1の完全なGoの実施を提供する
<code>gccgo</code>送り出すGCCの4.8.1を期待しています。
</p>

<h3 id="gc_flag">Command-line flag parsing
コマンドラインフラグのパース</h3>

<p>
In the gc tool chain, the compilers and linkers now use the
same command-line flag parsing rules as the Go flag package, a departure
from the traditional Unix flag parsing. This may affect scripts that invoke
the tool directly.
For example,
<code>go tool 6c -Fw -Dfoo</code> must now be written
<code>go tool 6c -F -w -D foo</code>. 
</p>
<p>
gcツールチェーンでは、コンパイラとリンカは、
Goのflagパッケージとしてのルールをパースする同じコマンドラインフラグを使用します。
伝統的なUNIXフラグパースからの新しい試みです。
これは、ツールを直接呼んでいるスクリプトに影響するかもしれません。
例えば<code>go tool 6c -Fw -Dfoo</code> は、<code>go tool 6c -F -w -D foo</code>と書かれなければなりません。
</p>

<h3 id="int">Size of int on 64-bit platforms</h3>

<p>
The language allows the implementation to choose whether the <code>int</code> type and
<code>uint</code> types are 32 or 64 bits. Previous Go implementations made <code>int</code>
and <code>uint</code> 32 bits on all systems. Both the gc and gccgo implementations
now make
<code>int</code> and <code>uint</code> 64 bits on 64-bit platforms such as AMD64/x86-64.
Among other things, this enables the allocation of slices with
more than 2 billion elements on 64-bit platforms.
</p>
<p>
言語は、<code>int</code>型と<code>uint</code>型が32bitか64bitかどうかを選択することができます。
以前のGo実装は、すべてのシステム上で32bitの<code>int</code>型と<code>uint</code>型を作った。
gcとgccgoの実装は、 AMD64/x86-64のような64bitプラットフォームでは64bitの
<code>int</code>型と<code>uint</code>型になります。
これは、64bitプラットフォーム上で20億以上の要素をもつスライスを割り当てることｄができます。
</p>
<p>
<em>Updating</em>:
Most programs will be unaffected by this change.
Because Go does not allow implicit conversions between distinct
<a href="/ref/spec/#Numeric_types">numeric types</a>,
no programs will stop compiling due to this change.
However, programs that contain implicit assumptions
that <code>int</code> is only 32 bits may change behavior.
For example, this code prints a positive number on 64-bit systems and
a negative one on 32-bit systems:
</p>
<p>
ほとんどのプログラムではこの変更は影響はないでしょう。
なぜなら、Goは、明示的な数値の型の間で暗黙的な変換を許容しないからです。
この変更によってプログラムがコンパイルを中止することはないでしょう。
しかしならが、<code>int</code>が32bitだけという過程を含むプログラムは、
 挙動が変更になるかもしれません。
 例えば、下記コードは64bit上では整数を出力し、32bit上では−1を出力します。
</p>
<pre>
x := ^uint32(0) // x is 0xffffffff
i := int(x)     // i is -1 on 32-bit systems, 0xffffffff on 64-bit
fmt.Println(i)
</pre>

<p>Portable code intending 32-bit sign extension (yielding <code>-1</code> on all systems)
would instead say:
</p>
<p>
32bitの符号拡張する（すべてのシステム上で-1が得られる）移植可能なコードは、以下のように表現します。
</p>

<pre>
i := int(int32(x))
</pre>

<h3 id="heap">Heap size on 64-bit architectures
64bitアーキテクチャ上のヒープサイズ</h3>

<p>
On 64-bit architectures only, the maximum heap size has been enlarged substantially,
from a few gigabytes to several tens of gigabytes.
(The exact details depend on the system and may change.)
</p>
<p>
64bitアーキテクチャ上だけでは、最大のヒープサイズは、
2,3GBから数十GBへ十分に大きくしてきました。
(詳細はシステムに依存し変更するかもしれません)
</p>	

<p>
On 32-bit architectures, the heap size has not changed.
</p>
<p>
32bitアーキテクチャでは、ヒープサイズは変更していません。
</p>

<p>
<em>Updating</em>:
This change should have no effect on existing programs beyond allowing them
to run with larger heaps.
</p>
<p>
この変更は既存のプログラムに影響はないはずです。
より大きいヒープで実行できます。
</p>

<h3 id="unicode">Unicode</h3>

<p>
To make it possible to represent code points greater than 65535 in UTF-16,
Unicode defines <em>surrogate halves</em>,
a range of code points to be used only in the assembly of large values, and only in UTF-16.
The code points in that surrogate range are illegal for any other purpose.
In Go 1.1, this constraint is honored by the compiler, libraries, and run-time:
a surrogate half is illegal as a rune value, when encoded as UTF-8, or when
encoded in isolation as UTF-16.
When encountered, for example in converting from a rune to UTF-8, it is
treated as an encoding error and will yield the replacement rune,
<a href="/pkg/unicode/utf8/#RuneError"><code>utf8.RuneError</code></a>,
U+FFFD.
</p>

<p>
This program,
</p>

<pre>
import "fmt"

func main() {
    fmt.Printf("%+q\n", string(0xD800))
}
</pre>

<p>
printed <code>"\ud800"</code> in Go 1.0, but prints <code>"\ufffd"</code> in Go 1.1.
</p>

<p>
Surrogate-half Unicode values are now illegal in rune and string constants, so constants such as
<code>'\ud800'</code> and <code>"\ud800"</code> are now rejected by the compilers.
When written explicitly as UTF-8 encoded bytes,
such strings can still be created, as in <code>"\xed\xa0\x80"</code>.
However, when such a string is decoded as a sequence of runes, as in a range loop, it will yield only <code>utf8.RuneError</code>
values.
</p>

<p>
The Unicode byte order marks U+FFFE and U+FEFF, encoded in UTF-8, are now permitted as the first
character of a Go source file.
Even though their appearance in the byte-order-free UTF-8 encoding is clearly unnecessary,
some editors add them as a kind of "magic number" identifying a UTF-8 encoded file.
</p>

<p>
<em>Updating</em>:
Most programs will be unaffected by the surrogate change.
Programs that depend on the old behavior should be modified to avoid the issue.
The byte-order-mark change is strictly backward-compatible.
</p>


<h3 id="race">Race detector</h3>

<p>
A major addition to the tools is a <em>race detector</em>, a way to find
bugs in programs caused by problems like concurrent changes to the same variable.
This new facility is built into the <code>go</code> tool.
For now, it is only available on Linux, Mac OS X, and Windows systems with
64-bit x86 processors.
To enable it, set the <code>-race</code> flag when building or testing your program 
(for instance, <code>go test -race</code>).
The race detector is documented in <a href="/doc/articles/race_detector.html">a separate article</a>.
</p>
<p>
ツールの主な追加は、<em>race detector</em>です。
これは、同じ変数を同時に変更するような問題で引き起こされる
プログラムのバグを見つける１つの方法です。
この新しい機能は、<code>go</code> ツールに組み込まれています。
今のところ、64bit x86プロセッサを持つLinux,Max OX X, Windows上で動きます。
これを可能にするには、プログラムをビルドするときやテストするときに、
<code>-race</code> フラグをセットします。
(例： <code>go test -race</code>)
このrace detectorのドキュメントは<a href="/doc/articles/race_detector.html">別の記事</a>にあります。

</p>
<h3 id="gc_asm">The gc assemblers</h3>

<p>
Due to the change of the <a href="#int"><code>int</code></a> to 64 bits and some other changes,
the arrangement of function arguments on the stack has changed in the gc tool chain.
Functions written in assembly will need to be revised at least
to adjust frame pointer offsets.
</p>
<p>
<a href="#int"><code>int</code></a>型の64bitへの変更といくつかの変更に伴い、
スタック上の関数の引数の配列は、gcツールチェーンに変更されました。
アセンブラで書かれた関数は、フレームポインタのオフセットを調整するために、
少なくとも改訂する必要があります。
</p>
<p>
<em>Updating</em>:
The <code>go vet</code> command now checks that functions implemented in assembly
match the Go function prototypes they implement.
</p>
<p>
 <code>go vet</code>コマンドは、アセンラブラで実装された関数が、
 Goの関数プロトタイプと一致することをチェックします。
</p>

<h3 id="gocmd">goコマンドの変更</h3>

<p>
The <a href="/cmd/go/"><code>go</code></a> command has acquired several
changes intended to improve the experience for new Go users.
</p>
<p>
<a href="/cmd/go/"><code>go</code></a>コマンドは、
新しいGoユーザーの体験を改善するための変更をしています。
</p>

<p>
First, when compiling, testing, or running Go code, the <code>go</code> command will now give more detailed error messages,
including a list of paths searched, when a package cannot be located.
</p>
<p>
一つ目に、Goのコードをコンパイル、テストあるいは実行するときに、
<code>go</code>コマンドは、
あるパッケージが配置されていなかったら、検索したパスのリストを含む
より詳細なエラーメッセージを与えてくれるでしょう。
</p>

<pre>
$ go build foo/quxx
can't load package: package foo/quxx: cannot find package "foo/quxx" in any of:
        /home/you/go/src/pkg/foo/quxx (from $GOROOT)
        /home/you/src/foo/quxx (from $GOPATH) 
</pre>

<p>
Second, the <code>go get</code> command no longer allows <code>$GOROOT</code>
as the default destination when downloading package source.
To use the <code>go get</code>
command, a valid <code>$GOPATH</code> is now required.
</p>
<p>
2番めに、<code>go get</code>コマンドは、
パッケージのソースをダウンロードするときに、
もはやデフォルトのインストール先として<code>$GOROOT</code>を許可しません。
<code>go get</code>コマンドを使うためには、有効な<code>$GOPATH</code>が必要になります。
</p>
<pre>
$ GOPATH= go get code.google.com/p/foo/quxx
package code.google.com/p/foo/quxx: cannot download, $GOPATH not set. For more details see: go help gopath 
</pre>

<p>
Finally, as a result of the previous change, the <code>go get</code> command will also fail
when <code>$GOPATH</code> and <code>$GOROOT</code> are set to the same value. 
</p>
<p>
最後に、先ほどの変更の結果として、<code>$GOPATH</code> と<code>$GOROOT</code>
に同じ値がセットされるとき、<code>go get</code>コマンドは失敗します。
</p>
<pre>
$ GOPATH=$GOROOT go get code.google.com/p/foo/quxx
warning: GOPATH set to GOROOT (/home/User/go) has no effect
package code.google.com/p/foo/quxx: cannot download, $GOPATH must not be set to $GOROOT. For more details see: go help gopath
</pre>

<h3 id="gotest">go testコマンドの変更</h3>

<p>
The <code>go test</code> command no longer deletes the binary when run with profiling enabled,
to make it easier to analyze the profile.
The implementation sets the <code>-c</code> flag automatically, so after running,
</p>
<p>
<code>go test</code>コマンドは、有効なプロファイルと一緒に実行するとき、
バイナリを削除しません。これは、プロファイルの解析を簡単にするためです。
実施するには、テストを実行後<code>-c</code> フラグを自動的にセットされます。
</p>

<pre>
$ go test -cpuprofile cpuprof.out mypackage
</pre>

<p>
the file <code>mypackage.test</code> will be left in the directory where <code>go test</code> was run.
</p>
<p>
<code>mypackage.test</code> は、<code>go test</code> が実行された
ディレクトリに残ったままとなっているでしょう。
</p>
<p>
The <code>go test</code> command can now generate profiling information
that reports where goroutines are blocked, that is,
where they tend to stall waiting for an event such as a channel communication.
The information is presented as a
<em>blocking profile</em>
enabled with the
<code>-blockprofile</code>
option of
<code>go test</code>.
Run <code>go help test</code> for more information.
</p>
<p>
<code>go test</code> コマンドは、
ゴルーチンがどこでブロックされるかなどのレポートするプロファイリング情報を生成することができます。
ブロックは、チャネル通信のようなイベントを機能停止する傾向にある場所です。
その情報は、<code>go test</code>オプションの<code>-blockprofile</code>
を使って<em>blocking profile</em>として表されます。
詳細は、<code>go help test</code>を実行してください。
</p>

<h3 id="gofix">go fixコマンドの変更</h3>

<p>
The <a href="/cmd/fix/"><code>fix</code></a> command, usually run as
<code>go fix</code>, no longer applies fixes to update code from
before Go 1 to use Go 1 APIs.
To update pre-Go 1 code to Go 1.1, use a Go 1.0 tool chain
to convert the code to Go 1.0 first.
</p>
<p>
通常は<code>go fix</code>として実行する
<a href="http://tip.golang.org/cmd/fix/"><code>fix</code></a>コマンドは、
Go1以前からGo1APIを使ったコードをアップデートする修正を適用しません。
</p>
<p>
Go 1以前のコードを Go 1.1へアップデートするには、Go1.0のツールチェーンを使い、
まずコードをGo1.0へ変換します。
</p>

<h3 id="gorun">go runコマンドの変更</h3>

<p>
The <code>go run</code> command now runs all files in the current working
directory if no file arguments are listed. Also, the <code>go run</code>
command now returns an error if test files are provided on the command line. In
this sense, "<code>go run</code>" replaces "<code>go run *.go</code>".
</p>
<p>
<code>go run</code>コマンドは、
ファイルの引数がなければ、カレントワーキングディレクトリにあるすべてのファイルを実行します。
<code>go run</code>コマンドは、
テストファイルがコマンドラインで提供されていれば、エラーを返します。
この意味では
<code>go run</code>" は"<code>go run *.go</code>を置き換えています。
</p>

<h3 id="platforms">追加プラットフォーム</h3>

<p>
The Go 1.1 tool chain adds experimental support for <code>freebsd/arm</code>,
<code>netbsd/386</code>, <code>netbsd/amd64</code>, <code>netbsd/arm</code>, 
<code>openbsd/386</code> and <code>openbsd/amd64</code> platforms.
</p>
<p>
Go 1.1 ツールチェーンは、下記プラットフォームを実験的なサポートを追加した。
 <code>freebsd/arm</code>,
<code>netbsd/386</code>, <code>netbsd/amd64</code>, <code>netbsd/arm</code>, 
<code>openbsd/386</code>,<code>openbsd/amd64</code> 
</p>
<p>
An ARMv6 or later processor is required for <code>freebsd/arm</code> or
<code>netbsd/arm</code>.
</p>
<p>
ARMv6あるいはそれ移行のプロセッサは<code>freebsd/arm</code> or
<code>netbsd/arm</code>が必要とされます。
</p>

<p>
Go 1.1 adds experimental support for <code>cgo</code> on <code>linux/arm</code>.
</p>
<p>
Go 1.1は <code>linux/arm</code>上の<code>cgo</code>を実験的にサポートします。
</p>

<h3 id="crosscompile">クロスコンパイル</h3>

<p>
When cross-compiling, the <code>go</code> tool will disable <code>cgo</code>
support by default.
</p>
<p>
クロスコンパイルするとき、<code>go</code>ツールはデフォルトでは
<code>cgo</code>をサポートしなくなります。
</p>

<p>
To explicitly enable <code>cgo</code>, set <code>CGO_ENABLED=1</code>.
</p>
<p>
<code>cgo</code>を明確に可能にするには、<code>CGO_ENABLED=1</code>を設定します。
</p>
<h2 id="performance">Performance</h2>

<p>
The performance of code compiled with the Go 1.1 gc tool suite should be noticeably
better for most Go programs.
Typical improvements relative to Go 1.0 seem to be about 30%-40%, sometimes
much more, but occasionally less or even non-existent.
There are too many small performance-driven tweaks through the tools and libraries
to list them all here, but the following major changes are worth noting:
</p>
<p>
Go 1.1のgcツールでコンパイルしたコードのパフォーマンスは、
ほとんどのプログラムにたいして目に見えてよくなっているはずです。
Go 1.0からの標準的な改善は、約30%から40%、ときおりそれ以上によくなると思われますが、
ときおり、改善がすくないか、改善されないこともあります。
ツールやライブラリを通して、多くの小さなパフォーマンスを手動とした調整がありますが、
注目に値する主な変更点は以下のとおりです。
</p>

<ul>
<li>The gc compilers generate better code in many cases, most noticeably for
floating point on the 32-bit Intel architecture.</li>
<li>The gc compilers do more in-lining, including for some operations
in the run-time such as <a href="/pkg/builtin/#append"><code>append</code></a>
and interface conversions.</li>
<li>There is a new implementation of Go maps with significant reduction in
memory footprint and CPU time.</li>
<li>The garbage collector has been made more parallel, which can reduce
latencies for programs running on multiple CPUs.</li>
<li>The garbage collector is also more precise, which costs a small amount of
CPU time but can reduce the size of the heap significantly, especially
on 32-bit architectures.</li>
<li>Due to tighter coupling of the run-time and network libraries, fewer
context switches are required on network operations.</li>
</ul>


<ul>
<li>gcコンパイラは、多くのケースでよりよいコードを生成します。
32bit intelアーキテクチャの浮動小数点に対しては、もっとも顕著です。</li>
<li>gcコンパイラは、内部で<a href="/pkg/builtin/#append"><code>append</code>
やインターフェースの変換のようなランタイムの作業を含む多くのことをしています。</li>
<li>メモリフットプリントとCPU時間の大幅な削減で、Goのマップの新しい実装があります。</li>
<li>ガベージコレクタは、複数のCPU上で動作するプログラムのための待ち時間を減らすことができ、
より多く並列化されています。</li>
<li>ガベージコレクタもより的確です。
少しの量のCPU時間のコストがかかりますが、
特に32bitアーキテクチャ上で、ヒープサイズを著しく減らすことができます。</li>
<li>ランタイムとネットワークライブラリの緊密な結合により
より少ないコンテキストの切り替えは、ネットワーク運用上必要とされます。</li>
</ul>
<h2 id="library">標準ライブラリの変更</h2>

<h3 id="bufio_scanner">bufio.Scanner</h3>

<p>
The various routines to scan textual input in the
<a href="/pkg/bufio/"><code>bufio</code></a>
package,
<a href="/pkg/bufio/#Reader.ReadBytes"><code>ReadBytes</code></a>,
<a href="/pkg/bufio/#Reader.ReadString"><code>ReadString</code></a>
and particularly
<a href="/pkg/bufio/#Reader.ReadLine"><code>ReadLine</code></a>,
are needlessly complex to use for simple purposes.
In Go 1.1, a new type,
<a href="/pkg/bufio/#Scanner"><code>Scanner</code></a>,
has been added to make it easier to do simple tasks such as
read the input as a sequence of lines or space-delimited words.
It simplifies the problem by terminating the scan on problematic
input such as pathologically long lines, and having a simple
default: line-oriented input, with each line stripped of its terminator.
Here is code to reproduce the input a line at a time:
</p>
<p>
<a href="http://tip.golang.org/pkg/bufio/"><code>bufio</code></a>パッケージでテキストインプットをスキャンする
様々な手順,
<a href="http://tip.golang.org/pkg/bufio/#Reader.ReadBytes"><code>ReadBytes</code></a>,
<a href="http://tip.golang.org/pkg/bufio/#Reader.ReadString"><code>ReadString</code></a>
特に
<a href="http://tip.golang.org/pkg/bufio/#Reader.ReadLine"><code>ReadLine</code></a>,
は、シンプルな目的に対して使うのが無駄に複雑です。
Go 1.1では新しい型<a href="http://tip.golang.org/pkg/bufio/#Scanner"><code>Scanner</code></a>
が追加されました。それは、線やスペースで区切られた一連の単語のようなインプットを読み取りやすくします。
異常に長い行のような不確かなインプットのスキャンを終えることによって問題を単純化します。
その終止符から取り除いた各行といっしょに行の方向づけされたインプットをもっています。
これは、一度に1行のインプットを再現するためのコードです。
</p>
<pre>
scanner := bufio.NewScanner(os.Stdin)
for scanner.Scan() {
    fmt.Println(scanner.Text()) // Println は最後に'\n'を追加します。
}
if err := scanner.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "reading standard input:", err)
}
</pre>

<p>
Scanning behavior can be adjusted through a function to control subdividing the input
(see the documentation for <a href="/pkg/bufio/#SplitFunc"><code>SplitFunc</code></a>),
but for tough problems or the need to continue past errors, the older interface
may still be required.
</p>
<p>
スキャンする動作は、インプットを分割制御する関数を通して調整されます
(<a href="/pkg/bufio/#SplitFunc"><code>SplitFunc</code></a>のドキュメントを見て下さい)。
しかし、頑固な問題や過去のエラーを継続するする必要があるため、
古いインターフェースはまだ必要かもしれません。	
</p>

<h3 id="net">net</h3>

<p>
The protocol-specific resolvers in the <a href="/pkg/net/"><code>net</code></a> package were formerly
lax about the network name passed in.
Although the documentation was clear
that the only valid networks for
<a href="/pkg/net/#ResolveTCPAddr"><code>ResolveTCPAddr</code></a>
are <code>"tcp"</code>,
<code>"tcp4"</code>, and <code>"tcp6"</code>, the Go 1.0 implementation silently accepted any string.
The Go 1.1 implementation returns an error if the network is not one of those strings.
The same is true of the other protocol-specific resolvers <a href="/pkg/net/#ResolveIPAddr"><code>ResolveIPAddr</code></a>,
<a href="/pkg/net/#ResolveUDPAddr"><code>ResolveUDPAddr</code></a>, and
<a href="/pkg/net/#ResolveUnixAddr"><code>ResolveUnixAddr</code></a>.
</p>
<p>
<a href="http://tip.golang.org/pkg/net/"><code>net</code></a>パッケージ内のプロトコル固有のレゾルバは、
ネットワーク名が渡されることについて以前はゆるかった。
ドキュメントにはっきりあるけれど、<a href="http://tip.golang.org/pkg/net/#ResolveTCPAddr"><code>ResolveTCPAddr</code></a>
の有効なネットワークは<code>"tcp"</code>,<code>"tcp4"</code>,<code>"tcp6"</code>
です。Go 1.0の実装は、任意の文字を受け入れていました。
Go 1.1の実装は、ネットワークがそれらの文字の１つになかったら、エラーを返します。
同じなのは、プロトコル固定のレゾルバ<a href="http://tip.golang.org/pkg/net/#ResolveIPAddr"><code>ResolveIPAddr</code></a>,
、<a href="http://tip.golang.org/pkg/net/#ResolveUDPAddr"><code>ResolveUDPAddr</code></a>、
<a href="http://tip.golang.org/pkg/net/#ResolveUnixAddr"><code>ResolveUnixAddr</code></a>です。
</p>

<p>
The previous implementation of
<a href="/pkg/net/#ListenUnixgram"><code>ListenUnixgram</code></a>
returned a
<a href="/pkg/net/#UDPConn"><code>UDPConn</code></a> as
a representation of the connection endpoint.
The Go 1.1 implementation instead returns a
<a href="/pkg/net/#UnixConn"><code>UnixConn</code></a>
to allow reading and writing
with its
<a href="/pkg/net/#UnixConn.ReadFrom"><code>ReadFrom</code></a>
and 
<a href="/pkg/net/#UnixConn.WriteTo"><code>WriteTo</code></a>
methods.
</p>
<p>
前の<a href="http://tip.golang.org/pkg/net/#ListenUnixgram"><code>ListenUnixgram</code></a>の実装は
コネクションエンドポイントの表現として<a href="http://tip.golang.org/pkg/net/#UDPConn"><code>UDPConn</code></a> 
を返していました。
Go 1.1の実装では、代わりに
<a href="http://tip.golang.org/pkg/net/#UnixConn.ReadFrom"><code>ReadFrom</code></a>と
<a href="http://tip.golang.org/pkg/net/#UnixConn.WriteTo"><code>WriteTo</code></a>メソッドを使って
読み込んだり書き込んだりできる
<a href="http://tip.golang.org/pkg/net/#UnixConn"><code>UnixConn</code></a>を返します。
</p>

<p>
The data structures
<a href="http://tip.golang.org/pkg/net/#IPAddr"><code>IPAddr</code></a>,
<a href="http://tip.golang.org/pkg/net/#TCPAddr"><code>TCPAddr</code></a>, and
<a href="http://tip.golang.org/pkg/net/#UDPAddr"><code>UDPAddr</code></a>
add a new string field called <code>Zone</code>.
Code using untagged composite literals (e.g. <code>net.TCPAddr{ip, port}</code>)
instead of tagged literals (<code>net.TCPAddr{IP: ip, Port: port}</code>)
will break due to the new field.
The Go 1 compatibility rules allow this change: client code must use tagged literals to avoid such breakages.
</p>
<p>
データ構造の
<a href="http://tip.golang.org/pkg/net/#IPAddr"><code>IPAddr</code></a>,
<a href="http://tip.golang.org/pkg/net/#TCPAddr"><code>TCPAddr</code></a>, 
<a href="http://tip.golang.org/pkg/net/#UDPAddr"><code>UDPAddr</code></a>
には、<code>Zone</code>と呼ぶ新しいstringフィールドを追加しています。
タグをつけたリテラル(例えば、<code>net.TCPAddr{IP: ip, Port: port}</code>)の代わりに
タグを付けていないリテラル(<code>net.TCPAddr{ip, port}</code>)を使っているコードは
新しいフィールドを追加したため、壊れます。
Go 1の互換性ルールはこの変化を許します。
クライアントコードは、このような破損を避けるためにタグをつけたリテラルを使わなければなりません。
</p>
<p>
<em>Updating</em>:
To correct breakage caused by the new struct field,
<code>go fix</code> will rewrite code to add tags for these types.
More generally, <code>go vet</code> will identify composite literals that
should be revised to use field tags.
</p>
<p>
<em>Updating</em>:
新しい構造体フィールドが原因による破損箇所を訂正することは、
<code>go fix</code> がそららの型にタグを追加して書きなおしてくれます。
より一般には、<code>go vet</code> は、フィールドタグを使うように修正すべき複合リテラルか確認します。
</p>
<h3 id="reflect">reflect</h3>

<p>
The <a href="/pkg/reflect/"><code>reflect</code></a> package has several significant additions.
</p>
<p>
<a href="http://tip.golang.org/pkg/reflect/"><code>reflect</code></a>パッケージは
いくつか重要な追加があります。
</p>

<p>
It is now possible to run a "select" statement using
the <code>reflect</code> package; see the description of
<a href="/pkg/reflect/#Select"><code>Select</code></a>
and
<a href="/pkg/reflect/#SelectCase"><code>SelectCase</code></a>
for details.
</p>
<p>
<code>reflect</code> パッケージを使って"select" ステートメントを実行することができるようになります。
詳細は、
<a href="/pkg/reflect/#Select"><code>Select</code></a>と
<a href="/pkg/reflect/#SelectCase"><code>SelectCase</code></a>
の説明を見て下さい。
</p>

<p>
The new method
<a href="/pkg/reflect/#Value.Convert"><code>Value.Convert</code></a>
(or
<a href="/pkg/reflect/#Type"><code>Type.ConvertibleTo</code></a>)
provides functionality to execute a Go conversion or type assertion operation
on a
<a href="/pkg/reflect/#Value"><code>Value</code></a>
(or test for its possibility).
</p>
<p>
新しいメソッド
<a href="/pkg/reflect/#Value.Convert"><code>Value.Convert</code></a>
(あるいは
<a href="/pkg/reflect/#Type"><code>Type.ConvertibleTo</code></a>)
は、<a href="/pkg/reflect/#Value"><code>Value</code></a>上で変換や
型アサーション（またはその可能性のためのテスト）を実行するための機能を提供します。
</p>
<p>
The new function
<a href="/pkg/reflect/#MakeFunc"><code>MakeFunc</code></a>
creates a wrapper function to make it easier to call a function with existing
<a href="/pkg/reflect/#Value"><code>Values</code></a>,
doing the standard Go conversions among the arguments, for instance
to pass an actual <code>int</code> to a formal <code>interface{}</code>.
</p>
<p>
新しい関数<a href="/pkg/reflect/#MakeFunc"><code>MakeFunc</code></a>は
存在する<a href="/pkg/reflect/#Value"><code>Values</code></a>を持つ関数をコールしやすくするための
ラッパー関数を作ります。
引数間で標準のGoの変換をします。
例えば、実際の<code>int</code>を表面的な<code>interface{}</code>へ渡します。
</p>

<p>
Finally, the new functions
<a href="/pkg/reflect/#ChanOf"><code>ChanOf</code></a>,
<a href="/pkg/reflect/#MapOf"><code>MapOf</code></a>
and
<a href="/pkg/reflect/#SliceOf"><code>SliceOf</code></a>
construct new
<a href="/pkg/reflect/#Type"><code>Types</code></a>
from existing types, for example to construct a the type <code>[]T</code> given
only <code>T</code>.
</p>
<p>
最後に、新しい関数
<a href="/pkg/reflect/#ChanOf"><code>ChanOf</code></a>,
<a href="/pkg/reflect/#MapOf"><code>MapOf</code></a>,
<a href="/pkg/reflect/#SliceOf"><code>SliceOf</code></a>
は、既存の型から新しい
<a href="/pkg/reflect/#Type"><code>Types</code></a>
を構築します。
例えば、<code>T</code>を与えた時、型<code>[]T</code>を構築します。
</p>

<h3 id="time">time</h3>
<p>
On FreeBSD, Linux, NetBSD, OS X and OpenBSD, previous versions of the
<a href="/pkg/time/"><code>time</code></a> package
returned times with microsecond precision.
The Go 1.1 implementation on these
systems now returns times with nanosecond precision.
Programs that write to an external format with microsecond precision
and read it back, expecting to recover the original value, will be affected
by the loss of precision.
There are two new methods of <a href="/pkg/time/#Time"><code>Time</code></a>,
<a href="/pkg/time/#Time.Round"><code>Round</code></a>
and
<a href="/pkg/time/#Time.Truncate"><code>Truncate</code></a>,
that can be used to remove precision from a time before passing it to
external storage.
</p>
<p>
FreeBSD, Linux, NetBSD, OS X, OpenBSD上では
以前のバージョンの<a href="/pkg/time/"><code>time</code></a>パッケージは、
マイクロ秒の精度の時間を返していました。
これらのシステム上でのGo 1.1の実装は、ナノ秒の精度の時間を返すようになります。
マイクロ秒の精度で外部形式への書き込み、
元の値に戻ることを期待し、それを読み戻すプログラムは、精度のロスによる影響を受けます。
<a href="/pkg/time/#Time"><code>Time</code></a>の新しい２つのメソッド
<a href="/pkg/time/#Time.Round"><code>Round</code></a>
と
<a href="/pkg/time/#Time.Truncate"><code>Truncate</code></a>
は、外部ストレージに渡す前に、時間から精度を削除するために使用することができます。
</p>

<p>
The new method
<a href="/pkg/time/#Time.YearDay"><code>YearDay</code></a>
returns the one-indexed integral day number of the year specified by the time value.
</p>

<p>
The
<a href="/pkg/time/#Timer"><code>Timer</code></a>
type has a new method
<a href="/pkg/time/#Timer.Reset"><code>Reset</code></a>
that modifies the timer to expire after a specified duration.
</p>
<p>
<a href="/pkg/time/#Timer"><code>Timer</code></a>型は
新しいメソッド
<a href="/pkg/time/#Timer.Reset"><code>Reset</code></a>
を持っています。これは指定した期間後に有効期限が切れるようにタイマーを変更します。
</p>
<p>
Finally, the new function
<a href="/pkg/time/#ParseInLocation"><code>ParseInLocation</code></a>
is like the existing
<a href="/pkg/time/#Parse"><code>Parse</code></a>
but parses the time in the context of a location (time zone), ignoring
time zone information in the parsed string.
This function addresses a common source of confusion in the time API.
</p>
<p>
最後に新しい関数
<a href="/pkg/time/#ParseInLocation"><code>ParseInLocation</code></a>
は既存の
<a href="/pkg/time/#Parse"><code>Parse</code></a>に似ています。
ロケーション（タイムゾーン）のコンテキスト内で時間をパースします。
パースされた文字列内のタイムゾーン情報は無視します。
この関数は、time APIの共通の混乱の原因に対処しています。
</p>
<p>
<em>Updating</em>:
Code that needs to read and write times using an external format with
lower precision should be modified to use the new methods.
</p>

<p>
<em>Updating</em>:
低い精度での外部フォーマットを使って時間を読んだり書いたりする必要があるコードは、
新しいメソッドを使って修正すべきです。
</p>

<h3 id="exp_old">Exp and old subtrees moved to go.exp and go.text subrepositories
Expと古いサブツリーはgo.exp と go.text のサブリポジトリに移動しました
</h3>

<p>
To make it easier for binary distributions to access them if desired, the <code>exp</code>
and <code>old</code> source subtrees, which are not included in binary distributions,
have been moved to the new <code>go.exp</code> subrepository at
<code>code.google.com/p/go.exp</code>. To access the <code>ssa</code> package,
for example, run
</p>
<p>
バイナリディストリビューションに含まれない<code>exp</code>と<code>old</code> ソースのサブツリーに
アクセスすること希望する場合、
バイナリディストリビューションに対してそれらにアクセスするには、
<code>code.google.com/p/go.exp</code>にある
新しい<code>go.exp</code> サブリポジトリに移動しました。
例えば、 <code>ssa</code>パッケージにアクセスするには、
</p>
<pre>
$ go get code.google.com/p/go.exp/ssa
</pre>

<p>
and then in Go source,
</p>
<p>
を実行して、それからGoのソースで
</p>

<pre>
import "code.google.com/p/go.exp/ssa"
</pre>
<p>
を記入します。
</p>

<p>
The old package <code>exp/norm</code> has also been moved, but to a new repository
<code>go.text</code>, where the Unicode APIs and other text-related packages will
be developed.
</p>
<p>
古いパッケージ<code>exp/norm</code>も移動しましたが、
<code>go.text</code>への移動です。
Unicode APIと他のテキストに関するパッケージが開発されます。

</p>
<h3 id="minor_library_changes">Minor changes to the library</h3>

<p>
The following list summarizes a number of minor changes to the library, mostly additions.
See the relevant package documentation for more information about each change.
</p>
<p>
以下のリストは、ライブラリへのマイナーチェンジ、おもに追加を要約したものです。
各変更についての詳細な情報は、関連するパッケージのドキュメントを見て下さい。
</p>

<ul>
<li> 
The <a href="/pkg/bytes/"><code>bytes</code></a> package has two new functions,
<a href="/pkg/bytes/#TrimPrefix"><code>TrimPrefix</code></a>
and
<a href="/pkg/bytes/#TrimSuffix"><code>TrimSuffix</code></a>,
with self-evident properties.
Also, the <a href="/pkg/bytes/#Buffer"><code>Buffer</code></a> type
has a new method
<a href="/pkg/bytes/#Buffer.Grow"><code>Grow</code></a> that
provides some control over memory allocation inside the buffer.
Finally, the
<a href="/pkg/bytes/#Reader"><code>Reader</code></a> type now has a
<a href="/pkg/strings/#Reader.WriteTo"><code>WriteTo</code></a> method
so it implements the 
<a href="/pkg/io/#WriterTo"><code>io.WriterTo</code></a> interface.

The <a href="/pkg/bytes/"><code>bytes</code></a> パッケージは２つの関数を追加しました。
わかりきったプロパティを持つ
<a href="/pkg/bytes/#TrimPrefix"><code>TrimPrefix</code></a>
と
<a href="/pkg/bytes/#TrimSuffix"><code>TrimSuffix</code></a>
です。

また、<a href="/pkg/bytes/#Buffer"><code>Buffer</code></a> 型は
新しいメソッド
<a href="/pkg/bytes/#Buffer.Grow"><code>Grow</code></a> 
を追加しました。
これはそのバッファ内のメモリの割り当てを拡張します。
最後に、
<a href="/pkg/bytes/#Reader"><code>Reader</code></a> 型は
<a href="/pkg/strings/#Reader.WriteTo"><code>WriteTo</code></a> メソッドを追加しました。
ですので、<a href="/pkg/bytes/#Reader"><code>Reader</code></a> 型は
<a href="/pkg/io/#WriterTo"><code>io.WriterTo</code></a> インターフェースを実装しています。

</li>

<li>
The <a href="/pkg/compress/gzip/"><code>compress/gzip</code></a> package has
a new <a href="/pkg/compress/gzip/#Writer.Flush"><code>Flush</code></a>
method for its
<a href="/pkg/compress/gzip/#Writer"><code>Writer</code></a>
type that flushes its underlying <code>flate.Writer</code>.

<a href="/pkg/compress/gzip/"><code>compress/gzip</code></a> パッケージは
<a href="/pkg/compress/gzip/#Writer"><code>Writer</code></a>型に対して
新しく
 <a href="/pkg/compress/gzip/#Writer.Flush"><code>Flush</code></a>
メソッドを追加しました。
that flushes its underlying <code>flate.Writer</code>.
</li>

<li>
The <a href="/pkg/crypto/hmac/"><code>crypto/hmac</code></a> package has a new function,
<a href="/pkg/crypto/hmac/#Equal"><code>Equal</code></a>, to compare two MACs.

The <a href="/pkg/crypto/hmac/"><code>crypto/hmac</code></a> パッケージは
2つのMACを比較する
<a href="/pkg/crypto/hmac/#Equal"><code>Equal</code></a>
を追加しました。
</li>

<li>
The <a href="/pkg/crypto/x509/"><code>crypto/x509</code></a> package
now supports PEM blocks (see
<a href="/pkg/crypto/x509/#DecryptPEMBlock"><code>DecryptPEMBlock</code></a> for instance),
and a new function
<a href="/pkg/crypto/x509/#ParseECPrivateKey"><code>ParseECPrivateKey</code></a> to parse elliptic curve private keys.

<a href="/pkg/crypto/x509/"><code>crypto/x509</code></a> パッケージは
PEM形式のブロックをサポートしました。
 (例えば、
<a href="/pkg/crypto/x509/#DecryptPEMBlock"><code>DecryptPEMBlock</code></a>
を見て下さい。 ),
それと、楕円曲線暗号のプライベートキーを解析する
<a href="/pkg/crypto/x509/#ParseECPrivateKey"><code>ParseECPrivateKey</code></a> .
をサポートしました。
</li>

<li>
The <a href="/pkg/database/sql/"><code>database/sql</code></a> package
has a new 
<a href="/pkg/database/sql/#DB.Ping"><code>Ping</code></a>
method for its
<a href="/pkg/database/sql/#DB"><code>DB</code></a>
type that tests the health of the connection.

<a href="/pkg/database/sql/"><code>database/sql</code></a> パッケージは
型<a href="/pkg/database/sql/#DB"><code>DB</code></a>に対して
<a href="/pkg/database/sql/#DB.Ping"><code>Ping</code></a>
メソッドを追加しました。
これは、接続の状態をテストします。
</li>

<li>
The <a href="/pkg/database/sql/driver/"><code>database/sql/driver</code></a> package
has a new
<a href="/pkg/database/sql/driver/#Queryer"><code>Queryer</code></a>
interface that a
<a href="/pkg/database/sql/driver/#Conn"><code>Conn</code></a>
may implement to improve performance.

<a href="/pkg/database/sql/driver/"><code>database/sql/driver</code></a> パッケージは
<a href="/pkg/database/sql/driver/#Queryer"><code>Queryer</code></a>
インターフェースを追加しました。
<a href="/pkg/database/sql/driver/#Conn"><code>Conn</code></a>
がパフォーマンスを改善するための実装となるかもしれません。
</li>

<li>
The <a href="/pkg/encoding/json/"><code>encoding/json</code></a> package's
<a href="/pkg/encoding/json/#Decoder"><code>Decoder</code></a>
has a new method
<a href="/pkg/encoding/json/#Decoder.Buffered"><code>Buffered</code></a>
to provide access to the remaining data in its buffer,
as well as a new method
<a href="/pkg/encoding/json/#Decoder.UseNumber"><code>UseNumber</code></a>
to unmarshal a value into the new type
<a href="/pkg/encoding/json/#Number"><code>Number</code></a>,
a string, rather than a float64.

<a href="/pkg/encoding/json/"><code>encoding/json</code></a> パッケージの
<a href="/pkg/encoding/json/#Decoder"><code>Decoder</code></a>は
<a href="/pkg/encoding/json/#Decoder.Buffered"><code>Buffered</code></a>
メソッドを追加しました。そのバッファ内に残っているデータへのアクセスを提供します。
同様に、値を新しい型
<a href="/pkg/encoding/json/#Number"><code>Number</code></a>(float64ではなくstring)
へ変換(unmarshal)する
<a href="/pkg/encoding/json/#Decoder.UseNumber"><code>UseNumber</code></a>
を追加しました。
</li>

<li>
The <a href="/pkg/encoding/xml/"><code>encoding/xml</code></a> package
has a new function,
<a href="/pkg/encoding/xml/#EscapeText"><code>EscapeText</code></a>,
which writes escaped XML output,
and a method on
<a href="/pkg/encoding/xml/#Encoder"><code>Encoder</code></a>,
<a href="/pkg/encoding/xml/#Encoder.Indent"><code>Indent</code></a>,
to specify indented output.

<a href="/pkg/encoding/xml/"><code>encoding/xml</code></a> パッケージ
はエスケープされたXML出力を書き出す<a href="/pkg/encoding/xml/#EscapeText"><code>EscapeText</code></a>
と
<a href="/pkg/encoding/xml/#Encoder"><code>Encoder</code></a>に対するインデントを指定する
<a href="/pkg/encoding/xml/#Encoder.Indent"><code>Indent</code></a>メソッドを追加しました。
</li>

<li>
<a href="/pkg/go/ast/"><code>go/ast</code></a> パッケージ内の
 <a href="/pkg/go/ast/#CommentMap"><code>CommentMap</code></a>メソッドと
 関連したメソッドは、Goプログラムにあるコメントを抽出と処理しやすくなります。
</li>

<li>
In the <a href="/pkg/go/doc/"><code>go/doc</code></a> package,
the parser now keeps better track of stylized annotations such as <code>TODO(joe)</code>
throughout the code,
information that the <a href="/cmd/godoc/"><code>godoc</code></a>
command can filter or present according to the value of the <code>-notes</code> flag.

<a href="/pkg/go/doc/"><code>go/doc</code></a>パッケージ内の
the parser now keeps better track of stylized annotations such as <code>TODO(joe)</code>
コード全体にわたって<code>TODO(joe)</code>のような形式化されたアノテーションを良くします。
<a href="/cmd/godoc/"><code>godoc</code></a>
コマンドは <code>-notes</code> フラグの値によってフィルタや表現することができます.
</li>

<li>
A new package, <a href="/pkg/go/format/"><code>go/format</code></a>, provides
a convenient way for a program to access the formatting capabilities of <code>gofmt</code>.
It has two functions,
<a href="/pkg/go/format/#Node"><code>Node</code></a> to format a Go parser
<a href="/pkg/go/ast/#Node"><code>Node</code></a>,
and
<a href="/pkg/go/format/#Source"><code>Source</code></a>
to format arbitrary Go source code.

新しいパッケージ<a href="/pkg/go/format/"><code>go/format</code></a>
は、<code>gofmt</code>のフォーマットするという特性にアクセスするための
プログラムの便利な方法を提供します。
このパッケージは２つの関数があります。
Go パーサー<a href="/pkg/go/ast/#Node"><code>Node</code></a>を
フォーマットする<a href="/pkg/go/format/#Node"><code>Node</code></a>関数と
任意のGo ソースコードをフォーマットする
<a href="/pkg/go/format/#Source"><code>Source</code></a>関数です。
</li>

<li>
The undocumented and only partially implemented "noescape" feature of the
<a href="/pkg/html/template/"><code>html/template</code></a>
package has been removed; programs that depend on it will break.

<a href="/pkg/html/template/"><code>html/template</code></a>パッケージの
ドキュメントされておらず、部分的にしか実装されていない"noescape"機能は削除されました。
これに依存するプログラムは壊れます。

</li>

<li>
The <a href="/pkg/image/jpeg/"><code>image/jpeg</code></a> package now
reads progressive JPEG files and handles a few more subsampling configurations.

<a href="/pkg/image/jpeg/"><code>image/jpeg</code></a> パッケージは
プログレッシブJPEGファイルを読み込んで、さらにいくつかのサブサンプリング構成を処理できるようになります。
</li>

<li>
The <a href="/pkg/io/"><code>io</code></a> package now exports the
<a href="/pkg/io/#ByteWriter"><code>io.ByteWriter</code></a> interface to capture the common
functionality of writing a byte at a time.

The <a href="/pkg/io/"><code>io</code></a> パッケージは
 一度に書き込む共通の機能を取り込むために
 <a href="/pkg/io/#ByteWriter"><code>io.ByteWriter</code></a>
インターフェースをエクスポートするようになりました。
</li>

<li>
The <a href="/pkg/log/syslog/"><code>log/syslog</code></a> package now provides better support
for OS-specific logging features.

<a href="/pkg/log/syslog/"><code>log/syslog</code></a> パッケージ
は、OS固有のログの特徴のよりよいサポートを提供します。
</li>

<li>
The <a href="/pkg/math/big/"><code>math/big</code></a> package's
<a href="/pkg/math/big/#Int"><code>Int</code></a> type now has
now has methods
<a href="/pkg/math/big/#Int.MarshalJSON"><code>MarshalJSON</code></a>
and
<a href="/pkg/math/big/#Int.UnmarshalJSON"><code>UnmarshalJSON</code></a>
to convert to and from a JSON representation.
Also,
<a href="/pkg/math/big/#Int"><code>Int</code></a>
can now convert directly to and from a <code>uint64</code> using
<a href="/pkg/math/big/#Int.Uint64"><code>Uint64</code></a>
and
<a href="/pkg/math/big/#Int.SetUint64"><code>SetUint64</code></a>,
while
<a href="/pkg/math/big/#Rat"><code>Rat</code></a>
can do the same with <code>float64</code> using
<a href="/pkg/math/big/#Rat.Float64"><code>Float64</code></a>
and
<a href="/pkg/math/big/#Rat.SetFloat64"><code>SetFloat64</code></a>.

<a href="/pkg/math/big/"><code>math/big</code></a> パッケージの
<a href="/pkg/math/big/#Int"><code>Int</code></a> 型は
JSON形式から変換するメソッド
<a href="/pkg/math/big/#Int.MarshalJSON"><code>MarshalJSON</code></a>
と
<a href="/pkg/math/big/#Int.UnmarshalJSON"><code>UnmarshalJSON</code></a>
を追加しました。

また、
<a href="/pkg/math/big/#Int"><code>Int</code></a>
は
<a href="/pkg/math/big/#Int.Uint64"><code>Uint64</code></a>
と
<a href="/pkg/math/big/#Int.SetUint64"><code>SetUint64</code></a>
を使って、
<code>uint64</code>へ直接変換できます。

<a href="/pkg/math/big/#Rat"><code>Rat</code></a>
は
<a href="/pkg/math/big/#Rat.Float64"><code>Float64</code></a>
と
<a href="/pkg/math/big/#Rat.SetFloat64"><code>SetFloat64</code></a>.
を使って、<code>float64</code>と同じことができます。
</li>

<li>
The <a href="/pkg/mime/multipart/"><code>mime/multipart</code></a> package
has a new method for its
<a href="/pkg/mime/multipart/#Writer"><code>Writer</code></a>,
<a href="/pkg/mime/multipart/#Writer.SetBoundary"><code>SetBoundary</code></a>,
to define the boundary separator used to package the output.

<a href="/pkg/mime/multipart/"><code>mime/multipart</code></a> パッケージは
<a href="/pkg/mime/multipart/#Writer"><code>Writer</code></a>,
に対する新しいメソッド
<a href="/pkg/mime/multipart/#Writer.SetBoundary"><code>SetBoundary</code></a>
を追加しました。
出力をパッケージ化するために使われる境界の区切り文字を定義します。
</li>

<li>
The
<a href="/pkg/net/"><code>net</code></a> package's
<a href="/pkg/net/#ListenUnixgram"><code>ListenUnixgram</code></a>
function has changed return types: it now returns a
<a href="/pkg/net/#UnixConn"><code>UnixConn</code></a>
rather than a
<a href="/pkg/net/#UDPConn"><code>UDPConn</code></a>, which was
clearly a mistake in Go 1.0.
Since this API change fixes a bug, it is permitted by the Go 1 compatibility rules.

<a href="/pkg/net/"><code>net</code></a> パッケージの
<a href="/pkg/net/#ListenUnixgram"><code>ListenUnixgram</code></a>
関数は型を返すように変更しました。
<a href="/pkg/net/#UDPConn"><code>UDPConn</code></a>
ではなく、
<a href="/pkg/net/#UnixConn"><code>UnixConn</code></a>
を返します。
これはGo1.0では明らかに間違いでした。
このAPIの修正はバグですので、
Go1互換性のルールによって許されています。
</li>

<li>
The <a href="/pkg/net/"><code>net</code></a> package includes a new type,
<a href="/pkg/net/#Dialer"><code>Dialer</code></a>, to supply options to
<a href="/pkg/net/#Dialer.Dial"><code>Dial</code></a>.

<a href="/pkg/net/"><code>net</code></a> パッケージに
新しい型
<a href="/pkg/net/#Dialer"><code>Dialer</code></a>
を追加しました。
<a href="/pkg/net/#Dialer.Dial"><code>Dial</code></a>
するオプションを提供するためです。
</li>

<li>
The <a href="/pkg/net/"><code>net</code></a> package adds support for
link-local IPv6 addresses with zone qualifiers, such as <code>fe80::1%lo0</code>.
The address structures <a href="/pkg/net/#IPAddr"><code>IPAddr</code></a>,
<a href="/pkg/net/#UDPAddr"><code>UDPAddr</code></a>, and
<a href="/pkg/net/#TCPAddr"><code>TCPAddr</code></a>
record the zone in a new field, and functions that expect string forms of these addresses, such as
<a href="/pkg/net/#Dial"><code>Dial</code></a>,
<a href="/pkg/net/#ResolveIPAddr"><code>ResolveIPAddr</code></a>,
<a href="/pkg/net/#ResolveUDPAddr"><code>ResolveUDPAddr</code></a>, and
<a href="/pkg/net/#ResolveTCPAddr"><code>ResolveTCPAddr</code></a>,
now accept the zone-qualified form.

<a href="/pkg/net/"><code>net</code></a> パッケージは
ゾーン修飾子を使用して<code>fe80::1%lo0</code>のような
リンクローカルIPv6アドレスのサポートが追加されます。

アドレスの構造<a href="/pkg/net/#IPAddr"><code>IPAddr</code></a>,
<a href="/pkg/net/#UDPAddr"><code>UDPAddr</code></a>
<a href="/pkg/net/#TCPAddr"><code>TCPAddr</code></a>
が、新しいフィールドのzoneを記録します。
functions that expect string forms of these addresses, such as
<a href="/pkg/net/#Dial"><code>Dial</code></a>,
<a href="/pkg/net/#ResolveIPAddr"><code>ResolveIPAddr</code></a>,
<a href="/pkg/net/#ResolveUDPAddr"><code>ResolveUDPAddr</code></a>, and
<a href="/pkg/net/#ResolveTCPAddr"><code>ResolveTCPAddr</code></a>,
now accept the zone-qualified form.


このような
<a href="/pkg/net/#Dial"><code>Dial</code></a>,
<a href="/pkg/net/#ResolveIPAddr"><code>ResolveIPAddr</code></a>,
<a href="/pkg/net/#ResolveUDPAddr"><code>ResolveUDPAddr</code></a>,
<a href="/pkg/net/#ResolveTCPAddr"><code>ResolveTCPAddr</code></a>,
として、これらのアドレスの文字列形式を期待する関数は、今ではゾーン修飾形式を受け入れます。
</li>

<li>
The <a href="/pkg/net/"><code>net</code></a> package adds
<a href="/pkg/net/#LookupNS"><code>LookupNS</code></a> to its suite of resolving functions.
<code>LookupNS</code> returns the <a href="/pkg/net/#NS">NS records</a> for a host name.

The <a href="/pkg/net/"><code>net</code></a> パッケージは
<a href="/pkg/net/#LookupNS"><code>LookupNS</code></a>
を追加しました。
<code>LookupNS</code>はホスト名に対する
<a href="/pkg/net/#NS">NS records</a>を返します。


</li>

<li>
The <a href="/pkg/net/"><code>net</code></a> package adds protocol-specific 
packet reading and writing methods to
<a href="/pkg/net/#IPConn"><code>IPConn</code></a>
(<a href="/pkg/net/#IPConn.ReadMsgIP"><code>ReadMsgIP</code></a>
and <a href="/pkg/net/#IPConn.WriteMsgIP"><code>WriteMsgIP</code></a>) and 
<a href="/pkg/net/#UDPConn"><code>UDPConn</code></a>
(<a href="/pkg/net/#UDPConn.ReadMsgUDP"><code>ReadMsgUDP</code></a> and
<a href="/pkg/net/#UDPConn.WriteMsgUDP"><code>WriteMsgUDP</code></a>).
These are specialized versions of <a href="/pkg/net/#PacketConn"><code>PacketConn</code></a>'s
<code>ReadFrom</code> and <code>WriteTo</code> methods that provide access to out-of-band data associated
with the packets.

<a href="/pkg/net/"><code>net</code></a>パッケージ
は読み書きできるメソッドを持ったプロトコル特有のパケットを追加しました。
<a href="/pkg/net/#IPConn"><code>IPConn</code></a>
(<a href="/pkg/net/#IPConn.ReadMsgIP"><code>ReadMsgIP</code></a>
、 <a href="/pkg/net/#IPConn.WriteMsgIP"><code>WriteMsgIP</code></a>)
と 
<a href="/pkg/net/#UDPConn"><code>UDPConn</code></a>
(<a href="/pkg/net/#UDPConn.ReadMsgUDP"><code>ReadMsgUDP</code></a> 、
<a href="/pkg/net/#UDPConn.WriteMsgUDP"><code>WriteMsgUDP</code></a>)
です。
これらはパケットと関連する帯域外のデータへのアクセスを提供する
<a href="/pkg/net/#PacketConn"><code>PacketConn</code></a>
<code>ReadFrom</code> and <code>WriteTo</code> メソッドのバージョンを
特殊化したものです。
 </li>
 
 <li>
The <a href="/pkg/net/"><code>net</code></a> package adds methods to
<a href="/pkg/net/#UnixConn"><code>UnixConn</code></a> to allow closing half of the connection 
(<a href="/pkg/net/#UnixConn.CloseRead"><code>CloseRead</code></a> and
<a href="/pkg/net/#UnixConn.CloseWrite"><code>CloseWrite</code></a>),
matching the existing methods of <a href="/pkg/net/#TCPConn"><code>TCPConn</code></a>.

<a href="/pkg/net/"><code>net</code></a> パッケージは 
<a href="/pkg/net/#TCPConn"><code>TCPConn</code></a>の既存のメソッドに合わせて、
<a href="/pkg/net/#UnixConn"><code>UnixConn</code></a> 
の接続の半分を閉じれるように、メソッド
(<a href="/pkg/net/#UnixConn.CloseRead"><code>CloseRead</code></a> と
<a href="/pkg/net/#UnixConn.CloseWrite"><code>CloseWrite</code></a>)
を追加しました。
</li>
 
<li>
The <a href="/pkg/net/http/"><code>net/http</code></a> package includes several new additions.
<a href="/pkg/net/http/#ParseTime"><code>ParseTime</code></a> parses a time string, trying
several common HTTP time formats.
The <a href="/pkg/net/http/#Request.PostFormValue">PostFormValue</a> method of
<a href="/pkg/net/http/#Request"><code>Request</code></a> is like
<a href="/pkg/net/http/#Request.FormValue"><code>FormValue</code></a> but ignores URL parameters.
The <a href="/pkg/net/http/#CloseNotifier"><code>CloseNotifier</code></a> interface provides a mechanism
for a server handler to discover when a client has disconnected.
The <code>ServeMux</code> type now has a
<a href="/pkg/net/http/#ServeMux.Handler"><code>Handler</code></a> method to access a path's
<code>Handler</code> without executing it.
The <code>Transport</code> can now cancel an in-flight request with
<a href="/pkg/net/http/#Transport.CancelRequest"><code>CancelRequest</code></a>.
Finally, the Transport is now more aggressive at closing TCP connections when
a <a href="/pkg/net/http/#Response"><code>Response.Body</code></a> is closed before
being fully consumed.

<a href="/pkg/net/http/"><code>net/http</code></a> 
パッケージは、新しくいくつか追加しました。
<a href="/pkg/net/http/#ParseTime"><code>ParseTime</code></a>
は、いくつかの共通のHTTPタイムフォーマットにしようと、時間文字列をパースします。
<a href="/pkg/net/http/#Request"><code>Request</code></a> の
<a href="/pkg/net/http/#Request.PostFormValue">PostFormValue</a> メソッドは
<a href="/pkg/net/http/#Request.FormValue"><code>FormValue</code></a>
に似ていますが、URLパラメータを無視します。

<a href="/pkg/net/http/#CloseNotifier"><code>CloseNotifier</code></a>
インターフェースは、クライアントが切断した時に検知するための
サーバハンドラに対する仕組みを提供します。

<code>ServeMux</code> 型はそれを実行せずにパスの
<code>Handler</code>にアクセスするための、
<a href="/pkg/net/http/#ServeMux.Handler"><code>Handler</code></a> メソッド
 を追加しました。

<code>Transport</code> can now cancel an in-flight request with
<a href="/pkg/net/http/#Transport.CancelRequest"><code>CancelRequest</code></a>
でin-flightリクエストをキャンセルできるようになりました。

最後に、
<a href="/pkg/net/http/#Response"><code>Response.Body</code></a>
が完全に使い終わる前にを閉じられるとき、
Transportは、TCP接続を閉じることにより積極的になりました。
</li>

<li>
新しい<a href="/pkg/net/http/cookiejar/"><code>net/http/cookiejar</code></a> 
パッケージは、HTTPクッキーを管理する基本的なもの提供します。
</li>

<li>
The <a href="/pkg/net/mail/"><code>net/mail</code></a> package has two new functions,
<a href="/pkg/net/mail/#ParseAddress"><code>ParseAddress</code></a> and
<a href="/pkg/net/mail/#ParseAddressList"><code>ParseAddressList</code></a>,
to parse RFC 5322-formatted mail addresses into
<a href="/pkg/net/mail/#Address"><code>Address</code></a> structures.

<a href="/pkg/net/mail/"><code>net/mail</code></a> 
パッケージは、
RFC 5322形式のメールアドレスを解析し、
<a href="/pkg/net/mail/#Address"><code>Address</code></a>
構造体へ格納する2つの関数
<a href="/pkg/net/mail/#ParseAddress"><code>ParseAddress</code></a>と
<a href="/pkg/net/mail/#ParseAddressList"><code>ParseAddressList</code></a>
を追加しました。
</li>

<li>
The <a href="/pkg/net/smtp/"><code>net/smtp</code></a> package's
<a href="/pkg/net/smtp/#Client"><code>Client</code></a> type has a new method,
<a href="/pkg/net/smtp/#Client.Hello"><code>Hello</code></a>,
which transmits a <code>HELO</code> or <code>EHLO</code> message to the server.
</li>

<li>
The <a href="/pkg/net/textproto/"><code>net/textproto</code></a> package
has two new functions,
<a href="/pkg/net/textproto/#TrimBytes"><code>TrimBytes</code></a> and
<a href="/pkg/net/textproto/#TrimString"><code>TrimString</code></a>,
which do ASCII-only trimming of leading and trailing spaces.
</li>

<li>
The new method <a href="/pkg/os/#FileMode.IsRegular"><code>os.FileMode.IsRegular</code></a> makes it easy to ask if a file is a plain file.
</li>

<li>
The <a href="/pkg/os/signal/"><code>os/signal</code></a> package has a new function,
<a href="/pkg/os/signal/#Stop"><code>Stop</code></a>, which stops the package delivering
any further signals to the channel.
</li>

<li>
The <a href="/pkg/regexp/"><code>regexp</code></a> package
now supports Unix-original leftmost-longest matches through the
<a href="/pkg/regexp/#Regexp.Longest"><code>Regexp.Longest</code></a>
method, while
<a href="/pkg/regexp/#Regexp.Split"><code>Regexp.Split</code></a> slices
strings into pieces based on separators defined by the regular expression.
</li>

<li>
The <a href="/pkg/runtime/debug/"><code>runtime/debug</code></a> package
has three new functions regarding memory usage.
The <a href="/pkg/runtime/debug/#FreeOSMemory"><code>FreeOSMemory</code></a>
function triggers a run of the garbage collector and then attempts to return unused
memory to the operating system;
the <a href="/pkg/runtime/debug/#ReadGCStats"><code>ReadGCStats</code></a>
function retrieves statistics about the collector; and
<a href="/pkg/runtime/debug/#SetGCPercent"><code>SetGCPercent</code></a>
provides a programmatic way to control how often the collector runs,
including disabling it altogether.
</li>

<li>
The <a href="/pkg/sort/"><code>sort</code></a> package has a new function,
<a href="/pkg/sort/#Reverse"><code>Reverse</code></a>.
Wrapping the argument of a call to 
<a href="/pkg/sort/#Sort"><code>sort.Sort</code></a>
with a call to <code>Reverse</code> causes the sort order to be reversed.
</li>

<li>
The <a href="/pkg/strings/"><code>strings</code></a> package has two new functions,
<a href="/pkg/strings/#TrimPrefix"><code>TrimPrefix</code></a>
and
<a href="/pkg/strings/#TrimSuffix"><code>TrimSuffix</code></a>
with self-evident properties, and the new method
<a href="/pkg/strings/#Reader.WriteTo"><code>Reader.WriteTo</code></a> so the
<a href="/pkg/strings/#Reader"><code>Reader</code></a>
type now implements the
<a href="/pkg/io/#WriterTo"><code>io.WriterTo</code></a> interface.
</li>

<li>
The <a href="/pkg/syscall/"><code>syscall</code></a> package has received many updates to make it more inclusive of constants and system calls for each supported operating system.
</li>

<li>
The <a href="/pkg/testing/"><code>testing</code></a> package now automates the generation of allocation
statistics in tests and benchmarks using the new
<a href="/pkg/testing/#AllocsPerRun"><code>AllocsPerRun</code></a> function. And the
<a href="/pkg/testing/#B.ReportAllocs"><code>ReportAllocs</code></a>
method on <a href="/pkg/testing/#B"><code>testing.B</code></a> will enable printing of
memory allocation statistics for the calling benchmark. It also introduces the
<a href="/pkg/testing/#BenchmarkResult.AllocsPerOp"><code>AllocsPerOp</code></a> method of
<a href="/pkg/testing/#BenchmarkResult"><code>BenchmarkResult</code></a>.
There is also a new
<a href="/pkg/testing/#Verbose"><code>Verbose</code></a> function to test the state of the <code>-v</code>
command-line flag,
and a new
<a href="/pkg/testing/#B.Skip"><code>Skip</code></a> method of
<a href="/pkg/testing/#B"><code>testing.B</code></a> and
<a href="/pkg/testing/#T"><code>testing.T</code></a>
to simplify skipping an inappropriate test.
</li>

<li>
In the <a href="/pkg/text/template/"><code>text/template</code></a>
and
<a href="/pkg/html/template/"><code>html/template</code></a> packages,
templates can now use parentheses to group the elements of pipelines, simplifying the construction of complex pipelines.
Also, as part of the new parser, the
<a href="/pkg/text/template/parse/#Node"><code>Node</code></a> interface got two new methods to provide
better error reporting.
Although this violates the Go 1 compatibility rules,
no existing code should be affected because this interface is explicitly intended only to be used
by the
<a href="/pkg/text/template/"><code>text/template</code></a>
and
<a href="/pkg/html/template/"><code>html/template</code></a>
packages and there are safeguards to guarantee that.
</li>

<li>
The implementation of the <a href="/pkg/unicode/"><code>unicode</code></a> package has been updated to Unicode version 6.2.0.
</li>

<li>
In the <a href="/pkg/unicode/utf8/"><code>unicode/utf8</code></a> package,
the new function <a href="/pkg/unicode/utf8/#ValidRune"><code>ValidRune</code></a> reports whether the rune is a valid Unicode code point.
To be valid, a rune must be in range and not be a surrogate half.
</li>
</ul>
