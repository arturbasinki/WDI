// Package z08 generuje zagadki sudoku z dokładnie jednym rozwiązaniem
// i maksymalną liczbą pustych pól, a także rozwiązuje dowolną planszę.
//
// Algorytm generatora:
//
//  1. Losowe wypełnienie całej planszy 9×9 (backtracking z losową
//     kolejnością kandydatów) — gwarantuje poprawną, kompletną grę.
//  2. „Kopanie dołków”: dla pól w losowej kolejności usuwamy wartość
//     i sprawdzamy, czy plansza nadal ma dokładnie jedno rozwiązanie
//     (zliczacz rozwiązań przerywa po 2). Jeśli unikalność znika,
//     przywracamy wartość. Powtarzamy, aż żadne pole nie da się usunąć.
//     Dzięki temu zagadka ma jak najwięcej pustych pól przy zachowaniu
//     unikatowego rozwiązania.
//
// Zliczanie rozwiązań (i solver) działa metodą backtrackingu z heurystyką
// MRV (następne pole = to o najmniejszej liczbie kandydatów) oraz
// maskami bitowymi wiersza/kolumny/bloku — typową trudną planszę
// rozwiązuje w dziesiątkach milisekund.
package z08

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// Grid — plansza 81 pól (0 = puste pole).
type Grid [81]int

// Result — wynik generacji zagadki.
type Result struct {
	Puzzle   Grid  `json:"puzzle"`
	Solution Grid  `json:"solution"`
	Blanks   int   `json:"blanks"`
	Removed  int   `json:"removed"`
	TimeMs   int64 `json:"timeMs"`
}

// SolveResult — wynik rozwiązania podanej planszy.
type SolveResult struct {
	Solved   bool  `json:"solved"`
	Solution Grid  `json:"solution,omitempty"`
	Count    int   `json:"count"`  // liczba rozwiązań (zliczane do cap)
	Capped   bool  `json:"capped"` // czy zliczanie urwano po cap
	TimeMs   int64 `json:"timeMs"`
}

// ValidateGrid sprawdza poprawność formatu planszy (81 wartości 0..9).
func ValidateGrid(g []int) error {
	if len(g) != 81 {
		return fmt.Errorf("plansza musi mieć 81 pól, podano %d", len(g))
	}
	for i, v := range g {
		if v < 0 || v > 9 {
			return fmt.Errorf("pole %d: wartość musi być z 0..9, podano %d", i+1, v)
		}
	}
	return nil
}

// rowColBox — pomocnicze indeksy.
func boxOf(r, c int) int { return (r/3)*3 + c/3 }

// CountSolutions zlicza rozwiązania planszy (do cap).
// Używa mask bitowych: 1<<v oznacza, że wartość v (1..9) zajęta.
func CountSolutions(g Grid, cap int) int {
	rowMask := make([]uint16, 9)
	colMask := make([]uint16, 9)
	boxMask := make([]uint16, 9)
	for i, v := range g {
		if v == 0 {
			continue
		}
		r, c := i/9, i%9
		bit := uint16(1) << v
		if rowMask[r]&bit != 0 || colMask[c]&bit != 0 || boxMask[boxOf(r, c)]&bit != 0 {
			return 0 // sprzeczna plansza
		}
		rowMask[r] |= bit
		colMask[c] |= bit
		boxMask[boxOf(r, c)] |= bit
	}

	count := 0
	// MRV: szukaj pustego pola z najmniej kandydatami.
	var findMRV func() (int, uint16)
	findMRV = func() (int, uint16) {
		best, bestCands := -1, uint16(0)
		for i, v := range g {
			if v != 0 {
				continue
			}
			r, c := i/9, i%9
			used := rowMask[r] | colMask[c] | boxMask[boxOf(r, c)]
			cands := ^used & 0b1111111110 // bity 1..9
			n := popcount16(cands)
			if best == -1 || n < popcount16(bestCands) {
				best, bestCands = i, cands
				if n == 1 {
					break
				}
			}
		}
		return best, bestCands
	}

	var dfs func()
	dfs = func() {
		if count >= cap {
			return
		}
		i, cands := findMRV()
		if i == -1 {
			count++
			return
		}
		r, c := i/9, i%9
		b := boxOf(r, c)
		for v := 1; v <= 9; v++ {
			bit := uint16(1) << v
			if cands&bit == 0 {
				continue
			}
			g[i] = v
			rowMask[r] |= bit
			colMask[c] |= bit
			boxMask[b] |= bit
			dfs()
			g[i] = 0
			rowMask[r] &^= bit
			colMask[c] &^= bit
			boxMask[b] &^= bit
			if count >= cap {
				return
			}
		}
	}
	dfs()
	return count
}

func popcount16(x uint16) int {
	n := 0
	for x != 0 {
		x &= x - 1
		n++
	}
	return n
}

