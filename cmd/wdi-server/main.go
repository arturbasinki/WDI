// Command wdi-server uruchamia jednorodną stronę WWW z ośmioma zadaniami
// WDI napisanymi w Go: statyczny frontend (wbudowany w binarię) oraz
// API JSON, w którym długie obliczenia (szukanie liczb Armstronga,
// benchmark sita, atak na Gronsfelda) działają jako joby w tle z
// postęgiem odpytywanym metodą GET.
//
// Uruchomienie:
//
//	go run ./cmd/wdi-server            # domyślnie http://localhost:8080
//	PORT=9000 go run ./cmd/wdi-server  # inny port
package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"wdi/internal/z01"
	"wdi/internal/z02"
	"wdi/internal/z03"
	"wdi/internal/z04"
	"wdi/internal/z05"
	"wdi/internal/z06"
	"wdi/internal/z07"
	"wdi/internal/z08"
)

//go:embed all:web
var webFS embed.FS

// ---------------------------------------------------------------------------
// Rejestr jobów (długie obliczenia z postęgiem)
// ---------------------------------------------------------------------------

// Job — stan zadania w tle.
//
// WSZYSTKIE pola (także cancel i endedAt) są chronione mutexem: handler POST
// zapisuje cancel, DELETE oraz zamknięcie serwera czytają je równolegle — bez
// blokady to klasyczny data race (wykrywany przez -race).
type Job struct {
	ID string

	mu     sync.Mutex
	status string // running | done | error | cancelled
	err    string
	// progress — ostatni stan postępu (dowolny typ JSON).
	progress any
	// result — wynik końcowy (dowolny typ JSON).
	result any
	// cancel — funkcja anulowania (wpisywana w handlerze POST, czytana w DELETE).
	cancel context.CancelFunc
	// endedAt — moment zakończenia joba (zero, gdy wciąż trwa); napędza TTL.
	endedAt time.Time
}

// setCancel — wpisuje funkcję anulowania pod mutexem.
func (j *Job) setCancel(c context.CancelFunc) {
	j.mu.Lock()
	j.cancel = c
	j.mu.Unlock()
}

// tryCancel — anuluje joba, jeśli wciąż trwa. false oznacza „już zakończony”
// (wtedy HTTP oddaje 409, a nie „cancelling”).
func (j *Job) tryCancel() bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.status != "running" || j.cancel == nil {
		return false
	}
	j.cancel()
	return true
}

func (j *Job) snapshot() map[string]any {
	j.mu.Lock()
	defer j.mu.Unlock()
	m := map[string]any{"jobId": j.ID, "status": j.status}
	if j.err != "" {
		m["error"] = j.err
	}
	if j.progress != nil {
		m["progress"] = j.progress
	}
	if j.result != nil {
		m["result"] = j.result
	}
	return m
}

func (j *Job) setProgress(p any) {
	j.mu.Lock()
	j.progress = p
	j.mu.Unlock()
}

// finish — kończy joba: status i endedAt ustawiane są w jednym krytycznym
// sectionie, żeby czystka TTL nigdy nie widziała „końca” bez daty.
func (j *Job) finish(now time.Time, result any, err error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.status != "running" {
		return
	}
	if err != nil {
		j.status = "error"
		j.err = err.Error()
	} else {
		j.status = "done"
		j.result = result
	}
	j.endedAt = now
}

// markCancelled — status „cancelled”, wyłącznie gdy job wciąż trwa (job, który
// zdążył się zakończyć, zachowuje swój stan).
func (j *Job) markCancelled(now time.Time) {
	j.mu.Lock()
	if j.status == "running" {
		j.status = "cancelled"
		j.endedAt = now
	}
	j.mu.Unlock()
}

// Polityka przechowalni jobów.
const (
	// jobTTL — jak długo trzymamy zakończone joby (pamięć: ~0,4 MB na job).
	jobTTL = 15 * time.Minute
	// jobMaxCount — twardy sufit liczby jobów w pamięci.
	jobMaxCount = 200
	// jobMaxParallel — ile obliczeń może trwać równolegle (semafor).
	jobMaxParallel = 2
	// z02MaxLimitAPI — twardy limit limitu w warstwie HTTP, spójny z selektorem
	// w UI (biblioteka z02 przyjmuje 10¹³, wystawiamy tylko 10¹²).
	z02MaxLimitAPI = int64(1e12)
)

