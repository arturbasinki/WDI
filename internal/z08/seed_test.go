package z08

import "testing"

// Regression: determinizm seeda (spec: „seed → deterministycznie”)
// nie był pilnowany żadnym testem.
func TestGenerateDeterministic(t *testing.T) {
	a, err := Generate(42)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Generate(42)
	if err != nil {
		t.Fatal(err)
	}
	if a.Puzzle != b.Puzzle || a.Solution != b.Solution {
		t.Fatal("Generate(42) nie jest deterministyczne")
	}
}
