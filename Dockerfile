# Build the frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/web

# Copy package files
COPY web/package.json web/package-lock.json* ./

# Install dependencies
RUN npm ci

# Copy source files
COPY web/ ./

# Build the frontend
RUN npm run build

# Build the Go backend
FROM golang:1.21-alpine AS go-builder

WORKDIR /app

# Install godotenv
RUN go install github.com/joho/godotenv@latest

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy server code
COPY server/ ./server/

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o waker ./server/main.go

# Final production image
FROM alpine:3.18

WORKDIR /app

# Copy built frontend from frontend-builder
COPY --from=frontend-builder /app/web/dist ./web/dist

# Copy built backend from go-builder
COPY --from=go-builder /app/waker .

# Copy .env.example for reference
COPY .env.example .

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Expose the port
EXPOSE 8080

# Run the server
CMD ["./waker"]
