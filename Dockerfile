# Use a Go version that supports net/netip (Go 1.18 or later)
FROM golang:1.19

# Define the working directory inside the container
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code and wait-for-it.sh
COPY . .

# Build the application
RUN go build -o main .

#get all depences
RUN go mod tidy

# Expose the port the app will run on
EXPOSE 3003

# Command to start the application
CMD ["./main"]