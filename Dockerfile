# Use a minimal Go base image
FROM golang:1.24

# Set working dir
WORKDIR /app

# Copy go mod and download deps
COPY go.mod ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the Go binary
RUN go build -o main ./cmd/api

# Run the binary
CMD ["./main"]
