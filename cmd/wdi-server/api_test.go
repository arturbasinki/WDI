// Testy warstwy HTTP dla zasad API wprowadzonych w P2 (dostępność i prawda
// w interfejsie/bibliotekach). Korzystają z pomocników z handlers_test.go.
//
// Cel: pilnują, że biblioteka nie chowa granic, a odpowiedź JSON mówi, co
// realnie użyto (n, samplesRequested, skipped, zwrócone vs count).
package main

import (
	"math"
	"testing"
)

// TestZ02DolnyProgJestWHTTP — biblioteka z02.Search przyjmuje każde limit >= 0
// (test jednostkowy TestSmallLimits), a warstwa HTTP nadal odrzuca limit
// poniżej 1000: klient nie może wystartować joba, który nic nie wnosi.
//
// UWAGA na granicę: `limit: 0` (i brak pola) znaczy „użyj domyślnego
// DefaultLimit" — to obsługa NIEOBECNEGO parametru, nie próg dolny, więc
// nie jest tu sprawdzane (startowałoby prawdziwy job na 10¹⁰).
func TestZ02DolnyProgJestWHTTP(t *testing.T) {
	srv, jobs, _ := newTestAPI(t)
	for _, limit := range []int{1, 100, 999} {
		code, body := do(t, "POST", srv.URL+"/api/z02/search", map[string]any{"limit": limit})
		if code != 400 {
			t.Errorf("limit=%d: HTTP %d, oczekiwano 400 (%s)", limit, code, errMsg(body))
		}
	}
	if n := runningCount(jobs); n != 0 {
		t.Errorf("odrzucone żądania nie mogą tworzyć jobów (jest %d)", n)
	}
}

// TestZ07LimitPrzycinaListeNieCount — limit z żądania: 0 = domyślnie 50,
// sufit 500 (przycięcie żądania), ujemny = 400. count jest pełny zawsze.
func TestZ07LimitPrzycinaListeNieCount(t *testing.T) {
	srv, _, _ := newTestAPI(t)

	// domyślnie: 50 zwróconych, pełny count 476 (AB - CD = EF)
	code, body := do(t, "POST", srv.URL+"/api/z07/solve", map[string]any{"equation": "AB - CD = EF"})
	if code != 200 {
		t.Fatalf("HTTP %d: %s", code, errMsg(body))
	}
	if got := len(arr(t, body, "solutions")); got != 50 {
		t.Errorf("bez limitu: zwrócono %d rozwiązań, oczekiwano domyślnych 50", got)
	}
	if got := num(body, "count"); got != 476 {
		t.Errorf("count=%v, oczekiwano pełnego 476", got)
	}

	// jawny mały limit
	code, body = do(t, "POST", srv.URL+"/api/z07/solve", map[string]any{"equation": "AB - CD = EF", "limit": 3})
	if code != 200 || len(arr(t, body, "solutions")) != 3 || num(body, "count") != 476 || body["truncated"] != true {
		t.Errorf("limit=3: HTTP %d solutions=%d count=%v truncated=%v",
			code, len(arr(t, body, "solutions")), num(body, "count"), body["truncated"])
	}

	// sufit 500: żądanie ponad sufit jest przycinane, count pozostaje pełny
	code, body = do(t, "POST", srv.URL+"/api/z07/solve", map[string]any{"equation": "ABC + DEF = GHI", "limit": 100000})
	if code != 200 {
		t.Fatalf("limit=100000: HTTP %d: %s", code, errMsg(body))
	}
	if got := len(arr(t, body, "solutions")); got != 500 {
		t.Errorf("limit=100000: zwrócono %d, oczekiwano 500 (sufit)", got)
	}
	if num(body, "count") != 1088 {
		t.Errorf("count=%v, oczekiwano 1088 (pełna liczba rozwiązań)", num(body, "count"))
	}

	// ujemny limit: błąd warstwy HTTP (biblioteka interpretuje <= 0 jako „bez limitu”)
	code, body = do(t, "POST", srv.URL+"/api/z07/solve", map[string]any{"equation": "AB - CD = EF", "limit": -5})
	if code != 400 {
		t.Errorf("limit=-5: HTTP %d, oczekiwano 400 (%s)", code, errMsg(body))
	}
}

