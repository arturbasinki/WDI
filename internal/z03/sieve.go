// Package z03 implementuje metodę Sita Eratostenesa trzema sposobami
// i porównuje czasy działania dla rosnących wartości n.
//
// W zadaniu (w Pythonie) użyto trzech struktur danych: listy, słownika
// i tablicy numpy. W Go odpowiednikami są:
//
//  1. wycinek []bool  — „lista”: każdy element to flaga „pierwszości”,
//  2. map[int]bool    — „słownik”: klucz to liczba, wartość to flaga,
//  3. bitmapa zbitowa — „tablica (numpy)”: gęsty, wektoryzowalny
//     zbiór bitów (analog tablicy bool numpy); 8 flag w 1 bajcie.
//
// Algorytm w każdym przypadku jest ten sam:
//
//   - zaznacz wszystkie liczby z [2, n] jako kandydaty,
//   - dla p = 2, 3, ..., ⌊√n⌋: jeśli p jest kandydatem, wykreśl
//     wielokrotności p·p, p·p+p, ..., ≤ n,
//   - pozostałe kandydaci to liczby pierwsze.
//
// Wersje z przydomkiem Ctx sprawdzają anulowanie WEWNĄTRZ pomiaru
// (co ctxTick iteracji), więc benchmark przerywa w trakcie wolnej
// struktury, a nie dopiero po jej zakończeniu.
package z03

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"time"
)

// Row to jeden wiersz tabeli benchmarkowej.
type Row struct {
	Structure string  `json:"structure"`
	N         int64   `json:"n"`
	Count     int     `json:"count"`
	Seconds   float64 `json:"seconds"`
}

// Progress opisuje postęp benchmarku (done/total = wykonanych pomiarów).
type Progress struct {
	Done  int `json:"done"`
	Total int `json:"total"`
	Last  Row `json:"last"`
}

// ctxTick — co ile iteracji wewnętrznej pętli sprawdzamy anulowanie.
// 2²⁰ ≈ milion: koszt wywołania ctx.Err() rozłożony na milion iteracji
// jest mierzalny jako zero, a najwolniejszy pojedynczy pomiar (mapa,
// n = 10⁸ — 4…15 s) urywa się w ułamek sekundy po anulowaniu.
const ctxTick = 1 << 20

// canceled — true, jeśli kontekst został anulowany (sprawdzenie co ctxTick
// iteracji: niska maska ustawiona w całości, czyli raz na 2²⁰ kroków).
func canceled(ctx context.Context, i int) bool {
	return i&(ctxTick-1) == ctxTick-1 && ctx.Err() != nil
}

// SieveSlice wykonuje sito na wycinku []bool („lista”).
func SieveSlice(n int) []int {
	sieve, _ := SieveSliceCtx(context.Background(), n)
	return sieve
}

// SieveSliceCtx — jak SieveSlice, z anulowaniem w środku pomiaru.
// Sprawdzenie siedzi w KAŻDEJ pętli (także w inicjalizacji flag), bo przy
// n = 10⁸ samo ustawienie tablicy to ~0,1 s bez -race i ~7 s pod -race:
// gdyby ctx był badany dopiero w pętli wykreślania, anulowanie czekałoby
// do końca inicjalizacji.
func SieveSliceCtx(ctx context.Context, n int) ([]int, error) {
	sieve := make([]bool, n+1)
	for i := range sieve {
		if canceled(ctx, i) {
			return nil, ctx.Err()
		}
		sieve[i] = true
	}
	p := 2
	for p*p <= n {
		if canceled(ctx, p) {
			return nil, ctx.Err()
		}
		if sieve[p] {
			for i := p * p; i <= n; i += p {
				if canceled(ctx, i) {
					return nil, ctx.Err()
				}
				sieve[i] = false
			}
		}
		p++
	}
	out := make([]int, 0)
	for p := 2; p <= n; p++ {
		if canceled(ctx, p) {
			return nil, ctx.Err()
		}
		if sieve[p] {
			out = append(out, p)
		}
	}
	return out, nil
}

// SieveMap wykonuje sito na mapie map[int]bool („słownik”).
func SieveMap(n int) []int {
	sieve, _ := SieveMapCtx(context.Background(), n)
	return sieve
}

