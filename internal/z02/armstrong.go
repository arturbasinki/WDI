// Package z02 znajduje liczby Armstronga (narcystyczne).
//
// Liczba Armstronga to N-cyfrowa liczba naturalna, która jest równa sumie
// swoich cyfr podniesionych do potęgi N, gdzie N to liczba jej cyfr.
// Przykład: 153 = 1^3 + 5^3 + 3^3.
//
// Program przeszukuje przedział [0, limit] metodą wyczerpującą, z
// wydajnym, przyrostowym obliczaniem sumy potęg cyfr:
//
//   - dla bloku liczb o tej samej liczbie cyfr d potęgi i^d (i = 0..9)
//     są preobliczone,
//   - przechodząc od n do n+1, gdy ostatnia cyfra b < 9, suma zmienia się
//     o stały przyrost (b+1)^d − b^d (odczyt z tablicy),
//   - gdy b = 9 (przeniesienie), suma jest przeliczana dla n/10+1 —
//     zdarza się to co 10 liczb, więc średni koszt pozostaje stały.
//
// Przeszukiwanie jest równoległe: przedział dzielimy na bloki po
// 4 000 000 liczb, które równolegle przetwarza pula gorutyn
// (runtime.GOMAXPROCS(0) workerów).
//
// Dzięki temu pojedynczy program potrafi przeanalizować 10^10 liczb
// w kilkanaście sekund na 12-rdzeniowej maszynie (Python z zadania
// przeszukiwał do 10^8).
package z02

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// Progress jest stanem postępu przeszukiwania.
type Progress struct {
	Processed int64   `json:"processed"`
	Limit     int64   `json:"limit"`
	Found     []int64 `json:"found"`
	Done      bool    `json:"done"`
}

const (
	blockSize = 4_000_000
	// maxLimit — twardy limit przeszukiwania (powyżej nie ma sensu).
	maxLimit = int64(1e13)
	// DefaultLimit — limit domyślny.
	DefaultLimit = int64(1e10)
)

// maxDigits — maksymalna liczba cyfr obsługiwana przez tablice potęg.
// 9^19 ≈ 1,35·10^18 mieści się w int64 (MaxInt64 ≈ 9,22·10^18).
const maxDigits = 19

// pow10[i] = 10^i. Ostatnia obsługiwana potęga to 10^18 — 10^19
// przepełniłaby int64.
var pow10 = func() []int64 {
	p := make([]int64, 20)
	p[0] = 1
	for i := 1; i < 19; i++ {
		p[i] = p[i-1] * 10
	}
	return p
}()

// exactPow — dokładna potęga w int64 (do 19. potęgi).
func exactPow(base, exp int64) int64 {
	r := int64(1)
	for i := int64(0); i < exp; i++ {
		r *= base
	}
	return r
}

// digitPow[d][i] = i^d dla d = 1..19, i = 0..9
var digitPow = func() [][]int64 {
	t := make([][]int64, maxDigits+1)
	for d := 1; d <= maxDigits; d++ {
		row := make([]int64, 10)
		for i := 0; i < 10; i++ {
			row[i] = exactPow(int64(i), int64(d))
		}
		t[d] = row
	}
	return t
}()

// digitPowerSum zwraca sumę potęg cyfr liczby n (potęga = digits).
func digitPowerSum(n int64, digits int) int64 {
	if n == 0 {
		return 0
	}
	s := int64(0)
	pow := digitPow[digits]
	for n > 0 {
		s += pow[n%10]
		n /= 10
	}
	return s
}

// chunkSearch przeszukuje przedział [start, end) i dopisuje znalezione
// liczby Armstronga do out (chronione mu). processed — licznik postępu.
//
// Optymalizacja: iterujemy „blokami” dziesięciu kolejnych liczb
// (o tym samej części starszej). Dla bloku o części starszej upper suma
// potęg cyfr to sUpper + k^d, gdzie k to cyfra jedności — sUpper liczymy
// tylko raz na blok (zamiast co 10 liczb w wersji przyrostowej, tu nie ma
// w ogóle drogich przeniesień cyfrowych). Dzięki temu na jedną liczbę
// przypada: jedno dodanie, jedno porównanie i ułamek jednej dzielenia.
func chunkSearch(start, end int64, digits int, processed *atomic.Int64, out *[]int64, mu *sync.Mutex) {
	pow := digitPow[digits]
	n := start
	for n < end {
		upper := n / 10
		first := n % 10
		sUpper := digitPowerSum(upper, digits)
		var cnt int64
		for k := first; k <= 9; k++ {
			nn := upper*10 + k
			if nn >= end {
				break
			}
			cnt++
			if sUpper+pow[k] == nn {
				mu.Lock()
				*out = append(*out, nn)
				mu.Unlock()
			}
		}
		processed.Add(cnt)
		n = (upper + 1) * 10
	}
}

