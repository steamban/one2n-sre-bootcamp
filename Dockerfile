# Stage 1: Build
FROM golang:1.26-alpine@sha256:c2a1f7b2095d046ae14b286b18413a05bb82c9bca9b25fe7ff5efef0f0826166 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o migrate cmd/migrate/main.go

# Stage 2: Runtime
FROM alpine:3.19@sha256:6baf43584bcb78f2e5847d1de515f23499913ac9f12bdf834811a3145eb11ca1

RUN addgroup appgroup && adduser -G appgroup -s /bin/sh -D appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /build/api .
COPY --from=builder --chown=appuser:appgroup /build/migrate .
COPY --from=builder --chown=appuser:appgroup /build/migrations ./migrations

USER appuser

EXPOSE 8080

CMD ["./api"]   