// jobStore — rejestr jobów z leniwą czystką (TTL + sufit liczby) i semaforem
// jednoczesnych obliczeń. Czystka wołana jest wyłącznie przy dostępie
// (create/get) — żadnych timerów globalnych.
type jobStore struct {
	mu    sync.Mutex
	seq   int
	jobs  map[string]*Job
	order []string // kolejność utworzenia — do wypychania najstarszych
	ttl   time.Duration
	max   int
	// sem — semafor jednoczesnych obliczeń (bufor o pojemności jobMaxParallel).
	sem chan struct{}
	// now — zegar wstrzykiwany: testy symulują upływ czasu bez spania.
	now func() time.Time
}

func newJobStore() *jobStore {
	return &jobStore{
		jobs: map[string]*Job{},
		ttl:  jobTTL,
		max:  jobMaxCount,
		sem:  make(chan struct{}, jobMaxParallel),
		now:  time.Now,
	}
}

// removeLocked — usuwa wskazane id z mapy i z kolejności. Wywołująca trzyma s.mu.
func (s *jobStore) removeLocked(ids []string) {
	if len(ids) == 0 {
		return
	}
	del := make(map[string]bool, len(ids))
	for _, id := range ids {
		del[id] = true
		delete(s.jobs, id)
	}
	keep := s.order[:0]
	for _, id := range s.order {
		if !del[id] {
			keep = append(keep, id)
		}
	}
	s.order = keep
}

// endedOf — data zakończenia joba (zero = wciąż trwa). Chwila trzymania
// mutexu samego joba, bez trzymania go podczas iteracji po mapie.
func endedOf(j *Job) time.Time {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.endedAt
}

// sweepLocked — leniwa czystka: najpierw usuwa zakończone joby po upływie TTL,
// a gdy nadal jest ich za dużo (powyżej s.max) — najstarsze zakończone, w
// kolejności utworzenia. Aktywne joby nigdy nie są usuwane.
// Wywołująca musi trzymać s.mu.
func (s *jobStore) sweepLocked() {
	now := s.now()
	var drop []string
	for _, id := range s.order {
		j := s.jobs[id]
		if j == nil {
			drop = append(drop, id)
			continue
		}
		if ended := endedOf(j); !ended.IsZero() && now.Sub(ended) > s.ttl {
			drop = append(drop, id)
		}
	}
	s.removeLocked(drop)

	if len(s.jobs) <= s.max {
		return
	}
	drop = drop[:0]
	for _, id := range s.order {
		if len(s.jobs)-len(drop) <= s.max {
			break
		}
		j := s.jobs[id]
		if j == nil {
			drop = append(drop, id)
			continue
		}
		if !endedOf(j).IsZero() {
			drop = append(drop, id)
		}
	}
	s.removeLocked(drop)
}

// sweep — czystka na żądanie (dostępna w testach; w produkcji wołana leniwie).
func (s *jobStore) sweep() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sweepLocked()
}

func (s *jobStore) create() (*Job, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sweepLocked()
	s.seq++
	id := fmt.Sprintf("j%d-%d", time.Now().UnixMilli(), s.seq)
	j := &Job{ID: id, status: "running"}
	s.jobs[id] = j
	s.order = append(s.order, id)
	return j, id
}

func (s *jobStore) get(id string) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sweepLocked()
	j, ok := s.jobs[id]
	return j, ok
}

// count — liczba trzymana jobów (testy progu pojemności).
func (s *jobStore) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.jobs)
}

