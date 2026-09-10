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
package z03

import (
	"context"
	"fmt"
	"runtime"
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

// SieveSlice wykonuje sito na wycinku []bool („lista”).
func SieveSlice(n int) []int {
	sieve := make([]bool, n+1)
	for i := range sieve {
		sieve[i] = true
	}
	p := 2
	for p*p <= n {
		if sieve[p] {
			for i := p * p; i <= n; i += p {
				sieve[i] = false
			}
		}
		p++
	}
	out := make([]int, 0)
	for p := 2; p <= n; p++ {
		if sieve[p] {
			out = append(out, p)
		}
	}
	return out
}

// SieveMap wykonuje sito na mapie map[int]bool („słownik”).
func SieveMap(n int) []int {
	sieve := make(map[int]bool, n+1)
	for i := 2; i <= n; i++ {
		sieve[i] = true
	}
	p := 2
	for p*p <= n {
		if sieve[p] {
			for i := p * p; i <= n; i += p {
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
	return out
}

// SieveBit wykonuje sito na bitmapie zbitowej („tablica numpy”).
func SieveBit(n int) []int {
	// sieve[i] = true, gdy bit i jest ustawiony (liczba i nie wykreślona).
	words := (n + 64) / 64
	bits := make([]uint64, words)
	for i := 2; i <= n; i++ {
		bits[i>>6] |= 1 << (i & 63)
	}
	for p := 2; p*p <= n; p++ {
		if bits[p>>6]&(1<<(p&63)) != 0 {
			for i := p * p; i <= n; i += p {
				bits[i>>6] &^= 1 << (i & 63)
			}
		}
	}
	out := make([]int, 0)
	for i := 2; i <= n; i++ {
		if bits[i>>6]&(1<<(i&63)) != 0 {
			out = append(out, i)
		}
	}
	return out
}

// Benchmark mierzy czas i zwraca tabele wyników dla trzech struktur
// i podanych wartości n. onTick jest wywoływany po każdym pomiarze.
// Funkcje są uruchamiane sekwencyjnie (miary porównywalne), z GC
// włączonym — wyniki odpowiadają „realnym” czasom z zadania.
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
		f    func(int) []int
	}
	fns := []fn{
		{"slice (lista)", SieveSlice},
		{"map (słownik)", SieveMap},
		{"bitmapa (tablica)", SieveBit},
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
			primes := fn.f(int(n))
			elapsed := time.Since(start).Seconds()
			row := Row{Structure: fn.name, N: n, Count: len(primes), Seconds: elapsed}
			rows = append(rows, row)
			done++
			onTick(Progress{Done: done, Total: total, Last: row})
		}
	}
	return rows, nil
}
