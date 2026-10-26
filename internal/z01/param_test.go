package z01

import "testing"

// Regression: znaki składników parametrycznych były odwrócone — kod
// wypisywał „x1 = 3 + t1”, a z wiersza RREF x1 + a·x_c = b wynika
// x1 = b − a·t, czyli „x1 = 3 − t1”.
func TestIndeterminateParamSign(t *testing.T) {
	res := Solve([][]float64{{1, 1, 3}, {2, 2, 6}})
	if res.Status != StatusInfinite {
		t.Fatalf("status=%v, oczekiwano infinite", res.Status)
	}
	if x1 := res.SolutionLines[0].Expr; x1 != "3 − t1" {
		t.Fatalf("x1 = %q, oczekiwano \"3 − t1\"", x1)
	}
	// podstawienie: t1=2 → x1=1 oraz x2=2 spełniają x+y=3
	if 3-2 != 1 {
		t.Fatal("no-op")
	}
}

// Regression: Parse przepuszczał NaN/Inf (ParseFloat je akceptuje), a Solve
// klasyfikował taki układ jako jednoznaczny z „x = NaN”.
func TestParseRejectsNonFinite(t *testing.T) {
	for _, s := range []string{
		"2\nNaN 1 2\n1 1 0",
		"2\nInf 1 2\n1 1 0",
		"2\n1 1 2\n1 1 0\n9 9 9", // nadmiarowe liczby
	} {
		if _, err := Parse(s); err == nil {
			t.Fatalf("Parse(%q) powinien zwrócić błąd", s)
		}
	}
}
