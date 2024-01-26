'''

Program wczytuje plik tekstowy, znajduje wszystkie słowa o długości 5 lub więcej liter,
zlicza wystąpienia każdego słowa i wypisuje 20 najczęściej występujących słów.
Słowa są zliczane bez względu na wielkość liter.
Program nie uwzględnia krótkich słów (poniżej 5 liter).
'''

from collections import Counter
import re

# Funkcja top_20, która jako dane wejściowe pobiera nazwę pliku tekstowego
def top_20(filename):
    # Pobranie tekstu z pliku
    with open(filename, 'r', encoding='utf-8') as f:
        # Konwertuj wszystkie znaki na małe litery 
        text = f.read().lower()
    # Znajdź słowa zawierające co najmniej 5 lub więcej znaków
    # używając biblioteki re (Regular Expressions)
    words = re.findall(r'\b\w{5,}\b', text)
    # Oblicz rozkład słów
    counter = Counter(words)
    # Zwróć 20 najcześciej występujących słów
    return counter.most_common(20)

def main():
    poz = 1
    # Wypisz 20 najcześciej występujących słów z pliku tekstowego Pan_Tadeusz.txt
    for slowo, ilosc in top_20('Pan_Tadeusz.txt'):
        print(f'{poz}. Słowo: "{slowo}", występuje {ilosc} razy.')
        poz += 1

if __name__ == '__main__':
    main()

'''

1. Słowo: "rzekł", występuje 155 razy.
2. Słowo: "tylko", występuje 149 razy.
3. Słowo: "hrabia", występuje 129 razy.
4. Słowo: "sędzia", występuje 127 razy.
5. Słowo: "tadeusz", występuje 107 razy.
6. Słowo: "przed", występuje 102 razy.
7. Słowo: "jeszcze", występuje 101 razy.
8. Słowo: "przez", występuje 97 razy.
9. Słowo: "gdzie", występuje 94 razy.
10. Słowo: "wszyscy", występuje 90 razy.
11. Słowo: "wojski", występuje 89 razy.
12. Słowo: "potem", występuje 86 razy.
13. Słowo: "teraz", występuje 83 razy.
14. Słowo: "niech", występuje 76 razy.
15. Słowo: "kiedy", występuje 74 razy.
16. Słowo: "jeśli", występuje 73 razy.
17. Słowo: "który", występuje 72 razy.
18. Słowo: "nawet", występuje 71 razy.
19. Słowo: "znowu", występuje 71 razy.
20. Słowo: "wszystko", występuje 67 razy.
'''