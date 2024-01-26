'''

Ten program zawiera trzy różne implementacje metody Sita Eratostenesa,
a także funkcję do obliczania ilości liczb pierwszych i mierzenia czasu
wykonania dla każdej z tych implementacji. Główna funkcja programu, main,
wywołuje te funkcje dla różnych wartości n.
'''

import time 
import numpy as np

# Implementacja Sita Eratostenesa używając listy
def sito_erastostenesa_list(n):
    # Inicjalizacja listy
    sito = [True] * (n+1)
    p = 2
    # Szukanie liczb pierwszych od 2 do n
    while p * p <= n:
        if sito[p]:
            for i in range(p * p, n+1, p):
                sito[i] = False
        p += 1
    return [p for p in range(2, n+1) if sito[p]]

# Implementacja Sita Eratostenesa używając słownika
def sito_erastostenesa_dict(n):
    # Inicjalizacja słownika
    sito = {i: True for i in range(2, n+1)}
    p = 2
    while p * p <= n:
        if sito[p]:
            for i in range(p * p, n+1, p):
                sito[i] = False
        p += 1
    return [p for p in sito if sito[p]]

# Implementacja Sita Eratostenesa używając tablicy numpy
def sito_erastostenesa_numpy(n):
    # Inicjalizacja tablicy numpy
    sito = np.ones(n+1, dtype=bool)
    sito[0:2] = False
    N_max = int(np.sqrt(len(sito) - 1))
    for i in range(2, N_max + 1):
        if sito[i]:
            sito[i**2::i] = False
    return np.nonzero(sito)[0]

# Funkcja do obliczania ilości liczb pierwszych i mierzenia czasu wykonania
def licz_pierwsze_i_czas(func, n_wartosci):
    for n in n_wartosci:
        start = time.time()
        pierwsze = func(n)
        stop = time.time()
        print(f"Dla n = {n}, znaleziono {len(pierwsze)} liczb pierwszych. Czas: {stop - start} sek.")


def main():
    n_wartosci = [10**i for i in range(3, 9)]

    print("Używając list:")
    licz_pierwsze_i_czas(sito_erastostenesa_list, n_wartosci)

    print("\nUżywając dictionary:")
    licz_pierwsze_i_czas(sito_erastostenesa_dict, n_wartosci)

    print("\nUżywając numpy array:")
    licz_pierwsze_i_czas(sito_erastostenesa_numpy, n_wartosci)

# Wywołanie głównej funkcji programu
if __name__ == "__main__":
    main()

'''

Używając list:
Dla n = 1000, znaleziono 168 liczb pierwszych. Czas: 0.0 sek.
Dla n = 10000, znaleziono 1229 liczb pierwszych. Czas: 0.0009992122650146484 sek.
Dla n = 100000, znaleziono 9592 liczb pierwszych. Czas: 0.009003639221191406 sek.
Dla n = 1000000, znaleziono 78498 liczb pierwszych. Czas: 0.10293865203857422 sek.
Dla n = 10000000, znaleziono 664579 liczb pierwszych. Czas: 1.5012073516845703 sek.
Dla n = 100000000, znaleziono 5761455 liczb pierwszych. Czas: 15.32214641571045 sek.

Używając dictionary:
Dla n = 1000, znaleziono 168 liczb pierwszych. Czas: 0.0 sek.
Dla n = 10000, znaleziono 1229 liczb pierwszych. Czas: 0.0019750595092773438 sek.
Dla n = 100000, znaleziono 9592 liczb pierwszych. Czas: 0.022000551223754883 sek.
Dla n = 1000000, znaleziono 78498 liczb pierwszych. Czas: 0.3162851333618164 sek.
Dla n = 10000000, znaleziono 664579 liczb pierwszych. Czas: 3.6641011238098145 sek.
Dla n = 100000000, znaleziono 5761455 liczb pierwszych. Czas: 42.909550189971924 sek.

Używając numpy array:
Dla n = 1000, znaleziono 168 liczb pierwszych. Czas: 0.0 sek.
Dla n = 10000, znaleziono 1229 liczb pierwszych. Czas: 0.0 sek.
Dla n = 100000, znaleziono 9592 liczb pierwszych. Czas: 0.0 sek.
Dla n = 1000000, znaleziono 78498 liczb pierwszych. Czas: 0.003000020980834961 sek.
Dla n = 10000000, znaleziono 664579 liczb pierwszych. Czas: 0.0409998893737793 sek.
Dla n = 100000000, znaleziono 5761455 liczb pierwszych. Czas: 0.7250897884368896 sek.
'''