# Этап 1: Сборка приложения
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем зависимости для сборки (protoc, если нужен для gRPC)
RUN apk add --no-cache protoc

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь проект (все файлы и папки)
COPY . .

# Генерируем gRPC код (если proto-файлы изменились, на случай если они не закоммичены)
RUN go generate ./...

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./main.go

# Этап 2: Создание минимального итогового образа
FROM alpine:latest

# Добавляем ca-certificates для HTTPS/gRPC и tzdata для временных зон, если нужно
# RUN apk add --no-cache ca-certificates tzdata

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/app .

# Копируем config.json
COPY --from=builder /app/config.json .

# Открываем порт (8080 для WebSocket, 50051 для gRPC)
EXPOSE 8080 50051

# Запускаем приложение
CMD ["./app"]