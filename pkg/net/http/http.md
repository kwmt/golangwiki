<h2 id="Overview">Overview</h2>
<p>
httpパッケージはHTTPクライアントとサーバを実装して提供します。
</p>
<p>
Get,Head,Post,PostFormはHTTPリクエストを生成します。
</p>
<pre class="go">
resp, err := http.Get("http://example.com/")
...
resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
...
resp, err := http.PostForm("http://example.com/form",
        url.Values{"key":{"Value"}, "id":{"123"}})
</pre>
<p>
クライアントは、レスポンスボディを使い終わったら、必ずクローズしなければなりません。
</p>
<pre class="go">
resp, err := http.Get("http://example.com/")
if err != nil {
    //handle error
}
defer resp.Body.Close()
body, err := ioutil.ReadAll(resp.Body)
//...
</pre>


<p>
ListenAndServeは与えられたアドレスとハンドラを使って、HTTPサーバーを開始します。
そのハンドラは普通はnilです。
nilの場合は、DefaultServeMuxを使うという意味です。
HandleとHandleFuncはDefaultServeMuxにハンドラを追加します。
</p>
<pre class="go">
http.Handle("/foo", fooHandler)
http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, &q", html.EscapeString(r.URL.Path))
})
</pre>


<h2 id="Redirect">func <a href="http://golang.org/src/pkg/net/http/server.go?s=21797:21865#L730">Redirect</a></h2>
<pre>func Redirect(w ResponseWriter, r *Request, urlStr string, code int)</pre>
<p>
Redirect replies to the request with a redirect to url,
which may be a path relative to the request path.


</p>
<h2 id="FileServer">func FileServer</h2>
<pre>func FileServer(root FileSystem) Handler</pre>
<p>
FileServerはrootで指定したファイルシステムをルートとして、その内容のHTTPリクエスト出すハンドラを返します。
</p>
<p>
オペレーティングシステムのファイルシステム実装を使うためには、http.Dirを使用します。
</p>
<pre>
http.Handle("/", http.FileServer(http.Dir("/tmp")))
</pre>
<h2 id="StripPrefix">func StripPrefix</h2>
<pre>func StripPrefix(prefix string, h Handler) Handler</pre>
<p>
StripPrefixはリクエストURLのパスから、与えられたprefixを取り除き、ハンドラhを呼ぶことでHTTPリクエストを出すハンドラを返します。
StripPrefixは、パスに対するリクエストをハンドルしますが、prefixで始まらない応答は404 not found エラーとなります。
</p>
<h2 id="EscapeString">

<h3 id="NewRequest">func <a href="/src/pkg/net/http/request.go?s=12527:12599#L398">NewRequest</a></h3>
<pre>func NewRequest(method, urlStr string, body io.Reader) (*Request, error)</pre>
<p>
NewRequest関数は、メソッド、URL、そしてオプションとしてbodyを与えることで、新しいRequestを返す関数です。
</p>

