# Build stage
FROM golang:1.21.6-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN sqlc generate
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /go/bin/sqlc ./sqlc
COPY app.env .
COPY start.sh .
COPY db/migration ./migration

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]