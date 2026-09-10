**Zadania do samodzielnego rozwiązania**

---

## Rozwiązanie w Go + jednorodna strona WWW

Wszystkie 8 zadań zostało przepisanych **od nowa w Go** (Go 1.25, wyłącznie biblioteka standardowa —
zero zależności zewnętrznych) i zaprezentowane jako **jedna strona WWW**: startowy ekran z 8
kafelkami, każdy prowadzący do graficznej, interaktywnej implementacji zadania. Całość (kod Go,
frontend HTML/CSS/JS, „Pan Tadeusz” i słownik angielski) jest wbudowana w **jeden plik binarny**
(`go:embed`) — nie ma plików konfiguracyjnych ani zasobów do kopирования.

### Uruchomienie

Natywne (Go ≥ 1.25.6):

```bash
go build -o bin/wdi-server ./cmd/wdi-server   # jeden binarny plik
./bin/wdi-server                              # http://localhost:8080
```

Docker (zalecane — zero zależności na hoście):

```bash
docker compose up --build -d    # buduje obraz (2 etapy) i startuje kontener
# → http://localhost:8080
docker compose down             # stop

# lub bez compose:
docker build -t wdi:latest .
docker run --rm -p 8080:8080 wdi:latest
```

Obraz: dwuwarstwowy (budowa `golang:1.25.6` → runtime `alpine:3.20`, użytkownik
nie-root, `CGO_ENABLED=0`), wszystko wbudowane w binarny plik przez `go:embed`,
healthcheck na `/api/tasks`. Zmiana portu: edycja mapowania w `docker-compose.yml`.

Testy (wszystkie pakiety, ~5 s; tryb `-short` ~1 s):

```bash
go test ./internal/...        # pełny zestaw
go test -short ./internal/... # szybki zestaw (bez ciękich przeszukiwań)
```

### Struktura

| katalog | zawartość |
|---|---|
| `internal/z01` | Gauss-Jordan z wyborem wiodącego + nagrywanie kroków (animacja RREF) |
| `internal/z02` | Armstronga: przeszukiwanie blokami, 12 gorutyn, do 10¹² |
| `internal/z03` | sito w `[]bool` / `map[int]bool` / bitmapie + benchmark |
| `internal/z04` | top-N słów (≥ 5 liter), „Pan Tadeusz” wbudowany (go:embed) |
| `internal/z05` | Gronsfeld: szyfruj/deszyfruj + atak słownikowy (go:embed słownik, 25 991 słów) |
| `internal/z06` | własny parser wyrażeń (wstęp rekurencyjny) + wykres z asymptotami |
| `internal/z07` | zagadki literowo-cyfrowe: wyczerpujące 10Pₖ z cięciami |
| `internal/z08` | sudoku: generacja + kopanie dołków z gwarancją unikalności |
| `cmd/wdi-server` | serwer HTTP (stdlib) + API JSON + wbudowany frontend |
| `cmd/wdi-server/web` | frontend: vanilla JS (hash routing) + SVG, zero build |

Ciężkie obliczenia (szukanie Armstronga, benchmark sita, atak na Gronsfelda) działają jako
**joby tła**: `POST` zwraca `jobId`, `GET` odpytuje postęp/wynik, `DELETE` anuluje.

### Zweryfikowane wyniki

- **z01**: przykład z zadania → `x=18, y=-58, z=35` (układ jednoznaczny); rozpoznawane układy
  sprzeczne i nieoznaczone (rozwiązanie parametryczne).
- **z02**: 41 liczb Armstronga do 10¹² — pełna lista zweryfikowana wyczerpującym przeszukiwaniem
  (31 min na 12 gorutynach); wartości ≤ 10⁹ pokrywają testy jednostkowe.
- **z03**: liczby pierwsze zgodne z wartościami referencyjnymi (168 / 1229 / 9592 / 78 498 /
  664 579 / 5 761 455); dla n = 10⁷: wycinek ~46 ms, bitmapa ~45 ms, mapa ~4 s.
- **z04**: top 3 dla „Pana Tadeusza”: *rzekł* 155, *tylko* 149, *hrabia* 129 (zgodne z Pythonem).
- **z05**: klucz 312063 → „computer science” = „fpop whsbsilhoee”; przysłowie z zadania złamane
  atakiem słownikowym: klucz **142807**, „there is no use crying over spilt milk” (253 918 kluczy).
- **z06**: `1/x` rysowany jako dwie gałęzie z asymptotą x ≈ 0; `tan(x)` — 4 asymptoty na [−6.3, 6.3].
- **z07**: `ABC·BD=EFGAH` — 8 rozwiązań (m.in. 623·29=18067); `SEND+MORE=MONEY` — jedno
  (9567+1085=10652).
- **z08**: generowane zagadki (np. seed 42: 56 pól pustych) mają dokładnie **jedno** rozwiązanie —
  weryfikowane niezależnie (Python) i przez zliczacz rozwiązań w Go.

Opis rozwiązania każdego zadania znajduje się na jego stronie WWW (panel „Opis rozwiązania”),
komentarze w kodzie są po polsku, a wyniki można wywołać interaktywnie.

---

