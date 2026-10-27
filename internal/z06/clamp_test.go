package z06

import (
	"math"
	"strings"
	"testing"
)

// TestClampKeepsPoints — „clamping Y” ma REALNIE przycinać do ±YClamp i
// ZACHOWYWAĆ punkt. Dawniej punkty ponad progiem były usuwane, więc gałęzie
// urywali się tuż pod progiem (dla 1/x na [-5,5] ~|y|=41) zamiast dojeżdżać
// do ramki wykresu.
func TestClampKeepsPoints(t *testing.T) {
	p, err := PlotFun("1/x", -5, 5, 1601)
	if err != nil {
		t.Fatal(err)
	}
	if p.YClamp <= 0 {
		t.Fatalf("YClamp = %v", p.YClamp)
	}
	var maxAbs float64
	var n int
	for _, seg := range p.Segments {
		for _, pt := range seg {
			n++
			if math.Abs(pt.Y) > p.YClamp {
				t.Fatalf("punkt (%v, %v) wyszedł poza przycięcie ±%v", pt.X, pt.Y, p.YClamp)
			}
			if math.Abs(pt.Y) > maxAbs {
				maxAbs = math.Abs(pt.Y)
			}
		}
	}
	if n == 0 {
		t.Fatal("brak punktów w segmentach")
	}
	// punkt leżący przy asymptocie dociera DO progu, nie zatrzymuje się pod nim
	if maxAbs != p.YClamp {
		t.Errorf("najdalszy |y| = %v, oczekiwano dokładnie YClamp = %v (gałąź musi dochodzić do ramki)", maxAbs, p.YClamp)
	}
	// 1/x: dwie gałęzie i jedna asymptota pionowa w x ≈ 0
	if len(p.Segments) != 2 {
		t.Errorf("segmentów: %d, oczekiwano 2", len(p.Segments))
	}
	if len(p.Asymptotes) != 1 {
		t.Errorf("asymptot: %d, oczekiwano 1", len(p.Asymptotes))
	}
	// gałęzie KOŃCZĄ się na progu przy asymptocie (x→0- daje -YClamp,
	// x→0+ daje +YClamp) — to dowód, że punkty przy asymptocie są zachowane
	left := p.Segments[0]
	leftEnd := left[len(left)-1]
	if leftEnd.Y != -p.YClamp {
		t.Errorf("koniec lewej gałęzi (x→0⁻): y=%v, oczekiwano %v", leftEnd.Y, -p.YClamp)
	}
	right := p.Segments[1]
	rightStart := right[0]
	if rightStart.Y != p.YClamp {
		t.Errorf("start prawej gałęzi (x→0⁺): y=%v, oczekiwano %v", rightStart.Y, p.YClamp)
	}
}

// TestInfinityEndsSegment — punkty ±Inf/NaN NADAL kończą segment
// (tylko „zwykłe” wyrzuty są przycinane i zachowywane).
func TestInfinityEndsSegment(t *testing.T) {
	// ln(x) dla x <= 0 to NaN: gałąź tylko dla x > 0
	p, err := PlotFun("ln(x)", -8, 8, 801)
	if err != nil {
		t.Fatal(err)
	}
	var nSegWithData int
	for _, seg := range p.Segments {
		if len(seg) > 0 {
			nSegWithData++
		}
		for _, pt := range seg {
			if math.IsNaN(pt.Y) || math.IsInf(pt.Y, 0) {
				t.Fatalf("punkt NaN/Inf w segmencie: (%v, %v)", pt.X, pt.Y)
			}
		}
	}
	if nSegWithData != 1 {
		t.Errorf("segmentów z danymi: %d, oczekiwano 1 (x>0)", nSegWithData)
	}
	if len(p.Asymptotes) == 0 {
		t.Error("oczekiwano asymptoty na granicy dziedziny ln(x)")
	}
}

// TestExponentialNotation — liczby przyjmują [eE][+-]cyfry: 2e3 = 2000.
// „3e” bez cyfr po e to dalej 3·e (nie wykładnik).
func TestExponentialNotation(t *testing.T) {
	cases := []struct {
		expr string
		x    float64
		want float64
	}{
		{"2e3", 0, 2000},
		{"2E3", 0, 2000},
		{"2e+3", 0, 2000},
		{"2e-3", 0, 0.002},
		{"1.5e2", 0, 150},
		{"3e", 0, 3 * math.E}, // stała e, nie wykładnik
		{"2e3x", 2, 4000},     // 2000 · x (implikowane mnożenie)
		{"1e-3+x", 1, 1.001},  // wykładnik ujemny w działaniu
		{"x^2+1e2", 3, 109},   // 9 + 100
	}
	for _, c := range cases {
		n, err := parseExpr(c.expr)
		if err != nil {
			t.Errorf("%q: %v", c.expr, err)
			continue
		}
		got := n.eval(c.x)
		if math.Abs(got-c.want) > 1e-9*math.Max(1, math.Abs(c.want)) {
			t.Errorf("%q dla x=%v = %v, oczekiwano %v", c.expr, c.x, got, c.want)
		}
	}
}

// TestRuneInErrorMessage — komunikat o nieznakach musi pokazywać RUNĘ.
// Dawniej „π” (dwa bajty) dawało mojibake: nieznany znak "Ï".
func TestRuneInErrorMessage(t *testing.T) {
	_, err := PlotFun("x + π", 0, 1, 601)
	if err == nil {
		t.Fatal("oczekiwano błędu dla „π”")
	}
	msg := err.Error()
	if !strings.Contains(msg, "π") {
		t.Errorf("komunikat %q nie zawiera „π”", msg)
	}
	if strings.Contains(msg, "Ï") {
		t.Errorf("komunikat %q to mojibake (bajt zamiast runy)", msg)
	}
	// emoji (4 bajty) też musi wyjść cała
	_, err = PlotFun("x + 😀", 0, 1, 601)
	if err == nil || !strings.Contains(err.Error(), "😀") {
		t.Errorf("komunikat o emoji: %v", err)
	}
}

// TestSamplesRequested — granice przycinane jawnie: Samples to to, co użyto,
// SamplesRequested to to, o co proszono (dawniej różnica znikała bez śladu).
func TestSamplesRequested(t *testing.T) {
	lo, err := PlotFun("x", 0, 1, 100)
	if err != nil {
		t.Fatal(err)
	}
	if lo.Samples != 501 || lo.SamplesRequested != 100 {
		t.Errorf("za mało próbek: Samples=%d SamplesRequested=%d, oczekiwano 501/100", lo.Samples, lo.SamplesRequested)
	}
	hi, err := PlotFun("x", 0, 1, 5_000_000)
	if err != nil {
		t.Fatal(err)
	}
	if hi.Samples != 20001 || hi.SamplesRequested != 5_000_000 {
		t.Errorf("za dużo próbek: Samples=%d SamplesRequested=%d, oczekiwano 20001/5000000", hi.Samples, hi.SamplesRequested)
	}
	ok, err := PlotFun("x", 0, 1, 1601)
	if err != nil {
		t.Fatal(err)
	}
	if ok.Samples != 1601 || ok.SamplesRequested != 1601 {
		t.Errorf("bez przycinania: Samples=%d SamplesRequested=%d", ok.Samples, ok.SamplesRequested)
	}
}
