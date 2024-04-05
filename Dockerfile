FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 go install github.com/mattn/go-sqlite3
RUN CGO_ENABLED=1 go build -o /app_linux ./cmd
RUN cp /app_linux /app/app_linux

COPY config /app/config
COPY web /app/web
VOLUME /app/db

CMD ["/app/app_linux"]
