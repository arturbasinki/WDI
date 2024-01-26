# Definiujemy alfabet, który będzie używany do deszyfrowania
alfabet = ' abcdefghijklmnopqrstuvwxyz'

# Funkcja szyfrująca
def gronsfeld_encrypt(wiadomosc, klucz):
    # Inicjalizujemy pusty ciąg, który będzie przechowywał zaszyfrowaną wiadomość
    zaszyfrowana_wiadomosc = ''
    # Inicjalizujemy indeks klucza na 0
    klucz_index = 0

    # Iterujemy przez każdy znak w wiadomości
    for char in wiadomosc:
        # Pobieramy wartość przesunięcia z klucza
        przesuniecie = int(klucz[klucz_index])
        # Obliczamy nowy indeks znaku po przesunięciu
        nowy_indeks = (alfabet.index(char) + przesuniecie) % len(alfabet)
        # Dodajemy nowy znak do zaszyfrowanej wiadomości
        zaszyfrowana_wiadomosc += alfabet[nowy_indeks]
        # Przechodzimy do następnego indeksu klucza
        klucz_index = (klucz_index + 1) % len(klucz)

    # Zwracamy zaszyfrowaną wiadomość
    return zaszyfrowana_wiadomosc

# Funkcja deszyfrująca
def gronsfeld_decrypt(zaszyfrowana_wiadomosc, klucz):
    # Inicjalizujemy pusty ciąg, który będzie przechowywał odszyfrowaną wiadomość
    rozszyfrowana_wiadomosc = ''
    # Inicjalizujemy indeks klucza na 0
    klucz_index = 0

    # Iterujemy przez każdy znak w zaszyfrowanej wiadomości
    for char in zaszyfrowana_wiadomosc:
        # Pobieramy wartość przesunięcia z klucza
        przesuniecie = int(klucz[klucz_index])
        # Obliczamy nowy indeks znaku po przesunięciu w przeciwnym kierunku
        nowy_indeks = (alfabet.index(char) - przesuniecie) % len(alfabet)
        # Dodajemy nowy znak do odszyfrowanej wiadomości
        rozszyfrowana_wiadomosc += alfabet[nowy_indeks]
        # Przechodzimy do następnego indeksu klucza
        klucz_index = (klucz_index + 1) % len(klucz)

    # Zwracamy odszyfrowaną wiadomość
    return rozszyfrowana_wiadomosc

def main():
    wiadomosc = input("Wprowadź wiadomość: ")
    klucz = input("Wprowadź klucz: ")

    zaszyfrowana_wiadomosc = gronsfeld_encrypt(wiadomosc, klucz)
    rozszyfrowana_wiadomosc = gronsfeld_decrypt(zaszyfrowana_wiadomosc, klucz)

    print("Zaszyfrowana wiadomość:", zaszyfrowana_wiadomosc)
    print("Rozszyfrowana wiadomość:", rozszyfrowana_wiadomosc)

if __name__ == "__main__":
    main()