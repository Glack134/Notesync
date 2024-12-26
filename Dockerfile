# Используем образ Golang для сборки
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем все файлы в рабочую директорию
COPY . .

# Сборка приложения
RUN go build -o main ./cmd/main.go

# Используем более легкий образ для запуска
FROM ubuntu:22.04

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем исполняемый файл из образа сборки
COPY --from=builder /app/main .

# Убедитесь, что файл исполняемый
RUN chmod +x main

# Команда для запуска
CMD ["./main"]


#Start Postgresql in Docker ---> docker run --name=notesync-db -e POSTGRES_PASSWORD='password' -p "port" -d --rm postgres 