package z07

import "testing"

// Regression: reguła braku zera wiodącego nie miała testu regulującego —
// dla „AB - CD = EF” daje 476 rozwiązań, po jej usunięciu byłoby 692.
func TestLeadingZeroRule(t *testing.T) {
	res, err := Solve("AB - CD = EF")
	if err != nil {
		t.Fatal(err)
	}
	if res.Count != 476 {
		t.Fatalf("Count=%d, oczekiwano 476", res.Count)
	}
	for _, s := range res.Solutions {
		for _, letter := range []string{"A", "C", "E"} { // pierwsze cyfry liczb AB, CD, EF
			if v, ok := s.Mapping[letter]; ok && v == 0 {
				t.Fatalf("zero wiodące %s=0 w %q", letter, s.Equation)
			}
		}
	}
}
