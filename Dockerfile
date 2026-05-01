FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# Build both binaries
RUN go build -o api ./cmd/api
RUN go build -o worker ./cmd/worker

# Final lightweight image
FROM alpine:latest

WORKDIR /root/

# Copy both binaries
COPY --from=build /app/api .
COPY --from=build /app/worker .

# Copy migrations if needed
COPY --from=build /app/migrations ./migrations

EXPOSE 8080

# Default (can be overridden by docker-compose)
CMD ["./api"]