// TestZ04DomyślneNUjawnione — handler stosuje DefaultTop JAWNIE i zwraca użyte n;
// TopWords zwraca błąd dla k <= 0, więc „20” nie bierze się znikąd.
func TestZ04DomyślneNUjawnione(t *testing.T) {
	srv, _, _ := newTestAPI(t)
	code, body := do(t, "GET", srv.URL+"/api/z04/top20", nil)
	if code != 200 {
		t.Fatalf("GET top20: HTTP %d: %s", code, errMsg(body))
	}
	if num(body, "n") != 20 {
		t.Errorf("GET: n=%v, oczekiwano 20", num(body, "n"))
	}
	if got := len(arr(t, body, "words")); got != 20 {
		t.Errorf("GET: %d słów, oczekiwano 20", got)
	}

	// n spoza zakresu => domyślna wartość (i jest widoczna w odpowiedzi)
	code, body = do(t, "POST", srv.URL+"/api/z04/top20", map[string]any{"text": "abcdef abcdef ghijkl", "n": 0})
	if code != 200 || num(body, "n") != 20 {
		t.Errorf("POST n=0: HTTP %d n=%v, oczekiwano 200/20 (%s)", code, num(body, "n"), errMsg(body))
	}
	// n = 1 jest legalne i zwraca dokładnie jedno słowo
	code, body = do(t, "POST", srv.URL+"/api/z04/top20", map[string]any{"text": "abcdef abcdef ghijkl", "n": 1})
	if code != 200 || len(arr(t, body, "words")) != 1 {
		t.Errorf("POST n=1: HTTP %d words=%d", code, len(arr(t, body, "words")))
	}
	// tekst, z którego nie da się nic wyciągnąć => pusta lista, nie błąd 500
	code, body = do(t, "POST", srv.URL+"/api/z04/top20", map[string]any{"text": "ala ma kot", "n": 5})
	if code != 200 || len(arr(t, body, "words")) != 0 {
		t.Errorf("POST tekst bez długich słów: HTTP %d words=%d", code, len(arr(t, body, "words")))
	}
}

// TestZ05SkippedWJSON — normalizacja nie jest cicha: pola „skipped” nie ma,
// gdy tekst mieści się w alfabecie; gdy nie mieści, są pominięte znaki.
func TestZ05SkippedWJSON(t *testing.T) {
	srv, _, _ := newTestAPI(t)
	code, body := do(t, "POST", srv.URL+"/api/z05/crypt",
		map[string]any{"mode": "enc", "text": "computer science", "key": "312063"})
	if code != 200 {
		t.Fatalf("enc: HTTP %d: %s", code, errMsg(body))
	}
	if _, ok := body["skipped"]; ok {
		t.Errorf("tekst w alfabecie nie powinien mieć pola skipped: %v", body["skipped"])
	}
	if body["out"] != "fpop whsbsilhoee" {
		t.Errorf("out=%v", body["out"])
	}

	code, body = do(t, "POST", srv.URL+"/api/z05/crypt",
		map[string]any{"mode": "enc", "text": "ZAŻÓŁĆ GĘŚLĄ", "key": "1"})
	if code != 200 {
		t.Fatalf("enc Zażółć: HTTP %d: %s", code, errMsg(body))
	}
	got := arr(t, body, "skipped")
	want := []string{"ó", "ą", "ć", "ę", "ł", "ś", "ż"}
	if len(got) != len(want) {
		t.Fatalf("skipped=%v, oczekiwano %v", got, want)
	}
	for i := range want {
		if s, _ := got[i].(string); s != want[i] {
			t.Errorf("skipped[%d]=%v, oczekiwano %q (posortowane po punkcie kodowym)", i, got[i], want[i])
		}
	}

	// Break: przycięcie capKeys do sufitu jest odrzucane w HTTP (400),
	// a nie udawane cichym okrojeniem pracy.
	code, body = do(t, "POST", srv.URL+"/api/z05/break",
		map[string]any{"cipher": "fpop whsbsilhoee", "maxLen": 6, "capKeys": 200_000_000})
	if code != 400 {
		t.Errorf("break capKeys>1e8: HTTP %d, oczekiwano 400 (%s)", code, errMsg(body))
	}
}

