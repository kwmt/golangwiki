<h2 id="introduction">Introduction to Go 1.2</h2>

2013年4月に<a href="/doc/go1.1.html">Go version 1.1</a>をリリースして以来、リリーススケジュールが
リリースのプロセスをより効率的にするために、短くされていました。
このリリースGo バージョン1.2 （短縮表記すると、Go1.2）では、だいたい1.1以降6ヶ月になります。
1.1は1.0がリリースされてから1年以上かかりました。
タイムスケールが短くなった理由は、1.2は1.0から1.1の段階より差分が少ないからです。
しかし、よりよいスケジューラと1つの新しい言語の特徴を含む、重要な変更があります。

もちろんGo1.2は<a href="/doc/go1compat.html">promise
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
    X [1<<24]byte
    Field int32
}

func main() {
    var x *T
    ...
}
</pre>

<p>
<code>nil</code>ポインタ<code>x</code>は間違ったメモリにアクセスできます:
<code>x.Field</code>はアドレス<code>1<<24</code>のメモリにアクセスできます。
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
Most code that depended on the old behavior is erroneous and will fail when run.
Such programs will need to be updated by hand.
古い仕様に依存するほとんどのコードはエラーとなり、実行すると失敗します。
そのようなプログラムは、手動で更新する必要があります。
</p>

<h3 id="three_index">Three-index slices</h3>

<p>
Go 1.2 adds the ability to specify the capacity as well as the length when using a slicing operation
on an existing array or slice.
A slicing operation creates a new slice by describing a contiguous section of an already-created array or slice:
</p>

<pre>
var array [10]int
slice := array[2:4]
</pre>

<p>
The capacity of the slice is the maximum number of elements that the slice may hold, even after reslicing;
it reflects the size of the underlying array.
In this example, the capacity of the <code>slice</code> variable is 8.
</p>

<p>
Go 1.2 adds new syntax to allow a slicing operation to specify the capacity as well as the length.
A second
colon introduces the capacity value, which must be less than or equal to the capacity of the
source slice or array, adjusted for the origin. For instance,
</p>

<pre>
slice = array[2:4:6]
</pre>

<p>
sets the slice to have the same length as in the earlier example but its capacity is now only 4 elements (6-2).
It is impossible to use this new slice value to access the last two elements of the original array.
</p>

<p>
In this three-index notation, a missing first index (<code>[:i:j]</code>) defaults to zero but the other
two indices must always be specified explicitly.
It is possible that future releases of Go may introduce default values for these indices.
</p>

<p>
Further details are in the
<a href="http://golang.org/s/go12slice">design document</a>.
</p>

<p>
<em>Updating</em>:
This is a backwards-compatible change that affects no existing programs.
</p>

<h2 id="impl">Changes to the implementations and tools</h2>

<h3 id="preemption">Pre-emption in the scheduler</h3>

<p>
In prior releases, a goroutine that was looping forever could starve out other
goroutines on the same thread, a serious problem when GOMAXPROCS
provided only one user thread.
In Go 1.2, this is partially addressed: The scheduler is invoked occasionally
upon entry to a function.
This means that any loop that includes a (non-inlined) function call can
be pre-empted, allowing other goroutines to run on the same thread.
</p>

<h3 id="cgo_and_cpp">Cgo and C++</h3>

<p>
The <a href="/cmd/cgo/"><code>cgo</code></a> command will now invoke the C++
compiler to build any pieces of the linked-to library that are written in C++; the
documentation has more detail.
</p>

<h3 id="go_tools_godoc">Godoc and vet moved to the go.tools subrepository</h3>

<p>
Both binaries are still included with the distribution, but the source code for the
godoc and vet commands has moved to the
<a href="http://code.google.com/p/go.tools">go.tools</a> subrepository.
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
<em>Updating</em>:
Since godoc and vet are not part of the library,
no client Go code depends on the their source and no updating is required.
</p>

<p>
The binary distributions available from <a href="http://golang.org">golang.org</a>
include these binaries, so users of these distributions are unaffected.
</p>

