// Package z04 wczytuje tekst z pliku i zlicza najczęściej występujące słowa.
//
// Program ignoruje krótkie słowa (krótsze niż 5 liter) typu „i”, „lub”,
// „się”, „aby”, „żeby” itd. Zwraca 20 najczęściej występujących słów
// (zliczanie bez różnicowania wielkości liter).
//
// Tekst „Pana Tadeusza” Adama Mickiewicza jest wbudowany w binarię
// (go:embed), dzięki czemu wynik dla domyślnych danych nie wymaga
// żadnego pliku na dysku. Można też przesłać dowolny własny tekst
// przez interfejs WWW i zliczyć jego słowa.
package z04

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

//go:embed pan-tadeusz.txt
var panTadeuszText string

// WordCount to słowo z liczbą wystąpień.
type WordCount struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// PanTadeusz zwraca pełny tekst „Pana Tadeusza” (bez BOM).
func PanTadeusz() string {
	return strings.TrimPrefix(panTadeuszText, "\uFEFF")
}

// MinLength — minimalna długość zliczanego słowa (w literach).
const MinLength = 5

// DefaultTop — domyślna liczba zwracanych słów (treść zadania: „20
// najczęściej występujących słów”). Warstwa wywołująca (HTTP/UI) stosuje
// tę wartość JAWNIE, gdy klient nie podaje swojego k; biblioteka nie
// podmienia jej w ukryciu (patrz TopWords).
const DefaultTop = 20

// ExtractWords dzieli tekst na słowa: ciągły fragment liter (Unicode),
// bez rozróżniania wielkości liter, o długości co najmniej MinLength.
// Znak apostrofa nie dzieli słowa (w tekście polskim to nie ma znaczenia,
// traktujemy go jako literę).
func ExtractWords(text string) []string {
	text = strings.ToLower(text)
	var words []string
	var cur []rune
	flush := func() {
		if len(cur) >= MinLength {
			words = append(words, string(cur))
		}
		cur = cur[:0]
	}
	for _, r := range text {
		if unicode.IsLetter(r) || r == '\'' {
			cur = append(cur, r)
			continue
		}
		flush()
	}
	flush()
	return words
}

// TopWords zwraca k najczęściej występujących słów w tekście.
// Sortowanie: malejąca liczba wystąpień, a przy remisie — alfabetycznie.
// Uwaga: Python (Counter.most_common) przy remisach zachowuje kolejność
// wstawiania, więc przy pozycjach „remisowych" pełne listy mogą się
// różnić — gwarantowane jest pokrywanie się rankingów bez remisów.
//
// k <= 0 jest BŁĘDEM, a nie cichym „20”: dawniej ujemne/zerowe k
// podmieniało domyślną wartość w środku biblioteki, więc wywołujący nie
// miał jak zapytać o „zero wyników" ani sprawdzić, jaka wartość weszła
// po drodze. Domyślną liczbę stosuje warstwa wywołująca (patrz DefaultTop).
func TopWords(text string, k int) ([]WordCount, error) {
	if k < 1 {
		return nil, fmt.Errorf("k musi być >= 1, podano %d", k)
	}
	counts := map[string]int{}
	for _, w := range ExtractWords(text) {
		counts[w]++
	}
	out := make([]WordCount, 0, len(counts))
	for w, c := range counts {
		out = append(out, WordCount{Word: w, Count: c})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Word < out[j].Word
	})
	if len(out) > k {
		out = out[:k]
	}
	return out, nil
}
