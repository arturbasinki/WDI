package z01

import (
	"math"
	"testing"
)

func TestParseExample(t *testing.T) {
	text := "3\n1 2 3 7\n-1 2 4 6\n2 1 1 13\n"
	A, err := Parse(text)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(A) != 3 || len(A[0]) != 4 {
		t.Fatalf("nieprawidłowy kształt macierzy: %v", A)
	}
	// wynik wersji Python: [18, -58, 35]
	res := Solve(A)
	if res.Status != StatusUnique {
		t.Fatalf("oczekiwano jednoznacznego, jest %s", res.Status)
	}
	want := []float64{18, -58, 35}
	for i, v := range want {
		if math.Abs(res.X[i]-v) > 1e-6 {
			t.Errorf("x[%d] = %v, oczekiwano %v", i, res.X[i], v)
		}
	}
}

func TestContradictory(t *testing.T) {
	// 2x+y=1 oraz 4x+2y=6 → sprzeczny
	A := [][]float64{
		{2, 1, 1},
		{4, 2, 6},
	}
	res := Solve(A)
	if res.Status != StatusNone {
		t.Fatalf("oczekiwano sprzecznego, jest %s", res.Status)
	}
}

func TestIndeterminate(t *testing.T) {
	// x+y=1 oraz 2x+2y=2 → nieoznaczony
	A := [][]float64{
		{1, 1, 1},
		{2, 2, 2},
	}
	res := Solve(A)
	if res.Status != StatusInfinite {
		t.Fatalf("oczekiwano nieoznaczonego, jest %s", res.Status)
	}
	if len(res.FreeVars) != 1 {
		t.Fatalf("oczekiwano 1 zmiennej wolnej, jest %v", res.FreeVars)
	}
}

func TestMinusWithSpace(t *testing.T) {
	// format z README: „- 1 2 4 6”
	text := "2\n1 1 3\n- 1 1 1\n"
	A, err := Parse(text)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	res := Solve(A)
	if res.Status != StatusUnique {
		t.Fatalf("oczekiwano jednoznacznego, jest %s", res.Status)
	}
	if math.Abs(res.X[0]-1) > 1e-9 || math.Abs(res.X[1]-2) > 1e-9 {
		t.Errorf("oczekiwano [1,2], jest %v", res.X)
	}
}

func TestStepsSanity(t *testing.T) {
	A, _ := Parse("3\n1 2 3 7\n-1 2 4 6\n2 1 1 13\n")
	res := Solve(A)
	if len(res.Steps) < 4 {
		t.Fatalf("za mało kroków: %d", len(res.Steps))
	}
	if len(res.Steps[0].Matrix) != 3 {
		t.Fatal("pierwszy krok musi mieć 3 wiersze")
	}
	// RREF: macierz jednostkowa + wyrazy wolne
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			want := 0.0
			if i == j {
				want = 1.0
			}
			if math.Abs(res.RREF[i][j]-want) > 1e-9 {
				t.Errorf("RREF[%d][%d] = %v, oczekiwano %v", i, j, res.RREF[i][j], want)
			}
		}
	}
}
