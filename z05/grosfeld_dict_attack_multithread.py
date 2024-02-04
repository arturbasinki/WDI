import itertools
import logging
import multiprocessing
import time

# Ustawiamy poziom logowania na INFO, aby wyświetlać wszystkie informacje o przebiegu programu
logging.basicConfig(level=logging.INFO)

# Wczytujemy słownik
with open('slownik.txt', 'r') as f:
    # Tworzymy zbiór wszystkich słów w słowniku, usuwając białe znaki na początku i na końcu każdego słowa
    dictionary = set(word.strip() for word in f)

# Funkcja deszyfrująca
def gronsfeld_decrypt(encrypted_message, key):
    """
    Deszyfruje zaszyfrowaną wiadomość za pomocą szyfru Gronsfelda.

    :param encrypted_message: str, Wiadomość do odszyfrowania.
    :param key: str, Klucz używany do deszyfrowania.

    :return: str, Odszyfrowana wiadomość.
    """
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

# Funkcja łamiąca szyfr dla danego klucza
def break_gronsfeld_key(encrypted_message, key):
    """
    Funkcja do złamania klucza Gronsfelda dla zaszyfrowanej wiadomości.
    
    :param encrypted_message (str): Zaszyfrowana wiadomość do odszyfrowania.
    :param key (int): Klucz używany do odszyfrowania.
    :returns: tuple: Krotka zawierająca odszyfrowaną wiadomość i klucz, jeśli co najmniej 80% słów w odszyfrowanej wiadomości znajduje się w słowniku, w przeciwnym razie zwraca None.
    """
    key = str(key)
    # Deszyfrujemy wiadomość za pomocą danego klucza
    decrypted_message = gronsfeld_decrypt(encrypted_message, key)
    # Sprawdzamy, czy przynajmniej 80% słów w odszyfrowanej wiadomości jest w słowniku
    # pozwoli to na złamanie szyfrogramu dla słów nie będących w słowniku, np z powodu deklinacji
    words = decrypted_message.split()
    if sum(word in dictionary for word in words) / len(words) > 0.8:
        # Jeśli tak, zwracamy odszyfrowaną wiadomość i klucz
        return decrypted_message, key
    # Jeśli mniej niż 80% słów pasuje, zwracamy None
    return None, None

# Generator wszystkich możliwych kluczy
def generate_keys():
    # Generujemy wszystkie możliwe klucze o długości mniejszej niż 10
    for i in range(1, 10):
        for p in itertools.product('0123456789', repeat=i):
            # Zwracamy kolejny klucz
            yield ''.join(p)

# Ten fragment kodu definiuje funkcję break_gronsfeld, która wykorzystuje przetwarzanie wieloprocesorowe
# do utworzenia puli procesów. Następnie iteruje po zakresie liczb i wykorzystuje pulę do obliczenia
# 1000 kluczy jednocześnie. Jeśli zostanie znaleziony pasujący klucz, zwraca wynik.
# Jeśli żaden klucz nie pasuje, zwraca None.
def break_gronsfeld(encrypted_message):
    """
    Funkcja do złamania szyfru Gronsfelda dla podanej wiadomości.

    :param encrypted_message: Zaszyfrowana wiadomość do odszyfrowania.
    :return: Krotka zawierająca odszyfrowaną wiadomość i klucz użyty do odszyfrowania. Jeśli żaden klucz nie pasuje, zwraca (None, None).
    """
    with multiprocessing.Pool() as pool:
        for i in range(1000000):
            #Obliczamy 1000 kluczy i po wykonaniu sprawdzamy wyniki
            for result in pool.starmap(break_gronsfeld_key, list(zip((encrypted_message for _ in range(1000)), range(1000*i, 1000*(i+1))))):
                if result[0] is not None:
                    return result

    # Jeśli żaden klucz nie pasuje, zwracamy None
    return None, None

def main():
    start = time.time()
    # Przykładowe użycie
    encrypted_message = 'axjmrlamuhnvayum jsbkvggpzgz zqmna tjpmh'
    # encrypted_message = 'fpop whsbsilhoee'
    decrypted_message, key = break_gronsfeld(encrypted_message)
    print(f'Szyfrogram: {encrypted_message}')
    print(f'Tekst jawny: {decrypted_message}, Klucz: {key}')
    print(f'Czas trwania: {time.time() - start}')

if __name__ == '__main__':
    main()

'''
    
Szyfrogram: fpop whsbsilhoee
Tekst jawny: computer science, Klucz: 312063
Czas trwania: 0.6394965648651123

Szyfrogram: axjmrlamuhnvayum jsbkvggpzgz zqmna tjpmh
Tekst jawny:  there is no use crying over spilt milk , Klucz: 142807
Czas trwania: 0.5416271686553955Czas trwania: 0.5453438758850098
'''