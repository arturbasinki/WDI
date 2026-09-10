// Package z07 rozwiązuje zagadki typu „alfanumeryka” (cryptarithm):
// w równaniu litery A–J zastępują cyfry. Przykład:
//
//	ABC * BD = EFGAH
//
// rozwiązuje się np. jako 623 * 29 = 18067 (A=6, B=2, C=3, D=9,
// E=1, F=8, G=0, H=7).
//
// Program wczytuje równanie w postaci <liczba> <operator> <liczba> =
// <liczba>, gdzie operatorem może być + - * /. Zasady:
//
//   - każda litera odpowiada dokładnie jednej cyfrze i odwrotnie,
//   - liczba wielocyfrowa nie może zaczynać się od zera,
//   - dzielenie jest dokładne (A / B = C ⟺ A = B·C w liczbach całkowitych),
//   - program wypisuje wszystkie rozwiązania (niektóre zagadki mają ich
//     więcej niż jedno).
//
// Metoda: wyczerpujące przeszukiwanie permutacji cyfr z nadanymi
// z góry ograniczeniami (liczba cyfr iloczynu/iloczynu, pierwsza cyfra
// ≠ 0) i weryfikacją równania. W Go przeszukanie 10P10 = 3 628 800
// permutacji zajmuje ułamek sekundy.
package z07

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// Solution — pojedyncze rozwiązanie zagadki.
type Solution struct {
	Mapping  map[string]int `json:"mapping"`
	Equation string         `json:"equation"`
}

// Result — wynik rozwiązania zagadki.
type Result struct {
	Formula   string     `json:"formula"`
	Op        string     `json:"op"`
	Letters   int        `json:"letters"`
	SpaceSize int64      `json:"spaceSize"`
	Count     int        `json:"count"`
	Truncated bool       `json:"truncated"`
	Solutions []Solution `json:"solutions"`
	TimeMs    int64      `json:"timeMs"`
}

const maxSolutions = 50 // ile przykładów zwracamy (liczba wszystkich)

// parseEquation parsuje „LHS op RHS = RES”.
func parseEquation(eq string) (lhs, op, rhs, res string, err error) {
	parts := strings.Fields(eq)
	if len(parts) != 5 {
		return "", "", "", "", fmt.Errorf("oczekiwano formatu <liczba> <operator> <liczba> = <liczba>, otrzymano %d tokenów", len(parts))
	}
	lhs, op, rhs, res = parts[0], parts[1], parts[2], parts[4]
	eqSign := parts[3]
	if eqSign != "=" {
		return "", "", "", "", fmt.Errorf("oczekiwano „=”, otrzymano %q", eqSign)
	}
	switch op {
	case "+", "-", "*", "/":
	default:
		return "", "", "", "", fmt.Errorf("nieznany operator %q (dozwolone: + - * /)", op)
	}
	check := func(s, name string) error {
		for _, r := range s {
			if r < 'A' || r > 'Z' {
				return fmt.Errorf("%s: dozwolone są tylko wielkie litery A–Z, znaleziono %q", name, s)
			}
		}
		if len(s) > 10 {
			return fmt.Errorf("%s: liczba nie może mieć więcej niż 10 cyfr", name)
		}
		return nil
	}
	if err := check(lhs, "lewą strona"); err != nil {
		return "", "", "", "", err
	}
	if err := check(rhs, "prawą stronę (czynnik 2)"); err != nil {
		return "", "", "", "", err
	}
	if err := check(res, "wynik"); err != nil {
		return "", "", "", "", err
	}
	return lhs, op, rhs, res, nil
}

