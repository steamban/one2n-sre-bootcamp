# Stage 1: Build
FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o migrate cmd/migrate/main.go

# Stage 2: Runtime
FROM alpine:3.19

RUN addgroup -g 1000 appgroup && adduser -u 1000 -G appgroup -s /bin/sh -D appuser

WORKDIR /app

COPY --from=builder /build/api .
COPY --from=builder /build/migrate .
COPY --from=builder /build/migrations ./migrations

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["./api"]