<p>
When building from source, users must use "go get" to install godoc and vet.
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

<h3 id="gc_changes">Changes to the gc compiler and linker</h3>

<p>
Go 1.2 has several semantic changes to the workings of the gc compiler suite.
Most users will be unaffected by them.
</p>

<p>
The <a href="/cmd/cgo/"><code>cgo</code></a> command now
works when C++ is included in the library being linked against.
See the <a href="/cmd/cgo/"><code>cgo</code></a> documentation
for details.
</p>

<p>
The gc compiler displayed a vestigial detail of its origins when
a program had no <code>package</code> clause: it assumed
the file was in package <code>main</code>.
The past has been erased, and a missing <code>package</code> clause
is now an error.
</p>

<p>
On the ARM, the toolchain supports "external linking", which
is a step towards being able to build shared libraries with the gc
tool chain and to provide dynamic linking support for environments
in which that is necessary.
</p>

<p>
In the runtime for the ARM, with <code>5a</code>, it used to be possible to refer
to the runtime-internal <code>m</code> (machine) and <code>g</code>
(goroutine) variables using <code>R9</code> and <code>R10</code> directly.
It is now necessary to refer to them by their proper names.
</p>

<p>
Also on the ARM, the <code>5l</code> linker (sic) now defines the
<code>MOVBS</code> and <code>MOVHS</code> instructions
as synonyms of <code>MOVB</code> and <code>MOVH</code>,
to make clearer the separation between signed and unsigned
sub-word moves; the unsigned versions already existed with a
<code>U</code> suffix.
</p>

<h3 id="cover">Test coverage</h3>

<p>
One major new feature of <a href="/pkg/go/"><code>go test</code></a> is
that it can now compute and, with help from a new, separately installed
"go tool cover" program, display test coverage results.
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
The cover tool does two things.
First, when "go test" is given the <code>-cover</code> flag, it is run automatically 
to rewrite the source for the package and insert instrumentation statements.
The test is then compiled and run as usual, and basic coverage statistics are reported:
</p>

<pre>
$ go test -cover fmt
ok  	fmt	0.060s	coverage: 91.4% of statements
$
</pre>

<p>
Second, for more detailed reports, different flags to "go test" can create a coverage profile file,
which the cover program, invoked with "go tool cover", can then analyze.
</p>

<p>
Details on how to generate and analyze coverage statistics can be found by running the commands
</p>

<pre>
$ go help testflag
$ go tool cover -help
</pre>

<h3 id="go_doc">The go doc command is deleted</h3>

<p>
The "go doc" command is deleted.
Note that the <a href="/cmd/godoc/"><code>godoc</code></a> tool itself is not deleted,
just the wrapping of it by the <a href="/cmd/go/"><code>go</code></a> command.
All it did was show the documents for a package by package path,
which godoc itself already does with more flexibility.
It has therefore been deleted to reduce the number of documentation tools and,
as part of the restructuring of godoc, encourage better options in future.
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

<pre>
$ godoc .
</pre>

<h3 id="gocmd">Changes to the go command</h3>

<p>
The <a href="/cmd/go/"><code>go get</code></a> command
now has a <code>-t</code> flag that causes it to download the dependencies
of the tests run by the package, not just those of the package itself.
By default, as before, dependencies of the tests are not downloaded.
</p>

<h2 id="performance">Performance</h2>

<p>
There are a number of significant performance improvements in the standard library; here are a few of them.
</p>

<ul> 

<li>
The <a href="/pkg/compress/bzip2/"><code>compress/bzip2</code></a>
decompresses about 30% faster.
</li>

<li>
The <a href="/pkg/crypto/des/"><code>crypto/des</code></a> package
is about five times faster.
</li>

<li>
The <a href="/pkg/encoding/json/"><code>encoding/json</code></a> package
encodes about 30% faster.
</li>

<li>
Networking performance on Windows and BSD systems is about 30% faster through the use
of an integrated network poller in the runtime, similar to what was done for Linux and OS X
in Go 1.1.
</li>

</ul>

