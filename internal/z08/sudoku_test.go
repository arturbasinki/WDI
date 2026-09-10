package z08

import (
	"testing"
)

// Plansza wygenerowana własnym generatorem (seed 42, 56 pustych pól)
// i niezależnie zweryfikowana (Python, brute force): dokładnie 1 rozwiązanie.
var verifiedPuzzle = []int{
	0, 0, 0, 7, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 5, 0, 0, 2,
	7, 0, 0, 9, 8, 0, 0, 0, 3,
	0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 1, 0, 5, 0, 0, 0, 9, 0,
	0, 0, 6, 4, 3, 0, 8, 0, 5,
	0, 9, 0, 6, 5, 0, 0, 0, 0,
	0, 0, 4, 0, 0, 9, 3, 6, 0,
	0, 0, 8, 0, 0, 3, 0, 1, 0,
}

func TestCountSolutionsVerified(t *testing.T) {
	var g Grid
	copy(g[:], verifiedPuzzle)
	n := CountSolutions(g, 3)
	if n != 1 {
		t.Fatalf("oczekiwano 1 rozwiązania, jest %d", n)
	}
}

func TestSolveVerified(t *testing.T) {
	res, err := Solve(verifiedPuzzle, 3)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	if !res.Solved || res.Count != 1 {
		t.Fatalf("oczekiwano 1 rozwiązania, jest %d", res.Count)
	}
	// weryfikacja poprawności rozwiązania
	if !isFull(res.Solution) || !isConsistent(res.Solution) {
		t.Fatal("rozwiązanie nie jest kompletną, poprawną planszą")
	}
	// wskazówki muszą się zgadzać
	for i, v := range verifiedPuzzle {
		if v != 0 && res.Solution[i] != v {
			t.Errorf("pole %d: wskazówka %d ≠ rozwiązanie %d", i+1, v, res.Solution[i])
		}
	}
}

func isFull(g Grid) bool {
	for _, v := range g {
		if v == 0 {
			return false
		}
	}
	return true
}

func isConsistent(g Grid) bool {
	for i := 0; i < 81; i++ {
		v := g[i]
		r, c := i/9, i%9
		for j := 0; j < 9; j++ {
			if j == i%9 || g[r*9+j] == v {
				continue
			}
			if g[r*9+j] == v {
				return false
			}
		}
		for j := 0; j < 9; j++ {
			if j == r || g[j*9+c] == v {
				continue
			}
			if g[j*9+c] == v {
				return false
			}
		}
		br, bc := (r/3)*3, (c/3)*3
		for a := 0; a < 3; a++ {
			for b := 0; b < 3; b++ {
				ii := (br+a)*9 + (bc + b)
				if ii != i && g[ii] == v {
					return false
				}
			}
		}
	}
	return true
}

func TestGenerate(t *testing.T) {
	res, err := Generate(42)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	// unikalność
	var pz Grid
	copy(pz[:], res.Puzzle[:])
	if n := CountSolutions(pz, 2); n != 1 {
		t.Fatalf("zagadka ma %d rozwiązań, oczekiwano 1", n)
	}
	// wskazówki muszą być spójne z rozwiązaniem
	for i := range 81 {
		if res.Puzzle[i] != 0 && res.Puzzle[i] != res.Solution[i] {
			t.Fatalf("pole %d: wskazówka %d ≠ rozwiązanie %d", i+1, res.Puzzle[i], res.Solution[i])
		}
	}
	if res.Blanks < 30 {
		t.Errorf("za mało pustych pól: %d", res.Blanks)
	}
	t.Logf("puste pola: %d, czas: %d ms", res.Blanks, res.TimeMs)
}

func TestValidateGrid(t *testing.T) {
	if err := ValidateGrid(make([]int, 80)); err == nil {
		t.Error("za krótki grid zaakceptowany")
	}
	g := make([]int, 81)
	g[5] = 11
	if err := ValidateGrid(g); err == nil {
		t.Error("wartość 11 zaakceptowana")
	}
}
