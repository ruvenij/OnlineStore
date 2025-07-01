# Build stage
FROM golang:1.24 AS builder

# Set container's working directory
WORKDIR /app

# Copy Go source and dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your project (including the external/data folder)
COPY . .

# Build the binary (adjust the path if needed)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd  # or "." if main.go is here

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary and CSV files from builder
COPY --from=builder /app/main .
COPY --from=builder /app/internal/config/config.json ./config.json

EXPOSE 8080
CMD ["./main"]