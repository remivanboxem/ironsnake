# Algorithmes et Complexité

## Introduction aux Algorithmes

Un algorithme est une suite finie et non ambiguë d'instructions permettant de résoudre un problème donné. Les algorithmes sont au cœur de l'informatique et constituent la base de tous les programmes informatiques.

## Complexité Algorithmique

### Notation Big O

La notation Big O permet d'exprimer la complexité temporelle et spatiale d'un algorithme en fonction de la taille de l'entrée `n`.

**Complexités courantes:**

- **O(1)** - Constant: L'opération prend toujours le même temps
- **O(log n)** - Logarithmique: Recherche binaire
- **O(n)** - Linéaire: Parcourir un tableau
- **O(n log n)** - Quasi-linéaire: Tri rapide, tri fusion
- **O(n²)** - Quadratique: Tri à bulles, tri par insertion
- **O(2^n)** - Exponentielle: Problèmes NP-complets

### Exemple: Analyse de Complexité

```py,playground,editable,docker-id=python3
def linear_search(arr, target):
    """Recherche linéaire - O(n)"""
    for i in range(len(arr)):
        if arr[i] == target:
            return i
    return -1

def binary_search(arr, target):
    """Recherche binaire - O(log n)"""
    left, right = 0, len(arr) - 1
    
    while left <= right:
        mid = (left + right) // 2
        if arr[mid] == target:
            return mid
        elif arr[mid] < target:
            left = mid + 1
        else:
            right = mid - 1
    
    return -1

# Test
sorted_array = [1, 3, 5, 7, 9, 11, 13, 15, 17, 19]
print(f"Linear search: {linear_search(sorted_array, 13)}")
print(f"Binary search: {binary_search(sorted_array, 13)}")
```

## Algorithmes de Tri

### Tri à Bulles (Bubble Sort)

Le tri à bulles compare des paires d'éléments adjacents et les échange s'ils sont dans le mauvais ordre.

**Complexité:** O(n²)

```py,playground,editable,docker-id=python3
def bubble_sort(arr):
    n = len(arr)
    for i in range(n):
        for j in range(0, n - i - 1):
            if arr[j] > arr[j + 1]:
                arr[j], arr[j + 1] = arr[j + 1], arr[j]
    return arr

# Test
numbers = [64, 34, 25, 12, 22, 11, 90]
print(f"Array trié: {bubble_sort(numbers.copy())}")
```

### Tri Rapide (Quick Sort)

Le tri rapide utilise la stratégie "diviser pour régner" avec un pivot.

**Complexité:** O(n log n) en moyenne

```py,playground,editable,docker-id=python3
def quick_sort(arr):
    if len(arr) <= 1:
        return arr
    
    pivot = arr[len(arr) // 2]
    left = [x for x in arr if x < pivot]
    middle = [x for x in arr if x == pivot]
    right = [x for x in arr if x > pivot]
    
    return quick_sort(left) + middle + quick_sort(right)

# Test
numbers = [64, 34, 25, 12, 22, 11, 90]
print(f"Quick sort: {quick_sort(numbers)}")
```

## Récursivité

La récursivité est une technique où une fonction s'appelle elle-même pour résoudre un problème.

### Exemple: Suite de Fibonacci

```py,playground,editable,docker-id=python3
def fibonacci(n):
    """Calcul récursif de Fibonacci - O(2^n)"""
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

def fibonacci_memo(n, memo={}):
    """Fibonacci avec mémoïsation - O(n)"""
    if n in memo:
        return memo[n]
    if n <= 1:
        return n
    memo[n] = fibonacci_memo(n-1, memo) + fibonacci_memo(n-2, memo)
    return memo[n]

# Test
print(f"Fibonacci(10): {fibonacci(10)}")
print(f"Fibonacci mémoïsé(30): {fibonacci_memo(30)}")
```

{{#task task04}}

<task id="task04"/>

## Exercices Pratiques

1. Implémentez le tri par insertion
2. Comparez les performances de différents algorithmes de tri
3. Analysez la complexité de vos algorithmes
