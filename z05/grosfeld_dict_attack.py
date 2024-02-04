import itertools
import time

# Wczytujemy słownik
with open('slownik.txt', 'r') as f:
    # Tworzymy zbiór wszystkich słów w słowniku, usuwając białe znaki na początku i na końcu każdego słowa
    dictionary = set(word.strip() for word in f)

# Funkcja deszyfrująca
def gronsfeld_decrypt(encrypted_message, key):
    # Definiujemy alfabet, który będzie używany do deszyfrowania
    alphabet = ' abcdefghijklmnopqrstuvwxyz'
    # Inicjalizujemy pusty ciąg, który będzie przechowywał odszyfrowaną wiadomość
    decrypted_message = ''
    # Inicjalizujemy indeks klucza na 0
    key_index = 0

    # Iterujemy przez każdy znak w zaszyfrowanej wiadomości
    for char in encrypted_message:
        # Pobieramy wartość przesunięcia z klucza
        shift = int(key[key_index])
        # Obliczamy nowy indeks znaku po przesunięciu w przeciwnym kierunku
        new_char_index = (alphabet.index(char) - shift) % len(alphabet)
        # Dodajemy nowy znak do odszyfrowanej wiadomości
        decrypted_message += alphabet[new_char_index]
        # Przechodzimy do następnego indeksu klucza
        key_index = (key_index + 1) % len(key)

    # Zwracamy odszyfrowaną wiadomość
    return decrypted_message

# Funkcja łamiąca szyfr
def break_gronsfeld(encrypted_message):
    # Generujemy wszystkie możliwe klucze o długości mniejszej niż 10
    possible_keys = (''.join(p) for i in range(1, 10) for p in itertools.product('0123456789', repeat=i))

    # Iterujemy przez wszystkie możliwe klucze
    for key in possible_keys:
        # Deszyfrujemy wiadomość za pomocą danego klucza
        decrypted_message = gronsfeld_decrypt(encrypted_message, key)
        words = decrypted_message.split()
        # Sprawdzamy, czy conajmniej 80% odszyfrowanych słóœ znajduje się w słowniku
        # co pozwala na złamanie szyfrogramu dla słów nie będąsych w słowniku z powodu np deklinacji 
        # czy casu innego niż teraźniejszy
        if sum(word in dictionary for word in words) / len(words) > 0.8:
            # Jeśli tak, zwracamy odszyfrowaną wiadomość i klucz
            return decrypted_message, key

    # Jeśli żaden klucz nie pasuje, zwracamy None
    return None, None

def main():
    start = time.time()
    # Przykładowe użycie
    encrypted_message = 'axjmrlamuhnvayum jsbkvggpzgz zqmna tjpmh'
    # encrypted_message = 'fpop whsbsilhoee'
    decrypted_message, key = break_gronsfeld(encrypted_message)
    print(f'Odszyfrowana wiadomość: {decrypted_message}, Klucz: {key}')
    print(f'Czas trwania: {time.time() - start}')

if __name__ == '__main__':
    main()

'''

Odszyfrowana wiadomość: computer science, Klucz: 312063
Czas trwania: 2.478550910949707

Odszyfrowana wiadomość:  there is no use crying over spilt milk , Klucz: 142807
Czas trwania: 3.5992813110351562
'''