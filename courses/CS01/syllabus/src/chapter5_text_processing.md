# Traitement de Texte et Analyse de Données

## Manipulation de Chaînes de Caractères

Les chaînes de caractères sont des séquences immuables en Python. Leur manipulation est essentielle pour le traitement de texte.

### Opérations de Base

```py,playground,editable,docker-id=python3
text = "  Hello, World! Welcome to Python.  "

# Nettoyage
print(f"Original: '{text}'")
print(f"Strip: '{text.strip()}'")
print(f"Lower: '{text.lower()}'")
print(f"Upper: '{text.upper()}'")

# Recherche
print(f"\nContient 'Python': {'Python' in text}")
print(f"Index de 'World': {text.find('World')}")
print(f"Commence par '  Hello': {text.startswith('  Hello')}")

# Remplacement
print(f"\nRemplacement: {text.replace('World', 'Universe')}")

# Division et jonction
words = text.strip().split()
print(f"Mots: {words}")
print(f"Jointure: {'-'.join(words)}")
```

### Formatage de Chaînes

```py,playground,editable,docker-id=python3
name = "Alice"
age = 25
height = 1.68

# f-strings (Python 3.6+)
print(f"Je m'appelle {name}, j'ai {age} ans")
print(f"Ma taille est {height:.2f}m")

# Format
print("Je m'appelle {}, j'ai {} ans".format(name, age))
print("Nom: {n}, Age: {a}".format(n=name, a=age))

# Alignement et padding
print(f"{'Gauche':<10}|{'Centre':^10}|{'Droite':>10}")
print(f"{42:05d}")  # Padding avec des zéros
```

## Expressions Régulières

Les expressions régulières permettent de rechercher et manipuler des motifs dans du texte.

```py,playground,editable,docker-id=python3
import re

text = "Contact: alice@example.com ou bob@test.org. Tél: 01-23-45-67-89"

# Recherche simple
emails = re.findall(r'\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b', text)
print(f"Emails trouvés: {emails}")

# Recherche de numéros de téléphone
phones = re.findall(r'\d{2}-\d{2}-\d{2}-\d{2}-\d{2}', text)
print(f"Téléphones: {phones}")

# Remplacement
censored = re.sub(r'\b\d{2}-\d{2}-\d{2}-\d{2}-\d{2}\b', '[REDACTED]', text)
print(f"Censuré: {censored}")

# Groupes de capture
match = re.search(r'(\w+)@(\w+\.\w+)', text)
if match:
    print(f"Utilisateur: {match.group(1)}, Domaine: {match.group(2)}")
```

## Analyse de Fréquence

### Comptage de Mots

```py,playground,editable,docker-id=python3
import string
from collections import Counter

text = """
Python is an interpreted, high-level programming language.
Python's design philosophy emphasizes code readability.
Python is dynamically typed and garbage-collected.
"""

def analyze_text(text):
    # Nettoyage
    text = text.lower()
    text = text.translate(str.maketrans('', '', string.punctuation))
    
    # Tokenisation
    words = text.split()
    
    # Statistiques
    word_count = len(words)
    unique_words = len(set(words))
    
    # Fréquence
    frequency = Counter(words)
    most_common = frequency.most_common(5)
    
    return {
        'total_words': word_count,
        'unique_words': unique_words,
        'most_common': most_common,
        'frequency': frequency
    }

# Analyse
stats = analyze_text(text)
print(f"Total de mots: {stats['total_words']}")
print(f"Mots uniques: {stats['unique_words']}")
print(f"\nMots les plus fréquents:")
for word, count in stats['most_common']:
    print(f"  '{word}': {count}")
```

{{#task task06}}

<task id="task06"/>

### Analyse Avancée

```py,playground,editable,docker-id=python3
def word_frequency_distribution(text):
    """Analyse la distribution de fréquence des mots"""
    words = text.lower().split()
    words = [w.strip(string.punctuation) for w in words]
    
    freq = {}
    for word in words:
        if word:
            freq[word] = freq.get(word, 0) + 1
    
    # Tri par fréquence décroissante
    sorted_freq = sorted(freq.items(), key=lambda x: x[1], reverse=True)
    
    return sorted_freq

def calculate_metrics(text):
    """Calcule diverses métriques sur le texte"""
    sentences = text.split('.')
    words = text.split()
    
    return {
        'sentences': len([s for s in sentences if s.strip()]),
        'words': len(words),
        'characters': len(text),
        'avg_word_length': sum(len(w) for w in words) / len(words) if words else 0,
        'avg_sentence_length': len(words) / len(sentences) if sentences else 0
    }

# Test
sample = "Python is great. Python is powerful. I love Python programming."
freq_dist = word_frequency_distribution(sample)
metrics = calculate_metrics(sample)

print("Distribution de fréquence:")
for word, count in freq_dist[:5]:
    print(f"  {word}: {count}")

print(f"\nMétriques:")
print(f"  Phrases: {metrics['sentences']}")
print(f"  Mots: {metrics['words']}")
print(f"  Longueur moyenne des mots: {metrics['avg_word_length']:.2f}")
```

## Traitement de Fichiers Texte

### Lecture et Écriture

```py,playground,editable,docker-id=python3
# Écriture dans un fichier
content = """Ligne 1: Introduction
Ligne 2: Développement
Ligne 3: Conclusion"""

# Simulation d'écriture (en mémoire pour l'exemple)
lines = content.split('\n')
print("Contenu écrit:")
for line in lines:
    print(f"  {line}")

# Lecture ligne par ligne
print("\nLecture ligne par ligne:")
for i, line in enumerate(lines, 1):
    print(f"  Ligne {i}: {line}")

# Lecture avec traitement
print("\nLignes contenant 'Ligne':")
filtered = [line for line in lines if 'Ligne' in line]
for line in filtered:
    print(f"  {line}")
```

### Parsing CSV

```py,playground,editable,docker-id=python3
import csv
from io import StringIO

# Données CSV simulées
csv_data = """nom,age,ville
Alice,25,Paris
Bob,30,Lyon
Charlie,28,Marseille"""

# Lecture
csvfile = StringIO(csv_data)
reader = csv.DictReader(csvfile)

print("Données CSV:")
for row in reader:
    print(f"  {row['nom']} ({row['age']} ans) - {row['ville']}")

# Statistiques
csvfile.seek(0)
reader = csv.DictReader(csvfile)
ages = [int(row['age']) for row in reader]
print(f"\nAge moyen: {sum(ages) / len(ages):.1f} ans")
```

## Traitement JSON

```py,playground,editable,docker-id=python3
import json

# Données structurées
data = {
    "cours": "CS01",
    "titre": "Introduction to Computer Science",
    "etudiants": [
        {"nom": "Alice", "note": 85},
        {"nom": "Bob", "note": 92},
        {"nom": "Charlie", "note": 78}
    ]
}

# Sérialisation
json_string = json.dumps(data, indent=2)
print("JSON:")
print(json_string)

# Désérialisation
parsed = json.loads(json_string)
print(f"\nCours: {parsed['cours']}")
print(f"Nombre d'étudiants: {len(parsed['etudiants'])}")

# Calcul de la moyenne
notes = [e['note'] for e in parsed['etudiants']]
print(f"Note moyenne: {sum(notes) / len(notes):.2f}")
```

## Exercices Pratiques

1. Créez un analyseur de log qui extrait les erreurs et avertissements
2. Implémentez un correcteur orthographique simple
3. Créez un générateur de résumés de texte
4. Développez un parser pour un format de configuration personnalisé