<h2 id="library">Changes to the standard library</h2>


<h3 id="archive_tar_zip">The archive/tar and archive/zip packages</h3>

<p>
The
<a href="/pkg/archive/tar/"><code>archive/tar</code></a>
and
<a href="/pkg/archive/zip/"><code>archive/zip</code></a>
packages have had a change to their semantics that may break existing programs.
The issue is that they both provided an implementation of the
<a href="/pkg/os/#FileInfo"><code>os.FileInfo</code></a>
interface that was not compliant with the specification for that interface.
In particular, their <code>Name</code> method returned the full
path name of the entry, but the interface specification requires that
the method return only the base name (final path element).
</p>

<p>
<em>Updating</em>: Since this behavior was newly implemented and
a bit obscure, it is possible that no code depends on the broken behavior.
If there are programs that do depend on it, they will need to be identified
and fixed manually.
</p>

<h3 id="encoding">The new encoding package</h3>

<p>
There is a new package, <a href="/pkg/encoding/"><code>encoding</code></a>,
that defines a set of standard encoding interfaces that may be used to
build custom marshalers and unmarshalers for packages such as
<a href="/pkg/encoding/xml/"><code>encoding/xml</code></a>,
<a href="/pkg/encoding/json/"><code>encoding/json</code></a>,
and
<a href="/pkg/encoding/binary/"><code>encoding/binary</code></a>.
These new interfaces have been used to tidy up some implementations in
the standard library.
</p>

<p>
The new interfaces are called
<a href="/pkg/encoding/#BinaryMarshaler"><code>BinaryMarshaler</code></a>,
<a href="/pkg/encoding/#BinaryUnmarshaler"><code>BinaryUnmarshaler</code></a>,
<a href="/pkg/encoding/#TextMarshaler"><code>TextMarshaler</code></a>,
and
<a href="/pkg/encoding/#TextUnmarshaler"><code>TextUnmarshaler</code></a>.
Full details are in the <a href="/pkg/encoding/">documentation</a> for the package
and a separate <a href="http://golang.org/s/go12encoding">design document</a>.
</p>

<h3 id="fmt_indexed_arguments">The fmt package</h3>

<p>
The <a href="/pkg/fmt/"><code>fmt</code></a> package's formatted print
routines such as <a href="/pkg/fmt/#Printf"><code>Printf</code></a>
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
For example, the normal <code>Printf</code> call
</p>

<pre>
fmt.Sprintf("%c %c %c\n", 'a', 'b', 'c')
</pre>

<p>
would create the string <code>"a b c"</code>, but with indexing operations like this,
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
The motivation for this feature is programmable format statements to access
the arguments in different order for localization, but it has other uses:
</p>

<pre>
log.Printf("trace: value %v of type %[1]T\n", expensiveFunction(a.b[c]))
</pre>

<p>
<em>Updating</em>: The change to the syntax of format specifications
is strictly backwards compatible, so it affects no working programs.
</p>

<h3 id="text_template">The text/template and html/template packages</h3>

<p>
The
<a href="/pkg/text/template/"><code>text/template</code></a> package
has a couple of changes in Go 1.2, both of which are also mirrored in the
<a href="/pkg/html/template/"><code>html/template</code></a> package.
</p>

<p>
First, there are new default functions for comparing basic types.
The functions are listed in this table, which shows their names and
the associated familiar comparison operator.
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

<pre>
{{if eq .A 1 2 3 }} equal {{else}} not equal {{end}}
</pre>

<p>
reports "equal" if <code>.A</code> is equal to <em>any</em> of 1, 2, or 3.
</p>

<p>
The second change is that a small addition to the grammar makes "if else if" chains easier to write.
Instead of writing,
</p>

<pre>
{{if eq .A 1}} X {{else}} {{if eq .A 2}} Y {{end}} {{end}} 
</pre>

<p>
one can fold the second "if" into the "else" and have only one "end", like this:
</p>

<pre>
{{if eq .A 1}} X {{else if eq .A 2}} Y {{end}}
</pre>

