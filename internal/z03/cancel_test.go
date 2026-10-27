package z03

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

// TestBenchmarkCancelInsideMeasurement — anulowanie musi przerywać POJEDYNCZY
// pomiar, nie czekać do jego końca. Przed zmianą ctx.Err() sprawdzano tylko
// MIĘDZY pomiarami, a najwolniejszy pomiar (n = 10⁸) trwa sekundy: DELETE
// na serwerze potrafiło czekać kilkanaście sekund na zwolnienie semafora.
func TestBenchmarkCancelInsideMeasurement(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	start := time.Now()
	rows, err := Benchmark(ctx, []int64{1e8}, nil)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatalf("oczekiwano przerwania pomiaru, zwrócono %d wierszy bez błędu", len(rows))
	}
	if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		t.Fatalf("oczekiwano błędu kontekstu, jest: %v", err)
	}
	// Dowód, że zadziałało sprawdzenie WEWNĄTRZ pomiaru: sam pomiar n = 10⁸
	// (100 MB flag) zajmuje co najmniej setki milisekund, a anulowanie
	// dotarło natychmiast po upływie terminu.
	if elapsed < 30*time.Millisecond {
		t.Errorf("przerwanie przed upływem terminu (%v) — pomiar w ogóle nie wystartował?", elapsed)
	}
	if elapsed > 3*time.Second {
		t.Errorf("anulowanie pojedynczego pomiaru zajęło %v — sprawdzanie ctx jest za rzadkie", elapsed)
	}
}

// TestSievesCtxCancel — sito mapy (najwolniejsze) przerywa w swoim wnętrzu.
func TestSievesCtxCancel(t *testing.T) {
	if testing.Short() {
		t.Skip("pomijane w -short: budowa mapy 10⁷ potrąca ~0,1 s CPU")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	start := time.Now()
	out, err := SieveMapCtx(ctx, 10_000_000)
	elapsed := time.Since(start)
	if err == nil {
		t.Fatalf("oczekiwano przerwania sita mapy, zwrócono %d liczb pierwszych", len(out))
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("oczekiwano context.DeadlineExceeded, jest: %v", err)
	}
	if out != nil {
		t.Errorf("przy przerwaniu zwrócono niepusty wynik: %d elementów", len(out))
	}
	if elapsed > 2*time.Second {
		t.Errorf("przerwanie sita mapy zajęło %v", elapsed)
	}
}

// TestSievesCtxAgreeWithPlain — ścieżka z kontekstem nie może zmieniać wyniku.
func TestSievesCtxAgreeWithPlain(t *testing.T) {
	for _, n := range []int{1000, 10_001, 99_733} {
		a, err := SieveSliceCtx(context.Background(), n)
		if err != nil {
			t.Fatalf("SieveSliceCtx(%d): %v", n, err)
		}
		b, err := SieveMapCtx(context.Background(), n)
		if err != nil {
			t.Fatalf("SieveMapCtx(%d): %v", n, err)
		}
		c, err := SieveBitCtx(context.Background(), n)
		if err != nil {
			t.Fatalf("SieveBitCtx(%d): %v", n, err)
		}
		if !reflect.DeepEqual(a, SieveSlice(n)) || !reflect.DeepEqual(b, SieveMap(n)) || !reflect.DeepEqual(c, SieveBit(n)) {
			t.Errorf("wersje z kontekstem rozmijają się z wersjami bez dla n=%d", n)
		}
		if !reflect.DeepEqual(a, b) || !reflect.DeepEqual(b, c) {
			t.Errorf("trzy struktury dają różne wyniki dla n=%d", n)
		}
	}
}

// TestBenchmarkWithoutCancel — Live kontekst: pełny zestaw wierszy, zero błędów.
func TestBenchmarkWithoutCancel(t *testing.T) {
	rows, err := Benchmark(context.Background(), []int64{1000, 10_000}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 6 { // 2 wartości n × 3 struktury
		t.Fatalf("wierszy: %d, oczekiwano 6", len(rows))
	}
	for _, r := range rows {
		if r.Count != 168 && r.Count != 1229 {
			t.Errorf("n=%d %s: count=%d (oczekiwano 168 dla 10³ i 1229 dla 10⁴)", r.N, r.Structure, r.Count)
		}
	}
}
