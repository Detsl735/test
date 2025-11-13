FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go build -o app ./cmd/main.go


FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose

ENV APP_PORT=:8080

CMD ["./app"]
