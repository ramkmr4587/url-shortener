# Start from the official Golang image
FROM golang:1.24.2

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o url-shortener .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./url-shortener"]
