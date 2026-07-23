# Этап сборки Go-приложения
FROM golang:1.26.2 AS builder

# Рабочая директория
WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник из main.go
RUN go build -o tracker ./cmd/main.go

# Этап финального образа
FROM ubuntu:latest

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/tracker .
# Копируем фронтенд
COPY web ./web

# Создаём папку для базы
RUN mkdir /dataBase

# Переменные окружения
ENV TODO_PORT=7540
ENV TODO_DBFILE=/dataBase/scheduler.db
ENV TODO_PASSWORD=12345

# Открываем порт
EXPOSE 7540

# Запуск приложения
CMD ["./tracker"]