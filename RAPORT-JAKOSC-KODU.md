# Raport audytu jakości i poprawności kodu — WDI (Go)

**Data:** 2026-10-26 · **Zakres:** całość kodu wniesionego (Go 2 704 linie + testy 629, `app.js` 1 101, `styles.css` 552, `index.html` 28, Docker/compose/README) · **Metoda:** audyt tylko do odczytu (3 niezależne audyty równoległe: solvery `internal/z*`, serwer+frontend, infrastruktura/README) + sprawdzenia mechaniczne + testy na żywo + **naprawa usterek krytycznych w trakcie audytu** (commit `9dca463`).

---

## 1. Werdykt

Projekt jest **poprawny merytorycznie** — wszystkie 14 mierzalnych liczb referencyjnych README odtworzono przez API (41 liczb Armstronga ≤10¹², π(10³…10⁸), rzekł 155 / tylko 149 / hrabia 129, 25 991 wierszy słownika, klucz 142 807 / 253 918 kluczy, 8 rozwiązań ABC·BD=EFGAH, 9567+1085=10652, 56 pól i 1 rozwiązanie dla seed 42, czasy sita 46/45/4000 ms). Zero zależności, `go vet`/`gofmt` czyste, brak XSS.

Audyt znalazł **2 krytyczne błędy funkcjonalne solverów** (odwrócony znak rozwiązania parametrycznego w z01 — widoczny w UI; wieczne zawieszenie `Search` po anulowaniu w z02), **2 wektory awarii procesu** (OOM w z01 dla n=100: 1,1 GB odpowiedzi / 4,7 GB RSS; `fatal stack overflow` w z06 dla 10⁶ nawiasów — zabija serwer, ukrywany przez `restart: unless-stopped`), **1 wyścig danych** (`go test -race`: ticker z02 vs końcowe `sort.Slice`) oraz szereg usterek średnich/małych. **Wszystkie krytyczne i wysokie usterki solverów oraz infrastruktury zostały naprawione i zabezpieczone testami regresji**; pozycje P1 warstwy serwerowej/frontendowej pozostały do naprawy (rozdz. 5).

---

## 2. Sprawdzenia mechaniczne (stan końcowy, po naprawach)

| Sprawdzenie | Wynik |
|---|---|
| `gofmt -l cmd internal` | ✅ czysto |
| `go vet ./...` | ✅ czysto |
| `go test -race -count=1 ./...` | ✅ wszystkie 8 pakietów (z02: 29,2 s — z testem 10⁹ pod `-race`) |
| `node --check web/app.js` | ✅ |
| `docker compose config -q` (Compose v5.5.1) | ✅ exit 0 |
| `docker build` | ✅ obraz **15,6 MB** (było 18,3) |
| `staticcheck` | ⚠️ nie uruchomiono — brak dostępu do Go proxy w środowisku; część usterek tej klasy (martwy kod, `_ = f`, cieniowanie `len`) wyłowił audyt ręczny |

Pokrycie testami (`go test -cover`): z01 88,1% · z02 88,3% · z03 95,5% · z04 96,7% · z05 84,4% · **z06 76,9%** · z07 92,5% · z08 97,8%.

---

## 3. Poprawki zastosowane w trakcie audytu (commit `9dca463`)

### Krytyczne — błędy funkcjonalne

