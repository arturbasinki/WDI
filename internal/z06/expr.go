// Package z06 rysuje wykresy funkcji jednej zmiennej.
//
// Wzór funkcji i przedział dla x są wczytywane z interfejsu WWW.
// Wzór jest parsowany własnym parserem wstępu rekurencyjnego
// (analog „sympy” z wersji Python) i numerycznie próbkowany
// w zadanym przedziale.
//
// Obsługiwane operatory i funkcje:
//
//	+  -  *  /  ^  (potęga, prawostronna)  ( )
//	funkcje: sin cos tan asin acos atan sinh cosh tanh
//	          sqrt cbrt ln log log10 exp abs floor ceil round
//	stałe:  x  pi  e
//	np.  y = x*x - 6*x + 3,  y = 1/x,  y = tan(x),  y = sin(x)/x
//
// Implikowane mnożenie jest dozwolone: „2x” oznacza „2*x”,
// „x sin(x)” — „x*sin(x)” (ale nie „sin(x)” — nazwy funkcji są
// rozpoznawane jako wywołania).
//
// Asymptoty pionowe wykrywane są numerycznie: punkt próbki jest
// „wyrzutem” gdy |y| przekracza próg zależny od rozkładu wartości
// (mediana + 15·IQR), a zmiana znaku między sąsiednimi próbkami
// z dużym skokiem oznacza asymptotę pionową. Wykres zwracany jest
// jako lista odcinków (segmentów), których nie przecinają asymptoty,
// więc interfejs rysuje je jako osobne ścieżki i przerywaną linię
// w miejscu asymptoty.
package z06

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// Point — punkt na wykresie.
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Plot — wynik próbkowania funkcji.
type Plot struct {
	Formula    string    `json:"formula"`
	XMin       float64   `json:"xmin"`
	XMax       float64   `json:"xmax"`
	Samples    int       `json:"samples"`
	Segments   [][]Point `json:"segments"`
	Asymptotes []float64 `json:"asymptotes"`
	YClamp     float64   `json:"yclamp"`
}

// ---------------------------------------------------------------------------
// Tokenizer
// ---------------------------------------------------------------------------

type tokenKind int

const (
	tNum tokenKind = iota
	tIdent
	tOp
)

type token struct {
	kind tokenKind
	text string
	num  float64
}

var tokenAliases = map[string]string{
	"*": "mul", "/": "div", "+": "add", "-": "sub", "^": "pow",
	"(": "lp", ")": "rp", ",": "comma",
}

func tokenize(s string) ([]token, error) {
	var toks []token
	i := 0
	for i < len(s) {
		c := s[i]
		switch {
		case c == ' ' || c == '\t' || c == '\n':
			i++
		case c >= '0' && c <= '9' || c == '.':
			j := i
			dots := 0
			for j < len(s) && ((s[j] >= '0' && s[j] <= '9') || s[j] == '.') {
				if s[j] == '.' {
					dots++
				}
				j++
			}
			if dots > 1 {
				return nil, fmt.Errorf("nieprawidłowy zapis liczby: %q", s[i:j])
			}
			v, err := strconv.ParseFloat(s[i:j], 64)
			if err != nil {
				return nil, fmt.Errorf("nieprawidłowa liczba: %q", s[i:j])
			}
			toks = append(toks, token{kind: tNum, text: s[i:j], num: v})
			i = j
		case isIdentStart(c):
			j := i
			for j < len(s) && isIdentPart(s[j]) {
				j++
			}
			toks = append(toks, token{kind: tIdent, text: s[i:j]})
			i = j
		case c == '*' || c == '/' || c == '+' || c == '-' || c == '^' || c == '(' || c == ')' || c == ',':
			toks = append(toks, token{kind: tOp, text: tokenAliases[string(c)]})
			i++
		default:
			return nil, fmt.Errorf("nieznany znak %q", string(c))
		}
	}
	return toks, nil
}

