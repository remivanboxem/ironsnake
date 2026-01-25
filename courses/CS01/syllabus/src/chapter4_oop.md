# Programmation Orientée Objet

## Introduction à la POO

La programmation orientée objet (POO) est un paradigme de programmation basé sur le concept d'objets, qui peuvent contenir des données (attributs) et du code (méthodes).

## Les Quatre Piliers de la POO

### 1. Encapsulation

L'encapsulation consiste à regrouper les données et les méthodes qui les manipulent dans une même unité (la classe), et à contrôler l'accès à ces données.

```py,playground,editable,docker-id=python3
class BankAccount:
    def __init__(self, owner, balance=0):
        self.owner = owner
        self._balance = balance  # Attribut "privé" par convention
    
    def deposit(self, amount):
        """Dépose de l'argent sur le compte"""
        if amount > 0:
            self._balance += amount
            return True
        return False
    
    def withdraw(self, amount):
        """Retire de l'argent du compte"""
        if 0 < amount <= self._balance:
            self._balance -= amount
            return True
        return False
    
    def get_balance(self):
        """Retourne le solde (accès contrôlé)"""
        return self._balance

# Test
account = BankAccount("Alice", 1000)
account.deposit(500)
account.withdraw(200)
print(f"Solde: {account.get_balance()}€")
```

### 2. Héritage

L'héritage permet de créer une nouvelle classe basée sur une classe existante, héritant de ses attributs et méthodes.

```py,playground,editable,docker-id=python3
class Animal:
    def __init__(self, name):
        self.name = name
    
    def speak(self):
        return "Un son quelconque"
    
    def move(self):
        return f"{self.name} se déplace"

class Dog(Animal):
    def speak(self):
        """Surcharge de la méthode speak"""
        return "Woof!"
    
    def fetch(self):
        """Méthode spécifique aux chiens"""
        return f"{self.name} rapporte la balle"

class Cat(Animal):
    def speak(self):
        return "Meow!"
    
    def scratch(self):
        return f"{self.name} fait ses griffes"

# Test
dog = Dog("Rex")
cat = Cat("Whiskers")

print(f"{dog.name} dit: {dog.speak()}")
print(f"{cat.name} dit: {cat.speak()}")
print(dog.fetch())
print(dog.move())
```

### 3. Polymorphisme

Le polymorphisme permet d'utiliser des objets de différentes classes à travers une interface commune.

```py,playground,editable,docker-id=python3
class Shape:
    def area(self):
        raise NotImplementedError("Cette méthode doit être implémentée")
    
    def perimeter(self):
        raise NotImplementedError("Cette méthode doit être implémentée")

class Rectangle(Shape):
    def __init__(self, width, height):
        self.width = width
        self.height = height
    
    def area(self):
        return self.width * self.height
    
    def perimeter(self):
        return 2 * (self.width + self.height)

class Circle(Shape):
    def __init__(self, radius):
        self.radius = radius
    
    def area(self):
        return 3.14159 * self.radius ** 2
    
    def perimeter(self):
        return 2 * 3.14159 * self.radius

# Polymorphisme en action
shapes = [
    Rectangle(5, 3),
    Circle(4),
    Rectangle(2, 7)
]

print("Calcul des aires:")
for shape in shapes:
    print(f"  {shape.__class__.__name__}: {shape.area():.2f}")
```

### 4. Abstraction

L'abstraction consiste à exposer uniquement les détails essentiels et à cacher les détails d'implémentation.

```py,playground,editable,docker-id=python3
from abc import ABC, abstractmethod

class Vehicle(ABC):
    def __init__(self, brand, model):
        self.brand = brand
        self.model = model
    
    @abstractmethod
    def start_engine(self):
        """Méthode abstraite - doit être implémentée par les sous-classes"""
        pass
    
    @abstractmethod
    def stop_engine(self):
        pass
    
    def display_info(self):
        """Méthode concrète disponible pour toutes les sous-classes"""
        return f"{self.brand} {self.model}"

class Car(Vehicle):
    def start_engine(self):
        return f"{self.display_info()}: Moteur démarré (clé)"
    
    def stop_engine(self):
        return f"{self.display_info()}: Moteur arrêté"

class ElectricCar(Vehicle):
    def start_engine(self):
        return f"{self.display_info()}: Système électrique activé"
    
    def stop_engine(self):
        return f"{self.display_info()}: Système électrique désactivé"

# Test
car = Car("Toyota", "Corolla")
electric = ElectricCar("Tesla", "Model 3")

print(car.start_engine())
print(electric.start_engine())
```

## Concepts Avancés

### Méthodes de Classe et Méthodes Statiques

```py,playground,editable,docker-id=python3
class MathOperations:
    pi = 3.14159
    
    def __init__(self, value):
        self.value = value
    
    @classmethod
    def from_string(cls, string_value):
        """Factory method - retourne une instance de la classe"""
        return cls(float(string_value))
    
    @staticmethod
    def add(a, b):
        """Méthode statique - pas d'accès à self ou cls"""
        return a + b
    
    def multiply_by_pi(self):
        """Méthode d'instance normale"""
        return self.value * MathOperations.pi

# Test
obj1 = MathOperations(5)
obj2 = MathOperations.from_string("10")

print(f"Addition statique: {MathOperations.add(3, 7)}")
print(f"Multiplication par pi: {obj1.multiply_by_pi():.2f}")
print(f"Objet créé depuis string: {obj2.value}")
```

### Properties (Getters/Setters)

```py,playground,editable,docker-id=python3
class Temperature:
    def __init__(self, celsius):
        self._celsius = celsius
    
    @property
    def celsius(self):
        """Getter pour celsius"""
        return self._celsius
    
    @celsius.setter
    def celsius(self, value):
        """Setter pour celsius avec validation"""
        if value < -273.15:
            raise ValueError("Température inférieure au zéro absolu!")
        self._celsius = value
    
    @property
    def fahrenheit(self):
        """Propriété calculée"""
        return self._celsius * 9/5 + 32
    
    @fahrenheit.setter
    def fahrenheit(self, value):
        self._celsius = (value - 32) * 5/9

# Test
temp = Temperature(25)
print(f"25°C = {temp.fahrenheit:.1f}°F")

temp.fahrenheit = 68
print(f"68°F = {temp.celsius:.1f}°C")
```

### Héritage Multiple

```py,playground,editable,docker-id=python3
class Flyable:
    def fly(self):
        return "Je peux voler!"

class Swimmable:
    def swim(self):
        return "Je peux nager!"

class Duck(Animal, Flyable, Swimmable):
    def __init__(self, name):
        Animal.__init__(self, name)
    
    def speak(self):
        return "Quack!"

# Test
duck = Duck("Donald")
print(duck.speak())
print(duck.fly())
print(duck.swim())
print(duck.move())
```

{{#task task08}}

<task id="task08"/>

## Design Patterns Courants

### Singleton

```py,playground,editable,docker-id=python3
class Singleton:
    _instance = None
    
    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
        return cls._instance

# Test
s1 = Singleton()
s2 = Singleton()
print(f"Même instance: {s1 is s2}")
```

## Exercices

1. Créez une hiérarchie de classes pour un système de gestion de bibliothèque
2. Implémentez le pattern Factory pour créer différents types de véhicules
3. Créez une classe avec validation des données utilisant properties
