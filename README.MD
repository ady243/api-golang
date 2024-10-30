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


# Authentification avec Google Cloud

gcloud auth login
gcloud config set project your-gcp-project-id

# Construire l'image Docker

docker build -t gcr.io/your-gcp-project-id/your-docker-image:latest -f Dockerfile.prod .

# Authentification avec Google Container Registry
gcloud auth configure-docker

# Pousser l'image Docker sur GCR
docker push gcr.io/your-gcp-project-id/your-docker-image:latest

# Naviguer vers le répertoire production
cd /Users/adymasivi/api-golang/production

# Créer un fichier .env
touch .env

# Ajouter les variables d'environnement dans le fichier .env
# (Ouvrez le fichier .env dans un éditeur de texte et ajoutez les variables)

# Installer les dépendances nécessaires
npm install cdktf constructs @cdktf/provider-google dotenv

# Déployer l'infrastructure avec CDKTF
cdktf deploy