# Importujemy potrzebne moduły
import sympy
import numpy as np
import matplotlib.pyplot as plt

# Pytamy użytkownika o wzór funkcji i przedział x
funkcja = input("Podaj wzór funkcji: ")
x_min = eval(input("Podaj dolną granicę przedziału x: "))
x_max = eval(input("Podaj górną granicę przedziału x: "))

# Definiujemy symbol x
x = sympy.Symbol('x')

# Próbujemy skonwertować wzór funkcji na wyrażenie sympy
try:
    expr = sympy.sympify(funkcja)
except Exception as e:
    # Jeśli wystąpi błąd, informujemy użytkownika i kończymy program
    print(f"Błąd: {e}")
    exit()

# Konwertujemy wyrażenie sympy na funkcję, używając biblioteki numpy
f = sympy.lambdify(x, expr, "numpy")

# Tworzymy tablicę wartości x z krokiem 0.01
x = np.arange(x_min, x_max + 0.01, 0.01)

# Obliczamy wartości y dla każdego x, używając funkcji f i funkcji np.nan_to_num
y = np.nan_to_num(f(x))

# Rysujemy wykres funkcji, używając funkcji plot
plt.plot(x, y, "b-")

# Dodajemy etykiety osi i tytuł wykresu
plt.xlabel("x")
plt.ylabel("y")
plt.title(f"y = {funkcja}")

# Ustawiamy skalę osi y, używając funkcji autoscale
plt.autoscale(axis="y")

# Pokazujemy wykres na ekranie
plt.show()
