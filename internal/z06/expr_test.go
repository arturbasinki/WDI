package z06

import (
	"math"
	"testing"
)

func TestParseAndEval(t *testing.T) {
	ast, err := parseExpr("x*x - 6*x + 3")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	for x, want := range map[float64]float64{0: 3, 1: -2, 3: -6, 6: 3} {
		if got := ast.eval(x); math.Abs(got-want) > 1e-9 {
			t.Errorf("f(%v) = %v, oczekiwano %v", x, got, want)
		}
	}
}

func TestImplicitMul(t *testing.T) {
	ast, err := parseExpr("2x")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := ast.eval(3); math.Abs(got-6) > 1e-12 {
		t.Errorf("2·3 = %v, oczekiwano 6", got)
	}
}

func TestPowerRightAssoc(t *testing.T) {
	ast, err := parseExpr("2^3^2")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := ast.eval(0); math.Abs(got-512) > 1e-9 {
		t.Errorf("2^3^2 = %v, oczekiwano 512", got)
	}
}

func TestParseErrors(t *testing.T) {
	// Uwaga: „x 3” i „(x)” są POPRAWNE (implikowane mnożenie / nawiasy).
	for _, s := range []string{"x*(*", "x +", "foo(x)", "(x", "x@2"} {
		if _, err := parseExpr(s); err == nil {
			t.Errorf("parse(%q) nie zwrócił błędu", s)
		}
	}
}

func TestPlotParabola(t *testing.T) {
	p, err := PlotFun("x*x-6*x+3", 0, 6, 1601)
	if err != nil {
		t.Fatalf("PlotFun: %v", err)
	}
	if len(p.Segments) != 1 {
		t.Fatalf("oczekiwano 1 segmentu, jest %d", len(p.Segments))
	}
	if len(p.Asymptotes) != 0 {
		t.Fatalf("oczekiwano 0 asymptot, jest %v", p.Asymptotes)
	}
}

func TestPlotReciprocal(t *testing.T) {
	p, err := PlotFun("1/x", -5, 5, 1601)
	if err != nil {
		t.Fatalf("PlotFun: %v", err)
	}
	if len(p.Asymptotes) != 1 {
		t.Fatalf("oczekiwano 1 asymptoty, jest %v", p.Asymptotes)
	}
	if math.Abs(p.Asymptotes[0]) > 0.05 {
		t.Errorf("asymptota w %v, oczekiwano ~0", p.Asymptotes[0])
	}
	// 2 gałęzie hiperboli
	if len(p.Segments) < 2 {
		t.Fatalf("oczekiwano >= 2 segmentów, jest %d", len(p.Segments))
	}
}

func TestPlotTan(t *testing.T) {
	p, err := PlotFun("tan(x)", -6.3, 6.3, 4001)
	if err != nil {
		t.Fatalf("PlotFun: %v", err)
	}
	// asymptoty przy ±π/2, ±3π/2
	if len(p.Asymptotes) != 4 {
		t.Fatalf("oczekiwano 4 asymptot, jest %v", p.Asymptotes)
	}
}

func TestPlotShiftedAsymptote(t *testing.T) {
	// asymptota przy x = 2
	p, err := PlotFun("1/(x-2)", 0, 4, 1601)
	if err != nil {
		t.Fatalf("PlotFun: %v", err)
	}
	if len(p.Asymptotes) != 1 || math.Abs(p.Asymptotes[0]-2) > 0.05 {
		t.Fatalf("asymptoty: %v, oczekiwano [2]", p.Asymptotes)
	}
}
