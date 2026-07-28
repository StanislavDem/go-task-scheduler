# Этап сборки Go-приложения
FROM golang:alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей и загружаем их (кеширование слоев)
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник: CGO_ENABLED=0 (Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o tracker ./cmd/main.go

# Этап финального образа
FROM alpine:latest

WORKDIR /app

# Копируем скомпилированный бинарник из builder
COPY --from=builder /app/tracker .

# Копируем фронтенд
COPY web ./web

# Переменные окружения
ENV TODO_DBFILE=/dataBase/scheduler.db

# Создаём директорию под базу данных
RUN mkdir -p /dataBase

# Запускаемый порт по умолчанию
EXPOSE 7540

# Запуск приложения
CMD ["./tracker"]