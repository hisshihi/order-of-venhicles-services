# Build stage
# Стадия сборки
FROM golang:1.23.6-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o service cmd/main.go

# Run stage
FROM alpine:3.21
WORKDIR /app

# Копируем бинарный файл из builder stage
COPY --from=builder /app/service .
# Копируем файл конфигурации в корневую директорию приложения
COPY --from=builder /app/cmd/app.env ./

# Открываем порт 8080
EXPOSE 8080

# Запускаем бинарный файл
CMD ["./service"]