// fillFull losowo wypełnia całą planszę (backtracking).
func fillFull(g *Grid, rng *rand.Rand) bool {
	rowMask := make([]uint16, 9)
	colMask := make([]uint16, 9)
	boxMask := make([]uint16, 9)

	var dfs func(i int) bool
	dfs = func(i int) bool {
		if i == 81 {
			return true
		}
		r, c := i/9, i%9
		b := boxOf(r, c)
		used := rowMask[r] | colMask[c] | boxMask[b]
		cands := make([]int, 0, 9)
		for v := 1; v <= 9; v++ {
			if used&(1<<v) == 0 {
				cands = append(cands, v)
			}
		}
		rng.Shuffle(len(cands), func(a, b int) { cands[a], cands[b] = cands[b], cands[a] })
		for _, v := range cands {
			bit := uint16(1) << v
			(*g)[i] = v
			rowMask[r] |= bit
			colMask[c] |= bit
			boxMask[b] |= bit
			if dfs(i + 1) {
				return true
			}
			(*g)[i] = 0
			rowMask[r] &^= bit
			colMask[c] &^= bit
			boxMask[b] &^= bit
		}
		return false
	}
	return dfs(0)
}

// Generate tworzy zagadkę z unikatowym rozwiązaniem i maksymalną
// liczbą pustych pól (dla losowej pełnej planszy bazowej).
func Generate(seed int64) (Result, error) {
	start := time.Now()
	rng := rand.New(rand.NewPCG(uint64(seed), 0))

	var full Grid
	if !fillFull(&full, rng) {
		return Result{}, fmt.Errorf("nie udało się wygenerować pełnej planszy (błąd wewnętrzny)")
	}

	// Kolejność kopania: losowa permutacja pól.
	order := make([]int, 81)
	for i := range order {
		order[i] = i
	}
	rng.Shuffle(len(order), func(a, b int) { order[a], order[b] = order[b], order[a] })

	puzzle := full
	removed := 0
	for _, i := range order {
		v := puzzle[i]
		puzzle[i] = 0
		if CountSolutions(puzzle, 2) != 1 {
			puzzle[i] = v // unikalność naruszona — przywracamy
		} else {
			removed++
		}
	}

	blanks := 0
	for _, v := range puzzle {
		if v == 0 {
			blanks++
		}
	}
	return Result{
		Puzzle:   puzzle,
		Solution: full,
		Blanks:   blanks,
		Removed:  removed,
		TimeMs:   time.Since(start).Milliseconds(),
	}, nil
}

// Solve rozwiązuje podaną planszę i zlicza jej rozwiązania (do cap).
func Solve(g []int, cap int) (SolveResult, error) {
	start := time.Now()
	if err := ValidateGrid(g); err != nil {
		return SolveResult{}, err
	}
	if cap < 1 {
		return SolveResult{}, fmt.Errorf("cap musi być >= 1, podano %d", cap)
	}
	var grid Grid
	copy(grid[:], g)
	count := CountSolutions(grid, cap)
	res := SolveResult{Solved: count > 0, Count: count, Capped: count >= cap, TimeMs: time.Since(start).Milliseconds()}
	if count > 0 {
		// odzyskaj pierwsze rozwiązanie (zresetuj grid i zagraj do końca)
		var grid2 Grid
		copy(grid2[:], g)
		if !fillFirst(&grid2) {
			return SolveResult{}, fmt.Errorf("błąd wewnętrzny: licznik znalazł rozwiązania, a fillFirst nie")
		}
		res.Solution = grid2
	}
	return res, nil
}

// fillFirst wypełnia planię pierwszym rozwiązaniem (backtracking bez liczenia).
func fillFirst(g *Grid) bool {
	rowMask := make([]uint16, 9)
	colMask := make([]uint16, 9)
	boxMask := make([]uint16, 9)
	for i, v := range *g {
		if v == 0 {
			continue
		}
		r, c := i/9, i%9
		bit := uint16(1) << v
		rowMask[r] |= bit
		colMask[c] |= bit
		boxMask[boxOf(r, c)] |= bit
	}

	// MRV jak w CountSolutions, ale zwracamy po pierwszym komplecie.
	var solve func() bool
	solve = func() bool {
		best, bestCands := -1, uint16(0)
		for i, v := range *g {
			if v != 0 {
				continue
			}
			r, c := i/9, i%9
			used := rowMask[r] | colMask[c] | boxMask[boxOf(r, c)]
			cands := ^used & 0b1111111110
			n := popcount16(cands)
			if n == 0 {
				return false
			}
			if best == -1 || n < popcount16(bestCands) {
				best, bestCands = i, cands
				if n == 1 {
					break
				}
			}
		}
		if best == -1 {
			return true // pełna plansza
		}
		r, c := best/9, best%9
		b := boxOf(r, c)
		for v := 1; v <= 9; v++ {
			bit := uint16(1) << v
			if bestCands&bit == 0 {
				continue
			}
			(*g)[best] = v
			rowMask[r] |= bit
			colMask[c] |= bit
			boxMask[b] |= bit
			if solve() {
				return true
			}
			(*g)[best] = 0
			rowMask[r] &^= bit
			colMask[c] &^= bit
			boxMask[b] &^= bit
		}
		return false
	}
	return solve()
}