func isIdentStart(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func isIdentPart(c byte) bool {
	return isIdentStart(c) || (c >= '0' && c <= '9')
}

// ---------------------------------------------------------------------------
// Parser (wstęp rekurencyjny)
//
// grammar:
//
//	expr   := term (('+'|'-') term)*
//	term   := factor (('*'|'/') factor | ident | '(')?   // implikowane mnożenie
//	factor := ('-'|'+')* power
//	power  := atom ('^' factor)?
//	atom   := NUM | IDENT | FUNC '(' expr (',' expr)* ')' | '(' expr ')'
// ---------------------------------------------------------------------------

// node — węzeł AST.
type node interface {
	eval(x float64) float64
}

type numNode float64

func (n numNode) eval(x float64) float64 { return float64(n) }

type varNode struct{}

func (varNode) eval(x float64) float64 { return x }

type binNode struct {
	op   string
	l, r node
}

func (n binNode) eval(x float64) float64 {
	a, b := n.l.eval(x), n.r.eval(x)
	switch n.op {
	case "add":
		return a + b
	case "sub":
		return a - b
	case "mul":
		return a * b
	case "div":
		return a / b
	case "pow":
		return math.Pow(a, b)
	}
	return math.NaN()
}

type unaryNode struct {
	op   string
	inner node
}

func (n unaryNode) eval(x float64) float64 {
	v := n.inner.eval(x)
	if n.op == "neg" {
		return -v
	}
	return v
}

type callNode struct {
	name string
	args []node
}

func (n callNode) eval(x float64) float64 {
	f, ok := funcs[n.name]
	if !ok {
		return math.NaN()
	}
	args := make([]float64, len(n.args))
	for i, a := range n.args {
		args[i] = a.eval(x)
	}
	v, err := f(args)
	if err != nil {
		return math.NaN()
	}
	return v
}

// isFinite — czy wartość jest skończona (ani NaN, ani Inf).
func isFinite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

var funcs = map[string]func([]float64) (float64, error){
	"sin":    func(a []float64) (float64, error) { return math.Sin(a[0]), nil },
	"cos":    func(a []float64) (float64, error) { return math.Cos(a[0]), nil },
	"tan":    func(a []float64) (float64, error) { return math.Tan(a[0]), nil },
	"asin":   func(a []float64) (float64, error) { return math.Asin(a[0]), nil },
	"acos":   func(a []float64) (float64, error) { return math.Acos(a[0]), nil },
	"atan":   func(a []float64) (float64, error) { return math.Atan(a[0]), nil },
	"sinh":   func(a []float64) (float64, error) { return math.Sinh(a[0]), nil },
	"cosh":   func(a []float64) (float64, error) { return math.Cosh(a[0]), nil },
	"tanh":   func(a []float64) (float64, error) { return math.Tanh(a[0]), nil },
	"sqrt":   func(a []float64) (float64, error) { return math.Sqrt(a[0]), nil },
	"cbrt":   func(a []float64) (float64, error) { return math.Cbrt(a[0]), nil },
	"ln":     func(a []float64) (float64, error) { return math.Log(a[0]), nil },
	"log":    func(a []float64) (float64, error) { return math.Log(a[0]), nil },
	"log10":  func(a []float64) (float64, error) { return math.Log10(a[0]), nil },
	"exp":    func(a []float64) (float64, error) { return math.Exp(a[0]), nil },
	"abs":    func(a []float64) (float64, error) { return math.Abs(a[0]), nil },
	"floor":  func(a []float64) (float64, error) { return math.Floor(a[0]), nil },
	"ceil":   func(a []float64) (float64, error) { return math.Ceil(a[0]), nil },
	"round":  func(a []float64) (float64, error) { return math.Round(a[0]), nil },
	"min":    func(a []float64) (float64, error) { m := a[0]; for _, v := range a[1:] { if v < m { m = v } }; return m, nil },
	"max":    func(a []float64) (float64, error) { m := a[0]; for _, v := range a[1:] { if v > m { m = v } }; return m, nil },
	"pow":    func(a []float64) (float64, error) { return math.Pow(a[0], a[1]), nil },
	"fmod":   func(a []float64) (float64, error) { return math.Mod(a[0], a[1]), nil },
	"sigmoid": func(a []float64) (float64, error) { return 1 / (1 + math.Exp(-a[0])), nil },
}

func unaryArity(name string) int {
	switch name {
	case "min", "max":
		return -1 // dowolna >= 2
	case "pow", "fmod":
		return 2
	}
	return 1
}

type parser struct {
	toks []token
	pos  int
}

func (p *parser) peek() *token {
	if p.pos < len(p.toks) {
		return &p.toks[p.pos]
	}
	return nil
}

func (p *parser) next() *token {
	t := p.peek()
	if t != nil {
		p.pos++
	}
	return t
}

func (p *parser) expect(text string) error {
	t := p.next()
	if t == nil || t.kind != tOp || t.text != text {
		return fmt.Errorf("oczekiwałem %q, a otrzymałem %v", text, t)
	}
	return nil
}

// parseparsuje całe wyrażenie.
func parseExpr(s string) (node, error) {
	toks, err := tokenize(s)
	if err != nil {
		return nil, err
	}
	p := &parser{toks: toks}
	n, err := p.expr()
	if err != nil {
		return nil, err
	}
	if p.pos != len(p.toks) {
		return nil, fmt.Errorf("nieoczekiwany token %q", p.toks[p.pos].text)
	}
	return n, nil
}

func (p *parser) expr() (node, error) {
	l, err := p.term()
	if err != nil {
		return nil, err
	}
	for {
		t := p.peek()
		if t == nil || t.kind != tOp || (t.text != "add" && t.text != "sub") {
			return l, nil
		}
		p.next()
		r, err := p.term()
		if err != nil {
			return nil, err
		}
		l = binNode{op: t.text, l: l, r: r}
	}
}

func (p *parser) term() (node, error) {
	l, err := p.factor()
	if err != nil {
		return nil, err
	}
	for {
		t := p.peek()
		if t == nil {
			return l, nil
		}
		if t.kind == tOp && (t.text == "mul" || t.text == "div") {
			p.next()
			r, err := p.factor()
			if err != nil {
				return nil, err
			}
			l = binNode{op: t.text, l: l, r: r}
			continue
		}
		// Implikowane mnożenie: liczba/litera/nawias bez operatora.
		if t.kind == tNum || t.kind == tIdent || (t.kind == tOp && t.text == "lp") {
			// ale tylko jeśli nie jest to nawias należący do funkcji —
			// wywołanie funkcji parsowane jest w atom(), więc tutaj
			// „(” zawsze oznacza nową liczbę mnożoną implikacyjnie.
			r, err := p.factor()
			if err != nil {
				return nil, err
			}
			l = binNode{op: "mul", l: l, r: r}
			continue
		}
		return l, nil
	}
}

func (p *parser) factor() (node, error) {
	neg := false
	for {
		t := p.peek()
		if t == nil || t.kind != tOp {
			break
		}
		if t.text == "sub" {
			neg = !neg
			p.next()
			continue
		}
		if t.text == "add" {
			p.next()
			continue
		}
		break
	}
	n, err := p.power()
	if err != nil {
		return nil, err
	}
	if neg {
		n = unaryNode{op: "neg", inner: n}
	}
	return n, nil
}

func (p *parser) power() (node, error) {
	base, err := p.atom()
	if err != nil {
		return nil, err
	}
	t := p.peek()
	if t != nil && t.kind == tOp && t.text == "pow" {
		p.next()
		// prawostronna asocjacja: x^2^3 = x^(2^3)
		exp, err := p.factor()
		if err != nil {
			return nil, err
		}
		return binNode{op: "pow", l: base, r: exp}, nil
	}
	return base, nil
}

func (p *parser) atom() (node, error) {
	t := p.next()
	if t == nil {
		return nil, fmt.Errorf("nieoczekiwany koniec wyrażenia")
	}
	switch t.kind {
	case tNum:
		return numNode(t.num), nil
	case tIdent:
		name := strings.ToLower(t.text)
		if f, ok := funcs[name]; ok {
			// wywołanie funkcji: nazwa '(' arg (',' arg)* ')'
			if err := p.expect("lp"); err != nil {
				return nil, fmt.Errorf("funkcja %s wymaga argumentów w nawiasach", name)
			}
			var args []node
			for {
				a, err := p.expr()
				if err != nil {
					return nil, err
				}
				args = append(args, a)
				nt := p.peek()
				if nt != nil && nt.kind == tOp && nt.text == "comma" {
					p.next()
					continue
				}
				break
			}
			if err := p.expect("rp"); err != nil {
				return nil, fmt.Errorf("funkcja %s: brak zamykającego nawiasu", name)
			}
			arity := unaryArity(name)
			if arity > 0 && len(args) != arity {
				return nil, fmt.Errorf("funkcja %s przyjmuje %d argumenty, podano %d", name, arity, len(args))
			}
			if arity == -1 && len(args) < 2 {
				return nil, fmt.Errorf("funkcja %s przyjmuje co najmniej 2 argumenty", name)
			}
			_ = f
			return callNode{name: name, args: args}, nil
		}
		switch name {
		case "x":
			return varNode{}, nil
		case "pi":
			return numNode(math.Pi), nil
		case "e":
			return numNode(math.E), nil
		default:
			return nil, fmt.Errorf("nieznana nazwa %q (dostępne: x, pi, e oraz funkcje %s)",
				t.text, funcNames())
		}
	case tOp:
		if t.text == "lp" {
			n, err := p.expr()
			if err != nil {
				return nil, err
			}
			if err := p.expect("rp"); err != nil {
				return nil, fmt.Errorf("brak zamykającego nawiasu")
			}
			return n, nil
		}
		return nil, fmt.Errorf("nieoczekiwany operator %q", t.text)
	}
	return nil, fmt.Errorf("nieoczekiwany token %q", t.text)
}

func funcNames() string {
	var names []string
	for n := range funcs {
		names = append(names, n)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

// ---------------------------------------------------------------------------
// Próbkowanie i wykrywanie asymptot
// ---------------------------------------------------------------------------

// PlotFun próbkuje funkcję wyrażoną wzorem w przedziale [xmin, xmax].
func PlotFun(formula string, xmin, xmax float64, samples int) (*Plot, error) {
	if samples < 501 {
		samples = 501
	}
	if samples > 20001 {
		samples = 20001
	}
	if xmin >= xmax {
		return nil, fmt.Errorf("dolna granica przedziału musi być mniejsza od górnej")
	}
	ast, err := parseExpr(formula)
	if err != nil {
		return nil, err
	}

	xs := make([]float64, samples)
	ys := make([]float64, samples)
	for i := 0; i < samples; i++ {
		x := xmin + float64(i)*(xmax-xmin)/float64(samples-1)
		xs[i] = x
		ys[i] = ast.eval(x)
	}

	// Percentyle |y| dla próbek skończonych.
	finite := []float64{}
	for i := 0; i < samples; i++ {
		if isFinite(ys[i]) {
			finite = append(finite, math.Abs(ys[i]))
		}
	}
	if len(finite) < 10 {
		return &Plot{Formula: formula, XMin: xmin, XMax: xmax, Samples: samples,
			Segments: [][]Point{}, Asymptotes: []float64{}, YClamp: 10}, nil
	}
	sort.Float64s(finite)
	quant := func(q float64) float64 {
		idx := q * float64(len(finite)-1)
		lo := int(idx)
		hi := lo + 1
		if hi >= len(finite) {
			return finite[lo]
		}
		frac := idx - float64(lo)
		return finite[lo]*(1-frac) + finite[hi]*frac
	}
	p25, p50, p75 := quant(0.25), quant(0.5), quant(0.75)
	iqr := p75 - p25
	if iqr < 1e-12 {
		iqr = math.Max(1, p50)
	}
	jumpThresh := 4 * iqr
	clampThresh := p50 + 15*iqr

	// Segmentacja: zerwij odcinek gdy
	//  - punkt nieskończony/NaN,
	//  - zmiana znaku z dużym skokiem (asymptota pionowa),
	//  - wejście/wyjście z obszaru „wyrzutów” (|y| > clampThresh).
	clamped := make([]bool, samples)
	for i := 0; i < samples; i++ {
		clamped[i] = isFinite(ys[i]) && math.Abs(ys[i]) > clampThresh
	}
	breakAt := make([]bool, samples) // breakAt[i] — przerwa między i-1 a i
	asymptotes := []float64{}
	for i := 1; i < samples; i++ {
		yi, yj := ys[i-1], ys[i]
		fi, fj := isFinite(yi), isFinite(yj)
		if !fi || !fj {
			// punkt nieskończony/NaN — przerwa; gdy po drugiej stronie
			// jest punkt skończony, w okolicy leży asymptota pionowa
			breakAt[i] = true
			if fi != fj {
				asymptotes = append(asymptotes, (xs[i-1]+xs[i])/2)
			}
			continue
		}
		if yi*yj < 0 && math.Abs(yj-yi) > jumpThresh {
			breakAt[i] = true
			asymptotes = append(asymptotes, (xs[i-1]+xs[i])/2)
			continue
		}
		if clamped[i-1] != clamped[i] {
			breakAt[i] = true
		}
	}
	// Deduplikacja asymptot bliższych niż 2 kroki.
	sort.Float64s(asymptotes)
	dedup := []float64{}
	step := (xmax - xmin) / float64(samples-1)
	for _, a := range asymptotes {
		if len(dedup) == 0 || a-dedup[len(dedup)-1] > 2*step {
			dedup = append(dedup, a)
		}
	}
	asymptotes = dedup

	segments := [][]Point{}
	cur := []Point{}
	flush := func() {
		if len(cur) > 0 {
			segments = append(segments, cur)
			cur = []Point{}
		}
	}
	for i := 0; i < samples; i++ {
		if isFinite(ys[i]) && !clamped[i] {
			cur = append(cur, Point{X: xs[i], Y: ys[i]})
			if i == samples-1 || (i < samples-1 && breakAt[i+1]) {
				flush()
			}
			continue
		}
		// punkt wyrzucony lub nieskończony: kończy segment
		flush()
	}

	return &Plot{
		Formula:    formula,
		XMin:       xmin,
		XMax:       xmax,
		Samples:    samples,
		Segments:   segments,
		Asymptotes: asymptotes,
		YClamp:     clampThresh,
	}, nil
}