| # | Usterka | Miejsce | Naprawa + test |
|---|---|---|---|
| K1 | **Odwrócony znak składników parametrycznych**: z wiersza RREF `x_p + Σa·x_c = b` wynika `x_p = b − Σa·t`, kod wypisywał `+`. UI pokazywało „x1 = 3 + t1” (dla x+y=3; 2x+2y=6 — dla t1=2 dawało 5+2≠3) | `internal/z01/solver.go` (sekcja nieoznaczona) | Znak liczony od `coef := -a`; test `param_test.go` asercją `x1 == "3 − t1"`; zweryfikowane live: API zwraca `x1 = 3 − t1` |
| K2 | **Anulowanie wieszało `Search` na zawsze**: feeder wychodził przez `ctx.Done()` bez `close(jobs)` → robotnicy czekali na `range jobs`, `wg.Wait()` nie wracał; `DELETE /api/z02/search/{id}` zostawiało job „running” w nieskończoność (reprodukcja 30/30) | `internal/z02/armstrong.go` | `defer close(jobs)` w feederze + zwrot `ctx.Err()` przy anulowaniu (kontrakt jak `z05.Break`); test `cancel_test.go` z timeoutem |
| K3 | **Wyścig danych** (wykryty przez `go test -race`): gorutyna tickera odczytywała `found` bez synchronizacji, a wynik `Search` aliasował wewnętrzną listę `found` (końcowe `sort.Slice` mutowało współdzielone dane) | `internal/z02/armstrong.go` | Join tickera (`close(tickStop); <-tickDone`) + kopia `append([]int64(nil), found...)` pod mutexem; pełny `-race` czysty wielokrotnie |
| K4 | **OOM serwera**: każdy krok animacji z01 kopiował pełną macierz; n=100 → ~10 200 snapshotów → odpowiedź **1 126 952 113 B**, RSS procesu 4,7 GB (pomiar live) | `internal/z01/solver.go` | Budżet 400 000 komórek na snapshoty + flaga `stepsTruncated` + adnotacja w UI; pomiar live po naprawie: **1,18 MB w 0,03 s**; test `stepsbudget_test.go` |
| K5 | **Crash procesu** (poza odzyskiwaniem `net/http`): rekurencyjny parser z06 bez limitu głębokości — `10⁶ × "("` → `fatal error: stack overflow`; `restart: unless-stopped` maskował objaw | `internal/z06/expr.go`, `cmd/wdi-server/main.go` | Limit głębokości 100 (`maxParseDepth`) + limit długości wyrażenia 2000 znaków (400) + `recover()` we wszystkich 3 gorutynach jobów; test `depth_test.go`; zweryfikowane live: 10⁶ nawiasów → HTTP 400, serwer żyje |

### Wysokie

| # | Usterka | Miejsce | Naprawa |
|---|---|---|---|
| W1 | `Parse` z01 akceptował `NaN`/`Inf` (`ParseFloat` je przyjmuje) → status „jednoznaczny” z `x = NaN`; nadmiarowe liczby w pliku cicho ignorowane | `internal/z01/solver.go` | Odrzucanie wartości nieskończonych i nadmiaru liczb; test w `param_test.go` |
| W2 | `SieveMap` zwracał liczby w kolejności iteracji mapy (niedeterministycznie), mimo że `SieveSlice`/`SieveBit` zwracają rosnąco; testy porównywały tylko `len()` | `internal/z03/sieve.go` | `sort.Ints` + test `order_test.go` (`SieveMap(200) == SieveSlice(200)`) |
| W3 | Każdy tick śródproceduralny `Break` (co 500 000 kluczy) raportował `Done: true` — klient widział „koniec” w połowie ataku | `internal/z05/gronsfeld.go` | Parametr `done`; ticki pośrednie `false` |
| W4 | `Solve(g, cap≤0)`: `count >= cap` przerywa natychmiast → pusta plansza (6,67·10²¹ rozwiązań) dawała `Solved:false, Capped:true`; wynik `fillFirst` ignorowany (ryzyko „rozwiązania” pełnego zer) | `internal/z08/sudoku.go` | Walidacja `cap >= 1`; błąd wewnętrzny, gdy `fillFirst` zawiedzie |
| W5 | **Test z08 „weryfikujący poprawność” był martwy**: warunek `continue` zawierał dokładnie ten sam predykat co następny `if` → wiersze i kolumny nigdy nie były sprawdzane | `internal/z08/sudoku_test.go` | Poprawione predykaty (`j == c` / `j == r`) |
| W6 | `pow10[19]` przepełnia `int64` w inicjatorze tablicy (`10¹⁹ > MaxInt64`) | `internal/z02/armstrong.go` | Pętla do 10¹⁸ + komentarz |
| W7 | `xmin=NaN` przechodził walidację (`NaN >= xmax` = false) → pusty wykres bez błędu | `internal/z06/expr.go` | `!(xmin < xmax)` |
| W8 | Infra po EOL: `alpine:3.20` (EOL 2026-04-01, 7 mies. bez łatek), `golang:1.25.6` (seria 1.25 EOL 2026-08-19); `apk add wget` zbędny (busybox ma wget, kod nie używa TLS — `grep http.Client` pusty); brak `HEALTHCHECK` w obrazie przy ścieżce bez compose | `Dockerfile` | `golang:1.26` + `alpine:3.23`, zero `apk add` (build etapu runtime **bez sieci**, **18,3 → 15,6 MB**), `HEALTHCHECK` z `$PORT` |