<p>
The two forms are identical in effect; the difference is just in the syntax.
</p>

<p>
<em>Updating</em>: Neither the "else if" change nor the comparison functions
affect existing programs. Those that
already define functions called <code>eq</code> and so on through a function
map are unaffected because the associated function map will override the new
default function definitions.
</p>

<h3 id="new_packages">New packages</h3>

<p>
There are two new packages.
</p>

<ul>
<li>
The <a href="/pkg/encoding/"><code>encoding</code></a> package is
<a href="#encoding">described above</a>.
</li>
<li>
The <a href="/pkg/image/color/palette/"><code>image/color/palette</code></a> package
provides standard color palettes.
</li>
</ul>

<h3 id="minor_library_changes">Minor changes to the library</h3>

<p>
The following list summarizes a number of minor changes to the library, mostly additions.
See the relevant package documentation for more information about each change.
</p>

<ul>

<li>
The <a href="/pkg/archive/zip/"><code>archive/zip</code></a> package
adds the
<a href="/pkg/archive/zip/#File.DataOffset"><code>DataOffset</code></a> accessor
to return the offset of a file's (possibly compressed) data within the archive.
</li>

<li>
The <a href="/pkg/bufio/"><code>bufio</code></a> package
adds <a href="/pkg/bufio/#Reader.Reset"><code>Reset</code></a>
methods to <a href="/pkg/bufio/#Reader"><code>Reader</code></a> and
<a href="/pkg/bufio/#Writer"><code>Writer</code></a>.
These methods allow the <a href="/pkg/Reader/"><code>Readers</code></a>
and <a href="/pkg/Writer/"><code>Writers</code></a>
to be re-used on new input and output readers and writers, saving
allocation overhead. 
</li>

<li>
The <a href="/pkg/compress/bzip2/"><code>compress/bzip2</code></a>
can now decompress concatenated archives.
</li>

<li>
The <a href="/pkg/compress/flate/"><code>compress/flate</code></a>
package adds a <a href="/pkg/compress/flate/#Reset"><code>Reset</code></a> 
method on the <a href="/pkg/compress/flate/#Writer"><code>Writer</code></a>,
to make it possible to reduce allocation when, for instance, constructing an
archive to hold multiple compressed files.
</li>

<li>
The <a href="/pkg/compress/gzip/"><code>compress/gzip</code></a> package's
<a href="/pkg/compress/gzip/#Writer"><code>Writer</code></a> type adds a
<a href="/pkg/compress/gzip/#Writer.Reset"><code>Reset</code></a>
so it may be reused.
</li>

<li>
The <a href="/pkg/compress/zlib/"><code>compress/zlib</code></a> package's
<a href="/pkg/compress/zlib/#Writer"><code>Writer</code></a> type adds a
<a href="/pkg/compress/zlib/#Writer.Reset"><code>Reset</code></a>
so it may be reused.
</li>

<li>
The <a href="/pkg/container/heap/"><code>container/heap</code></a> package
adds a <a href="/pkg/container/heap/#Fix"><code>Fix</code></a>
method to provide a more efficient way to update an item's position in the heap.
</li>

<li>
The <a href="/pkg/container/list/"><code>container/list</code></a> package
adds the <a href="/pkg/container/list/#MoveBefore"><code>MoveBefore</code></a>
and
<a href="/pkg/container/list/#MoveAfter"><code>MoveAfter</code></a>
methods, which implement the obvious rearrangement.
</li>

<li>
The <a href="/pkg/crypto/cipher/"><code>crypto/cipher</code></a> package
adds the a new GCM mode (Galois Counter Mode), which is almost always
used with AES encryption.
</li>

<li>
The 
<a href="/pkg/crypto/md5/"><code>crypto/md5</code></a> package
adds a new <a href="/pkg/crypto/md5/#Sum"><code>Sum</code></a> function
to simplify hashing without sacrificing performance.
</li>

<li>
Similarly, the 
<a href="/pkg/crypto/md5/"><code>crypto/sha1</code></a> package
adds a new <a href="/pkg/crypto/sha1/#Sum"><code>Sum</code></a> function.
</li>

