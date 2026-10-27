// Package z05 implementuje szyfr Gronsfelda (modyfikacja szyfru Cezara
// z kluczem liczbowym), jego deszyfrowanie oraz atak słownikowy
// pozwalający złamać szyfrogram.
//
// Alfabet szyfru zawiera 27 znaków: spację i 26 małych liter
// angielskiego alfabetu: „ abcdefghijklmnopqrstuvwxyz”.
// Kolejne litery wiadomości są przesuwane o kolejne cyfry klucza
// (cyklujące się po jego długości). Przykład:
//
//	klucz:       312063 312063 3120
//	jawny:       computer science
//	szyfrogram:  fpop whsbsilhoee
//
// Normalizacja (małe litery + usuwanie znaków spoza alfabetu) NIE jest
// niema: Encrypt/Decrypt/Break zwracają zbiór unikalnych znaków, które
// zostały pominięte, więc „ZAŻÓŁĆ GĘŚLĄ” nie znika bez śladu.
//
// Atak słownikowy wylicza kolejne kandydackie klucze (1-cyfrowe,
// 2-cyfrowe, ...), odszyfrowuje nimi szyfrogram i sprawdza, czy ponad
// 80% słów wyniku znajduje się w słowniku angielskim. Słownik
// (wbudowany, plik slownik.txt) ma 25 991 wierszy = 25 839 unikalnych
// haseł; hasła ze znakami spoza „a-z" (1 862) nigdy nie mogą zostać
// dopasowane, bo tekst jawny jest normalizowany do „a-z" + spacja.
// Tick postępu wysyłany jest co 500 000 wypróbowanych kluczy.
// Pierwszy pasujący klucz jest zwracany wraz z tekstem jawnym.
package z05

import (
	"context"
	_ "embed"
	"fmt"
	"sort"
	"strings"
	"time"
)

//go:embed slownik.txt
var slownikText string

// Alphabet — alfabet szyfru: spacja + 26 małych liter.
const Alphabet = " abcdefghijklmnopqrstuvwxyz"

// MaxKeyLen — maksymalna długość klucza (w zadaniu: dłużność < 10).
const MaxKeyLen = 9

// BreakCapDefault — capKeys stosowany, gdy żądający poda wartość <= 0.
const BreakCapDefault = int64(2_000_000)

// BreakCapMax — twardy sufit liczby wypróbowanych kluczy w ataku.
// Powyżej niego capKeys jest PRZYCINANY, a wynik ustawia CapApplied=true
// (dawniej cięcie było nieme: caller nie wiedział, że atak skończył się
// na suficie, a nie na wyczerpaniu przestrzeni kluczy).
const BreakCapMax = int64(100_000_000)

// Dictionary — słownik angielskich słów (zbiór).
var Dictionary = func() map[string]struct{} {
	m := make(map[string]struct{}, 30000)
	for _, line := range strings.Split(slownikText, "\n") {
		w := strings.TrimSpace(line)
		if w != "" {
			m[w] = struct{}{}
		}
	}
	return m
}()

// normalizeText przygotowuje tekst do szyfrowania: zamienia na małe litery
// i usuwa znaki spoza alfabetu (czyli wszystko poza „a–z” i spacją).
//
// Drugi wynik to POSORTOWANY po punkcie kodowym Unicode zbiór unikalnych
// znaków pominiętych — w pisowni po normalizacji (małe litery), bo na tym
// alfabecie pracuje szyfr. Gdy tekst mieścił się w alfabecie, zbiór jest
// pusty. Dawniej normalizacja była cicha: „ZAŻÓŁĆ GĘŚLĄ” dawało 200 OK z
// wynikiem „  ” i siedem polskich liter zniknęło bez śladu.
//
// Trzeci wynik to błąd, jeśli po normalizacji nie zostało nic do zaszyfrowania.
func normalizeText(text string) (string, []string, error) {
	var b strings.Builder
	seen := map[rune]struct{}{}
	var dropped []rune
	for _, r := range strings.ToLower(text) {
		if r >= 'a' && r <= 'z' || r == ' ' {
			b.WriteRune(r)
			continue
		}
		if _, ok := seen[r]; !ok {
			seen[r] = struct{}{}
			dropped = append(dropped, r)
		}
	}
	out := b.String()
	if out == "" {
		return "", nil, fmt.Errorf("tekst nie zawiera żadnych liter a–z")
	}
	sort.Slice(dropped, func(i, j int) bool { return dropped[i] < dropped[j] })
	skipped := make([]string, len(dropped))
	for i, r := range dropped {
		skipped[i] = string(r)
	}
	return out, skipped, nil
}

