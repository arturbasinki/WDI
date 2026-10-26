package z01

import "testing"

// Regression: każdy krok animacji kopiował pełną macierz; dla n=100
// dawało to ~10 000 snapshotów (~1,1 GB JSON, ~4,7 GB RSS — OOM).
func TestSolveStepsBudget(t *testing.T) {
	n := 100
	A := make([][]float64, n)
	for i := range A {
		A[i] = make([]float64, n+1)
		for j := 0; j < n; j++ {
			A[i][j] = float64((i*7+j*13)%5) - 2
		}
		A[i][i] += 20 // dominacja przekątnej — układ jednoznaczny
		A[i][n] = float64(i)
	}
	res := Solve(A)
	if !res.StepsTruncated {
		t.Fatal("dla n=100 oczekiwano stepsTruncated=true")
	}
	if cells := n * (n + 1) * len(res.Steps); cells > 400_000 {
		t.Fatalf("przekroczono budżet kroków: %d komórek", cells)
	}
	if res.Status != StatusUnique {
		t.Fatalf("wynik musi pozostać poprawny, status=%v", res.Status)
	}
	res2 := Solve([][]float64{{1, 2, 3, 7}, {-1, 2, 4, 6}, {2, 1, 1, 13}})
	if res2.StepsTruncated || len(res2.Steps) < 4 {
		t.Fatalf("mały układ ma mieć pełną animację: trunc=%v steps=%d", res2.StepsTruncated, len(res2.Steps))
	}
}
