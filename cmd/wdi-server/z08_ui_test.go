// Testy jednego UI sudoku (z08) — pilnują decyzji projektowej: generowanie,
// ręczne wpisywanie, sprawdzanie i rozwiązywanie dzieje się na JEDNEJ
// graficznej planszy. „Pokaż rozwiązanie” i „Sprawdź / rozwiąż” były dwiema
// nazwami tej samej akcji, a tekstowy solver był drugim, osobnym UI.
package main

import (
	"strings"
	"testing"
)

func TestZ08JednaPlanszaJednoUI(t *testing.T) {
	js := appJS(t)

	if n := strings.Count(js, "VIEWS.z08 ="); n != 1 {
		t.Fatalf("widok z08: %d definicje, oczekiwano 1", n)
	}

	// Jedna plansza i jedna ścieżka do każdej akcji.
	for _, want := range []string{
		`id="grid"`,
		`id="check"`,
		`id="solve"`,
		`id="clear"`,
		`id="blank"`,
		`id="own-run"`,
		`/api/z08/generate`,
		`/api/z08/solve`,
		"Sprawdź",
		"Rozwiąż",
	} {
		if !strings.Contains(js, want) {
			t.Errorf("app.js nie zawiera %q — brakuje elementu zunifikowanego UI z08", want)
		}
	}

	// Stare, zdublowane funkcje i osobny solver tekstowy nie mogą wrócić.
	for _, banned := range []string{
		"Pokaż rozwiązanie",
		"Sprawdź / rozwiąż",
		`id="show"`,
		`id="own-out"`,
		`id="stats"`,
	} {
		if strings.Contains(js, banned) {
			t.Errorf("app.js zawiera %q — to stary, zdublowany fragment UI z08", banned)
		}
	}
}