// TestZ06PrawdaOPrzycieciu — samplesRequested obok samples, przycięte Y
// zachowuje punkt (dociera do ±yclamp) i notacja wykładnicza.
func TestZ06PrawdaOPrzycieciu(t *testing.T) {
	srv, _, _ := newTestAPI(t)

	code, body := do(t, "POST", srv.URL+"/api/z06/plot",
		map[string]any{"expr": "x", "xmin": 0, "xmax": 1, "samples": 5_000_000})
	if code != 200 {
		t.Fatalf("plot: HTTP %d: %s", code, errMsg(body))
	}
	if num(body, "samples") != 20001 || num(body, "samplesRequested") != 5_000_000 {
		t.Errorf("samples=%v samplesRequested=%v, oczekiwano 20001/5000000",
			num(body, "samples"), num(body, "samplesRequested"))
	}

	code, body = do(t, "POST", srv.URL+"/api/z06/plot",
		map[string]any{"expr": "1/x", "xmin": -5, "xmax": 5, "samples": 1601})
	if code != 200 {
		t.Fatalf("plot 1/x: HTTP %d: %s", code, errMsg(body))
	}
	clamp := num(body, "yclamp")
	maxAbs, pts := 0.0, 0
	for _, seg := range arr(t, body, "segments") {
		s, ok := seg.([]any)
		if !ok {
			continue
		}
		for _, p := range s {
			pt, ok := p.(map[string]any)
			if !ok {
				continue
			}
			y := num(pt, "y")
			if math.Abs(y) > maxAbs {
				maxAbs = math.Abs(y)
			}
			pts++
		}
	}
	if pts == 0 {
		t.Fatal("brak punktów w segmentach")
	}
	if math.Abs(maxAbs-clamp) > 1e-12 {
		t.Errorf("najdalszy |y| = %v, yclamp = %v — gałąź musi DOCIERAĆ do progu (punkt zachowany)", maxAbs, clamp)
	}

	// notacja wykładnicza: 2e3 to liczba 2000
	code, body = do(t, "POST", srv.URL+"/api/z06/plot",
		map[string]any{"expr": "2e3", "xmin": 0, "xmax": 1, "samples": 601})
	if code != 200 {
		t.Fatalf("plot 2e3: HTTP %d: %s", code, errMsg(body))
	}
	segs := arr(t, body, "segments")
	if len(segs) == 0 {
		t.Fatal("brak segmentów dla 2e3")
	}
	first, ok := segs[0].([]any)
	if !ok || len(first) == 0 {
		t.Fatal("pusty segment dla 2e3")
	}
	pt, ok := first[0].(map[string]any)
	if !ok || math.Abs(num(pt, "y")-2000) > 1e-9 {
		t.Errorf("2e3: pierwszy y = %v, oczekiwano 2000", pt["y"])
	}
}

// TestZ06KomunikatNaRunie — 400 musi pokazywać znak, nie jego bajt (π, nie Ï).
func TestZ06KomunikatNaRunie(t *testing.T) {
	srv, _, _ := newTestAPI(t)
	code, body := do(t, "POST", srv.URL+"/api/z06/plot",
		map[string]any{"expr": "x + π", "xmin": 0, "xmax": 1, "samples": 601})
	if code != 400 {
		t.Fatalf("HTTP %d, oczekiwano 400", code)
	}
	msg := errMsg(body)
	if !stringsContains(msg, "π") || stringsContains(msg, "Ï") {
		t.Errorf("komunikat %q nie pokazuje runy „π”", msg)
	}
}

// ---------------------------------------------------------------------------
// pomocniki odczytu JSON (map[string]any z encoding/json)
// ---------------------------------------------------------------------------

func num(m map[string]any, key string) float64 {
	v, ok := m[key]
	if !ok {
		return -1
	}
	f, ok := v.(float64)
	if !ok {
		return -1
	}
	return f
}

func arr(t *testing.T, m map[string]any, key string) []any {
	t.Helper()
	v, ok := m[key]
	if !ok {
		t.Fatalf("brak pola %q w odpowiedzi: %v", key, m)
	}
	a, ok := v.([]any)
	if !ok {
		t.Fatalf("pole %q nie jest tablicą: %T", key, v)
	}
	return a
}

func stringsContains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