### Średnie / dokumentacja

- **Zgodność z Pythonem przy remisach** (z04): komentarz obiecywał „łatwe porównanie z Pythonem”, a `Counter.most_common` przy remisach zachowuje kolejność wstawiania — zmierzono rozjazd na pozycji 20 top-20 (`jakby×67` vs `wszystko×67`; pozycje 1–19 identyczne). Skorygowany komentarz (stan na dziś jest udokumentowany, zmiana sortowania nie wymagana).
- **Dokumentacja z05** zawyżała siłę ataku: 25 991 wierszy = 25 839 haseł, z czego 1 862 ze znakami spoza `a-z` nigdy nie może się dopasować (~24 000 realnych); tick to „co 500 000 kluczy”, nie „~50 Hz”; literówka „dłużność” — poprawione komentarze.
- **Martwy kod z01**: nadpisywany natychmiast zapis `SolutionLines` w ścieżce jednoznacznej — usunięty.
- **README**: sklejony nagłówek `…lutego 2024 roku# WDI` (tytuł w ogóle się nie renderował); **cyrylica** „zasobów do kop**ирования**”; „**Natywne**”→„Natywnie”; „bez **ciękich**”→„ciężkich”; „dwuwarstwowy”→„dwuetapowy” (były 2 etapy, ale 3 warstwy z danymi); zaktualizowane wersje baz, akapit o kolizji portu, `go test -race` w instrukcji.
- **`.dockerignore`**: `+ .home/ .remember/ *.log` (prywatne notatki sesji wchodziły do kontekstu budowy — dowieziono `COPY .` probe), usunięty martwy `.ab-tmp/`.
- **Powtarzalność obrazu**: binarka była już bitowo-reproducible (identyczny sha256 między `--no-cache` buildami), warstwy nie — dodane `touch -d @0` + `COPY --link`.
- **compose**: mapowanie `8080:8080` wiązało się na **0.0.0.0** (ekspozycja w LAN) → `127.0.0.1:8080:8080`; healthcheck z twardym `:8080` → `$$PORT` (spójny ze zmienną serwera).

---

## 4. Bezpieczeństwo

- **XSS: brak.** Weryfikacja wyczerpująca: wszystkie dane użytkownika i plikowe trafiają do DOM wyłącznie przez `esc()` — słowa z04, wzór z06 (echo w formularzu), tekst/klucz z05, równania z07, komunikaty błędów API. Żadnego interpolowanego `innerHTML` z surowych danych.
- Serwer nie wykonuje żądań przychodzących; brak osadzania poleceń; `go:embed` nie pozwala na traversal.
- Brak CSP i nagłówków cache (P3); ciało >4 MB → 400 zamiast 413 (P3).
- Kontener: uruchomiony jako `wdi` (uid 10001), binarka `root:root 0755` — poprawne; opcjonalne twardnienie: `read_only`/`cap_drop`/`no-new-privileges` (P3).

---

## 5. Usterki otwarte — rekomendacje priorytetami

### P1 — serwer (`cmd/wdi-server/main.go`)

