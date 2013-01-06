<h2 id="Overview">Overview</h2>
<p>
htmlパッケージは、HTMLテキストをエスケープすることとアンエスケープする関数を提供します。
</p>

<h2 id="EscapeString">
func <a href="/src/pkg/html/escape.go?s=5523:5557#L229">EscapeString</a></h2>
<pre class="go">
func EscapeString(s string) string
</pre>
<p>
EscapeString escapes special characters like "<" to become "&lt;".
 It escapes only five such characters: <, >, &, ' and ".
 UnescapeString(EscapeString(s)) == s always holds, but the converse isn't always true.

EscapeStringは¥"<¥"が¥"&lt;¥"になるような特殊文字をエスケープします。
</p>


<h2 id="UnescapeString">
func <a href="/src/pkg/html/escape.go?s=5986:6022#L243">UnescapeString</a></h2>
<pre class="go">
func UnescapeString(s string) string
</pre>
<p>
</p>
