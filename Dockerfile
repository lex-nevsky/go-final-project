# вместо убунту возмём alpine, он легче
FROM golang:1.26-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# режем образ до минимума, выкидываем дебаг
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

# фиксируем сборку
FROM alpine:3.23

RUN apk --no-cache add ca-certificates

WORKDIR /app

# копируем бинарник и фронтенд
COPY --from=builder /app/server .
COPY --from=builder /app/web ./web

# порт сервера (по умолчанию 7540)
EXPOSE 7540

# переменные окружения (можно переопределить при запуске)
ENV TODO_DBFILE=/app/scheduler.db
ENV TODO_PASSWORD=""
ENV TODO_PORT=""
ENV JWT_SECRET=""


CMD ["./server"]
