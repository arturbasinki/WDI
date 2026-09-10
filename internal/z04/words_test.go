package z04

import "testing"

// Wyniki wersji Python dla Pana Tadeusza (top 20):
// rzekł 155, tylko 149, hrabia 129, sędzia 127, tadeusz 107,
// przed 102, jeszcze 101, przez 97, gdzie 94, wszyscy 90,
// wojski 89, potem 86, teraz 83, niech 76, kiedy 74,
// jeśli 73, który 72, nawet 71, znowu 71, wszystko 67.
func TestTop20PanTadeusz(t *testing.T) {
	top := TopWords(PanTadeusz(), 20)
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
