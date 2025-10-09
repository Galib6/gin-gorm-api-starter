# ========== Development ==========
FROM golang:1.25-alpine AS dev

# Install git & air
RUN apk add --no-cache git && go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go.mod first for caching deps
COPY go.mod go.sum ./
RUN go mod download

# Copy all files
COPY . .

# Expose default Gin port
EXPOSE 8080

# Run Air for live reload
CMD ["air", "-c", ".air.toml"]


# ========== Production (optional) ==========
# FROM golang:1.25-alpine AS prod
# WORKDIR /app
# COPY . .
# RUN go build -o main .
# EXPOSE 8080
# CMD ["./main"]
