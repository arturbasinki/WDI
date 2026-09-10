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
// Atak słownikowy wylicza kolejne kandydackie klucze (1-cyfrowe,
// 2-cyfrowe, ...), odszyfrowuje nimi szyfrogram i sprawdza, czy ponad
// 80% słów wyniku znajduje się w słowniku angielskim (wbudowany,
// 25 991 słów, z pliku slownik.txt). Pierwszy pasujący klucz jest
// zwracany wraz z tekstem jawnym.
package z05

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"time"
)

//go:embed slownik.txt
var slownikText string

// Alphabet — alfabet szyfru: spacja + 26 małych liter.
const Alphabet = " abcdefghijklmnopqrstuvwxyz"

// MaxKeyLen — maksymalna długość klucza (w zadaniu: dłużność < 10).
const MaxKeyLen = 9

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

// normalizeText przygotowuje tekst do szyfrowania: zamienia małe litery,
// usuwa znaki niebędące w alfabecie. Zwraca błąd, jeśli po normalizacji
// tekst jest pusty.
func normalizeText(text string) (string, error) {
	var b strings.Builder
	for _, r := range strings.ToLower(text) {
		if r >= 'a' && r <= 'z' || r == ' ' {
			b.WriteRune(r)
		}
	}
	out := b.String()
	if out == "" {
		return "", fmt.Errorf("tekst nie zawiera żadnych liter a–z")
	}
	return out, nil
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
func Encrypt(text, key string) (string, error) {
	t, err := normalizeText(text)
	if err != nil {
		return "", err
	}
	if err := validateKey(key); err != nil {
		return "", err
	}
	out := make([]rune, len(t))
	for i, r := range t {
		iAl := strings.IndexRune(Alphabet, r)
		out[i] = rune(Alphabet[(iAl+int(key[i%len(key)]-'0'))%len(Alphabet)])
	}
	return string(out), nil
}

// Decrypt deszyfruje szyfrogram kluczem cyfrowym.
func Decrypt(cipher, key string) (string, error) {
	t, err := normalizeText(cipher)
	if err != nil {
		return "", err
	}
	if err := validateKey(key); err != nil {
		return "", err
	}
	out := make([]rune, len(t))
	for i, r := range t {
		iAl := strings.IndexRune(Alphabet, r)
		out[i] = rune(Alphabet[(iAl-int(key[i%len(key)]-'0')+len(Alphabet))%len(Alphabet)])
	}
	return string(out), nil
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
// capKeys ogranicza całkowitą liczbę wypróbowanych kluczy (bezpieczeństwo).
// onTick jest wywoływany z postęgiem ~50 razy na sekundę.
func Break(ctx context.Context, cipher string, maxLen int, capKeys int64, onTick func(BreakProgress)) (BreakResult, error) {
	t, err := normalizeText(cipher)
	if err != nil {
		return BreakResult{}, err
	}
	if maxLen < 1 || maxLen > MaxKeyLen {
		return BreakResult{}, fmt.Errorf("maxLen musi być z 1..%d", MaxKeyLen)
	}
	if capKeys <= 0 {
		capKeys = 2_000_000
	}
	if capKeys > 100_000_000 {
		capKeys = 100_000_000
	}
	if onTick == nil {
		onTick = func(BreakProgress) {}
	}

	start := time.Now()
	keysTried := int64(0)
	notify := func(len int, found bool, plain, key string) {
		onTick(BreakProgress{KeysTried: keysTried, Len: len, Found: found, Plain: plain, Key: key, Done: true})
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
				return BreakResult{}, ctx.Err()
			}
			if keysTried >= capKeys {
				notify(l, false, "", "")
				return BreakResult{Found: false, KeysTried: keysTried, TimeMs: timeMs()}, nil
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
				notify(l, true, plain, string(buf))
				return BreakResult{Found: true, Plain: plain, Key: string(buf), KeysTried: keysTried, TimeMs: timeMs()}, nil
			}
			if keysTried%500_000 == 0 {
				notify(l, false, "", "")
			}
		}
	}
	notify(maxLen, false, "", "")
	return BreakResult{Found: false, KeysTried: keysTried, TimeMs: timeMs()}, nil
}
