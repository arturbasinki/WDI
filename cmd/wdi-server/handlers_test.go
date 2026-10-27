// Testy warstwy HTTP (cmd/wdi-server).httptest, bez portów live.
//
// Długie obliczenia są podstawiane przez szwy testowe (searchArmstrong,
// runBenchmark, breakGronsfeld), dzięki czemu testy są deterministyczne
// i nie czekają na realne obliczenia (żaden sleep > 200 ms).
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"wdi/internal/z02"
)

// ---------------------------------------------------------------------------
// pomocniki
// ---------------------------------------------------------------------------

// fakeClock — wstrzykiwany zegar przechowalni: symuluje upływ TTL bez spania.
type fakeClock struct {
	mu   sync.Mutex
	base time.Time
	off  time.Duration
}

func (c *fakeClock) now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.base.Add(c.off)
}

func (c *fakeClock) advance(d time.Duration) {
	c.mu.Lock()
	c.off += d
	c.mu.Unlock()
}

// newTestAPI — aplikacja na httptest z wstrzyknietym zegarem.
func newTestAPI(t *testing.T) (*httptest.Server, *jobStore, *fakeClock) {
	t.Helper()
	clock := &fakeClock{base: time.Unix(1_700_000_000, 0)}
	jobs := newJobStore()
	jobs.now = clock.now
	srv := httptest.NewServer(newAPI(jobs))
	t.Cleanup(srv.Close)
	return srv, jobs, clock
}

// useSeams — test przejmuje szwy obliczeń. Po teście: najpierw WYPOMPOWANIE
// wszystkich gorutyn jobów (czyszczenie rejestru + oczekiwanie na stan
// końcowy), dopiero potem przywrócenie prawdziwych obliczeń — bez tej
// kolejności globalne szwy są czytane przez gorutynę ze starego testu
// w momencie zapisu (data race pod -race).
func useSeams(t *testing.T, jobs *jobStore) {
	t.Helper()
	s, b, k := searchArmstrong, runBenchmark, breakGronsfeld
	t.Cleanup(func() {
		searchArmstrong, runBenchmark, breakGronsfeld = s, b, k
	})
	t.Cleanup(func() {
		deadline := time.Now().Add(2 * time.Second)
		for {
			jobs.cancelAll()
			if runningCount(jobs) == 0 {
				return
			}
			if time.Now().After(deadline) {
				t.Errorf("gorutyny jobów nie zakończyły się w 2 s")
				return
			}
			time.Sleep(2 * time.Millisecond)
		}
	})
}

// runningCount — liczba wciąż aktywnych jobów w rejestrze.
func runningCount(s *jobStore) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, j := range s.jobs {
		j.mu.Lock()
		if j.status == "running" {
			n++
		}
		j.mu.Unlock()
	}
	return n
}

// blockSearch — symuluje długi job: sygnalizuje start, czeka na anulowanie
// (DELETE / cancelAll) albo na zamknięcie testu (żeby nie zostawiać gorutyn).
func blockSearch(started chan<- struct{}, stop <-chan struct{}) func(
	context.Context, int64, func(z02.Progress)) (z02.Progress, error) {
	return func(ctx context.Context, _ int64, _ func(z02.Progress)) (z02.Progress, error) {
		started <- struct{}{}
		select {
		case <-ctx.Done():
			return z02.Progress{}, ctx.Err()
		case <-stop:
			return z02.Progress{}, context.Canceled
		}
	}
}

// okSearch — tani job kończący się natychmiast (bez realnych obliczeń).
func okSearch(ctx context.Context, limit int64, onTick func(z02.Progress)) (z02.Progress, error) {
	p := z02.Progress{Limit: limit, Processed: limit, Done: true, Found: []int64{0, 1, 153}}
	onTick(p)
	return p, nil
}

// call — niskopoziomowe żądanie (bezpieczne do użycia z gorutyn).
func call(method, url string, body any) (int, map[string]any, error) {
	var rd io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return 0, nil, err
		}
		rd = strings.NewReader(string(b))
	}
	req, err := http.NewRequest(method, url, rd)
	if err != nil {
		return 0, nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	out := map[string]any{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &out) // nie-JSON (np. strona) jest dozwolony
	}
	return resp.StatusCode, out, nil
}

