**Zadania do samodzielnego rozwiązania**

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
    Biorąc klucz o wartości 31206 3 i niezaszyfrowany tekst „computer science”, uzyskujemy następujący
    szyfrogram:
    3120633120633120
    computer science
    fpop whsbsilhoee
    Kolejne litery są przesuwane o kolejne wartości z klucza. Alfabet użyty w systemie zawiera spacje i 26
    małych liter alfabetu angielskiego. Proszę napisać programy dokonujące szyfrowania i deszyfrowania
    pliku tekstowego zadanym kluczem. Następnie proszę napisać program łamiący ten szyfr. Korzystając
    ze słownika wyrazów języka angielskiego (plik słownik.txt) proszę odszyfrować angielskie przysłowie:
    „axjmrlamuhnvayum jsbkvggpzgz zqmna tjpmh”. Długość klucza jest mniejsza od 10.


6. Używając biblioteki matplotlib, napisać program rysujący wykresy funkcji jednej zmiennej (na
    przykład: y=x*x- 6 *x+3). Jako dane należy wczytać wzór funkcji oraz przedział dla zmiennej x.
    Wskazówka: przydatna będzie funkcja eval. Proszę uwzględnić istnienie asymptot pionowych np. w
    funkcji y = 1/x.
7. Proszę napisać program, rozwiązujący zagadki polegające na znalezieniu cyfr, które w równaniu
    zostały zastąpione literami A-J. Na przykład dla równania ABC * BD = EFGAH, rozwiązaniem jest:
    A=6, B=2, C=3, D=9, E=1, F=8, G=0, H=7, czyli jest to równanie 623 * 29 = 18067. Program powinien
    wczytywać z klawiatury równanie w postaci literowej i wypisywać rozwiązanie w postaci cyfrowej.
    Można założyć, że działania będą w postaci <liczba> <operator> <liczba> = <liczba>, gdzie
    operatorem może być + - * /. Proszę uwzględnić, że niektóre zagadki mogą mieć więcej rozwiązań.
8. Proszę napisać program generujący zagadki sudoku. Wygenerowana zagadka powinna posiadać tylko
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
