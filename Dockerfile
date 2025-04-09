# Use a minimal base image for Go
FROM golang:1.24.1 AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/countries-dashboard-service/main.go

# --- Use a smaller image for the final binary ---
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/server .

# Set the default command
CMD ["./server"]
