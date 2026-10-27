package z07

import "testing"

// TestLimitCapsSolutionsKeepsCount — limit przycina TYLKO zwracane Solutions.
// Count jest zawsze pełną liczbą rozwiązań (żądanie P2), a Truncated mówi,
// czy lista jest okrojona. limit <= 0 = bez limitu (wszystkie rozwiązania).
func TestLimitCapsSolutionsKeepsCount(t *testing.T) {
	full, err := Solve("AB - CD = EF", 0)
	if err != nil {
		t.Fatal(err)
	}
	if full.Count != 476 {
		t.Fatalf("Count=%d, oczekiwano 476", full.Count)
	}
	if len(full.Solutions) != full.Count {
		t.Fatalf("bez limitu: Solutions=%d, Count=%d (lista powinna być pełna)", len(full.Solutions), full.Count)
	}
	if full.Truncated {
		t.Error("bez limitu: Truncated=true, oczekiwano false")
	}

	for _, limit := range []int{1, 7, 100} {
		res, err := Solve("AB - CD = EF", limit)
		if err != nil {
			t.Fatal(err)
		}
		if res.Count != 476 {
			t.Errorf("limit=%d: Count=%d, oczekiwano pełnego 476", limit, res.Count)
		}
		if len(res.Solutions) != limit {
			t.Errorf("limit=%d: zwrócono %d rozwiązań", limit, len(res.Solutions))
		}
		if !res.Truncated {
			t.Errorf("limit=%d: Truncated=false, oczekiwano true", limit)
		}
		// przycięta lista to prefix pełnej — ten sam porządek przeszukiwania
		for i := range res.Solutions {
			if res.Solutions[i].Equation != full.Solutions[i].Equation {
				t.Fatalf("limit=%d: Solutions[%d]=%q, w pełnej liście %q",
					limit, i, res.Solutions[i].Equation, full.Solutions[i].Equation)
			}
		}
	}

	// limit większy niż liczba rozwiązań: nic nie jest ucinane
	res, err := Solve("AB - CD = EF", 500)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Solutions) != 476 || res.Truncated {
		t.Errorf("limit=500: Solutions=%d, Truncated=%v", len(res.Solutions), res.Truncated)
	}

	// limit ujemny = bez limitu (kontrakt: limit <= 0)
	neg, err := Solve("SEND + MORE = MONEY", -1)
	if err != nil {
		t.Fatal(err)
	}
	if neg.Count != 1 || len(neg.Solutions) != 1 {
		t.Errorf("limit=-1: Count=%d, Solutions=%d", neg.Count, len(neg.Solutions))
	}
}

// TestDefaultSolutions — jawna wartość domyślna stosowana przez handler.
func TestDefaultSolutions(t *testing.T) {
	if DefaultSolutions != 50 || MaxSolutions != 500 {
		t.Fatalf("DefaultSolutions=%d MaxSolutions=%d, oczekiwano 50/500", DefaultSolutions, MaxSolutions)
	}
	res, err := Solve("AB - CD = EF", DefaultSolutions)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Solutions) != DefaultSolutions {
		t.Errorf("przy domyślnym limicie zwrócono %d rozwiązań", len(res.Solutions))
	}
}