// SieveMapCtx — jak SieveMap, z anulowaniem w środku pomiaru.
// To ta struktura jest najwolniejsza (wstawienia do mapy dla 10⁸ kluczy),
// więc to głównie dla niej istnieje sprawdzanie wewnątrz pomiaru.
func SieveMapCtx(ctx context.Context, n int) ([]int, error) {
	sieve := make(map[int]bool, n+1)
	for i := 2; i <= n; i++ {
		if canceled(ctx, i) {
			return nil, ctx.Err()
		}
		sieve[i] = true
	}
	p := 2
	for p*p <= n {
		if canceled(ctx, p) {
			return nil, ctx.Err()
		}
		if sieve[p] {
			for i := p * p; i <= n; i += p {
				if canceled(ctx, i) {
					return nil, ctx.Err()
				}
				sieve[i] = false
			}
		}
		p++
	}
	out := make([]int, 0)
	for p, isP := range sieve {
		if isP {
			out = append(out, p)
		}
	}
	// Iteracja mapy jest niedeterministyczna — zwracamy listę rosnąco,
	// spójnie z SieveSlice/SieveBit (test TestSieveMapOrdered).
	sort.Ints(out)
	return out, nil
}

// SieveBit wykonuje sito na bitmapie zbitowej („tablica numpy”).
func SieveBit(n int) []int {
	sieve, _ := SieveBitCtx(context.Background(), n)
	return sieve
}

// SieveBitCtx — jak SieveBit, z anulowaniem w środku pomiaru.
func SieveBitCtx(ctx context.Context, n int) ([]int, error) {
	// sieve[i] = true, gdy bit i jest ustawiony (liczba i nie wykreślona).
	words := (n + 64) / 64
	bits := make([]uint64, words)
	for i := 2; i <= n; i++ {
		if canceled(ctx, i) {
			return nil, ctx.Err()
		}
		bits[i>>6] |= 1 << (i & 63)
	}
	for p := 2; p*p <= n; p++ {
		if canceled(ctx, p) {
			return nil, ctx.Err()
		}
		if bits[p>>6]&(1<<(p&63)) != 0 {
			for i := p * p; i <= n; i += p {
				if canceled(ctx, i) {
					return nil, ctx.Err()
				}
				bits[i>>6] &^= 1 << (i & 63)
			}
		}
	}
	out := make([]int, 0)
	for i := 2; i <= n; i++ {
		if canceled(ctx, i) {
			return nil, ctx.Err()
		}
		if bits[i>>6]&(1<<(i&63)) != 0 {
			out = append(out, i)
		}
	}
	return out, nil
}

// Benchmark mierzy czas i zwraca tabele wyników dla trzech struktur
// i podanych wartości n. onTick jest wywoływany po każdym pomiarze.
// Funkcje są uruchamiane sekwencyjnie (miary porównywalne), z GC
// włączonym — wyniki odpowiadają „realnym” czasom z zadania.
//
// Anulowanie kontekstu przerywa TAKŻE wewnątrz pojedynczego pomiaru
// (patrz Sieve*Ctx): zwracane są dotychczasowe wiersze oraz ctx.Err().
func Benchmark(ctx context.Context, ns []int64, onTick func(Progress)) ([]Row, error) {
	if len(ns) == 0 {
		return nil, fmt.Errorf("pusta lista wartości n")
	}
	for _, n := range ns {
		if n < 1000 || n > 1e8 {
			return nil, fmt.Errorf("n musi być z przedziału 1000..100000000, podano %d", n)
		}
	}
	if onTick == nil {
		onTick = func(Progress) {}
	}
	// Ustaw GC na stałą częstotliwość dla powtarzalnych pomiarów.
	runtime.GC()

	type fn struct {
		name string
		f    func(context.Context, int) ([]int, error)
	}
	fns := []fn{
		{"slice (lista)", SieveSliceCtx},
		{"map (słownik)", SieveMapCtx},
		{"bitmapa (tablica)", SieveBitCtx},
	}

	rows := []Row{}
	total := len(ns) * len(fns)
	done := 0
	for _, n := range ns {
		for _, fn := range fns {
			if ctx.Err() != nil {
				return rows, ctx.Err()
			}
			runtime.GC()
			start := time.Now()
			primes, err := fn.f(ctx, int(n))
			elapsed := time.Since(start).Seconds()
			if err != nil {
				// Przerwany ŚRODEK pomiaru: nie raportujemy wiersza z
				// uciętym wynikiem — zwracamy to, co zmierzone, i błąd.
				return rows, err
			}
			row := Row{Structure: fn.name, N: n, Count: len(primes), Seconds: elapsed}
			rows = append(rows, row)
			done++
			onTick(Progress{Done: done, Total: total, Last: row})
		}
	}
	return rows, nil
}
