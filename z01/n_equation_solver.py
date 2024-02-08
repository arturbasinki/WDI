'''
Program służy do rozwiązywania układów równań liniowych o N niewiadomych i korzysta z biblioteki numpy.
Dane do programu są wczytywane z pliku tekstowego. Oto szczegółowy opis działania programu:

main(): Jest to główna funkcja programu, która jest wywoływana, gdy program jest uruchamiany.
Sprawdza ona, czy podano odpowiednią liczbę argumentów wiersza poleceń (nazwę pliku tekstowego).
Jeśli tak, to wczytuje dane z pliku za pomocą funkcji wczytaj_dane() i następnie rozwiązuje układ równań
za pomocą funkcji rozwiaz_uklad().

wczytaj_dane(nazwa_pliku): Ta funkcja otwiera plik tekstowy o podanej nazwie i wczytuje z niego dane.
W pierwszym wierszu pliku powinna znajdować się liczba równań N. Kolejne wiersze powinny zawierać
współczynniki równań i wyrazy wolne. Jeśli dane w pliku są niepoprawne (na przykład, jeśli wiersz
nie zawiera odpowiedniej liczby wartości), funkcja wyświetla komunikat o błędzie i kończy działanie programu.

rozwiaz_uklad(A, B): Ta funkcja rozwiązuje układ równań liniowych o N niewiadomych.
Sprawdza ona, czy macierz A jest kwadratowa i czy wymiary macierzy A i wektora B są zgodne.
Jeśli tak, to rozwiązuje układ równań za pomocą funkcji numpy.linalg.solve().
Jeśli układ równań jest sprzeczny lub nieoznaczony, funkcja wyświetla odpowiedni komunikat.

Program uruchamia się z wiersza poleceń, podając nazwę pliku jako argument, na przykład:
$ python nazwa_programu.py nazwa_pliku.txt.
Należy się upewnić, że plik istnieje w tym samym katalogu, co skrypt Pythona, lub podać pełną ścieżkę do pliku.

'''

import numpy as np
import sys

def wczytaj_dane(nazwa_pliku):
    # Otwórz plik i wczytaj dane
    try:
        with open(nazwa_pliku, 'r') as plik:
            # Wczytaj liczbę równań
            N = int(plik.readline())

            # Wczytaj macierz współczynników i wyrazy wolne
            A = np.zeros((N, N))
            B = np.zeros(N)
            for i in range(N):
                row = list(map(float, plik.readline().split()))
                if len(row) != N+1:
                    raise ValueError("Niepoprawny format danych w pliku.")
                A[i] = row[:-1]
                B[i] = row[-1]
        return A, B
    except Exception as e:
        print(f"Wystąpił błąd podczas wczytywania danych z pliku: {e}")
        sys.exit(1)

def rozwiaz_uklad(A, B):
    # Rozwiąż układ równań
    try:
        # Sprawdź, czy macierz A jest kwadratowa
        if A.shape[0] != A.shape[1]:
            raise ValueError("Macierz A musi być macierzą kwadratową.")
        # Sprawdź, czy wymiary macierzy A i wektora B są zgodne
        if A.shape[0] != B.shape[0]:
            raise ValueError("Wymiary macierzy A i wektora B muszą być zgodne.")
        x = np.linalg.solve(A, B)
        print("Rozwiązanie układu równań to:", x)
    except ValueError as ve:
        # Dane w macierzy A lub/i wektorze B są niepoprawne
        print(f"Wystąpił błąd podczas rozwiązywania układu równań: {ve}")
    except np.linalg.LinAlgError:
        # Układ równań jest sprzeczny lub nieoznaczony
        print("Układ równań jest sprzeczny lub nieoznaczony.")

def main():
    # Sprawdź, czy podano argumenty
    if len(sys.argv) != 2:
        print("Użycie: python nazwa_programu.py nazwa_pliku.txt")
        return

    # Pobierz nazwę pliku z argumentów
    nazwa_pliku = sys.argv[1]

    A, B = wczytaj_dane(nazwa_pliku)
    rozwiaz_uklad(A, B)

if __name__ == "__main__":
    main()