// do — jak call, ale wywala błąd testu (woławane tylko z gorutyny testu).
func do(t *testing.T, method, url string, body any) (int, map[string]any) {
	t.Helper()
	code, out, err := call(method, url, body)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	return code, out
}

func errMsg(body map[string]any) string {
	s, _ := body["error"].(string)
	return s
}

func jobID(t *testing.T, body map[string]any) string {
	t.Helper()
	id, ok := body["jobId"].(string)
	if !ok || id == "" {
		t.Fatalf("brak jobId w odpowiedzi: %v", body)
	}
	return id
}

// waitStatus — czeka (bez realnego obliczenia: status ustawiany jest w
// gorutynie joba), aż job przyjmie oczekiwany stan końcowy.
func waitStatus(t *testing.T, url string, want string) map[string]any {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for {
		code, body := do(t, "GET", url, nil)
		if code != 200 {
			t.Fatalf("GET %s: kod %d, spodziewany 200", url, code)
		}
		if s, _ := body["status"].(string); s == want {
			return body
		}
		if time.Now().After(deadline) {
			t.Fatalf("GET %s: status %v, spodziewany %q", url, body["status"], want)
		}
		time.Sleep(2 * time.Millisecond)
	}
}

// waitPostAccepted — powtarza POST, aż zwolni się miejsce w semaforze
// (zwolnienie następuje w gorutynie joba po ustawieniu statusu).
func waitPostAccepted(t *testing.T, url string) map[string]any {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for {
		code, body := do(t, "POST", url, map[string]any{"limit": int64(1e6)})
		if code == http.StatusAccepted {
			return body
		}
		if code != http.StatusTooManyRequests {
			t.Fatalf("POST %s: kod %d, spodziewany 202/429", url, code)
		}
		if time.Now().After(deadline) {
			t.Fatalf("POST %s: semafor nie zwolnił (kod %d)", url, code)
		}
		time.Sleep(2 * time.Millisecond)
	}
}

const z02URL = "/api/z02/search"

// ---------------------------------------------------------------------------
// 3. DELETE na zakończonym jobie → 409 + komunikat PL
// ---------------------------------------------------------------------------

func TestDeleteNaZakonczonymZwraca409(t *testing.T) {
	srv, jobs, _ := newTestAPI(t)
	useSeams(t, jobs)
	searchArmstrong = okSearch

	code, body := do(t, "POST", srv.URL+z02URL, map[string]any{"limit": int64(1e6)})
	if code != http.StatusAccepted {
		t.Fatalf("POST: kod %d, spodziewany 202 (%v)", code, errMsg(body))
	}
	id := jobID(t, body)
	url := srv.URL + z02URL + "/" + id

	if s := waitStatus(t, url, "done"); s["result"] == nil {
		t.Fatalf("brak wyniku w snapse: %v", s)
	}

	code, body = do(t, "DELETE", url, nil)
	if code != http.StatusConflict {
		t.Fatalf("DELETE zakończonego: kod %d, spodziewany 409", code)
	}
	if msg := errMsg(body); !strings.Contains(msg, "zakończy") {
		t.Fatalf("komunikat 409 po polsku oczekiwany, jest: %q", msg)
	}
	// Stan joba nietknięty: wciąż „done”, nie „cancelling”.
	_, body = do(t, "GET", url, nil)
	if st, _ := body["status"].(string); st != "done" {
		t.Fatalf("status po nieudanym DELETE: %q, spodziewany done", st)
	}
}

func TestDeleteNaNieistniejacym404(t *testing.T) {
	srv, jobs, _ := newTestAPI(t)
	useSeams(t, jobs)
	if code, _ := do(t, "DELETE", srv.URL+z02URL+"/j0-0", nil); code != http.StatusNotFound {
		t.Fatalf("DELETE nieznaczngo: kod %d, spodziewany 404", code)
	}
	if code, _ := do(t, "GET", srv.URL+z02URL+"/j0-0", nil); code != http.StatusNotFound {
		t.Fatalf("GET nieznaczngo: kod %d, spodziewany 404", code)
	}
}

