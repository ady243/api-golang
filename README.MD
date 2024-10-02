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