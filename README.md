# Forum

Un projet de forum web inspiré de Reddit, développé en Go avec une interface web moderne.

## Fonctionnalités

- Inscription et connexion utilisateur
- Création et affichage de posts avec images
- Stockage des utilisateurs et posts en SQLite
- Interface responsive avec sidebar dynamique

## Structure du projet

- `/main` : Code principal du serveur web (Go)
- `/tables` : Fonctions d'accès aux bases de données (Go)
- `/BDD` : Fichiers de base de données SQLite (`accounts.db`, `posts.db`)
- `/html` : Pages HTML principales (accueil, création de post, etc.)
- `/Login` : Pages d'inscription et de connexion
- `/css` : Feuilles de style CSS
- `/img` : Images utilisées dans l'interface

## Installation

> **Important** : Le serveur doit être lancé depuis un terminal Ubuntu (WSL) sous Windows. Si vous lancez le serveur depuis PowerShell ou CMD, certaines fonctionnalités risquent de ne pas fonctionner correctement.

1. **Cloner le dépôt**
   ```bash
   git clone <url-du-repo>
   cd Forum
   ```

2. **Initialiser les bases de données**
   ```bash
   go run initdb.go
   ```

3. **Lancer le serveur**
   ```bash
   cd main
   go run main.go
   ```

4. **Accéder à l'application**
   Ouvrir [http://localhost:8080](http://localhost:8080) dans votre navigateur.

## Dépendances

- Go 1.18+
- [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
- [github.com/google/uuid](https://github.com/google/uuid)
- [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)

## Commandes utiles

- Installer Sqlite 3:
  ```bash
  sudo apt install sqlite3
  ```
- Voir les posts dans la BDD :
  ```bash
  sqlite3 ./BDD/posts.db "SELECT * FROM posts;"
  ```
- Voir les utilisateurs :
  ```bash
  sqlite3 ./BDD/accounts.db "SELECT * FROM people;"
  ```

## Auteurs

- Projet réalisé par Kévin, Thomas et Victor

---