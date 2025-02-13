# Используем официальный образ Go как базовый
FROM golang:1.23.6-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем исполняемый файл
RUN go build -o server ./cmd/server/main.go

# Используем минимальный образ Alpine для финального контейнера
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем исполняемый файл из предыдущего этапа
COPY --from=builder /app/server .

# Копируем файл конфигурации
COPY --from=builder /app/config.yml . 
COPY --from=builder /app/migrations ./migrations
ENV ENV="docker"

# Указываем команду для запуска контейнера
CMD ["./server"]