func TestDeleteNaAktywnymAnuluje(t *testing.T) {
	srv, jobs, _ := newTestAPI(t)
	useSeams(t, jobs)
	started := make(chan struct{}, 4)
	stop := make(chan struct{})
	t.Cleanup(func() { close(stop) })
	searchArmstrong = blockSearch(started, stop)

	code, body := do(t, "POST", srv.URL+z02URL, map[string]any{"limit": int64(1e6)})
	if code != http.StatusAccepted {
		t.Fatalf("POST: kod %d", code)
	}
	id := jobID(t, body)
	<-started
	url := srv.URL + z02URL + "/" + id

	code, body = do(t, "DELETE", url, nil)
	if code != http.StatusOK {
		t.Fatalf("DELETE aktywnego: kod %d, spodziewany 200", code)
	}
	if st, _ := body["status"].(string); st != "cancelling" {
		t.Fatalf("status: %q, spodziewany cancelling", st)
	}
	waitStatus(t, url, "cancelled")
	// Po anulowaniu DELETE oddaje już 409.
	if code, _ = do(t, "DELETE", url, nil); code != http.StatusConflict {
		t.Fatalf("DELETE po anulowaniu: kod %d, spodziewany 409", code)
	}
}

// ---------------------------------------------------------------------------
// 4. Walidacja synchroniczna w POST → 400, bez utworzenia joba
// ---------------------------------------------------------------------------

func TestWalidacjaSynchroniczna400(t *testing.T) {
	srv, jobs, _ := newTestAPI(t)
	useSeams(t, jobs)
	// Gdyby któryś POST utworzył joba, zobaczymy to w jobs.count().
	searchArmstrong = func(ctx context.Context, _ int64, _ func(z02.Progress)) (z02.Progress, error) {
		t.Error("obliczenia nie powinny wystartować przy nieprawidłowym parametrze")
		return z02.Progress{}, nil
	}

	cases := []struct {
		name string
		url  string
		body any
		want string
	}{
		{"z02 limit 1e13", z02URL, map[string]any{"limit": int64(1e13)}, "limit"},
		{"z02 limit 10", z02URL, map[string]any{"limit": int64(10)}, "limit"},
		{"z02 limit ujemny", z02URL, map[string]any{"limit": int64(-5)}, "limit"},
		{"z03 n za duże", "/api/z03/benchmark", map[string]any{"ns": []int64{1e9}}, "n musi"},
		{"z03 n za małe", "/api/z03/benchmark", map[string]any{"ns": []int64{7}}, "n musi"},
		{"z05 maxLen", "/api/z05/break", map[string]any{"cipher": "abc", "maxLen": 99}, "maxLen"},
		{"z05 pusty szyfrogram", "/api/z05/break", map[string]any{"cipher": "   "}, "szyfrogram"},
		{"z05 capKeys", "/api/z05/break", map[string]any{"cipher": "abc", "maxLen": 6, "capKeys": int64(1e12)}, "capKeys"},
		{"z06 zła składnia", "/api/z06/plot", map[string]any{"expr": "(", "xmin": 0, "xmax": 1}, ""},
	}
	for _, tc := range cases {
		code, body := do(t, "POST", srv.URL+tc.url, tc.body)
		if code != http.StatusBadRequest {
			t.Errorf("%s: kod %d, spodziewany 400 (%v)", tc.name, code, errMsg(body))
		}
		if tc.want != "" && !strings.Contains(errMsg(body), tc.want) {
			t.Errorf("%s: komunikat %q nie zawiera %q", tc.name, errMsg(body), tc.want)
		}
	}
	if n := jobs.count(); n != 0 {
		t.Fatalf("utworzono %d jobów przy samych błędach wejścia", n)
	}
	// Limit graniczny (1e12) musi przechodzić walidację — był testowany
	// realnie, tu tylko ścieżka HTTP: podstawiamy tani seam.
	searchArmstrong = okSearch
	code, _ := do(t, "POST", srv.URL+z02URL, map[string]any{"limit": z02MaxLimitAPI})
	if code != http.StatusAccepted {
		t.Fatalf("POST limit=1e12: kod %d, spodziewany 202", code)
	}
}

