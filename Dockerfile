# Use a Go version that supports net/netip (Go 1.18 or later)
FROM golang:1.19

# Define the working directory inside the container
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Expose the port the app will run on
EXPOSE 3003

# Command to start the application
CMD ["./main"]
