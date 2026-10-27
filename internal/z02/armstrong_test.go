package z02

import (
	"context"
	"testing"
)

func TestSearchUpTo1e6(t *testing.T) {
	var last Progress
	res, err := Search(context.Background(), 1_000_000, func(pr Progress) { last = pr })
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if !last.Done {
		t.Error("ostatni tick postępu nie oznaczył Done")
	}
	// Wyniki Python do 10^6 (i weryfikacja własna): 0-9, 153, 370, 371,
	// 407, 1634, 8208, 9474, 54748, 92727, 93084, 548834
	want := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 153, 370, 371, 407, 1634, 8208, 9474, 54748, 92727, 93084, 548834}
	if len(res.Found) != len(want) {
		t.Fatalf("znaleziono %d liczb: %v", len(res.Found), res.Found)
	}
	for i, v := range want {
		if res.Found[i] != v {
			t.Errorf("Found[%d] = %d, oczekiwano %d", i, res.Found[i], v)
		}
	}
	if !res.Done {
		t.Fatal("res.Done musi być true")
	}
}

// TestSmallLimits — próg „limit >= 1000” należał do warstwy HTTP, nie do
// algorytmu: biblioteka musi przyjmować dowolny limit >= 0.
func TestSmallLimits(t *testing.T) {
	cases := []struct {
		limit int64
		want  []int64
	}{
		{0, []int64{0}},
		{20, []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{100, []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{152, []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{153, []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 153}},
	}
	for _, tc := range cases {
		res, err := Search(context.Background(), tc.limit, nil)
		if err != nil {
			t.Errorf("limit=%d: %v", tc.limit, err)
			continue
		}
		if !res.Done {
			t.Errorf("limit=%d: Done=false", tc.limit)
		}
		if len(res.Found) != len(tc.want) {
			t.Errorf("limit=%d: Found=%v, oczekiwano %v", tc.limit, res.Found, tc.want)
			continue
		}
		for i, v := range tc.want {
			if res.Found[i] != v {
				t.Errorf("limit=%d: Found[%d]=%d, oczekiwano %d", tc.limit, i, res.Found[i], v)
			}
		}
	}
	// Ujemny limit pozostaje błędem (przedział [0, limit] byłby pusty).
	if _, err := Search(context.Background(), -1, nil); err == nil {
		t.Error("limit ujemny musi zwracać błąd")
	}
}

func TestKnownUpTo1e9(t *testing.T) {
	if testing.Short() {
		t.Skip("pomijane w trybie -short (przeszukiwanie do 10^9 trwa kilka sekund)")
	}
	// Szybsza weryfikacja do 10^9 (Python: dodatkowo 24678050, 24678051, 88593477)
	res, err := Search(context.Background(), 1_000_000_000, nil)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	want := []int64{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 153, 370, 371, 407, 1634,
		8208, 9474, 54748, 92727, 93084, 548834, 1741725, 4210818,
		9800817, 9926315, 24678050, 24678051, 88593477,
		146511208, 472335975, 534494836, 912985153,
	}
	if len(res.Found) != len(want) {
		t.Fatalf("znaleziono %d liczb: %v", len(res.Found), res.Found)
	}
	for i, v := range want {
		if res.Found[i] != v {
			t.Errorf("Found[%d] = %d, oczekiwano %d", i, res.Found[i], v)
		}
	}
}
