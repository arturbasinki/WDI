// Testy katalogu zadań (GET /api/tasks) — pilnują, że rejestr w Go jest
// JEDYNYM źródłem tytułów i opisów: frontend nie przechowuje kopii, więc nie
// ma się co rozjechać przy zmianie opisu (P3.1).
package main

import (
	"strings"
	"testing"
)

// TestTasksZwracaOsemZadan — spis zwraca wszystkie osiem zadań, z kompletnymi
// polami i bez duplikatów id (na tym rejestrze budowane są kafelki i nagłówki).
func TestTasksOsobistyRejestr(t *testing.T) {
	srv, _, _ := newTestAPI(t)
	code, body := do(t, "GET", srv.URL+"/api/tasks", nil)
	if code != 200 {
		t.Fatalf("GET /api/tasks: kod %d, oczekiwano 200", code)
	}
	raw, _ := body["tasks"].([]any)
	if len(raw) != 8 {
		t.Fatalf("zadania: %d wpisów, oczekiwano 8", len(raw))
	}
	want := []string{"z01", "z02", "z03", "z04", "z05", "z06", "z07", "z08"}
	seen := map[string]bool{}
	for i, e := range raw {
		task, _ := e.(map[string]any)
		id, _ := task["id"].(string)
		if id != want[i] {
			t.Errorf("pozycja %d: id %q, oczekiwano %q", i, id, want[i])
		}
		if seen[id] {
			t.Errorf("duplikat id %q w rejestrze", id)
		}
		seen[id] = true
		for _, field := range []string{"num", "title", "tag"} {
			if s, _ := task[field].(string); s == "" {
				t.Errorf("%s: pole %q jest puste", id, field)
			}
		}
	}
}

// TestZ08TagBezObietnicyMaksimum — opis z08 nie może obiecywać „maksymalnej
// liczby pustych pól”: generator kopiuje dołki monotonicznie i dla seed 1–8
// oraz 42 daje 55–58 pustych, a teoretyczny sufit jest wyższy (~64).
func TestZ08TagBezObietnicyMaksimum(t *testing.T) {
	srv, _, _ := newTestAPI(t)
	_, body := do(t, "GET", srv.URL+"/api/tasks", nil)
	raw, _ := body["tasks"].([]any)
	for _, e := range raw {
		task, _ := e.(map[string]any)
		if id, _ := task["id"].(string); id != "z08" {
			continue
		}
		tag, _ := task["tag"].(string)
		if strings.Contains(strings.ToLower(tag), "maksymaln") {
			t.Errorf("tag z08 obiecuje maksimum: %q", tag)
		}
		for _, want := range []string{"wiele pustych pól", "55–58", "jednoznaczność"} {
			if !strings.Contains(tag, want) {
				t.Errorf("tag z08 powinien zawierać %q, jest: %q", want, tag)
			}
		}
		return
	}
	t.Fatalf("brak z08 w /api/tasks")
}

// appJS — treść frontendu wbudowanego w binarię (go:embed).
func appJS(t *testing.T) string {
	t.Helper()
	app, err := webFS.ReadFile("web/app.js")
	if err != nil {
		t.Fatalf("odczyt web/app.js: %v", err)
	}
	return string(app)
}

// TestFrontendNieKopiujeKatalogu — kafelki i nagłówki podstron budowane są z
// rejestru Go (GET /api/tasks). Szkielet nagłówka może występować w app.js
// DOKŁADNIE raz — w helperze taskHead(), który wstawia tytuł z katalogu.
// Osiem wywołań taskHead("zNN") to osiem podstron: każde odchylenie znaczy,
// że ktoś z powrotem wpisał tytuł w tekst sztywny. (Proza opisów — „Szyfr
// Gronsfelda to…” — jest celowa i nie jest kopiowaniem katalogu.)
func TestFrontendNieKopiujeKatalogu(t *testing.T) {
	js := appJS(t)
	if !strings.Contains(js, `"/api/tasks"`) {
		t.Errorf("app.js nie woła GET /api/tasks — skąd bierze spis zadań?")
	}
	if n := strings.Count(js, `<div class="task-head"`); n != 1 {
		t.Errorf("szkielet nagłówka w app.js: %d miejsca, oczekiwano 1 (wyłącznie taskHead)", n)
	}
	if n := strings.Count(js, `taskHead("z`); n != len(taskCatalog) {
		t.Errorf("wywołania taskHead: %d podstron, oczekiwano %d", n, len(taskCatalog))
	}
	for _, task := range taskCatalog {
		if strings.Contains(js, `<h2>`+task.Title+`</h2>`) ||
			strings.Contains(js, `<h3>`+task.Title+`</h3>`) {
			t.Errorf("app.js kopiuje tytuł %s (%q) w nagłówku — jedno źródło jest w rejestrze Go",
				task.ID, task.Title)
		}
	}
}

// TestFrontendJednoOknoDeweloperskie — stan pomocniczy widoków (sufit ataku,
// ostatnie rozwiązanie sudoku) mieszka w domknięciu modułu; na oknie
// przeglądarki wisi JEDNO nazwane API do odczytu: window.wdiDebug. Dawniej
// były to dwa anonimowe globalne kanały zapisywalne z konsoli.
func TestFrontendJednoOknoDeweloperskie(t *testing.T) {
	js := appJS(t)
	if strings.Contains(js, "window.__") {
		t.Errorf("app.js wystawia globalne okna window.__… — stan ma być w domknięciu modułu")
	}
	if !strings.Contains(js, "window.wdiDebug") {
		t.Errorf("app.js nie wystawia window.wdiDebug — ma być jedno nazwane API deweloperskie")
	}
}
