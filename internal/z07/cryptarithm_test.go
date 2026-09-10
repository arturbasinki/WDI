package z07

import (
	"testing"
)

// Przykład z zadania: ABC * BD = EFGAH, m.in. 623 * 29 = 18067.
func TestExample(t *testing.T) {
	res, err := Solve("ABC * BD = EFGAH")
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	if res.Count < 1 {
		t.Fatal("brak rozwiązań")
	}
	found := false
	for _, s := range res.Solutions {
		if s.Equation == "623 * 29 = 18067" {
			found = true
			if s.Mapping["A"] != 6 || s.Mapping["B"] != 2 || s.Mapping["C"] != 3 ||
				s.Mapping["D"] != 9 || s.Mapping["E"] != 1 || s.Mapping["F"] != 8 ||
				s.Mapping["G"] != 0 || s.Mapping["H"] != 7 {
				t.Errorf("nieprawidłowe przypisanie: %v", s.Mapping)
			}
		}
	}
	if !found {
		t.Fatalf("rozwiązanie 623 * 29 = 18067 nie występuje; pierwsze: %v", res.Solutions[0].Equation)
	}
	t.Logf("znaleziono %d rozwiązań (przestrzeń %d)", res.Count, res.SpaceSize)
}

func TestAddition(t *testing.T) {
	// SEND + MORE = MONEY — klasyczny przykład, 1 rozwiązanie: 9567 + 1085 = 10652
	res, err := Solve("SEND + MORE = MONEY")
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	if res.Count != 1 {
		t.Fatalf("oczekiwano 1 rozwiązania, jest %d", res.Count)
	}
	if res.Solutions[0].Equation != "9567 + 1085 = 10652" {
		t.Fatalf("wynik: %q, oczekiwano %q", res.Solutions[0].Equation, "9567 + 1085 = 10652")
	}
}

func TestDivision(t *testing.T) {
	// AB / C = DE z dokładnym dzieleniem
	res, err := Solve("AB / C = DE")
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	if res.Count < 1 {
		t.Fatal("brak rozwiązań")
	}
	// weryfikacja: każde rozwiązanie musi spełniać A*10+B = C*(10*D+E)
	for _, s := range res.Solutions {
		m := s.Mapping
		a, b, c, d, e := m["A"], m["B"], m["C"], m["D"], m["E"]
		if c == 0 || (a*10+b)%c != 0 || (a*10+b)/c != d*10+e {
			t.Errorf("błędne rozwiązanie dzielenia: %q", s.Equation)
		}
	}
	t.Logf("AB/C=DE: %d rozwiązań", res.Count)
}

func TestErrors(t *testing.T) {
	for _, eq := range []string{
		"ABC = BD",        // za mało tokenów
		"A ? B = C",       // zły operator
		"AB1 * C = DE",    // cyfra zamiast litery
		"ABCDEFGHIJKLMNOP * C = DE", // > 10 liter
	} {
		if _, err := Solve(eq); err == nil {
			t.Errorf("Solve(%q) nie zwrócił błędu", eq)
		}
	}
}