// ---------------------------------------------------------------------------
// 5. Semafor: maksymalnie jobMaxParallel jobów, nadmiar → 429
// ---------------------------------------------------------------------------

func TestSemaforDwaRunningTrzeci429(t *testing.T) {
	srv, jobs, _ := newTestAPI(t)
	useSeams(t, jobs)
	started := make(chan struct{}, 8)
	stop := make(chan struct{})
	t.Cleanup(func() { close(stop) })
	searchArmstrong = blockSearch(started, stop)

	url := srv.URL + z02URL
	ids := make([]string, 0, jobMaxParallel)
	for i := 0; i < jobMaxParallel; i++ {
		code, body := do(t, "POST", url, map[string]any{"limit": int64(1e6)})
		if code != http.StatusAccepted {
			t.Fatalf("POST #%d: kod %d, spodziewany 202", i+1, code)
		}
		ids = append(ids, jobID(t, body))
		<-started // potwierdzenie, że job realnie trwa
	}

	code, body := do(t, "POST", url, map[string]any{"limit": int64(1e6)})
	if code != http.StatusTooManyRequests {
		t.Fatalf("POST #3: kod %d, spodziewany 429", code)
	}
	if !strings.Contains(errMsg(body), "zajęte") {
		t.Fatalf("komunikat 429 po polsku oczekiwany, jest: %q", errMsg(body))
	}
	if n := jobs.count(); n != jobMaxParallel {
		t.Fatalf("store: %d jobów, spodziewany %d (429 nie tworzy joba)", n, jobMaxParallel)
	}

	// Zwolnienie przy zakończeniu: anulowane joby oddają miejsca semafora.
	for _, id := range ids {
		if code, _ = do(t, "DELETE", url+"/"+id, nil); code != http.StatusOK {
			t.Fatalf("DELETE %s: kod %d", id, code)
		}
	}
	for _, id := range ids {
		waitStatus(t, url+"/"+id, "cancelled")
	}
	body = waitPostAccepted(t, url)
	if code, _ = do(t, "DELETE", url+"/"+jobID(t, body), nil); code != http.StatusOK {
		t.Fatalf("DELETE po zwolnieniu: kod %d", code)
	}
}

// ---------------------------------------------------------------------------
// 2. Czystka jobStore: TTL od zakończenia + sufit liczby → GET 404
// ---------------------------------------------------------------------------

func TestCzystkaTTLPo15Minutach404(t *testing.T) {
	srv, jobs, clock := newTestAPI(t)
	useSeams(t, jobs)
	searchArmstrong = okSearch

	code, body := do(t, "POST", srv.URL+z02URL, map[string]any{"limit": int64(1e6)})
	if code != http.StatusAccepted {
		t.Fatalf("POST: kod %d", code)
	}
	url := srv.URL + z02URL + "/" + jobID(t, body)
	waitStatus(t, url, "done")

	// Zaraz po zakończeniu: 200.
	if code, _ = do(t, "GET", url, nil); code != http.StatusOK {
		t.Fatalf("GET zaraz po done: kod %d", code)
	}
	// Po_ttl: czystka przy dostępie usuwa joba → 404.
	clock.advance(jobTTL + time.Minute)
	if code, _ = do(t, "GET", url, nil); code != http.StatusNotFound {
		t.Fatalf("GET po TTL: kod %d, spodziewany 404", code)
	}
	if code, _ = do(t, "DELETE", url, nil); code != http.StatusNotFound {
		t.Fatalf("DELETE po TTL: kod %d, spodziewany 404", code)
	}
	if n := jobs.count(); n != 0 {
		t.Fatalf("store: %d jobów, spodziewany 0 po czystce TTL", n)
	}
}

