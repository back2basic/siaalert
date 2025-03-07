# Build Stage
FROM golang:1.23.5-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app
RUN go build -o scanner

# Run Stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the build stage
COPY --from=build /app/scanner .

# Start the Go app
CMD ["./scanner"]

# Expose the port the app runs on
# EXPOSE 8000