// Solve rozwiązuje zagadkę i zwraca wszystkie rozwiązania.
func Solve(eq string) (*Result, error) {
	start := time.Now()
	lhs, op, rhs, res, err := parseEquation(eq)
	if err != nil {
		return nil, err
	}

	// Zbiór liter.
	letterSet := map[rune]struct{}{}
	for _, s := range []string{lhs, rhs, res} {
		for _, r := range s {
			letterSet[r] = struct{}{}
		}
	}
	letters := make([]rune, 0, len(letterSet))
	for r := range letterSet {
		letters = append(letters, r)
	}
	sort.Slice(letters, func(i, j int) bool { return letters[i] < letters[j] })
	if len(letters) > 10 {
		return nil, fmt.Errorf("zagadka ma %d liter, a dostępnych jest tylko 10 cyfr", len(letters))
	}

	// Litery, które nie mogą być zerem (pierwsza litera liczby >1 cyfra).
	noZero := map[rune]bool{}
	for _, s := range []string{lhs, rhs, res} {
		if len(s) > 1 {
			noZero[rune(s[0])] = true
		}
	}

	// Przeszukiwanie permutacji cyfr (rekursja z cięciami).
	digitOf := make([]int32, 26) // -1 = niesprzydzielone
	used := [10]bool{}

	// Cięcia wstępne: liczba cyfr iloczynu/ilorazu musi się zgadzać.
	okSpace := true
	if op == "*" {
		if len(res) < len(lhs)+len(rhs)-1 || len(res) > len(lhs)+len(rhs) {
			okSpace = false
		}
	} else if op == "/" {
		// A / B = C  ⟺  A = B*C: długość A należy do [d(B)+d(C)-1, d(B)+d(C)]
		if len(lhs) < len(rhs)+len(res)-1 || len(lhs) > len(rhs)+len(res) {
			okSpace = false
		}
	}

	solutions := []Solution{}
	count := 0

	var dfs func(k int)
	dfs = func(k int) {
		if k == len(letters) {
			// Weryfikacja pełnego przypisania.
			var A, B, C int64
			for i := range lhs {
				A = A*10 + int64(digitOf[lhs[i]-'A'])
			}
			for i := range rhs {
				B = B*10 + int64(digitOf[rhs[i]-'A'])
			}
			for i := range res {
				C = C*10 + int64(digitOf[res[i]-'A'])
			}
			ok := false
			switch op {
			case "+":
				ok = A+B == C
			case "-":
				ok = A-B == C
			case "*":
				if A != 0 && B > math.MaxInt64/A {
					ok = false // iloczyn przekracza int64, więc na pewno ≠ C
				} else {
					ok = A*B == C
				}
			case "/":
				ok = B != 0 && A%B == 0 && A/B == C
			}
			if ok {
				count++
				if len(solutions) < maxSolutions {
					mapping := map[string]int{}
					for _, r := range letters {
						mapping[string(r)] = int(digitOf[r-'A'])
					}
					numEq := mapDigits(eq, digitOf)
					solutions = append(solutions, Solution{Mapping: mapping, Equation: numEq})
				}
			}
			return
		}
		r := letters[k]
		for d := 0; d < 10; d++ {
			if used[d] {
				continue
			}
			if noZero[r] && d == 0 {
				continue
			}
			used[d] = true
			digitOf[r-'A'] = int32(d)
			dfs(k + 1)
			digitOf[r-'A'] = -1
			used[d] = false
		}
	}

	if okSpace {
		for i := range digitOf {
			digitOf[i] = -1
		}
		dfs(0)
	}

	// Rozmiar przestrzeni: 10Pk permutacji.
	space := int64(10)
	for i := 1; i < len(letters); i++ {
		space *= int64(10 - i)
	}

	return &Result{
		Formula:   strings.Join([]string{lhs, op, rhs, "=", res}, " "),
		Op:        op,
		Letters:   len(letters),
		SpaceSize: space,
		Count:     count,
		Truncated: count > len(solutions),
		Solutions: solutions,
		TimeMs:    time.Since(start).Milliseconds(),
	}, nil
}

// mapDigits podmienia litery na cyfry w równaniu (do wyświetlenia).
func mapDigits(eq string, digitOf []int32) string {
	var b strings.Builder
	for i := 0; i < len(eq); i++ {
		c := eq[i]
		if c >= 'A' && c <= 'Z' {
			b.WriteByte(byte('0' + digitOf[c-'A']))
		} else {
			b.WriteByte(c)
		}
	}
	return b.String()
}
