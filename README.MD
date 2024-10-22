# Pour se connecter à la base de donnée en local 
- Host : localhost ou 127.0.0.1
- Hassword:postgres
- Hser:postgres
- Hame:postgres

# Pour se connecter à la base de donnée Adminer 

- Systme : PostgreSQL
- Host : postgres 
- Utilisateur : postgres 
- Password : postgres 
- DataBase : postgres

# Build le project 

- Docker compose buil 

# Pour run le project 

- Docker compose up 


# Pour créer un user 

- localhost:3003/api/register 
- Les information require pour créer un user sont : 
{
    "username"
    "password"
    "email"
}

# Pour se connecter et avoir le token 

- localhost:3003/api/login
- Les information pour se connecter : 
{
    "email"
    "password"
}

# Les informations à mettre dans le .env => créez un .env si vous l'avez pas 

POSTGRES_DATABASE=postgres
POSTGRES_HOST=postgres
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_PORT=5432

API_PORT=3003
DRAGONFLY_PORT=6379
DRAGONFLY_HOST=dragonfly
SECRET_KEY=secret


# NB: quand vous pushez faites attention à ne pas push les fichiez inutile

# Pour créer un match 

- Route : `/api/matches/`
- **Authorization** : Nécessite un token utilisateur dans l'en-tête `Authorization` (Bearer Token).
- Modèle de données du JSON à envoyer : 
```json
{
    "organizer_id": "string",
    "referee_id": "string", // optionnel
    "description": "string", // optionnel
    "match_date": "YYYY-MM-DD",
    "match_time": "HH:MM",
    "address": "string",
    "number_of_players": integer
}
```

# Gérer les joueurs d'un match (matchPlayers)

## **Récupérer les joueurs d'un match**
- **Route**: `/api/matchesPlayers/:match_id`
- **Méthode**: `GET`
- **Description**: Retourne la liste des joueurs associés à un match donné.

## **Ajouter un joueur à un match**
- **Route**: `/api/matchesPlayers/`
- **Méthode**: `POST`
- **Requête**:
    ```json
    {
      "match_id": "string",
      "player_id": "string"
    }
    ```
- **Description**: Associe un joueur à un match en spécifiant l'ID du match et l'ID du joueur.

## **Attribuer une équipe à un joueur (organisateur uniquement)**
- **Route**: `/api/matchesPlayers/assignTeam`
- **Méthode**: `PUT`
- **Authorization**: Seul l'organisateur du match peut effectuer cette opération.
- **Requête**:
    ```json
    {
      "match_player_id": "string",
      "team_number": 1 | 2
    }
    ```
- **Description**: Permet à l'organisateur d'assigner une équipe (1 ou 2) à un joueur pour un match spécifique.
