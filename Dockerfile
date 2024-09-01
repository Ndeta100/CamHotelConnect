# Step 1: Build the Go application
FROM golang:1.20 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

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

# Copy the entrypoint script into the container
COPY entrypoint.sh /app/entrypoint.sh

# Ensure the binary and entrypoint script have execute permissions
RUN chmod +x ./main /app/entrypoint.sh

# Expose the application on port 5000
EXPOSE 5000

# Set the entrypoint to the shell script
ENTRYPOINT ["/app/entrypoint.sh"]

# Command to run the executable
CMD ["./main"]