1. Proszę napisać program rozwiązujący układ równań N równań liniowych o N niewiadomych.
    Dane dla problemu należy wczytać z pliku tekstowego. W pierwszym wierszu znajduje się liczba
    równań N, kolejne wiersze zawierają macierz współczynników oraz wyrazy wolne, na przykład plik:
    3
    1 2 3 7
    - 1 2 4 6
    2 1 1 13
    Odpowiada układowi 3 równań o 3 niewiadomych w postaci:
    X+2Y+3Z=
    - X+2Y+4Z=
    2X+Y+Z=
    Program powinien uwzględnić przypadki układu nieoznaczonego i sprzecznego. Wskazówka:
    rozważyć zastosowanie biblioteki numpy.
2. Liczby Armstronga to N-cyfrowa liczba naturalna która jest sumą swoich cyfr podniesionych do
    potęgi N. Na przykład: 153 = 1^3 +5^3 +3^3. Proszę napisać program znajdujący jak najwięcej takich liczb.
3. Metoda Sita Eratostenesa. Ze zbioru liczb naturalnych z przedziału [2,n], wybieramy najmniejszą,
    czyli 2, i wykreślamy wszystkie jej wielokrotności większe od niej samej. Z pozostałych liczb
    wybieramy najmniejszą niewykreśloną liczbę (3) i usuwamy wszystkie jej wielokrotności większe od
    niej samej. Według tej samej procedury postępujemy dla kolejnych liczb. Proces ten pozostawia
    nieskreślone wyłącznie liczby pierwsze. Proszę napisać program wyszukujący liczby pierwsze w
    zadanym zakresie. Proszę zaimplementować program na 3 sposoby: używając listy, słownika i tablicy
    array z pakietu numpy. Dla każdej ze struktur obliczyć ilość liczb pierwszych dla n równego
    103 ,10^4 ,10^5 ,10^6 ,10^7 ,10^8 , a następnie porównać ich czasy działania.
4. Proszę napisać program, który wczytuje tekst z podanego pliku i wypisuje 20 najczęściej
    występujących w nim słów. Program powinien ignorować krótkie słowa (krótsze niż 5 liter) typu: i,
    lub, się, aby, żeby , itp. Proszę podać wyniki dla tekstu „Pana Tadeusza”.
5. Komputer jest doskonałym narzędziem służącym do szyfrowania i deszyfrowania tajnych
    wiadomości. W metodzie Gronsfelda, będącą modyfikacją szyfru Cezara, stosuje się klucz liczbowy.
    Biorąc klucz o wartości 312063 i niezaszyfrowany tekst „computer science”, uzyskujemy następujący
    szyfrogram:   
    3120633120633120   
    computer science  
    fpop whsbsilhoee  
    Kolejne litery są przesuwane o kolejne wartości z klucza. Alfabet użyty w systemie zawiera spacje i 26
    małych liter alfabetu angielskiego. Proszę napisać programy dokonujące szyfrowania i deszyfrowania
    pliku tekstowego zadanym kluczem. Następnie proszę napisać program łamiący ten szyfr. Korzystając
    ze słownika wyrazów języka angielskiego (plik słownik.txt) proszę odszyfrować angielskie przysłowie:
    „axjmrlamuhnvayum jsbkvggpzgz zqmna tjpmh”. Długość klucza jest mniejsza od 10.


7. Używając biblioteki matplotlib, napisać program rysujący wykresy funkcji jednej zmiennej (na
    przykład: y=x*x- 6 *x+3). Jako dane należy wczytać wzór funkcji oraz przedział dla zmiennej x.
    Wskazówka: przydatna będzie funkcja eval. Proszę uwzględnić istnienie asymptot pionowych np. w
    funkcji y = 1/x.
8. Proszę napisać program, rozwiązujący zagadki polegające na znalezieniu cyfr, które w równaniu
    zostały zastąpione literami A-J. Na przykład dla równania ABC * BD = EFGAH, rozwiązaniem jest:
    A=6, B=2, C=3, D=9, E=1, F=8, G=0, H=7, czyli jest to równanie 623 * 29 = 18067. Program powinien
    wczytywać z klawiatury równanie w postaci literowej i wypisywać rozwiązanie w postaci cyfrowej.
    Można założyć, że działania będą w postaci <liczba> <operator> <liczba> = <liczba>, gdzie
    operatorem może być + - * /. Proszę uwzględnić, że niektóre zagadki mogą mieć więcej rozwiązań.
9. Proszę napisać program generujący zagadki sudoku. Wygenerowana zagadka powinna posiadać tylko
    jedno rozwiązanie i maksymalną liczbę pustych pól.

Uwagi:

- Należy rozwiązać minimum 4 zadania.
- Maksymalna ocena za rozwiązanie N zadań nie może przekroczyć N/2+
- Rozwiązanie zadania powinno zawierać: krótki opis rozwiązania, kod programu wraz z komentarzami,
    wyniki programu dla przykładowych danych.
- Rozwiązania zadań wraz z opisami proszę umieścić w jednym plik typu DOC lub PDF.
- Rozwiązania zadań proszę umieścić w kolejności rosnących numerów.
- Plik z rozwiązaniami proszę umieścić w systemie Moodle do końca lutego 2024 roku# WDI
Rozwiązania zadań WDI
