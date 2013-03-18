<h2 id="Exit">func Exit</h2>
<pre class="go">
func Exit(code int)
</pre>
<p>
Exitは、与えられたステータスコードで、現在のプログラムを終了させます。
慣習的にコード0は成功を表し、0でないコードはエラーを表します。
</p>

<h2 id="Getenv">func <a href="http://golang.org/src/pkg/os/env.go?s=2363:2393#L69">Getenv</a></h2>
<pre>func Getenv(key string) string</pre>
<p>
Getenvは指定したキーの環境変数の値を取得します。
変数が存在しない場合、空の値を返します。
</p
<h2 id="TempDir">func <a href="http://golang.org/src/pkg/os/file_unix.go?s=8657:8678#L287">TempDir</a></h2>
<pre>func TempDir() string</pre>
<p>
TempDirは、一時ファイル用に使用するデフォルトのディレクトリを返します。
</p>

<h2 id="OpenFile">func OpenFile</h2>
<pre class="go">
func OpenFile(name string, flag int, perm FileMode) (file *File, err error)
</pre>
<p>
OpenFileは、一般化されたファイルを開く関数です。
ほとんどのユーザーはOpenやCloseの代わりに使うでしょう。
指定されたフラグ(O_RDONLYなど)とパーミッション(0666など)を使って、
指定したファイル名のファイルを開きます。
もし成功したら、I/Oに対してFileに関するメソッド郡が使えます。
エラーであれば、型*PathErrorになります。
</p>
