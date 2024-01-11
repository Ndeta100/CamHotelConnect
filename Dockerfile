FROM golang:1.20.3-alpine

# Set the working dir
WORKDIR /app

# Copy the go.mod and go.sum files to the working dir
COPY go.mod go.sum ./

# Download and install any requires Go dependencies
RUN go mod download

# Copy the entire source to the working dir
COPY . .

# Build the go application
RUN go build -o main .

# Expose the port specified by the PORT env variiable
EXPOSE 3000

# Set the entry point of the container to the executable
CMD ["./main"]