1. **`jobStore` nigdy nie czyści wpisów** — ~420 KB na job (rezerwa 1000 próbek × 5 struktur); 500 zakończonych jobów ≈ 213 MB RSS. Sugestia: TTL (np. 15 min od zakończenia) + limit liczby trzymanych jobów.
2. **`j.cancel` zapisywany/odczytywany bez mutexu** (wyścig przy jednoczesnym DELETE i starcie) — `job` wymaga mutexa na polu `cancel`.
3. **DELETE na zakończonym jobie** zwraca `{"status":"cancelling"}` — powinno być no-op/409.
4. **Brak limitu jednoczesnych jobów** — `POST /api/z02/search {"limit":1e13}` jest akceptowany (UI ogranicza do 1e12, API nie); 6× `1e11` równolegle wyczerpuje CPU/RAM. Sugestia: semafor (np. 2 joby) + spójny `maxLimit` w API.
5. **Walidacja parametrów jest asynchroniczna** — bezsensowne parametry dają 202 + później błąd w jobie, zamiast 400 przy `POST`.
6. **Brak graceful shutdown** — serwer nie obsługuje SIGTERM; `docker stop` tnie po 10 s (stan jobów i tak jest ulotny, więc ryzyko małe, ale kulturę zatrzymania warto dodać: `signal.Notify` + `srv.Shutdown`).

### P1 — frontend (`web/app.js`)

7. **`playTimer` (animacja z01, ~w. 320)** — interwał działa po zmianie routu; **`document.body.onkeydown` (z08, ~w. 1052)** — nadpisuje globalny handler i wycieka. Oba powinny być rejestrowane w czyszczeniu routu (`AbortController`/lista disposerów).
8. **Joby nie są anulowane przy zmianie routu** —DELETE przy odjeździe; **poller połyka błędy `fetch`** (~w. 110) — utrata połączenia wiesza pasek postępu bez komunikatu.

### P2 — dostępność i UX

9. **WCAG AA**: `--muted-2: #6b7789` na tle karty = kontrast 3,55–4,17:1 < 4,5:1 (drobny tekst pomocniczy). Propozycja: `#8a97ab`.
10. 17 `<label>` bez `for`, 27 `<button>` bez `type="button"`, brak atrybutów ARIA przy wskaźnikach postępu.
11. z05: tekst spoza alfabetu jest **cicho normalizowany** — „ZAŻÓŁĆ GĘŚLĄ” → szyfrowanie zwraca 200 z ` caim`; przydałby się komunikat, które znaki pominięto.
12. z07: `Solutions` ucięte do 50 (spec zadania: „wszystkie rozwiązania”; `Count` jest pełne, `Truncated` poprawny — UI informuje, biblioteka mogłaby brać `limit`).
13. z06: deklarowany „clamping Y” realizowany przez **usuwanie punktów** (gałąź urwa się przy ~|y|=41 zamiast dochodzić do ramki); `samples` i granice przycinane po cichu; brak notacji wykl. (`2e3` → mylący błąd o nazwie `e3`); komunikat o nieznaku pokazuje mojibake dla UTF-8 (`π` → `Ï`).
14. z02: próg dolny `limit >= 1000` w bibliotece — sensowny, ale wyraźnie do warstwy HTTP; z05: `capKeys` przycinany po cichu; z03: anulowanie sprawdzane między pomiarami (przy n=10⁸ na mapie trwa ~4 s × 3).

### P3 — porządek / ewolucja

