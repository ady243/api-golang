# Use a Go version that supports Go 1.22
FROM golang:1.22

# Define the working directory inside the container
WORKDIR /app

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

# Copy dependency files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code, wait-for-it.sh, and air.conf
COPY . .


#get all depences
RUN go mod tidy

# Expose the port the app will run on
EXPOSE 3003

# Command to start the application
CMD ["air"]