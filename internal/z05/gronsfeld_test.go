package z05

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

// Przykład z zadania: klucz 312063, „computer science” → „fpop whsbsilhoee”.
func TestEncryptExample(t *testing.T) {
	out, skipped, err := Encrypt("computer science", "312063")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if out != "fpop whsbsilhoee" {
		t.Fatalf("Encrypt = %q, oczekiwano %q", out, "fpop whsbsilhoee")
	}
	if len(skipped) != 0 {
		t.Errorf("tekst w alfabecie: skipped=%v, oczekiwano pustego", skipped)
	}
}

func TestDecryptRoundtrip(t *testing.T) {
	plain := "the quick brown fox"
	cipher, _, err := Encrypt(plain, "42719")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	back, skipped, err := Decrypt(cipher, "42719")
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if back != plain {
		t.Fatalf("roundtrip: %q != %q", back, plain)
	}
	if len(skipped) != 0 {
		t.Errorf("skipped=%v, oczekiwano pustego", skipped)
	}
}

// TestSkippedSigns — normalizacja nie może być cicha (żądanie P2).
// „ZAŻÓŁĆ GĘŚLĄ” nie zawiera ani jednej litery a–z poza „a g l”, więc
// dawniej szyfrowanie zwracało 200 OK z niemal pustym wynikiem i nikt
// nie mówił, które znaki zniknęły.
func TestSkippedSigns(t *testing.T) {
	out, skipped, err := Encrypt("ZAŻÓŁĆ GĘŚLĄ", "1")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	// porządek: rosnąco po punkcie kodowym Unicode (stąd „ó” na początku)
	want := []string{"ó", "ą", "ć", "ę", "ł", "ś", "ż"}
	if !reflect.DeepEqual(skipped, want) {
		t.Errorf("skipped = %q, oczekiwano %q", skipped, want)
	}
	// po normalizacji zostaje „za gl” (Z i A to litery a–z po obniżeniu
	// wielkości; spacja jest w alfabecie) — 5 znaków, wszystkie przesunięte o 1
	if len(out) != 5 {
		t.Errorf("długość szyfrogramu = %d (%q), oczekiwano 5", len(out), out)
	}
	// szyfrogram zwrócony przez Encrypt da się odszyfrować tym samym kluczem
	// i wychodzi dokładnie znormalizowany tekst; Decrypt na surowym tekście
	// przesuwa w drugą stronę, więc nie jest równy Encrypt — sprawdzasz pętlę
	back, skipped2, err := Decrypt(out, "1")
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if back != "za gl" {
		t.Errorf("Decrypt(szyfrogram) = %q, oczekiwano %q", back, "za gl")
	}
	// Decrypt na surowym, niezaszyfrowanym tekście raportuje te same znaki
	_, skipped2, err = Decrypt("ZAŻÓŁĆ GĘŚLĄ", "1")
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !reflect.DeepEqual(skipped2, want) {
		t.Errorf("Decrypt skipped = %q, oczekiwano %q", skipped2, want)
	}

	// Break też musi raportować pominięte znaki (szyfrogram po normalizacji
	// może być krótszy, niż myśli użytkownik).
	res, err := Break(context.Background(), "fpop whsbsilhoee, 2024!", 6, 2_000_000, nil)
	if err != nil {
		t.Fatalf("Break: %v", err)
	}
	// pominięte: spacja jest w alfabecie — zostaje przecinek i cyfry
	wantBreak := []string{"!", ",", "0", "2", "4"}
	if !reflect.DeepEqual(res.Skipped, wantBreak) {
		t.Errorf("Break.Skipped = %q, oczekiwano %q", res.Skipped, wantBreak)
	}
}

// TestSkippedAllNonAlphabet — tekst, z którego po normalizacji nie zostało
// NIC (żadnej litery a–z i żadnej spacji), to błąd, nie pusty wynik.
func TestSkippedAllNonAlphabet(t *testing.T) {
	out, skipped, err := Encrypt("12345!!!", "1")
	if err == nil {
		t.Errorf("oczekiwano błędu, jest %q (skipped %q)", out, skipped)
	}
	// spacja NALEŻY do alfabetu, więc „12345 !!!” to nie pusty tekst: zostaje
	// jedna spacja (rozdzielająca cyfry od wykrzykników) i operacja ma sens.
	out2, skipped2, err := Encrypt("12345 !!!", "1")
	if err != nil {
		t.Errorf("tekst ze spacją nie powinien dawać błędu: %v", err)
	}
	if out2 != "a" {
		t.Errorf("wynik = %q, oczekiwano %q (spacja przesunięta o 1)", out2, "a")
	}
	if len(skipped2) != 6 {
		t.Errorf("skipped = %q, oczekiwano 6 znaków (pięć cyfr i wykrzyknik)", skipped2)
	}
}

// TestBreakCapApplied — przycięcie capKeys do sufitu musi być widoczne
// w wyniku (dawniej cięcie było nieme).
func TestBreakCapApplied(t *testing.T) {
	// maxLen=1 → przestrzeń to 10 kluczy; capKeys ponad sufit zostaje przycięty
	res, err := Break(context.Background(), "fpop whsbsilhoee", 1, BreakCapMax+1, nil)
	if err != nil {
		t.Fatalf("Break: %v", err)
	}
	if !res.CapApplied {
		t.Error("CapApplied=false, oczekiwano true (capKeys był > BreakCapMax)")
	}
	if res.KeysTried != 10 {
		t.Errorf("KeysTried = %d, oczekiwano 10 (cała przestrzeń maxLen=1)", res.KeysTried)
	}
	// capKeys w granicach: flaga nie ustawiona
	res2, err := Break(context.Background(), "fpop whsbsilhoee", 6, 2_000_000, nil)
	if err != nil {
		t.Fatalf("Break: %v", err)
	}
	if res2.CapApplied {
		t.Error("CapApplied=true przy capKeys w granicach")
	}
	// capKeys <= 0 → domyślny, bez flagi
	res3, err := Break(context.Background(), "fpop whsbsilhoee", 1, 0, nil)
	if err != nil {
		t.Fatalf("Break: %v", err)
	}
	if res3.CapApplied {
		t.Error("CapApplied=true przy capKeys=0 (stosowany domyślny)")
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
	if _, _, err := Encrypt("abc", "12a"); err == nil {
		t.Error("klucz z literami nie został odrzucony")
	}
	if _, _, err := Encrypt("abc", "1234567890"); err == nil {
		t.Error("klucz 10-cyfrowy nie został odrzucony")
	}
}