<li>
Also, the
<a href="/pkg/crypto/sha256/"><code>crypto/sha256</code></a> package
adds <a href="/pkg/crypto/sha256/#Sum256"><code>Sum256</code></a>
and <a href="/pkg/crypto/sha256/#Sum224"><code>Sum224</code></a> functions.
</li>

<li>
Finally, the <a href="/pkg/crypto/sha512/"><code>crypto/sha512</code></a> package
adds <a href="/pkg/crypto/sha512/#Sum512"><code>Sum512</code></a> and
<a href="/pkg/crypto/sha512/#Sum384"><code>Sum384</code></a> functions.
</li>

<li>
The <a href="/pkg/crypto/x509/"><code>crypto/x509</code></a> package
adds support for reading and writing arbitrary extensions.
</li>

<li>
The <a href="/pkg/crypto/tls/"><code>crypto/tls</code></a> package adds
support for TLS 1.1, 1.2 and AES-GCM.
</li>

<li>
The <a href="/pkg/database/sql/"><code>database/sql</code></a> package adds a
<a href="/pkg/database/sql/#DB.SetMaxOpenConns"><code>SetMaxOpenConns</code></a>
method on <a href="/pkg/database/sql/#DB"><code>DB</code></a> to limit the
number of open connections to the database.
</li>

<li>
The <a href="/pkg/encoding/csv/"><code>encoding/csv</code></a> package
now always allows trailing commas on fields.
</li>

<li>
The <a href="/pkg/encoding/gob/"><code>encoding/gob</code></a> package
now treats channel and function fields of structures as if they were unexported,
even if they are not. That is, it ignores them completely. Previously they would
trigger an error, which could cause unexpected compatibility problems if an
embedded structure added such a field.
The package also now supports the generic encoding interfaces of the
<a href="/pkg/encoding/"><code>encoding</code></a> package
described above.
</li>

<li>
The <a href="/pkg/encoding/json/"><code>encoding/json</code></a> package
now will always escape ampersands as "\u0026" when printing strings.
It will now accept but correct invalid UTF-8 in
<a href="/pkg/encoding/json/#Marshal"><code>Marshal</code></a>
(such input was previously rejected).
Finally, it now supports the generic encoding interfaces of the
<a href="/pkg/encoding/"><code>encoding</code></a> package
described above.
</li>

<li>
The <a href="/pkg/encoding/xml/"><code>encoding/xml</code></a> package
now allows attributes stored in pointers to be marshaled.
It also supports the generic encoding interfaces of the
<a href="/pkg/encoding/"><code>encoding</code></a> package
described above through the new
<a href="/pkg/encoding/xml/#Marshaler"><code>Marshaler</code></a>,
<a href="/pkg/encoding/xml/#Unmarshaler"><code>Unmarshaler</code></a>,
and related
<a href="/pkg/encoding/xml/#MarshalerAttr"><code>MarshalerAttr</code></a> and
<a href="/pkg/encoding/xml/#UnmarshalerAttr"><code>UnmarshalerAttr</code></a>
interfaces.
The package also adds a
<a href="/pkg/encoding/xml/#Encoder.Flush"><code>Flush</code></a> method
to the
<a href="/pkg/encoding/xml/#Encoder"><code>Encoder</code></a>
type for use by custom encoders. See the documentation for
<a href="/pkg/encoding/xml/#Encoder.EncodeToken"><code>EncodeToken</code></a>
to see how to use it.
</li>

<li>
The <a href="/pkg/flag/"><code>flag</code></a> package now
has a <a href="/pkg/flag/#Getter"><code>Getter</code></a> interface
to allow the value of a flag to be retrieved. Due to the
Go 1 compatibility guidelines, this method cannot be added to the existing
<a href="/pkg/flag/#Value"><code>Value</code></a>
interface, but all the existing standard flag types implement it.
The package also now exports the <a href="/pkg/flag/#CommandLine"><code>CommandLine</code></a>
flag set, which holds the flags from the command line.
</li>

