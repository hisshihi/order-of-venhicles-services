# Build stage
# Стадия сборки
FROM golang:1.23.6-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o service cmd/main.go

# Используем официальный образ Go
FROM golang:1.23.6-alpine3.21

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем текущий каталог в образ
COPY . .

# Собираем бинарный файл
RUN go build -o service cmd/main.go

# Run stage
FROM alpine:3.21

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из builder stage
COPY --from=builder /app/service .
COPY cmd/app.env ./cmd

# Открываем порт 8080
EXPOSE 8080

# Запускаем бинарный файл
CMD ["./service"]