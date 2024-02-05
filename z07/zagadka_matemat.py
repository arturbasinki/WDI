import itertools

def solve_equation(equation):
    # Znajdź wszystkie unikalne litery w równaniu
    letters = set(equation) - set("+-*/= ")

    # Sprawdź, czy liczba liter nie przekracza 10
    if len(letters) > 10:
        print("Za dużo liter, nie można rozwiązać zagadki.")
        return

    # Stwórz słownik, który będzie mapował litery na cyfry
    mapping = dict()

    # Przypisz każdej literze cyfrę od 0 do 9, z wykluczeniem przypadków,
    # gdzie liczba nie może zaczynać się od zera w wyniku
    for digits in itertools.permutations(range(10), len(letters)):
        # Utwórz parę litera-cyfra dla każdej litery
        for letter, digit in zip(letters, digits):
            mapping[letter] = digit

        # Podziel równanie na trzy części: lewa strona, znak równości i prawa strona
        left_side, equal_sign, right_side = equation.partition("=")

        # Zamień litery na cyfry w równaniu
        numeric_equation = ''.join(str(mapping.get(char, char)) for char in equation)

        # Sprawdź, czy jakakolwiek liczba w rozwiązaniu zaczyna się od zera
        if any(num.startswith('0') for num in numeric_equation.split()):
            continue

        # Sprawdź, czy równanie jest prawdziwe
        if eval(numeric_equation.replace("=", "=="), {"__builtins__": None}):
            # Ustaw zmienną valid na True
            valid = True
            # Sprawdź, czy każda liczba w równaniu ma tyle samo cyfr co odpowiadająca jej litera
            for letter, number in zip(equation, numeric_equation):
                # Jeśli litera i cyfra są różne
                if letter != number:
                    # Sprawdź, czy litera jest w zbiorze liter
                    if letter in letters:
                        # Sprawdź, czy liczba ma tyle samo cyfr co litera
                        if len(number) != len(letter):
                            # Jeśli nie, to ustaw zmienną valid na False
                            valid = False
                            # Przerwij pętlę
                            break
            # Jeśli zmienna valid jest True, to wypisz rozwiązanie w postaci cyfrowej
            if valid:
                print(numeric_equation)


# Wczytaj równanie z klawiatury
equation = input("Podaj równanie w postaci literowej: ")
solve_equation(equation)