<li>
The <a href="/pkg/go/build/"><code>go/build</code></a> package adds
the <a href="/pkg/go/build/#Package.AllTags"><code>AllTags</code></a> field
to the <a href="/pkg/go/build/#Package"><code>Package</code></a> type,
to make it easier to process build tags.
</li>

<li>
The <a href="/pkg/image/draw/"><code>image/draw</code></a> package now
exports an interface, <a href="/pkg/image/draw/#Drawer"><code>Drawer</code></a>,
that wraps the standard <a href="/pkg/image/draw/#Draw"><code>Draw</code></a> method.
The Porter-Duff operators now implement this interface, in effect binding an operation to
the draw operator rather than providing it explicitly.
Given a paletted image as its destination, the new
<a href="/pkg/image/draw/#FloydSteinberg"><code>FloydSteinberg</code></a>
implementation of the
<a href="/pkg/image/draw/#Drawer"><code>Drawer</code></a>
interface will use the Floyd-Steinberg error diffusion algorithm to draw the image.
To create palettes suitable for such processing, the new
<a href="/pkg/image/draw/#Quantizer"><code>Quantizer</code></a> interface
represents implementations of quantization algorithms that choose a palette
given a full-color image.
There are no implementations of this interface in the library.
</li>

<li>
The <a href="/pkg/image/gif/"><code>image/gif</code></a> package
can now create GIF files using the new
<a href="/pkg/image/gif/#Encode"><code>Encode</code></a>
and <a href="/pkg/image/gif/#EncodeAll"><code>EncodeAll</code></a>
functions.
Their options argument allows specification of an image
<a href="/pkg/image/draw/#Quantizer"><code>Quantizer</code></a> to use;
if it is <code>nil</code>, the generated GIF will use the 
<a href="/pkg/image/color/palette/#Plan9"><code>Plan9</code></a>
color map (palette) defined in the new
<a href="/pkg/image/color/palette/"><code>image/color/palette</code></a> package.
The options also specify a
<a href="/pkg/image/draw/#Drawer"><code>Drawer</code></a>
to use to create the output image;
if it is <code>nil</code>, Floyd-Steinberg error diffusion is used.
</li>

<li>
The <a href="/pkg/io/#Copy"><code>Copy</code></a> method of the
<a href="/pkg/io/"><code>io</code></a> package now prioritizes its
arguments differently.
If one argument implements <a href="/pkg/io/#WriterTo"><code>WriterTo</code></a>
and the other implements <a href="/pkg/io/#ReaderFrom"><code>ReaderFrom</code></a>,
<a href="/pkg/io/#Copy"><code>Copy</code></a> will now invoke
<a href="/pkg/io/#WriterTo"><code>WriterTo</code></a> to do the work,
so that less intermediate buffering is required in general.
</li>

<li>
The <a href="/pkg/net/"><code>net</code></a> package requires cgo by default
because the host operating system must in general mediate network call setup.
On some systems, though, it is possible to use the network without cgo, and useful
to do so, for instance to avoid dynamic linking.
The new build tag <code>netgo</code> (off by default) allows the construction of a
<code>net</code> package in pure Go on those systems where it is possible.
</li>

<li>
The <a href="/pkg/net/"><code>net</code></a> package adds a new field
<code>DualStack</code> to the <a href="/pkg/net/#Dialer"><code>Dialer</code></a>
struct for TCP connection setup using a dual IP stack as described in
<a href="http://tools.ietf.org/html/rfc6555">RFC 6555</a>.
</li>

<li>
The <a href="/pkg/net/http/"><code>net/http</code></a> package will no longer
transmit cookies that are incorrect according to
<a href="http://tools.ietf.org/html/rfc6265">RFC 6265</a>.
It just logs an error and sends nothing.
Also,
the <a href="/pkg/net/http/"><code>net/http</code></a> package's
<a href="/pkg/net/http/#ReadResponse"><code>ReadResponse</code></a>
function now permits the <code>*Request</code> parameter to be <code>nil</code>,
whereupon it assumes a GET request.
Finally, an HTTP server will now serve HEAD
requests transparently, without the need for special casing in handler code.
While serving a HEAD request, writes to a 
<a href="/pkg/net/http/#Handler"><code>Handler</code></a>'s
<a href="/pkg/net/http/#ResponseWriter"><code>ResponseWriter</code></a>
are absorbed by the
<a href="/pkg/net/http/#Server"><code>Server</code></a>
and the client receives an empty body as required by the HTTP specification.
</li>