// Search przeszukuje [0, limit] w poszukiwaniu liczb Armstronga.
// Wywoływany jest w osobnej gorutynie serwera (job); onTick jest
// wywoływane ~10 razy na sekundę z aktualnym stanem.
func Search(ctx context.Context, limit int64, onTick func(Progress)) (Progress, error) {
	if limit < 1000 {
		return Progress{}, fmt.Errorf("limit musi być >= 1000")
	}
	if limit > maxLimit {
		return Progress{}, fmt.Errorf("limit nie może przekraczać %d", maxLimit)
	}
	if onTick == nil {
		onTick = func(Progress) {}
	}

	var (
		processed atomic.Int64
		mu        sync.Mutex
		found     []int64
	)

	// Zbierz bloki: dla każdej liczby cyfr d przedział [10^(d-1), 10^d)
	// (dla d=1 zaczynamy od 0, bo 0 jest liczbą Armstronga).
	type chunk struct {
		start, end int64
		digits     int
	}
	chunks := []chunk{}
	for d := 1; d <= maxDigits; d++ {
		lo := pow10[d-1]
		if d == 1 {
			lo = 0
		}
		hi := pow10[d] // exclusive
		if hi > limit+1 {
			hi = limit + 1
		}
		if lo >= limit+1 {
			break
		}
		for s := lo; s < hi; s += blockSize {
			e := s + blockSize
			if e > hi {
				e = hi
			}
			chunks = append(chunks, chunk{s, e, d})
		}
	}

	workers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	jobs := make(chan chunk, workers)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for c := range jobs {
				if ctx.Err() != nil {
					return
				}
				chunkSearch(c.start, c.end, c.digits, &processed, &found, &mu)
			}
		}()
	}
	go func() {
		// close NIGDY nie może zostać pominięte: bez niego robotnicy
		// czekają na `range jobs` w nieskończoność po anulowaniu.
		defer close(jobs)
		for _, c := range chunks {
			select {
			case jobs <- c:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Ticker postępu (~10 Hz). Musi zostać zatrzymany i dołączony (join)
	// przed fazą końcową — bez tego finalny sort kopia/found współzawodniczy
	// z odczytem tablicy przez tickera (wykryte przez `go test -race`).
	tickStop := make(chan struct{})
	tickDone := make(chan struct{})
	go func() {
		defer close(tickDone)
		t := time.NewTicker(100 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-t.C:
			case <-tickStop:
				return
			case <-ctx.Done():
				return
			}
			mu.Lock()
			f := append([]int64(nil), found...)
			mu.Unlock()
			sort.Slice(f, func(i, j int) bool { return f[i] < f[j] })
			onTick(Progress{Processed: processed.Load(), Limit: limit, Found: f})
		}
	}()

	wg.Wait()
	// Zatrzymaj tickera i poczekaj na jego zakończenie: od tej chwili
	// `found` mutuje wyłącznie ta gorutyna.
	close(tickStop)
	<-tickDone
	// Przerwane przeszukiwanie: zwracamy dotychczasowy stan BŁĘDEM
	// contextu (kontrakt jak w z05.Break) — nie udajemy sukcesu.
	if cerr := ctx.Err(); cerr != nil {
		mu.Lock()
		f := append([]int64(nil), found...)
		mu.Unlock()
		sort.Slice(f, func(i, j int) bool { return f[i] < f[j] })
		return Progress{Processed: processed.Load(), Limit: limit, Found: f}, cerr
	}
	// Ostateczny stan: uzupełnij licznik do limitu (przybliżenie).
	processed.Store(limit + 1)
	mu.Lock()
	f := append([]int64(nil), found...) // kopia: zwracany stan nie aliasuje wewnętrznego found
	mu.Unlock()
	sort.Slice(f, func(i, j int) bool { return f[i] < f[j] })
	res := Progress{Processed: limit + 1, Limit: limit, Found: f, Done: true}
	onTick(res)
	return res, nil
}
