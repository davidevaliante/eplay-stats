# Start from the official Golang base image
FROM golang:1.20-alpine

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files for dependency management
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the Go app
CMD ["./main"]
