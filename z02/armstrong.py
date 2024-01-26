'''

Ta funkcja sprawdza, czy dana liczba n jest liczbą Armstronga.
Liczba Armstronga to taka liczba, która jest równa sumie swoich cyfr podniesionych
do potęgi równej liczbie jej cyfr. Na przykład, 153 jest liczbą Armstronga,
ponieważ 153 = 1^3 + 5^3 + 3^3.
Funkcja ta jest używana w programie do znalezienia wszystkich liczb Armstronga w określonym zakresie.
'''

def jest_armstrong(n):
    # Oblicz liczbę cyfr w liczbie n
    liczba_cyfr = len(str(n))

    # Zainicjuj zmienną suma na 0
    suma = 0

    # Skopiuj wartość n do zmiennej tymczasowej
    temp = n

    # Pętla wykonuje się, dopóki temp jest większe od 0
    while temp > 0:
        # Oblicz ostatnią cyfrę liczby temp
        cyfra = temp % 10

        # Dodaj do sumy cyfrę podniesioną do potęgi liczba_cyfr
        suma += cyfra ** liczba_cyfr

        # Usuń ostatnią cyfrę z liczby temp
        temp //= 10

    # Zwróć prawdę, jeśli n jest równa sumie, w przeciwnym razie zwróć fałsz
    return n == suma

def znajdz_liczby_armstrong():
    # Znajduje i drukuje wszystkie liczby Armstronga w określonym zakresie.
    for i in range(0, 10**7):  # Zmienić zakres według potrzeb
        if jest_armstrong(i):
            print(i)

def main():
    znajdz_liczby_armstrong()

if __name__ == "__main__":
    main()
