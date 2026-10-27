package z04

import "testing"

// Wyniki wersji Python dla Pana Tadeusza (top 20):
// rzekł 155, tylko 149, hrabia 129, sędzia 127, tadeusz 107,
// przed 102, jeszcze 101, przez 97, gdzie 94, wszyscy 90,
// wojski 89, potem 86, teraz 83, niech 76, kiedy 74,
// jeśli 73, który 72, nawet 71, znowu 71, wszystko 67.
func TestTop20PanTadeusz(t *testing.T) {
	top, err := TopWords(PanTadeusz(), 20)
	if err != nil {
		t.Fatalf("TopWords: %v", err)
	}
	if len(top) != 20 {
		t.Fatalf("oczekiwano 20 słów, jest %d", len(top))
	}
	if top[0].Word != "rzekł" || top[0].Count != 155 {
		t.Errorf("pierwsze słowo: %v, oczekiwano rzekł×155", top[0])
	}
	if top[1].Word != "tylko" || top[1].Count != 149 {
		t.Errorf("drugie słowo: %v, oczekiwano tylko×149", top[1])
	}
	want5 := "tadeusz"
	if top[4].Word != want5 {
		t.Errorf("piąte słowo: %v, oczekiwano %q", top[4].Word, want5)
	}
}

// TestTopWordsRejectsNonPositiveK — k <= 0 jest błędem, a nie cichym „20”
// (poprzednio TopWords(tekst, 0) i TopWords(tekst, -5) zwracały 20 wyników).
func TestTopWordsRejectsNonPositiveK(t *testing.T) {
	for _, k := range []int{0, -1, -100} {
		out, err := TopWords("abcdef abcdef ghijkl", k)
		if err == nil {
			t.Errorf("k=%d: oczekiwano błędu, zwrócono %d wyników", k, len(out))
		}
		if out != nil {
			t.Errorf("k=%d: przy błędzie spodziewano się nil, jest %v", k, out)
		}
	}
	// k = 1 jest legalne i zwraca dokładnie jeden wynik.
	one, err := TopWords("abcdef abcdef ghijkl", 1)
	if err != nil || len(one) != 1 || one[0].Word != "abcdef" {
		t.Errorf("k=1: %v, %v (oczekiwano [abcdef])", one, err)
	}
}

// TestDefaultTop — jawna stała domyślna (20) zgodna z treścią zadania.
func TestDefaultTop(t *testing.T) {
	if DefaultTop != 20 {
		t.Fatalf("DefaultTop = %d, oczekiwano 20", DefaultTop)
	}
	a, err := TopWords(PanTadeusz(), DefaultTop)
	if err != nil {
		t.Fatalf("TopWords(DefaultTop): %v", err)
	}
	b, err := TopWords(PanTadeusz(), 20)
	if err != nil {
		t.Fatalf("TopWords(20): %v", err)
	}
	if len(a) != len(b) {
		t.Fatalf("różna liczba wyników: %d vs %d", len(a), len(b))
	}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("różnica na pozycji %d: %v vs %v", i, a[i], b[i])
		}
	}
}

func TestExtractWords(t *testing.T) {
	// „lis” (3 litery), „się”, „i”, „lub” są za krótkie — pomijane.
	words := ExtractWords("Szybki brzydkły lis. się i lub abcdef, ABCDEF!")
	want := []string{"szybki", "brzydkły", "abcdef", "abcdef"}
	if len(words) != len(want) {
		t.Fatalf("words = %v, oczekiwano %v", words, want)
	}
	for i := range want {
		if words[i] != want[i] {
			t.Errorf("words[%d] = %q, oczekiwano %q", i, words[i], want[i])
		}
	}
}
