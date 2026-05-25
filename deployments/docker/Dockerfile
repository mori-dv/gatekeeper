# syntax=docker/dockerfile:1
# ---- build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git curl

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o gatekeeper \
    ./cmd/api

# ---- runtime stage ----
FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/gatekeeper .

USER app

EXPOSE 8080

CMD ["./gatekeeper"]