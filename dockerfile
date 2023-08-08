# Start from the official Go image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project into the container
COPY . .

# Build the Go application
RUN go build -o todo-api ./cmd/api/main.go

# Expose the port that your Go application listens on
EXPOSE 8080

# Run the Go application
CMD ["./todo-api"]