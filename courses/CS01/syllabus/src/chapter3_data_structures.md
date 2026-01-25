# Structures de Données

## Introduction

Les structures de données sont des moyens d'organiser et de stocker des données de manière efficace. Le choix de la bonne structure de données est crucial pour la performance d'un programme.

## Structures Linéaires

### Listes (Arrays)

Les listes sont des collections ordonnées d'éléments accessibles par index.

```py,playground,editable,docker-id=python3
# Création et manipulation de listes
numbers = [1, 2, 3, 4, 5]

# Accès par index - O(1)
print(f"Premier élément: {numbers[0]}")

# Ajout à la fin - O(1) amorti
numbers.append(6)
print(f"Après append: {numbers}")

# Insertion au début - O(n)
numbers.insert(0, 0)
print(f"Après insert: {numbers}")

# Recherche - O(n)
index = numbers.index(3)
print(f"Index de 3: {index}")
```

### Listes Chaînées (Linked Lists)

Une liste chaînée est une collection d'éléments où chaque élément pointe vers le suivant.

```py,playground,editable,docker-id=python3
class Node:
    def __init__(self, data):
        self.data = data
        self.next = None

class LinkedList:
    def __init__(self):
        self.head = None
    
    def append(self, data):
        """Ajoute un élément à la fin - O(n)"""
        new_node = Node(data)
        
        if not self.head:
            self.head = new_node
            return
        
        current = self.head
        while current.next:
            current = current.next
        current.next = new_node
    
    def prepend(self, data):
        """Ajoute un élément au début - O(1)"""
        new_node = Node(data)
        new_node.next = self.head
        self.head = new_node
    
    def display(self):
        """Affiche la liste"""
        elements = []
        current = self.head
        while current:
            elements.append(str(current.data))
            current = current.next
        return " -> ".join(elements)

# Test
ll = LinkedList()
ll.append(1)
ll.append(2)
ll.append(3)
ll.prepend(0)
print(f"Liste chaînée: {ll.display()}")
```

### Piles (Stacks)

Une pile suit le principe LIFO (Last In, First Out).

```py,playground,editable,docker-id=python3
class Stack:
    def __init__(self):
        self.items = []
    
    def push(self, item):
        """Ajoute un élément au sommet - O(1)"""
        self.items.append(item)
    
    def pop(self):
        """Retire et retourne l'élément au sommet - O(1)"""
        if not self.is_empty():
            return self.items.pop()
        return None
    
    def peek(self):
        """Retourne l'élément au sommet sans le retirer - O(1)"""
        if not self.is_empty():
            return self.items[-1]
        return None
    
    def is_empty(self):
        return len(self.items) == 0
    
    def size(self):
        return len(self.items)

# Test
stack = Stack()
stack.push(1)
stack.push(2)
stack.push(3)
print(f"Sommet: {stack.peek()}")
print(f"Pop: {stack.pop()}")
print(f"Taille: {stack.size()}")
```

### Files (Queues)

Une file suit le principe FIFO (First In, First Out).

```py,playground,editable,docker-id=python3
from collections import deque

class Queue:
    def __init__(self):
        self.items = deque()
    
    def enqueue(self, item):
        """Ajoute un élément à la fin - O(1)"""
        self.items.append(item)
    
    def dequeue(self):
        """Retire et retourne le premier élément - O(1)"""
        if not self.is_empty():
            return self.items.popleft()
        return None
    
    def front(self):
        """Retourne le premier élément sans le retirer - O(1)"""
        if not self.is_empty():
            return self.items[0]
        return None
    
    def is_empty(self):
        return len(self.items) == 0
    
    def size(self):
        return len(self.items)

# Test
queue = Queue()
queue.enqueue("Premier")
queue.enqueue("Deuxième")
queue.enqueue("Troisième")
print(f"Front: {queue.front()}")
print(f"Dequeue: {queue.dequeue()}")
print(f"Nouveau front: {queue.front()}")
```

## Structures Non-Linéaires

### Arbres Binaires

Un arbre binaire est une structure hiérarchique où chaque nœud a au plus deux enfants.

```py,playground,editable,docker-id=python3
class TreeNode:
    def __init__(self, value):
        self.value = value
        self.left = None
        self.right = None

class BinaryTree:
    def __init__(self, root_value):
        self.root = TreeNode(root_value)
    
    def inorder_traversal(self, node, result=None):
        """Parcours en ordre (gauche-racine-droite)"""
        if result is None:
            result = []
        
        if node:
            self.inorder_traversal(node.left, result)
            result.append(node.value)
            self.inorder_traversal(node.right, result)
        
        return result
    
    def preorder_traversal(self, node, result=None):
        """Parcours préfixé (racine-gauche-droite)"""
        if result is None:
            result = []
        
        if node:
            result.append(node.value)
            self.preorder_traversal(node.left, result)
            self.preorder_traversal(node.right, result)
        
        return result

# Test
tree = BinaryTree(1)
tree.root.left = TreeNode(2)
tree.root.right = TreeNode(3)
tree.root.left.left = TreeNode(4)
tree.root.left.right = TreeNode(5)

print(f"Inorder: {tree.inorder_traversal(tree.root)}")
print(f"Preorder: {tree.preorder_traversal(tree.root)}")
```

### Tables de Hachage (Hash Tables)

Les dictionnaires Python sont implémentés comme des tables de hachage.

```py,playground,editable,docker-id=python3
# Dictionnaires Python - Table de hachage
hash_table = {}

# Insertion - O(1) en moyenne
hash_table["nom"] = "Alice"
hash_table["age"] = 25
hash_table["ville"] = "Paris"

# Accès - O(1) en moyenne
print(f"Nom: {hash_table.get('nom')}")

# Suppression - O(1) en moyenne
del hash_table["ville"]

# Itération
print("\nContenu:")
for key, value in hash_table.items():
    print(f"  {key}: {value}")

# Vérification d'existence - O(1)
print(f"\n'age' existe: {'age' in hash_table}")
```

{{#task task05}}

<task id="task05"/>

## Exercices

1. Implémentez une file de priorité
2. Créez un arbre binaire de recherche (BST)
3. Implémentez une table de hachage simple avec gestion des collisions
