"""
Ten skrypt generuje i rozwiązuje planszę Sudoku za pomocą algorytmu backtracking. 

Na początku generuje planszę Sudoku z 17 wskazówkami, a następnie renderuje ją na ekranie. 

Następnie rozwiązuje planszę Sudoku i zlicza liczbę rozwiązań. 

Jeśli istnieje przynajmniej jedno rozwiązanie, wyświetla pierwsze z nich.

Funkcje używane w skrypcie to:

- `generuj_sudoku`: Generuje planszę Sudoku z 17 wskazówkami.
- `jest_poprawne`: Sprawdza, czy podana liczba jest poprawna w danym wierszu i kolumnie.
- `sprawdz_sudoku `: Sprawdza, czy siatka jest kompletna.
- `rozwiaz_sudoku`: Rozwiązuje siatkę Sudoku za pomocą algorytmu backtracking.
- `renderuj_sudoku`: Renderuje siatkę sudoku jako łańcuch znaków z odpowiednim formatowaniem.
"""


import numpy as np

def generuj_sudoku():
    """
    Generuje planszę Sudoku z 17 wskazówkami.
    """
    sudoku = np.zeros((9, 9), dtype=int)
    for _ in range(23):
        while True:
            row = np.random.randint(0, 9)
            col = np.random.randint(0, 9)
            num = np.random.randint(1, 10)
            if sudoku[row][col] == 0 and jest_poprawne(sudoku, row, col, num):
                sudoku[row][col] = num
                break
    return sudoku

def jest_poprawne(sudoku, row, col, num):
    """
    Sprawdza, czy podana liczba jest poprawna w danym wierszu i kolumnie.
    """
    # Sprawdź unikalność w wierszu
    if num in sudoku[row]:
        return False

    # Sprawdź unikalność w kolumnie
    if num in sudoku[:, col]:
        return False

    # Sprawdź unikalność w podsiatce 3x3
    start_row, start_col = 3 * (row // 3), 3 * (col // 3)
    if num in sudoku[start_row:start_row+3, start_col:start_col+3].flatten():
        return False

    return True

def sprawdz_sudoku(siatka):
    """
    Sprawdza, czy siatka jest kompletna.
    """
    return not any(0 in row for row in siatka)

def rozwiaz_sudoku(grid):
    """
    Rozwiązuje siatkę Sudoku za pomocą algorytmu backtracking i zlicza unikatowe rozwiązania.
    """
    def _rozwiaz_sudoku(grid):
        nonlocal count
        for i in range(81):
            row, col = divmod(i, 9)
            if grid[row][col] == 0:
                for num in range(1, 10):
                    if jest_poprawne(grid, row, col, num):
                        grid[row][col] = num
                        if sprawdz_sudoku(grid):
                            count += 1
                            return True  # Znaleziono rozwiązanie, przerywamy się
                        else:
                            if _rozwiaz_sudoku(grid):
                                return True  # Znaleziono rozwiązanie, przerywamy się
                grid[row][col] = 0  # Zresetuj komórkę do pustej
                return False  # Nie znaleziono rozwiązania
        return True  # Cała siatka jest wypełniona, znaleziono rozwiązanie

    count = 0
    _rozwiaz_sudoku(grid)
    return count


def renderuj_sudoku(sudoku):
    """
    Renderuje siatkę sudoku jako łańcuch znaków z odpowiednim formatowaniem.
    """
    pierwsza_linia = "╔═══╤═══╤═══╦═══╤═══╤═══╦═══╤═══╤═══╗"
    separator = "╟───┼───┼───╫───┼───┼───╫───┼───┼───╢"
    wew_linia = "╠═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╣"
    pionowa_linia = "║ {} │ {} │ {} ║ {} │ {} │ {} ║ {} │ {} │ {} ║"
    ostatnia_linia = "╚═══╧═══╧═══╩═══╧═══╧═══╩═══╧═══╧═══╝"

    print(pierwsza_linia)
    for i in range(9):
        if i in (3, 6):
            print(wew_linia)
        if i in (1, 2, 4, 5, 7, 8):
            print(separator)
        row_values = []
        for j in range(9):
            if sudoku[i, j] == 0:
                row_values.append(" ")
            else:
                row_values.append(str(sudoku[i][j]))
        print(pionowa_linia.format(*row_values))        
    print(ostatnia_linia)



def main():
    # Generuj planszę Sudoku z 17 wskazówkami
    sudoku = generuj_sudoku()

    # Renderuj planszę Sudoku
    renderuj_sudoku(sudoku)

    # Rozwiąż planszę Sudoku i policz liczbę unikatowych rozwiązań
    liczba_rozwiazan = rozwiaz_sudoku(sudoku.copy())
    print(f"Liczba unikatowych rozwiązań: {liczba_rozwiazan}")

if __name__ == "__main__":
    main()


"""
╔═══╤═══╤═══╦═══╤═══╤═══╦═══╤═══╤═══╗
║   │ 3 │ 6 ║   │   │ 8 ║   │   │   ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ 4 │ 5 │   ║   │ 7 │   ║ 1 │   │   ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║   │   │   ║   │   │   ║ 2 │ 3 │   ║
╠═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╣
║   │   │ 5 ║   │ 6 │   ║   │   │   ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ 9 │   │   ║   │   │   ║   │   │ 8 ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║   │   │ 8 ║   │   │ 1 ║   │   │   ║
╠═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╣
║   │   │   ║   │ 9 │   ║   │ 1 │   ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║   │ 1 │   ║   │   │   ║   │   │ 7 ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║   │ 9 │   ║ 3 │   │   ║ 6 │ 4 │   ║
╚═══╧═══╧═══╩═══╧═══╧═══╩═══╧═══╧═══╝
Liczba unikatowych rozwiązań: 1
"""