// tryAcquire — zajmuje jedno miejsce w semaforze; false = brak wolnych (429).
func (s *jobStore) tryAcquire() bool {
	select {
	case s.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

// release — zwalnia miejsce w semaforze (wołane po ustawieniu statusu końcowego).
func (s *jobStore) release() {
	select {
	case <-s.sem:
	default:
	}
}

// cancelAll — anuluje wszystkie aktywne joby (przy zamknięciu serwera).
func (s *jobStore) cancelAll() {
	s.mu.Lock()
	live := make([]*Job, 0, len(s.order))
	for _, id := range s.order {
		if j := s.jobs[id]; j != nil {
			live = append(live, j)
		}
	}
	s.mu.Unlock()
	for _, j := range live {
		j.tryCancel()
	}
}

// startJob — start joba w tle: najpierw semafor (429, gdy brak miejsc), potem
// rejestracja i gorutyna. Miejsce w semaforze zwalniamy dopiero PO ustawieniu
// statusu końcowego, żeby „429” nie przeplatało się z jobem, który już trwa.
// run zwraca wynik/błąd — obsługuje panic, anulowanie i czystkę TTL.
func startJob(jobs *jobStore, w http.ResponseWriter, run func(ctx context.Context, j *Job) (any, error)) {
	if !jobs.tryAcquire() {
		httpErr(w, http.StatusTooManyRequests,
			fmt.Sprintf("zajęte: trwają już %d obliczenia równolegle — spróbuj za chwilę", jobMaxParallel))
		return
	}
	j, id := jobs.create()
	// Kontekst odłączony od żądania HTTP: job trwa dalej po zamknięciu
	// połączenia POST; anuluje go DELETE albo zamknięcie serwera.
	ctx, cancel := context.WithCancel(context.Background())
	j.setCancel(cancel)
	go func() {
		defer func() {
			// panic w tej gorutynie zabiłby cały proces (leży poza
			// odzyskiwaniem net/http) — zamieniamy go w stan „error”.
			if r := recover(); r != nil {
				j.finish(jobs.now(), nil, fmt.Errorf("panic w obliczeniach: %v", r))
			}
			jobs.release()
		}()
		res, err := run(ctx, j)
		if errors.Is(ctx.Err(), context.Canceled) {
			j.markCancelled(jobs.now())
		} else {
			j.finish(jobs.now(), res, err)
		}
	}()
	writeJSON(w, http.StatusAccepted, map[string]string{"jobId": id})
}

// Szwy testowe: warstwa HTTP woła obliczenia przez te funkcje, dzięki czemu
// testy httptest symulują długie joby bez realnych obliczeń (bez portów live).
var (
	searchArmstrong = z02.Search
	runBenchmark    = z03.Benchmark
	breakGronsfeld  = z05.Break
)

// ---------------------------------------------------------------------------
// Pomocnicze funkcje HTTP
// ---------------------------------------------------------------------------

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode: %v", err)
	}
}

func httpErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func decodeBody(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 4<<20) // 4 MB
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(v); err != nil {
		httpErr(w, http.StatusBadRequest, "nieprawidłowe dane JSON: "+err.Error())
		return false
	}
	return true
}

// ---------------------------------------------------------------------------
// Handler API
// ---------------------------------------------------------------------------

// jobGet / jobDelete — jedna obsługa GET i DELETE dla wszystkich trzech
// endpointów jobowych (z02, z03, z05).
//
// GET usuniętego joba (TTL / sufit liczby) → 404.
// DELETE: aktywny → anuluj i 202 „cancelling”; zakończony → 409 z komunikatem
// po polsku (nie ma czego przerywać); nieznany → 404.
func jobGet(jobs *jobStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		j, ok := jobs.get(r.PathValue("id"))
		if !ok {
			httpErr(w, http.StatusNotFound, "nieznany job — wygasł albo nigdy nie istniał")
			return
		}
		writeJSON(w, http.StatusOK, j.snapshot())
	}
}

func jobDelete(jobs *jobStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		j, ok := jobs.get(r.PathValue("id"))
		if !ok {
			httpErr(w, http.StatusNotFound, "nieznany job — wygasł albo nigdy nie istniał")
			return
		}
		if !j.tryCancel() {
			httpErr(w, http.StatusConflict, "job już się zakończył — nie ma czego przerywać")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "cancelling"})
	}
}

// taskDef — wpis katalogu zadań.
//
// Rejestr taskCatalog jest JEDYNYM źródłem tytułów i opisów: frontend nie
// przechowuje własnej kopii (kafelki startowe i nagłówki podstron czytają je z
// GET /api/tasks), więc opis nie może się rozjechać między Go a JS.
type taskDef struct {
	ID    string `json:"id"`
	Num   string `json:"num"`
	Title string `json:"title"`
	Tag   string `json:"tag"`
}

