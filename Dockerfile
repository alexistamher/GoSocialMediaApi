# Build stage
FROM golang:1.26.5-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Install dependencies
COPY go.mod go.sum .
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/api .

COPY .env .

EXPOSE 3000

CMD ["./api"]
