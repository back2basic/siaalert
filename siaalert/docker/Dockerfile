# Use the official Golang image as the base image
FROM golang:1.23.5-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app
RUN go build -o myapp .

# Start the Go app
CMD ["./myapp"]

# Expose the port the app runs on
# EXPOSE 8000