// taskCatalog — spis zadań serwujący kafelki i nagłówki podstron.
var taskCatalog = []taskDef{
	{ID: "z01", Num: "01", Title: "Układ równań liniowych", Tag: "N równań, N niewiadomych — Gauss-Jordan, animacja redukcji"},
	{ID: "z02", Num: "02", Title: "Liczby Armstronga", Tag: "wyszukiwanie do 10¹², 12 gorutyn, postęgi w czasie rzeczywistym"},
	{ID: "z03", Num: "03", Title: "Sito Eratostenesa", Tag: "wycinek vs mapa vs bitmapa — benchmark 10³…10⁸"},
	{ID: "z04", Num: "04", Title: "Częstotliwość słów", Tag: "top 20 dla „Pana Tadeusza” — tekst wbudowany w binarię"},
	{ID: "z05", Num: "05", Title: "Szyfr Gronsfelda", Tag: "szyfrowanie, deszyfrowanie i atak słownikowy (26 tys. słów)"},
	{ID: "z06", Num: "06", Title: "Wykresy funkcji", Tag: "własny parser wyrażeń, detekcja asymptot pionowych"},
	{ID: "z07", Num: "07", Title: "Zagadki literowo-cyfrowe", Tag: "cryptarithm: wszystkie rozwiązania, m.in. ABC·BD = EFGAH"},
	// „Maksymalna liczba pustych pól” było sformułowaniem umownym: generator
	// kopiuje dołki monotonicznie i dla seed 1–8 oraz 42 daje 55–58 pustych pól,
	// a teoretyczne maksimum przy jednoznacznym rozwiązaniu jest wyższe (~64).
	{ID: "z08", Num: "08", Title: "Generator sudoku", Tag: "unikatowe rozwiązanie, wiele pustych pól (typowo 55–58; gwarantowana jednoznaczność)"},
}

