FROM golang:1.24.4-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Correct build command — no runtime flags here
RUN go build -o server ./cmd/students-api/main.go

EXPOSE 9090

# Pass the config file at runtime if needed
CMD ["./server", "--config", "./config/mysql.yaml"]