// validateKey sprawdza klucz: cyfry 0–9, długość 1..MaxKeyLen.
func validateKey(key string) error {
	if len(key) < 1 || len(key) > MaxKeyLen {
		return fmt.Errorf("klucz musi mieć 1–%d cyfr, podano %d", MaxKeyLen, len(key))
	}
	for _, r := range key {
		if r < '0' || r > '9' {
			return fmt.Errorf("klucz może zawierać tylko cyfry 0–9")
		}
	}
	return nil
}

// Encrypt szyfruje tekst kluczem cyfrowym.
//
// Zwraca szyfrogram, posortowany zbiór znaków POMINIĘTYCH przy normalizacji
// (pusty, gdy tekst był w całości w alfabecie) oraz błąd.
func Encrypt(text, key string) (string, []string, error) {
	t, skipped, err := normalizeText(text)
	if err != nil {
		return "", skipped, err
	}
	if err := validateKey(key); err != nil {
		return "", skipped, err
	}
	out := make([]rune, len(t))
	for i, r := range t {
		iAl := strings.IndexRune(Alphabet, r)
		out[i] = rune(Alphabet[(iAl+int(key[i%len(key)]-'0'))%len(Alphabet)])
	}
	return string(out), skipped, nil
}

// Decrypt deszyfruje szyfrogram kluczem cyfrowym.
// Zwraca tekst jawny, posortowany zbiór znaków pominiętych przy normalizacji
// oraz błąd.
func Decrypt(cipher, key string) (string, []string, error) {
	t, skipped, err := normalizeText(cipher)
	if err != nil {
		return "", skipped, err
	}
	if err := validateKey(key); err != nil {
		return "", skipped, err
	}
	out := make([]rune, len(t))
	for i, r := range t {
		iAl := strings.IndexRune(Alphabet, r)
		out[i] = rune(Alphabet[(iAl-int(key[i%len(key)]-'0')+len(Alphabet))%len(Alphabet)])
	}
	return string(out), skipped, nil
}

// BreakProgress — postęp ataku słownikowego.
type BreakProgress struct {
	KeysTried int64  `json:"keysTried"`
	Len       int    `json:"len"`
	Found     bool   `json:"found"`
	Plain     string `json:"plain,omitempty"`
	Key       string `json:"key,omitempty"`
	Done      bool   `json:"done"`
}

// BreakResult — wynik ataku.
type BreakResult struct {
	Found     bool   `json:"found"`
	Plain     string `json:"plain,omitempty"`
	Key       string `json:"key,omitempty"`
	KeysTried int64  `json:"keysTried"`
	TimeMs    int64  `json:"timeMs"`
	// Skipped — unikalne znaki pominięte przy normalizacji szyfrogramu
	// (posortowane po punkcie kodowym), nieobecne gdy tekst był w alfabecie.
	Skipped []string `json:"skipped,omitempty"`
	// CapApplied — prawda, gdy żądany capKeys był > BreakCapMax i został
	// przycięty do sufitu; atak skończył się wtedy na suficie, nie na
	// wyczerpaniu przestrzeni kluczy.
	CapApplied bool `json:"capApplied,omitempty"`
}