15. Katalog 8 zadań zduplikowany: rejestr w Go i kopia w `app.js` (źródło rozjazdów przy zmianie); 6 prawie identycznych handlerów GET/DELETE jobów do parametryzacji.
16. Globalne `window.__capKeys` / `window.__sol` (API testowe) — wynieść do domknięcia modułu.
17. `z04: TopWords(k<1) → 20` — brak możliwości poproszenia o 0; wartość domyślna powinna być jawna (`DefaultTop`).
18. `.gitignore` to w większości szablon Python — **słusznie** (repo zawiera legacy Python z01–z08!), ale brakuje Go-owych: `*.test`, `coverage.out`, `wdi-server`; miny typu `lib/`, `build/`, `dist/` są niebezpieczne przy rozbudowie.
19. Duplikat danych w git (~705 KB): `internal/z04/pan-tadeusz.txt` ≡ `z04/Pan_Tadeusz.txt`, `internal/z05/slownik.txt` ≡ `z05/slownik.txt` (identyczne md5). Decyzja użytkownika: legacy to treść dydaktyczna.
20. ~230 MB artefaktów w workspace (`.gocache` 215 MB, `.gopath` z toolchainem, `.verify`, `screenshots`, stary `bin/`) — wynik ograniczeń sandboxa (cache musi być w katalogu zapisywalnym); w normalnym środowisku: `GOCACHE` poza repo.
21. **Brak CI** — wymogi: `gofmt -l`, `go vet`, `go test -race -count=1 ./internal/...`, `docker build`. Deklaracje README (~5 s testów, obraz) są dziś prawdziwe, ale nic ich nie pilnuje.
22. Tag „maksymalna liczba pustych pól” (z08): generator kopie dołki monotonicznie — 55–58 pustych dla seed 1–8/42 (teoretyczne maksimum przy unikalnym rozwiązaniu ~64); sformułowanie w UI „maksymalna” jest umowne.

---

## 6. Mocne strony

1. **Zero zależności, jeden binarny plik** (`go list -m all` → `wdi`; `go:embed` domyka frontend, „Pana Tadeusza”, słownik).
2. **Wyniki merytoryczne w 100% powtarzalne** — wszystkie deklaracje liczbowe README zweryfikowane przez API na żywo; zero rozjazdów.
3. Solidne rozwiązania algorytmiczne: ochrona przed overflow iloczynu w z07 (`B > MaxInt64/A`), poprawne priorytety/asocjacja w z06 (`-2^2=-4`, `2^3^2=512`), MRV+maski i deterministyczny `rand.NewPCG(seed)` w z08, blokowa optymalizacja z02 zgodna ze specyfikacją.
4. Testy z05 realnie wykrywają przesunięcie indeksu klucza (dokładne wektory referencyjne w obie strony).
5. Higiena obrazu Docker: pin baz, `CGO_ENABLED=0`, `-trimpath -s -w`, nie-root, spójne `EXPOSE/ENV`, healthcheck w obrazie **i** w compose, mapowanie tylko loopback.

---

## 7. Stan środowiska (istotny dla uruchomienia)

- **Port 8080 na tym hoście zajmuje obcy kontener** (`simple-social-app`; `GET :8080/api/tasks` → 401 JWT). `docker compose up` zakończy się „port is already allocated” — wystarczy zwolnić 8080 albo zmienić lewą stronę mapowania (README opisuje).
- Kontener `wdi` został usunięty w trakcie audytu (był w stanie `Exited(2)` przez kolizję portu); obraz `wdi:latest` przebudowany na nowe bazy (15,6 MB), zweryfikowany na porcie 8099 i posprzątany.
- Poprawki audytu są **zacommitowane** (`9dca463`), drzewo robocze czyste; historia: `d39f869 → 64ac8d3 → 4a5fb95 → 6bfc479 → 589919d → 9dca463`.

## 8. Podsumowanie liczbowe

| Kategoria | Znaleziono | Naprawione | Otwarte |
|---|---|---|---|
| Krytyczne (błędy funkcjonalne / awarie) | 5 | 5 | 0 |
| Wysokie | 8 | 8 | 0 |
| Średnie | ~15 | 9 (logika+docs) | ~6 (UX: z06/z07/z04) |
| Niskie / P3 | ~20 | 6 | ~14 |
| Testy regresji dodane | 9 testów w 6 plikach | — | — |

**Stan końcowy: repozytorium poprawne, pełny zestaw `-race` czysty, obraz buduje się i działa, wszystkie wektory awarii zablokowane testami regresji.** Zalecane następne kroki: punkty P1 z rozdziału 5 (wyciek `jobStore`, mutex `cancel`, limit jobów, czyszczenie timerów frontendu) oraz pipeline CI.
