# Start from a Debian-based image with Go installed
FROM golang:1.21.7

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/api 

# Expose port 8080 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./main"]
