# ---------- Stage 1: Build ----------
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Install git (required for go mod if private dependencies exist)
RUN apk add --no-cache git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
# RUN CGO_ENABLED=0  go build -o server ./cmd/students-api/main.go
RUN CGO_ENABLED=0  go build -ldflags="-s -w" -o server ./cmd/students-api/main.go


# ---------- Stage 2: Run ----------
FROM scratch

WORKDIR /app

# Install required system dependencies if any (optional)
# RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/server .

# Copy config if needed inside image
COPY --from=builder /app/config ./config

EXPOSE 9090

CMD ["./server", "--config", "./config/mysql.yaml"]
