<h3 id="Method_values">Method values</h3>

<p>
If the expression <code>x</code> has static type <code>T</code> and
<code>M</code> is in the <a href="#Method_sets">method set</a> of type <code>T</code>,
<code>x.M</code> is called a <i>method value</i>.
The method value <code>x.M</code> is a function value that is callable
with the same arguments as a method call of <code>x.M</code>.
The expression <code>x</code> is evaluated and saved during the evaluation of the
method value; the saved copy is then used as the receiver in any calls,
which may be executed later.
</p>

<p>
式<code>x</code>がスタティックな型<code>T</code>を持っていて、
<code>M</code>が型<code>T</code>の<a href="#Method_sets">method set</a>である場合、
<code>x.M</code>を<i>method value</i>と呼びます。
このmethod valueである<code>x.M</code> は、
<code>x.M</code>のメソッド呼び出しとして同じ引数で呼び出しできる関数の値です。
式<code>x</code>は評価され、method valueの評価している間保存されます。
保存されたコピーは、あとで実行されるかもしれないコールでのレシーバとして使われます。
</p>

<p>
The type <code>T</code> may be an interface or non-interface type.
</p>
<p>
型<code>T</code>は、interface型か非interface型かもしれません。
</p>

<p>
As in the discussion of <a href="#Method_expressions">method expressions</a> above,
consider a struct type <code>T</code> with two methods,
<code>Mv</code>, whose receiver is of type <code>T</code>, and
<code>Mp</code>, whose receiver is of type <code>*T</code>.
</p>
<p>
上述の<a href="#Method_expressions">method expressions</a> の議論において、
２つのメソッド
：型<code>T</code>のレシーバを持つ
<code>Mv</code>と
型<code>*T</code>のレシーバをもつ
<code>Mp</code>
、構造体を考えます。
</p>
<pre>
type T struct {
	a int
}
func (tv  T) Mv(a int) int         { return 0 }  // 値レシーバ
func (tp *T) Mp(f float32) float32 { return 1 }  // ポインタレシーバ

var t T
var pt *T
func makeT() T
</pre>

<p>
式
</p>

<pre>
t.Mv
</pre>

<p>
は、型
</p>

<pre>
func(int) int
</pre>
<p>
の関数値を生成します。
</p>

<p>
These two invocations are equivalent:
</p>
<p>
これらの２つの呼び出しは等しいです。
</p>

<pre>
t.Mv(7)
f := t.Mv; f(7)
</pre>

<p>
Similarly, the expression
</p>
<p>
同じように、式
</p>

<pre>
pt.Mp
</pre>

<p>
yields a function value of type
</p>
<p>
は、型
</p>	
<pre>
func(float32) float32
</pre>
<p>
の関数値を生成します。
</p>

<p>
As with <a href="#Selectors">selectors</a>, a reference to a non-interface method with a value receiver
using a pointer will automatically dereference that pointer: <code>pt.Mv</code> is equivalent to <code>(*pt).Mv</code>.
</p>
<p>
<a href="#Selectors">selectors</a>を持つように、
ポインタを使っている値レシーバをもっている非インターフェースなメソッドへの参照は
参照先の値を自動的に取得します。
<code>pt.Mv</code>は<code>(*pt).Mv</code>と同じです。

</p>
<p>
As with <a href="#Calls">method calls</a>, a reference to a non-interface method with a pointer receiver
using an addressable value will automatically take the address of that value: <code>t.Mv</code> is equivalent to <code>(&amp;t).Mv</code>.
</p>
<p>
<a href="#Calls">method calls</a>のように、アドレス指定可能な値を使って
ポインタレシーバを持つ非インターフェースなメソッドへの参照は、
値のアドレスを自動的に取得します。
<code>t.Mv</code>は、<code>(&amp;t).Mv</code>と同じです。
</p>

<pre>
f := t.Mv; f(7)   // like t.Mv(7)
f := pt.Mp; f(7)  // like pt.Mp(7)
f := pt.Mv; f(7)  // like (*pt).Mv(7)
f := t.Mp; f(7)   // like (&amp;t).Mp(7)
f := makeT().Mp   // 無効です：makeT()の結果はアドレス指定できません
</pre>

<p>
Although the examples above use non-interface types, it is also legal to create a method value
from a value of interface type.
</p>
<p>
上の例は、非インターフェース型を使っていますが、インターフェース型の値からmethod valueを作ることもできます。
</p>

<pre>
var i interface { M(int) } = myVal
f := i.M; f(7)  // like i.M(7)
</pre>