func newAPI(jobs *jobStore) http.Handler {
	getJob := jobGet(jobs)
	delJob := jobDelete(jobs)
	mux := http.NewServeMux()

	// --- Spis zadań: jedno źródło dla kafelków startowych i nagłówków podstron ---
	mux.HandleFunc("GET /api/tasks", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"tasks": taskCatalog})
	})

	// --- z01: układ równań ---
	mux.HandleFunc("POST /api/z01/solve", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Text string `json:"text"`
		}
		if !decodeBody(w, r, &req) {
			return
		}
		A, err := z01.Parse(req.Text)
		if err != nil {
			httpErr(w, 400, err.Error())
			return
		}
		writeJSON(w, 200, z01.Solve(A))
	})

	// --- z02: liczby Armstronga ---
	mux.HandleFunc("POST /api/z02/search", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Limit int64 `json:"limit"`
		}
		if !decodeBody(w, r, &req) {
			return
		}
		if req.Limit == 0 {
			req.Limit = z02.DefaultLimit
		}
		// Walidacja synchroniczna: bezsensowny limit = 400 przy POST, bez
		// tworzenia joba i bez zajmowania miejsca w semaforze.
		if req.Limit < 1000 || req.Limit > z02MaxLimitAPI {
			httpErr(w, http.StatusBadRequest,
				fmt.Sprintf("limit musi być z przedziału 1000…%d (10¹²)", z02MaxLimitAPI))
			return
		}
		startJob(jobs, w, func(ctx context.Context, j *Job) (any, error) {
			return searchArmstrong(ctx, req.Limit, func(p z02.Progress) { j.setProgress(p) })
		})
	})

	mux.HandleFunc("GET /api/z02/search/{id}", getJob)
	mux.HandleFunc("DELETE /api/z02/search/{id}", delJob)

	mux.HandleFunc("GET /api/z02/known", func(w http.ResponseWriter, r *http.Request) {
		// Lista zweryfikowana lokalnie przeszukiwaniem wyczerpującym
		// (wartości do 10^9 pokryte testami jednostkowymi, całość — jobem
		// weryfikacyjnym; zakres podany w verifiedUpTo).
		writeJSON(w, 200, map[string]any{
			"verifiedUpTo": int64(1e12),
			"note":         "Lista zweryfikowana przeszukiwaniem wyczerpującym w Go (12 gorutyn).",
			"numbers":      knownArmstrong,
		})
	})

	// --- z03: sito Eratostenesa ---
	mux.HandleFunc("POST /api/z03/benchmark", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			NS []int64 `json:"ns"`
		}
		if !decodeBody(w, r, &req) {
			return
		}
		if len(req.NS) == 0 {
			req.NS = []int64{1e3, 1e4, 1e5, 1e6, 1e7, 1e8}
		}
		// Walidacja synchroniczna: złe n = 400 przy POST, nie 202 + późny błąd.
		for _, n := range req.NS {
			if n < 1000 || n > 100_000_000 {
				httpErr(w, http.StatusBadRequest,
					fmt.Sprintf("n musi być z przedziału 1000…100000000, podano %d", n))
				return
			}
		}
		startJob(jobs, w, func(ctx context.Context, j *Job) (any, error) {
			return runBenchmark(ctx, req.NS, func(p z03.Progress) { j.setProgress(p) })
		})
	})

	mux.HandleFunc("GET /api/z03/benchmark/{id}", getJob)
	mux.HandleFunc("DELETE /api/z03/benchmark/{id}", delJob)

	// --- z04: słowa ---
	// Domyślną liczbę wyników (DefaultTop) stosuje TUTOR handler — biblioteka
	// dla k <= 0 zwraca błąd, więc „20” nie pojawia się w odpowiedzi znikąd.
	mux.HandleFunc("GET /api/z04/top20", func(w http.ResponseWriter, r *http.Request) {
		n := z04.DefaultTop
		if s := r.URL.Query().Get("n"); s != "" {
			if v, err := strconv.Atoi(s); err == nil && v >= 1 && v <= 100 {
				n = v
			}
		}
		words, err := z04.TopWords(z04.PanTadeusz(), n)
		if err != nil {
			httpErr(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, 200, map[string]any{
			"source": "pan-tadeusz",
			"words":  words,
			"n":      n,
		})
	})

	mux.HandleFunc("POST /api/z04/top20", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Text string `json:"text"`
			N    int    `json:"n"`
		}
		if !decodeBody(w, r, &req) {
			return
		}
		if req.N < 1 || req.N > 100 {
			req.N = z04.DefaultTop
		}
		words, err := z04.TopWords(req.Text, req.N)
		if err != nil {
			httpErr(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, 200, map[string]any{
			"source": "wlasny-tekst",
			"words":  words,
			"n":      req.N,
		})
	})

	// --- z05: Gronsfeld ---
	mux.HandleFunc("POST /api/z05/crypt", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Mode string `json:"mode"` // enc | dec
			Text string `json:"text"`
			Key  string `json:"key"`
		}
		if !decodeBody(w, r, &req) {
			return
		}
		var (
			out     string
			skipped []string
			err     error
		)
		switch req.Mode {
		case "enc":
			out, skipped, err = z05.Encrypt(req.Text, req.Key)
		case "dec":
			out, skipped, err = z05.Decrypt(req.Text, req.Key)
		default:
			httpErr(w, 400, "mode musi być „enc” lub „dec”")
			return
		}
		if err != nil {
			httpErr(w, 400, err.Error())
			return
		}
		// „skipped” jest zwracane TYLKO gdy normalizacja coś usunęła: alfabet
		// szyfru to a–z + spacja, więc polskie litery, cyfry i interpunkcja
		// znikają z wyniku — klient musi to widzieć, a nie zgadywać.
		resp := map[string]any{"out": out}
		if len(skipped) > 0 {
			resp["skipped"] = skipped
		}
		writeJSON(w, 200, resp)
	})

	mux.HandleFunc("POST /api/z05/break", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Cipher  string `json:"cipher"`
			MaxLen  int    `json:"maxLen"`
			CapKeys int64  `json:"capKeys"`
		}
		if !decodeBody(w, r, &req) {
			return
		}
		if req.MaxLen == 0 {
			req.MaxLen = 6
		}
		if req.CapKeys == 0 {
			req.CapKeys = z05.BreakCapDefault
		}
		// Walidacja synchroniczna — 400 przy POST, bez tworzenia joba.
		if strings.TrimSpace(req.Cipher) == "" {
			httpErr(w, http.StatusBadRequest, "szyfrogram jest pusty")
			return
		}
		if req.MaxLen < 1 || req.MaxLen > z05.MaxKeyLen {
			httpErr(w, http.StatusBadRequest, fmt.Sprintf("maxLen musi być z przedziału 1…%d", z05.MaxKeyLen))
			return
		}
		if req.CapKeys < 1 || req.CapKeys > z05.BreakCapMax {
			httpErr(w, http.StatusBadRequest,
				fmt.Sprintf("capKeys musi być z przedziału 1…%d", z05.BreakCapMax))
			return
		}
		startJob(jobs, w, func(ctx context.Context, j *Job) (any, error) {
			return breakGronsfeld(ctx, req.Cipher, req.MaxLen, req.CapKeys,
				func(p z05.BreakProgress) { j.setProgress(p) })
		})
	})

	mux.HandleFunc("GET /api/z05/break/{id}", getJob)
	mux.HandleFunc("DELETE /api/z05/break/{id}", delJob)

	// --- z06: wykresy ---
	mux.HandleFunc("POST /api/z06/plot", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Expr    string  `json:"expr"`
			XMin    float64 `json:"xmin"`
			XMax    float64 `json:"xmax"`
			Samples int     `json:"samples"`
		}
		if !decodeBody(w, r, &req) {
			return
		}
		if len(req.Expr) > 2000 {
			httpErr(w, 400, "wyrażenie jest za długie (maksymalnie 2000 znaków)")
			return
		}
		if req.Samples == 0 {
			req.Samples = 1601
		}
		p, err := z06.PlotFun(req.Expr, req.XMin, req.XMax, req.Samples)
		if err != nil {
			httpErr(w, 400, err.Error())
			return
		}
		writeJSON(w, 200, p)
	})

	// --- z07: alfanumeryka ---
	// limit liczby ZWRACANYCH rozwiączeń jest parametrem żądania:
	// domyślnie 50, sufit 500 (przycinamy, nie odrzucamy — Count jest pełny).
	mux.HandleFunc("POST /api/z07/solve", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Equation string `json:"equation"`
			Limit    int    `json:"limit"`
		}
		if !decodeBody(w, r, &req) {
			return
		}
		limit := req.Limit
		if limit == 0 {
			limit = z07.DefaultSolutions
		}
		if limit < 0 {
			httpErr(w, http.StatusBadRequest, "limit musi być >= 0 (0 = wartość domyślna, -1 = bez limitu w bibliotece)")
			return
		}
		if limit > z07.MaxSolutions {
			limit = z07.MaxSolutions
		}
		res, err := z07.Solve(req.Equation, limit)
		if err != nil {
			httpErr(w, 400, err.Error())
			return
		}
		writeJSON(w, 200, res)
	})

	// --- z08: sudoku ---
	mux.HandleFunc("POST /api/z08/generate", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Seed int64 `json:"seed"`
		}
		if !decodeBody(w, r, &req) {
			return
		}
		if req.Seed == 0 {
			req.Seed = time.Now().UnixNano()
		}
		res, err := z08.Generate(req.Seed)
		if err != nil {
			httpErr(w, 500, err.Error())
			return
		}
		writeJSON(w, 200, res)
	})

	mux.HandleFunc("POST /api/z08/solve", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Grid []int `json:"grid"`
		}
		if !decodeBody(w, r, &req) {
			return
		}
		res, err := z08.Solve(req.Grid, 100)
		if err != nil {
			httpErr(w, 400, err.Error())
			return
		}
		writeJSON(w, 200, res)
	})

	return mux
}