// breakKey odszyfrowuje szyfrogram podanym kluczem i sprawdza
// czy > 80% słów znajduje się w słowniku.
func breakKey(cipher string, key string) (string, bool) {
	out := make([]rune, len(cipher))
	for i, r := range cipher {
		iAl := strings.IndexRune(Alphabet, r)
		out[i] = rune(Alphabet[(iAl-int(key[i%len(key)]-'0')+len(Alphabet))%len(Alphabet)])
	}
	plain := string(out)
	words := strings.Fields(plain)
	if len(words) == 0 {
		return plain, false
	}
	bad := 0
	for _, w := range words {
		if _, ok := Dictionary[w]; !ok {
			bad++
			// wczesne zakończenie: > 20% słów poza słownikiem
			if 5*bad > len(words) {
				return plain, false
			}
		}
	}
	return plain, true
}

// Break atakuje szyfrogram klucze o długościach 1..maxLen, wyliczane w
// kolejności leksykograficznej (tak jak itertools.product w Pythonie).
// capKeys ogranicza całkowitą liczbę wypróbowanych kluczy (bezpieczeństwo):
// wartość <= 0 daje BreakCapDefault, a wartość > BreakCapMax jest przycinana
// do sufitu — wtedy wynik ma CapApplied=true (nic nie jest przycinane po cichu).
// onTick jest wywoływany z postęgiem ~50 razy na sekundę.
func Break(ctx context.Context, cipher string, maxLen int, capKeys int64, onTick func(BreakProgress)) (BreakResult, error) {
	t, skipped, err := normalizeText(cipher)
	if err != nil {
		return BreakResult{Skipped: skipped}, err
	}
	if maxLen < 1 || maxLen > MaxKeyLen {
		return BreakResult{Skipped: skipped}, fmt.Errorf("maxLen musi być z 1..%d", MaxKeyLen)
	}
	capApplied := false
	if capKeys <= 0 {
		capKeys = BreakCapDefault
	}
	if capKeys > BreakCapMax {
		capKeys = BreakCapMax
		capApplied = true
	}
	if onTick == nil {
		onTick = func(BreakProgress) {}
	}

	start := time.Now()
	keysTried := int64(0)
	notify := func(length int, found bool, plain, key string, done bool) {
		onTick(BreakProgress{KeysTried: keysTried, Len: length, Found: found, Plain: plain, Key: key, Done: done})
	}
	timeMs := func() int64 { return time.Since(start).Milliseconds() }

	// Generator kluczy: długość 1..maxLen, każda długość — 0..10^len-1
	// z zerami na początku (np. „007”), w kolejności leksykograficznej.
	for l := 1; l <= maxLen; l++ {
		limit := int64(1)
		for i := 0; i < l; i++ {
			limit *= 10
		}
		buf := make([]byte, l)
		for v := int64(0); v < limit; v++ {
			if ctx.Err() != nil {
				return BreakResult{Skipped: skipped, CapApplied: capApplied}, ctx.Err()
			}
			if keysTried >= capKeys {
				notify(l, false, "", "", true)
				return BreakResult{Found: false, KeysTried: keysTried, TimeMs: timeMs(),
					Skipped: skipped, CapApplied: capApplied}, nil
			}
			// zapisz v jako l-cyfrowy łańcuch z zerami
			x := v
			for i := l - 1; i >= 0; i-- {
				buf[i] = byte('0' + x%10)
				x /= 10
			}
			plain, ok := breakKey(t, string(buf))
			keysTried++
			if ok {
				notify(l, true, plain, string(buf), true)
				return BreakResult{Found: true, Plain: plain, Key: string(buf), KeysTried: keysTried,
					TimeMs: timeMs(), Skipped: skipped, CapApplied: capApplied}, nil
			}
			if keysTried%500_000 == 0 {
				notify(l, false, "", "", false)
			}
		}
	}
	notify(maxLen, false, "", "", true)
	return BreakResult{Found: false, KeysTried: keysTried, TimeMs: timeMs(),
		Skipped: skipped, CapApplied: capApplied}, nil
}
