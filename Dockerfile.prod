# Étape 1 : Construction de l'application dans une image Golang
FROM golang:1.22 AS builder

# Définir le répertoire de travail à l'intérieur du conteneur
WORKDIR /app

# Copier les fichiers de dépendances (go.mod et go.sum) dans l'image
COPY go.mod go.sum ./

# Télécharger les dépendances nécessaires
RUN go mod download

# Installer swag pour générer la documentation Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Ajouter le répertoire GOPATH/bin au PATH
ENV PATH="/go/bin:${PATH}"

# Copier tout le code source de l'application
COPY . .

# Générer la documentation Swagger
RUN swag init

# Compiler l'application en mode production
# CGO_ENABLED=0 permet de désactiver les dépendances au système pour rendre l'image plus légère
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/api ./main.go

# Étape 2 : Créer une image minimale pour exécuter l'application
FROM alpine:latest

# Installer les certificats SSL pour que Go puisse effectuer des requêtes HTTPS
RUN apk --no-cache add ca-certificates

# Définir le répertoire de travail
WORKDIR /root/

# Copier l'exécutable construit depuis l'image builder
COPY --from=builder /app/api .

# Copier la documentation Swagger générée depuis l'image builder
COPY --from=builder /app/docs ./docs

# Lister les fichiers dans le répertoire docs pour vérifier qu'ils sont bien copiés
RUN ls -la ./docs

# Exposer le port sur lequel l'application tourne
EXPOSE 3003

# Démarrer l'application
CMD ["./api"]