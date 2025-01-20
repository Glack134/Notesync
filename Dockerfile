# Используем базовый образ с Go для архитектуры amd64
FROM --platform=linux/amd64 golang:1.20-buster

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы в контейнер
COPY ./ ./

# Устанавливаем клиент PostgreSQL
RUN apt-get update && apt-get -y install postgresql-client

# Скачиваем зависимости
RUN go mod download

# Компилируем приложение для Linux amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notesync .

# Устанавливаем права на выполнение
RUN chmod +x notesync

# Указываем команду для запуска приложения
CMD ["./notesync"]