<li>
The <a href="/pkg/runtime/"><code>runtime</code></a> package relaxes
the constraints on finalizer functions in
<a href="/pkg/runtime/#SetFinalizer"><code>SetFinalizer</code></a>: the
actual argument can now be any type that is assignable to the formal type of
the function, as is the case for any normal function call in Go.
</li>

<li>
The <a href="/pkg/sort/"><code>sort</code></a> package has a new
<a href="/pkg/sort/#Stable"><code>Stable</code></a> function that implements
stable sorting. It is less efficient than the normal sort algorithm, however.
</li>

<li>
The <a href="/pkg/strings/"><code>strings</code></a> package adds
an <a href="/pkg/strings/#IndexByte"><code>IndexByte</code></a>
function for consistency with the <a href="/pkg/bytes/"><code>bytes</code></a> package.
</li>

<li>
The <a href="/pkg/sync/atomic/"><code>sync/atomic</code></a> package
adds a new set of swap functions that atomically exchange the argument with the
value stored in the pointer, returning the old value.
The functions are
<a href="/pkg/sync/atomic/#SwapInt32"><code>SwapInt32</code></a>,
<a href="/pkg/sync/atomic/#SwapInt64"><code>SwapInt64</code></a>,
<a href="/pkg/sync/atomic/#SwapUint32"><code>SwapUint32</code></a>,
<a href="/pkg/sync/atomic/#SwapUint64"><code>SwapUint64</code></a>,
<a href="/pkg/sync/atomic/#SwapUintptr"><code>SwapUintptr</code></a>,
and
<a href="/pkg/sync/atomic/#SwapPointer"><code>SwapPointer</code></a>,
which swaps an <code>unsafe.Pointer</code>.
</li>

<li>
syscall: implemented Sendfile for Darwin, added Syscall9 for Darwin/amd64 (CL 10980043).
</li>

<li>
The <a href="/pkg/testing/"><code>testing</code></a> package
now exports the <a href="/pkg/testing/#TB"><code>TB</code></a> interface.
It records the methods in common with the
<a href="/pkg/testing/#T"><code>T</code></a>
and
<a href="/pkg/testing/#B"><code>B</code></a> types,
to make it easier to share code between tests and benchmarks.
Also, the
<a href="/pkg/testing/#AllocsPerRun"><code>AllocsPerRun</code></a>
function now quantizes the return value to an integer (although it
still has type <code>float64</code>), to round off any error caused by
initialization and make the result more repeatable. 
</li>

<li>
The <a href="/pkg/text/template/"><code>text/template</code></a> package
now automatically dereferences pointer values when evaluating the arguments
to "escape" functions such as "html", to bring the behavior of such functions
in agreement with that of other printing functions such as "printf".
</li>

<li>
In the <a href="/pkg/time/"><code>time</code></a> package, the
<a href="/pkg/time/#Parse"><code>Parse</code></a> function
and
<a href="/pkg/time/#Format"><code>Format</code></a>
method
now handle time zone offsets with seconds, such as in the historical
date "1871-01-01T05:33:02+00:34:08".
Also, pattern matching in the formats for those routines is stricter: a non-lowercase letter
must now follow the standard words such as "Jan" and "Mon".
</li>

<li>
The <a href="/pkg/unicode/"><code>unicode</code></a> package
adds <a href="/pkg/unicode/#In"><code>In</code></a>,
a nicer-to-use but equivalent version of the original
<a href="/pkg/unicode/#IsOneOf"><code>IsOneOf</code></a>,
to see whether a character is a member of a Unicode category.
</li>

</ul>