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
	"strconv"
	"sync"
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
type Job struct {
	ID string

	mu     sync.Mutex
	status string // running | done | error | cancelled
	err    string
	// progress — ostatni stan postępu (dowolny typ JSON).
	progress any
	// result — wynik końcowy (dowolny typ JSON).
	result any

	cancel context.CancelFunc
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

func (j *Job) finish(result any, err error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if err != nil {
		j.status = "error"
		j.err = err.Error()
		return
	}
	j.status = "done"
	j.result = result
}

type jobStore struct {
	mu   sync.Mutex
	seq  int
	jobs map[string]*Job
}

func (s *jobStore) create() (*Job, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	id := fmt.Sprintf("j%d-%d", time.Now().UnixMilli(), s.seq)
	j := &Job{ID: id, status: "running"}
	s.jobs[id] = j
	return j, id
}

func (s *jobStore) get(id string) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[id]
	return j, ok
}

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

func newAPI(jobs *jobStore) http.Handler {
	mux := http.NewServeMux()

	// --- Spis zadań (dla kafelków) ---
	mux.HandleFunc("GET /api/tasks", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{
			"tasks": []map[string]string{
				{"id": "z01", "num": "01", "title": "Układ równań liniowych", "tag": "N równań, N niewiadomych — Gauss-Jordan, animacja redukcji"},
				{"id": "z02", "num": "02", "title": "Liczby Armstronga", "tag": "wyszukiwanie do 10¹², 12 gorutyn, postęgi w czasie rzeczywistym"},
				{"id": "z03", "num": "03", "title": "Sito Eratostenesa", "tag": "wycinek vs mapa vs bitmapa — benchmark 10³…10⁸"},
				{"id": "z04", "num": "04", "title": "Częstotliwość słów", "tag": "top 20 dla „Pana Tadeusza” — tekst wbudowany w binarię"},
				{"id": "z05", "num": "05", "title": "Szyfr Gronsfelda", "tag": "szyfrowanie, deszyfrowanie i atak słownikowy (26 tys. słów)"},
				{"id": "z06", "num": "06", "title": "Wykresy funkcji", "tag": "własny parser wyrażeń, detekcja asymptot pionowych"},
				{"id": "z07", "num": "07", "title": "Zagadki literowo-cyfrowe", "tag": "cryptarithm: wszystkie rozwiązania, m.in. ABC·BD = EFGAH"},
				{"id": "z08", "num": "08", "title": "Generator sudoku", "tag": "unikatowe rozwiązanie, maksymalna liczba pustych pól"},
			},
		})
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
		j, id := jobs.create()
		// Kontekst odłączony od żądania HTTP: job trwa dalej po zamknięciu
		// połączenia POST; anuluje go wyłącznie jawnie DELETE.
		ctx, cancel := context.WithCancel(context.Background())
		j.cancel = cancel
		go func() {
			defer func() {
				// panic w tej gorutynie zabiłby cały proces (leży poza
				// odzyskiwaniem net/http) — zamieniamy go w stan „error”.
				if r := recover(); r != nil {
					j.finish(nil, fmt.Errorf("panic w obliczeniach: %v", r))
				}
			}()
			res, err := z02.Search(ctx, req.Limit, func(p z02.Progress) { j.setProgress(p) })
			if errors.Is(ctx.Err(), context.Canceled) {
				j.mu.Lock()
				j.status = "cancelled"
				j.mu.Unlock()
				return
			}
			j.finish(res, err)
		}()
		writeJSON(w, 202, map[string]string{"jobId": id})
	})

	mux.HandleFunc("GET /api/z02/search/{id}", func(w http.ResponseWriter, r *http.Request) {
		j, ok := jobs.get(r.PathValue("id"))
		if !ok {
			httpErr(w, 404, "nieznany job")
			return
		}
		writeJSON(w, 200, j.snapshot())
	})

	mux.HandleFunc("DELETE /api/z02/search/{id}", func(w http.ResponseWriter, r *http.Request) {
		j, ok := jobs.get(r.PathValue("id"))
		if !ok {
			httpErr(w, 404, "nieznany job")
			return
		}
		if j.cancel != nil {
			j.cancel()
		}
		writeJSON(w, 200, map[string]string{"status": "cancelling"})
	})

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
		j, id := jobs.create()
		ctx, cancel := context.WithCancel(context.Background())
		j.cancel = cancel
		go func() {
			defer func() {
				if r := recover(); r != nil {
					j.finish(nil, fmt.Errorf("panic w obliczeniach: %v", r))
				}
			}()
			rows, err := z03.Benchmark(ctx, req.NS, func(p z03.Progress) { j.setProgress(p) })
			if errors.Is(ctx.Err(), context.Canceled) {
				j.mu.Lock()
				j.status = "cancelled"
				j.mu.Unlock()
				return
			}
			j.finish(rows, err)
		}()
		writeJSON(w, 202, map[string]string{"jobId": id})
	})

	mux.HandleFunc("GET /api/z03/benchmark/{id}", func(w http.ResponseWriter, r *http.Request) {
		j, ok := jobs.get(r.PathValue("id"))
		if !ok {
			httpErr(w, 404, "nieznany job")
			return
		}
		writeJSON(w, 200, j.snapshot())
	})

	mux.HandleFunc("DELETE /api/z03/benchmark/{id}", func(w http.ResponseWriter, r *http.Request) {
		j, ok := jobs.get(r.PathValue("id"))
		if !ok {
			httpErr(w, 404, "nieznany job")
			return
		}
		if j.cancel != nil {
			j.cancel()
		}
		writeJSON(w, 200, map[string]string{"status": "cancelling"})
	})

	// --- z04: słowa ---
	mux.HandleFunc("GET /api/z04/top20", func(w http.ResponseWriter, r *http.Request) {
		n := 20
		if s := r.URL.Query().Get("n"); s != "" {
			if v, err := strconv.Atoi(s); err == nil && v >= 1 && v <= 100 {
				n = v
			}
		}
		writeJSON(w, 200, map[string]any{
			"source": "pan-tadeusz",
			"words":  z04.TopWords(z04.PanTadeusz(), n),
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
			req.N = 20
		}
		writeJSON(w, 200, map[string]any{
			"source": "wlasny-tekst",
			"words":  z04.TopWords(req.Text, req.N),
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
			out string
			err error
		)
		switch req.Mode {
		case "enc":
			out, err = z05.Encrypt(req.Text, req.Key)
		case "dec":
			out, err = z05.Decrypt(req.Text, req.Key)
		default:
			httpErr(w, 400, "mode musi być „enc” lub „dec”")
			return
		}
		if err != nil {
			httpErr(w, 400, err.Error())
			return
		}
		writeJSON(w, 200, map[string]string{"out": out})
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
			req.CapKeys = 2_000_000
		}
		j, id := jobs.create()
		ctx, cancel := context.WithCancel(context.Background())
		j.cancel = cancel
		go func() {
			defer func() {
				if r := recover(); r != nil {
					j.finish(nil, fmt.Errorf("panic w obliczeniach: %v", r))
				}
			}()
			res, err := z05.Break(ctx, req.Cipher, req.MaxLen, req.CapKeys, func(p z05.BreakProgress) { j.setProgress(p) })
			if errors.Is(ctx.Err(), context.Canceled) {
				j.mu.Lock()
				j.status = "cancelled"
				j.mu.Unlock()
				return
			}
			j.finish(res, err)
		}()
		writeJSON(w, 202, map[string]string{"jobId": id})
	})

	mux.HandleFunc("GET /api/z05/break/{id}", func(w http.ResponseWriter, r *http.Request) {
		j, ok := jobs.get(r.PathValue("id"))
		if !ok {
			httpErr(w, 404, "nieznany job")
			return
		}
		writeJSON(w, 200, j.snapshot())
	})

	mux.HandleFunc("DELETE /api/z05/break/{id}", func(w http.ResponseWriter, r *http.Request) {
		j, ok := jobs.get(r.PathValue("id"))
		if !ok {
			httpErr(w, 404, "nieznany job")
			return
		}
		if j.cancel != nil {
			j.cancel()
		}
		writeJSON(w, 200, map[string]string{"status": "cancelling"})
	})

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
	mux.HandleFunc("POST /api/z07/solve", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Equation string `json:"equation"`
		}
		if !decodeBody(w, r, &req) {
			return
		}
		res, err := z07.Solve(req.Equation)
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
	jobs := &jobStore{jobs: map[string]*Job{}}
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
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
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
