# Step 1: Build the Go application
FROM golang:1.20 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Copy the .env file to the working directory
COPY .env ../

# Build the Go app with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Step 2: Create a minimal image to run the built Go binary
FROM alpine:latest

# Install certificates for HTTPS communication
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Ensure binary has executable permissions
RUN chmod +x ./main

# Expose the application on port 5000
EXPOSE 5000

# Command to run the executable
CMD ["./main"]
