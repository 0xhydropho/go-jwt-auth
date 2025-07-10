# BUILD STAGE
FROM golang:1.24-alpine AS build

RUN apk add --no-cache build-base sqlite-dev

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod tidy

COPY .env.example .env
COPY . .

ENV GOMAXPROCS=1
RUN go build -o go-jwt-auth ./cmd/api/main.go

# FINAL IMAGE
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache sqlite-libs

COPY --from=build /build/go-jwt-auth .
COPY --from=build /build/.env .

CMD ["./go-jwt-auth"]
