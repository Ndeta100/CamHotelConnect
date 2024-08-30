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

# Use build arguments for environment variables
ARG HTTP_LISTEN_ADDRESS
ARG JWT_SECRET
ARG MONGO_DB_URL
ARG MONGO_DB_URL_TEST
ARG MONGO_DB_NAME
ARG CLOUDINARY_CLOUD_NAME
ARG CLOUDINARY_API_KEY
ARG CLOUDINARY_API_SECRET
ARG CLOUDINARY_UPLOAD_FOLDER

# Set environment variables in the Docker image using the build arguments
ENV HTTP_LISTEN_ADDRESS=${HTTP_LISTEN_ADDRESS}
ENV JWT_SECRET=${JWT_SECRET}
ENV MONGO_DB_URL=${MONGO_DB_URL}
ENV MONGO_DB_URL_TEST=${MONGO_DB_URL_TEST}
ENV MONGO_DB_NAME=${MONGO_DB_NAME}
ENV CLOUDINARY_CLOUD_NAME=${CLOUDINARY_CLOUD_NAME}
ENV CLOUDINARY_API_KEY=${CLOUDINARY_API_KEY}
ENV CLOUDINARY_API_SECRET=${CLOUDINARY_API_SECRET}
ENV CLOUDINARY_UPLOAD_FOLDER=${CLOUDINARY_UPLOAD_FOLDER}

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

# Ensure the binary has execute permissions (usually this is not necessary, but let's add it for safety)
RUN chmod +x ./main

# Expose the application on port 5000
EXPOSE 5000

# Command to run the executable
CMD ["./main"]