func TestCzystkaCap500JobowDo200(t *testing.T) {
	srv, jobs, _ := newTestAPI(t)
	useSeams(t, jobs)
	searchArmstrong = okSearch

	first := ""
	for i := 0; i < 500; i++ {
		code, body := do(t, "POST", srv.URL+z02URL, map[string]any{"limit": int64(1e6)})
		if code != http.StatusAccepted {
			t.Fatalf("POST #%d: kod %d (%v)", i, code, errMsg(body))
		}
		if i == 0 {
			first = jobID(t, body)
		}
	}
	// Joby kończą się w gorutynach; czystka usuwa tylko zakończone.
	deadline := time.Now().Add(2 * time.Second)
	for jobs.count() > jobMaxCount {
		jobs.sweep()
		if time.Now().After(deadline) {
			break
		}
		time.Sleep(2 * time.Millisecond)
	}
	if n := jobs.count(); n > jobMaxCount {
		t.Fatalf("store: %d jobów, sufit to %d", n, jobMaxCount)
	}
	// Najstarszy job wyleciał → GET 404.
	if code, _ := do(t, "GET", srv.URL+z02URL+"/"+first, nil); code != http.StatusNotFound {
		t.Fatalf("GET najstarszego po czystce: kod %d, spodziewany 404", code)
	}
}

func TestCzystkaNieUsunAktywnegoJoba(t *testing.T) {
	srv, jobs, clock := newTestAPI(t)
	useSeams(t, jobs)
	started := make(chan struct{}, 4)
	stop := make(chan struct{})
	t.Cleanup(func() { close(stop) })
	searchArmstrong = blockSearch(started, stop)

	code, body := do(t, "POST", srv.URL+z02URL, map[string]any{"limit": int64(1e6)})
	if code != http.StatusAccepted {
		t.Fatalf("POST: kod %d", code)
	}
	url := srv.URL + z02URL + "/" + jobID(t, body)
	<-started

	// Aktywny job przeżywa czystkę (nie ma daty zakończenia).
	clock.advance(24 * time.Hour)
	jobs.sweep()
	if code, _ = do(t, "GET", url, nil); code != http.StatusOK {
		t.Fatalf("GET aktywnego po czystce: kod %d, spodziewany 200", code)
	}
}

// ---------------------------------------------------------------------------
// 1. Równoległe POST + DELETE (pod -race: wyścig na polu cancel)
// ---------------------------------------------------------------------------

func TestRownoleglePostDeleteBezWyscigu(t *testing.T) {
	srv, jobs, _ := newTestAPI(t)
	useSeams(t, jobs)
	started := make(chan struct{}, 128)
	stop := make(chan struct{})
	t.Cleanup(func() { close(stop) })
	searchArmstrong = blockSearch(started, stop)

	const n = 24
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []string
	)
	fail := func(format string, args ...any) {
		mu.Lock()
		errs = append(errs, fmt.Sprintf(format, args...))
		mu.Unlock()
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			code, body, err := call(http.MethodPost, srv.URL+z02URL, map[string]any{"limit": int64(1e6)})
			if err != nil {
				fail("POST #%d: %v", i, err)
				return
			}
			if code == http.StatusTooManyRequests {
				return // miejsce zajęte — poprawna odpowiedź
			}
			if code != http.StatusAccepted {
				fail("POST #%d: kod %d, spodziewany 202/429", i, code)
				return
			}
			id, _ := body["jobId"].(string)
			// Krótka sekwencja GET/DELETE/GET na tym samym jobie — to ściga
			// zapis cancel (gorutyna joba) z odczytem w DELETE.
			for k := 0; k < 3; k++ {
				if _, _, err := call(http.MethodGet, srv.URL+z02URL+"/"+id, nil); err != nil {
					fail("GET %s: %v", id, err)
				}
			}
			dcode, _, err := call(http.MethodDelete, srv.URL+z02URL+"/"+id, nil)
			if err != nil {
				fail("DELETE %s: %v", id, err)
			} else if dcode != http.StatusOK && dcode != http.StatusConflict {
				fail("DELETE %s: kod %d, spodziewany 200/409", id, dcode)
			}
			if _, _, err := call(http.MethodGet, srv.URL+z02URL+"/"+id, nil); err != nil {
				fail("GET po DELETE %s: %v", id, err)
			}
		}(i)
	}
	wg.Wait()

	if len(errs) > 0 {
		t.Fatalf("błędy równoległego POST/DELETE:\n  %s", strings.Join(errs, "\n  "))
	}
}
