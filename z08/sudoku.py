"""
Ten program generuje zagadkę Sudoku, która ma tylko jedno rozwiązanie. 

1. Najpierw, program generuje pustą planszę Sudoku.
2. Następnie, program wypełnia planszę losowymi liczbami od 1 do 9, upewniając się, że każda liczba jest unikalna
    w swoim wierszu, kolumnie i podsiatce 3x3.
3. Program kontynuuje wypełnianie planszy, dopóki nie osiągnie stanu, w którym każde kolejne pole może być
    wypełnione tylko jedną liczbą.
4. Jeśli taki stan jest osiągnięty, program zwraca wygenerowaną planszę Sudoku jako zagadkę do rozwiązania.
5. W trakcie generowania planszy, program wyświetla pasek postępu, który pokazuje, ile prób generowania planszy
    zostało już wykonanych. Do renderowania paska postępu program korzysta z biblioteki tqdm.
6. Po wygenerowaniu planszy, program sprawdza, czy wygenerowana zagadka Sudoku ma tylko jedno rozwiązanie.
    Jeśli tak, program zwraca planszę jako wynik. W przeciwnym razie, proces generowania jest powtarzany.
"""

import os
import random
from tqdm import tqdm

# Sprawdza, czy liczba jest w wierszu
def w_wierszu(liczba, wiersz, sudoku):
    # Zwraca prawdę, jeśli liczba jest w danym wierszu sudoku
    return liczba in sudoku[wiersz]

# Sprawdza, czy liczba jest w kolumnie
def w_kolumnie(liczba, kolumna, sudoku):
    # Przechodzi przez każdy wiersz w danej kolumnie
    for wiersz in range(9):
        # Jeśli liczba jest w danym wierszu, zwraca prawdę
        if sudoku[wiersz][kolumna] == liczba:
            return True
    # Jeśli liczba nie jest w kolumnie, zwraca fałsz
    return False

# Sprawdza, czy liczba jest w podsiatce
def w_podsiatce(liczba, wiersz, kolumna, sudoku):
    # Oblicza początkowy wiersz i kolumnę dla podsiatki 3x3
    wiersz_start = (wiersz // 3) * 3
    kolumna_start = (kolumna // 3) * 3
    # Przechodzi przez każde pole w podsiatce
    for i in range(wiersz_start, wiersz_start + 3):
        for j in range(kolumna_start, kolumna_start + 3):
            # Jeśli liczba jest w podsiatce, zwraca prawdę
            if sudoku[i][j] == liczba:
                return True
    # Jeśli liczba nie jest w podsiatce, zwraca fałsz
    return False

# Generuje liczbę
def generuj_liczbe(wiersz, kolumna, sudoku):
    # Generuje losową liczbę między 1 a 9, dopóki liczba ta nie jest unikalna w danym wierszu, kolumnie i podsiatce
    while True:
        liczba = random.randint(1, 9)
        if not (w_wierszu(liczba, wiersz, sudoku) or w_kolumnie(liczba, kolumna, sudoku) or w_podsiatce(liczba, wiersz, kolumna, sudoku)):
            return liczba

# Generuje zagadkę Sudoku
def generuj_sudoku():
    # Inicjalizuje pustą planszę sudoku
    sudoku = [[0 for _ in range(9)] for _ in range(9)]
    liczba_pustych_pol = 0

    # Wypełnianie sudoku
    for i in range(9):
        for j in range(9):
            # Z prawdopodobieństwem 30% wypełnia pole losową liczbą
            if random.random() < 0.21:
                liczba = generuj_liczbe(i, j, sudoku)
                sudoku[i][j] = liczba
            else:
                # Zwiększa liczbę pustych pól
                liczba_pustych_pol += 1

    # Zwraca wygenerowane sudoku i liczbę pustych pól
    return sudoku, liczba_pustych_pol

# Sprawdza, czy istnieje tylko jedno rozwiązanie
def ma_jedno_rozwiazanie(plansza):
    # Kopiowanie planszy
    kopia_planszy = [row[:] for row in plansza]

    def _rozw_sudoku(plansza):
        # Przechodzi przez każde pole na planszy Sudoku (od lewej do prawej, od góry do dołu).
        for i in range(9):
            for j in range(9):
                # Gdy napotka puste pole (z wartością 0), próbuje wstawić każdą liczbę od 1 do 9.
                if plansza[i][j] == 0:
                    for liczba in range(1, 10):
                        # Sprawdza, czy liczba jest unikalna w danym wierszu, kolumnie i podsiatce.
                        if not (w_wierszu(liczba, i, plansza) or w_kolumnie(liczba, j, plansza) or w_podsiatce(liczba, i, j, plansza)):
                            # Jeśli liczba jest unikalna, umieszcza ją na planszy.
                            plansza[i][j] = liczba
                            # Następnie próbuje wypełnić resztę planszy. Jeśli udaje się to zrobić, zwraca prawdę.
                            if _rozw_sudoku(plansza):
                                return True
                            # Jeśli liczba nie pasuje, usuwa ją z planszy (cofanie kroku).
                            plansza[i][j] = 0
                    # Jeśli żadna liczba nie pasuje do danego pola, zwraca fałsz. To powoduje "cofnięcie" poprzednich kroków.
                    return False
        # Jeśli cała plansza jest wypełniona, zwraca prawdę.
        return True

    # Zwraca prawdę, jeśli sudoku ma tylko jedno rozwiązanie, w przeciwnym razie zwraca fałsz
    return _rozw_sudoku(kopia_planszy)


# Główna funkcja
def main():
    # Czyszczenie ekranu przed wygenerowaniem planszy i paska postępu w zależności od sytemu operacyjnego
    os.system('cls' if os.name == 'nt' else 'clear')
    print("Generowanie zagadki Sudoku. Pasek postępu pokazuje liczbę prób generowania planszy z jednym rozwiązaniem.")
    # Generowanie sudoku
    liczba_prob = 0
    # Inicjalizacja paska postępu
    pbar = tqdm(total=100)
    while True:
        liczba_prob += 1
        # Aktualizacja paska postępu
        pbar.update(1)
        # Wymusza aktualizację paska postępu
        pbar.refresh()
        sudoku, liczba_pustych_pol = generuj_sudoku()
        kopia_planszy = [row[:] for row in sudoku]
        # Jeśli sudoku ma tylko jedno rozwiązanie, przerywa pętlę
        if ma_jedno_rozwiazanie(kopia_planszy):
            # i zamyka pasek postępu
            pbar.close()  
            break


    print("\nWygenerowana plansza:\n")
    # Wyświetlanie sudoku
    for wiersz in sudoku:
        print(" ".join(str(liczba) if liczba != 0 else "." for liczba in wiersz))

    print(f"\nLiczba pustych pól: {liczba_pustych_pol}")
    print(f"Liczba prób: {liczba_prob}")
    print("Wygenerowana zagadka sudoku posiada tylko jedno rozwiązanie.")

if __name__ == "__main__":
    main()