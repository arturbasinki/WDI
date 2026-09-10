// Package z01 rozwiązuje układy N równań liniowych o N niewiadomych.
//
// Dane układu wczytywane są z pliku tekstowego (lub z polika w interfejsie WWW)
// w następującym formacie:
//
//	wiersz 1: liczba równań N
//	wiersze 2..N+1: współczynniki równania oraz wyraz wolny (N+1 liczb)
//
// np. plik:
//
//	3
//	1 2 3 7
//	-1 2 4 6
//	2 1 1 13
//
// odpowiada układowi:
//
//	X1 + 2X2 + 3X3 = 7
//	-X1 + 2X2 + 4X3 = 6
//	2X1 + X2 + X3 = 13
//
// Układ rozwiązujemy metodą eliminacji Gaussa z wyborem wiodącego (Gauss-Jordan),
// co daje macierz w postaci schodkowej zredukowanej (RREF). Dzięki temu program
// poprawnie klasyfikuje trzy przypadki:
//
//  1. układ jednoznaczny    — dokładnie jedno rozwiązanie,
//  2. układ sprzeczny       — brak rozwiązań (wiersz: 0 = c ≠ 0),
//  3. układ nieoznaczony    — nieskończenie wiele rozwiązań
//     (podawane parametrycznie przez niewiadome wolne).
//
// Dodatkowo program zwraca poszczególnie kroki eliminacji (zamiany wierszy,
// dzielenie wiersza, odejmowanie wielokrotności), które interfejs WWW może
// zaprezentować jako animację redukcji macierzy.
package z01

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Epsilon — próg, poniżej którego wartość uznajemy za zero (dane zadania
// mają rzędy wielkości ~1, więc stały próg jest wystarczający).
const Epsilon = 1e-9

// Status opisuje charakter rozwiązania układu.
type Status string

const (
	// StatusUnique — układ jednoznaczny (dokładnie jedno rozwiązanie).
	StatusUnique Status = "unique"
	// StatusNone — układ sprzeczny (brak rozwiązań).
	StatusNone Status = "none"
	// StatusInfinite — układ nieoznaczony (nieskończenie wiele rozwiązań).
	StatusInfinite Status = "infinite"
)

// Step to pojedyncza operacja na wierszach dla animacji eliminacji.
type Step struct {
	Description string      `json:"description"`
	Matrix      [][]float64 `json:"matrix"`
}

// SolutionLine to jeden wiersz wyniku: wartość niewiadomej albo parametr.
type SolutionLine struct {
	Var  string `json:"var"`
	Expr string `json:"expr"`
}

// Result jest kompletnym wynikiem rozwiązania układu.
type Result struct {
	N             int            `json:"n"`
	Status        Status         `json:"status"`
	X             []float64      `json:"x,omitempty"`
	RREF          [][]float64    `json:"rref"`
	FreeVars      []string       `json:"freeVars,omitempty"`
	SolutionLines []SolutionLine `json:"solutionLines"`
	Steps         []Step         `json:"steps"`
	Rank          int            `json:"rank"`
}

// Parse parsuje tekst w formacie pliku zadania: w pierwszym wierszu N,
// następnie N wierszy po N+1 liczb (współczynniki i wyraz wolny).
// Rozdziałem mogą być spacje, tabulator, przecinki i znaki nowej linii.
// Tolerowany jest zapis „- 1" (minus oddzielony spacją od cyfry).
func Parse(text string) (A [][]float64, err error) {
	// Zastąp przecinki spacjami i podziel na tokeny.
	text = strings.ReplaceAll(text, ",", " ")
	fields := strings.Fields(text)
	// Złącz pojedyncze „-" z następnym tokenem (zapis „- 1 2 4 6").
	merged := make([]string, 0, len(fields))
	for i := 0; i < len(fields); i++ {
		if fields[i] == "-" && i+1 < len(fields) && fields[i+1] != "-" {
			merged = append(merged, "-"+fields[i+1])
			i++
			continue
		}
		merged = append(merged, fields[i])
	}
	if len(merged) == 0 {
		return nil, fmt.Errorf("brak danych wejściowych")
	}
	n, err := strconv.Atoi(merged[0])
	if err != nil {
		return nil, fmt.Errorf("pierwsza wartość (liczba równań N) nie jest liczbą całkowitą: %q", merged[0])
	}
	if n < 1 || n > 100 {
		return nil, fmt.Errorf("N musi być z przedziału 1..100, podano %d", n)
	}
	need := n * (n + 1)
	if len(merged) < 1+need {
		return nil, fmt.Errorf("oczekiwałem %d liczb (N wierszy po N+1), znaleziono %d", need, len(merged)-1)
	}
	A = make([][]float64, n)
	for i := 0; i < n; i++ {
		A[i] = make([]float64, n+1)
		for j := 0; j < n+1; j++ {
			v, err := strconv.ParseFloat(merged[1+i*(n+1)+j], 64)
			if err != nil {
				return nil, fmt.Errorf("niepoprawna liczba %q w wierszu %d", merged[1+i*(n+1)+j], i+1)
			}
			A[i][j] = v
		}
	}
	return A, nil
}

// snapshot zwraca kopię macierzy (do zapisu kroku animacji).
func snapshot(m [][]float64) [][]float64 {
	c := make([][]float64, len(m))
	for i := range m {
		c[i] = append([]float64(nil), m[i]...)
	}
	return c
}

// Solve rozwiązuje macierz poszerzoną A [N][N+1] metodą Gaussa-Jordana
// i zwraca wynik wraz z krokami eliminacji.
func Solve(A [][]float64) Result {
	n := len(A)
	M := snapshot(A)
	res := Result{N: n, RREF: M, Steps: []Step{}}
	steps := &res.Steps
	*steps = append(*steps, Step{Description: "Macierz poszerzona [A|b] na start.", Matrix: snapshot(M)})

	pivotCols := []int{}
	row := 0
	for col := 0; col < n && row < n; col++ {
		// Wybór wiodącego: wiersz o największej |wartości| w kolumnie col.
		p := row
		for r := row + 1; r < n; r++ {
			if math.Abs(M[r][col]) > math.Abs(M[p][col]) {
				p = r
			}
		}
		if math.Abs(M[p][col]) < Epsilon {
			continue // kolumna liniowo zależna — bez wiodącego
		}
		if p != row {
			M[p], M[row] = M[row], M[p]
			*steps = append(*steps, Step{
				Description: fmt.Sprintf("Zamiana wierszy W%d i W%d (wybór wiodącego w kolumnie %d).", row+1, p+1, col+1),
				Matrix:      snapshot(M),
			})
		}
		// Normalizacja wiersza wiodącego.
		div := M[row][col]
		for j := col; j <= n; j++ {
			M[row][j] /= div
		}
		*steps = append(*steps, Step{
			Description: fmt.Sprintf("W%d := W%d / %s", row+1, row+1, fmtNum(div)),
			Matrix:      snapshot(M),
		})
		// Eliminacja reszty wierszy (góra i dół) — pełny Gauss-Jordan.
		for r := 0; r < n; r++ {
			if r == row {
				continue
			}
			f := M[r][col]
			if math.Abs(f) < Epsilon {
				continue
			}
			for j := col; j <= n; j++ {
				M[r][j] -= f * M[row][j]
			}
			*steps = append(*steps, Step{
				Description: fmt.Sprintf("W%d := W%d − %s·W%d", r+1, r+1, fmtNum(f), row+1),
				Matrix:      snapshot(M),
			})
		}
		pivotCols = append(pivotCols, col)
		row++
	}

	res.Rank = len(pivotCols)
	res.RREF = snapshot(M)

	// Klasyfikacja: szukaj wiersza sprzecznego (0 = c, c ≠ 0).
	for i := 0; i < n; i++ {
		allZero := true
		for j := 0; j < n; j++ {
			if math.Abs(M[i][j]) >= Epsilon {
				allZero = false
				break
			}
		}
		if allZero && math.Abs(M[i][n]) >= Epsilon {
			res.Status = StatusNone
			res.SolutionLines = []SolutionLine{{Var: "—", Expr: "Układ sprzeczny — nie ma rozwiązania."}}
			return res
		}
	}

	pivotSet := make([]bool, n)
	for _, c := range pivotCols {
		pivotSet[c] = true
	}

	if res.Rank == n {
		// Układ jednoznaczny: x_j = wartość z wiersza, którego wiodący to j.
		res.Status = StatusUnique
		res.X = make([]float64, n)
		res.SolutionLines = make([]SolutionLine, n)
		for i, c := range pivotCols {
			v := M[i][n]
			if math.Abs(v) < Epsilon {
				v = 0
			} else {
				// usuń szum float64: zaokrąglenie do 9 miejsc po przecinku
				v = math.Round(v*1e9) / 1e9
			}
			res.X[c] = v
			res.SolutionLines[c] = SolutionLine{Var: varName(c), Expr: fmtNum(v)}
		}
		for c := 0; c < n; c++ {
			res.SolutionLines[c] = SolutionLine{Var: varName(c), Expr: fmtNum(res.X[c])}
		}
		return res
	}

	// Układ nieoznaczony: niewiadome wolne (kolumny bez wiodącego) to parametry.
	res.Status = StatusInfinite
	free := []int{}
	for c := 0; c < n; c++ {
		if !pivotSet[c] {
			free = append(free, c)
		}
	}
	res.FreeVars = make([]string, len(free))
	for i := range free {
		res.FreeVars[i] = "t" + strconv.Itoa(i+1)
	}
	// x_wiodący = (wyraz wolny) − Σ M[w][c_wolny]·t_k
	res.SolutionLines = make([]SolutionLine, n)
	for i := range res.SolutionLines {
		res.SolutionLines[i] = SolutionLine{Var: varName(i), Expr: "?"}
	}
	for i, p := range pivotCols {
		expr := strings.Builder{}
		expr.WriteString(fmtNum(M[i][n]))
		for k, c := range free {
			a := M[i][c]
			if math.Abs(a) < Epsilon {
				continue
			}
			if a > 0 {
				expr.WriteString(" + ")
				if math.Abs(a-1) < Epsilon {
					expr.WriteString(res.FreeVars[k])
				} else {
					expr.WriteString(fmtNum(a) + "·" + res.FreeVars[k])
				}
			} else {
				expr.WriteString(" − ")
				aa := -a
				if math.Abs(aa-1) < Epsilon {
					expr.WriteString(res.FreeVars[k])
				} else {
					expr.WriteString(fmtNum(aa) + "·" + res.FreeVars[k])
				}
			}
		}
		res.SolutionLines[p] = SolutionLine{Var: varName(p), Expr: expr.String()}
	}
	for _, c := range free {
		idx := 0
		for _, f := range free {
			if f == c {
				break
			}
			idx++
		}
		res.SolutionLines[c] = SolutionLine{Var: varName(c), Expr: res.FreeVars[idx]}
	}
	return res
}

// varName zwraca nazwę niewiadomej: x1, x2, ...
func varName(col int) string { return "x" + strconv.Itoa(col+1) }

// fmtNum ładnie formatuje liczbę zmiennoprzecinkową (bez zbędnych zer).
func fmtNum(v float64) string {
	if math.Abs(v) < Epsilon {
		return "0"
	}
	s := strconv.FormatFloat(v, 'g', 10, 64)
	return s
}
