package z05

import (
	"context"
	"strings"
	"testing"
)

// Przykład z zadania: klucz 312063, „computer science” → „fpop whsbsilhoee”.
func TestEncryptExample(t *testing.T) {
	out, err := Encrypt("computer science", "312063")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if out != "fpop whsbsilhoee" {
		t.Fatalf("Encrypt = %q, oczekiwano %q", out, "fpop whsbsilhoee")
	}
}

func TestDecryptRoundtrip(t *testing.T) {
	plain := "the quick brown fox"
	cipher, err := Encrypt(plain, "42719")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	back, err := Decrypt(cipher, "42719")
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if back != plain {
		t.Fatalf("roundtrip: %q != %q", back, plain)
	}
}

// Złamanie szyfrogramu z przykładu (klucz 6-cyfrowy).
func TestBreakExample(t *testing.T) {
	res, err := Break(context.Background(), "fpop whsbsilhoee", 6, 2_000_000, nil)
	if err != nil {
		t.Fatalf("Break: %v", err)
	}
	if !res.Found {
		t.Fatal("atak nie znalazł klucza")
	}
	if res.Key != "312063" {
		t.Errorf("klucz = %q, oczekiwano 312063 (znaleziony szyfrogram może mieć inny klucz o tej samej dekodowalności)", res.Key)
	}
	if res.Plain != "computer science" {
		t.Errorf("plain = %q, oczekiwano %q", res.Plain, "computer science")
	}
}

// Zagadka z zadania: „axjmrlamuhnvayum jsbkvggpzgz zqmna tjpmh”, klucz < 10 cyfr.
// Weryfikacja ręczna (cyfra po cyfrze) potwierdza: klucz 142807,
// tekst jawny: „there is no use crying over spilt milk”.
func TestBreakPuzzle(t *testing.T) {
	res, err := Break(context.Background(), "axjmrlamuhnvayum jsbkvggpzgz zqmna tjpmh", 9, 100_000_000, nil)
	if err != nil {
		t.Fatalf("Break: %v", err)
	}
	if !res.Found {
		t.Fatal("atak nie znalazł klucza dla zagadki")
	}
	if res.Key != "142807" {
		t.Errorf("klucz = %q, oczekiwano 142807", res.Key)
	}
	if plain := strings.TrimSpace(res.Plain); plain != "there is no use crying over spilt milk" {
		t.Errorf("plain = %q, oczekiwano %q", plain, "there is no use crying over spilt milk")
	}
	t.Logf("rozwiązanie zagadki: plain=%q key=%s (wypróbowano %d kluczy, %d ms)", res.Plain, res.Key, res.KeysTried, res.TimeMs)
}

func TestValidateKey(t *testing.T) {
	if _, err := Encrypt("abc", "12a"); err == nil {
		t.Error("klucz z literami nie został odrzucony")
	}
	if _, err := Encrypt("abc", "1234567890"); err == nil {
		t.Error("klucz 10-cyfrowy nie został odrzucony")
	}
}
