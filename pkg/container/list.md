<h2 id="Element">type <a href="/src/pkg/container/list/list.go?s=408:709#L5">Element</a></h2>
<pre>type Element struct {
    <span class="comment">// リストの要素の値です。</span>
    Value interface{}
    <span class="comment">// フィルタしたフィールドや公開されないフィールドも含んでいます。</span>
}</pre>
<p>
Element は連結リストのひとつの要素です。
</p>

<h2 id="List">type <a href="/src/pkg/container/list/list.go?s=1012:1071#L25">List</a></h2>
<pre>type List struct {
    <span class="comment">// フィルタしたフィールドや公開されないフィールドも含んでいます。</span>
}</pre>
<p>
Listは双方向リストを表します。
Listのゼロは、空リストです。
</p>

<h3 id="New">fun New</h3>
<pre class="go">
func New() *List
</pre>
<p>
Newは、初期化したリストを返します。
</p>

<h3 id="List.Back">func (*List) Back</h3>
<pre class="go">
func (l *List) Back() *Element
</pre>
<p>
Backはリストの最後の要素を返します。
</p>

<h3 id="List.Front">func (*List) Front</h3>
<pre class="go">
func (l *List) Front() *Element
</pre>
<p>
Frontはリストの最初の要素を返します。
</p>
<h3 id="List.Init">func (*List) Init</h3>
<pre class="go">
func (l *List) Init() *List
</pre>
<p>
InitはListを初期化あるいはクリアします。
</p>