// ---------------------------------------------------------------------------
// main
// ---------------------------------------------------------------------------

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	jobs := newJobStore()
	api := newAPI(jobs)

	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", api)
	mux.Handle("/", http.FileServer(http.FS(sub)))

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 2 * time.Minute, // z08 solve + z07 mogą dłużej
		IdleTimeout:  120 * time.Second,
	}
	log.Printf("WDI w Go — serwer WWW na http://localhost:%s", port)

	// Graceful shutdown: SIGTERM/SIGINT → anuluj aktywne joby → srv.Shutdown
	// (maksymalnie 5 s na dokończenie obsługiwanych żądań) → wyjście z kodem 0.
	// Dzięki temu `docker stop` kończy kontener od razu, bez cięcia po 10 s.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)
		<-ctx.Done()
		log.Printf("sygnał zamknięcia — anuluję aktywne joby…")
		jobs.cancelAll()
		sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(sctx); err != nil {
			log.Printf("shutdown: %v", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	<-shutdownDone
	log.Printf("zamknięto")
}

// knownArmstrong — liczby Armstronga zweryfikowane lokalnym przeszukiwaniem
// wyczerpującym do 10^12 (zadanie z02; pomiar: 31 min na 12 gorutynach).
// Wartości ≤10^9 pokrywają testy jednostkowe (internal/z02).
var knownArmstrong = []int64{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
	153, 370, 371, 407,
	1634,
	8208, 9474,
	54748, 92727, 93084,
	548834,
	1741725, 4210818, 9800817, 9926315,
	24678050, 24678051,
	88593477,
	146511208, 472335975, 534494836, 912985153,
	4679307774,
	32164049650, 32164049651, 40028394225, 42678290603,
	44708635679, 49388550606, 82693916578, 94204591914,
}
