package z03

import (
	"context"
	"testing"
)

// Liczby pierwsze do 10^n — wartości referencyjne (zgodne z Python
// z zadania): 168, 1229, 9592, 78498, 664579, 5761455.
func TestPrimeCounts(t *testing.T) {
	cases := []struct {
		n    int
		want int
	}{
		{1000, 168},
		{10_000, 1229},
		{100_000, 9592},
		{1_000_000, 78498},
	}
	for _, tc := range cases {
		if got := len(SieveSlice(tc.n)); got != tc.want {
			t.Errorf("SieveSlice(%d) = %d, oczekiwano %d", tc.n, got, tc.want)
		}
		if got := len(SieveMap(tc.n)); got != tc.want {
			t.Errorf("SieveMap(%d) = %d, oczekiwano %d", tc.n, got, tc.want)
		}
		if got := len(SieveBit(tc.n)); got != tc.want {
			t.Errorf("SieveBit(%d) = %d, oczekiwano %d", tc.n, got, tc.want)
		}
	}
	// duże n — tylko spójność między strukturami (ciężki test: kilka sekund)
	if !testing.Short() {
		n := 10_000_000
		a, b, c := len(SieveSlice(n)), len(SieveMap(n)), len(SieveBit(n))
		if a != b || b != c {
			t.Errorf("niezgodne wyniki: slice=%d map=%d bit=%d", a, b, c)
		}
		if a != 664579 {
			t.Errorf("Sieve(%d) = %d, oczekiwano 664579", n, a)
		}
	}
}

func TestBenchmark(t *testing.T) {
	rows, err := Benchmark(context.Background(), []int64{1000, 10_000}, nil)
	if err != nil {
		t.Fatalf("Benchmark: %v", err)
	}
	if len(rows) != 6 {
		t.Fatalf("oczekiwano 6 wierszy, jest %d", len(rows))
	}
	for _, r := range rows {
		if r.Seconds < 0 {
			t.Errorf("ujemny czas: %+v", r)
		}
		if r.Count == 0 {
			t.Errorf("zero liczb pierwszych: %+v", r)
		}
	}
}
