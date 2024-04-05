# Используем образ Golang для сборки приложения
FROM golang:1.20 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем все файлы с исходным кодом в контейнер
COPY . .

# Устанавливаем CGO_ENABLED и устанавливаем go-sqlite3 внутри контейнера
RUN CGO_ENABLED=1 go install github.com/mattn/go-sqlite3

# Собираем исполняемый файл из исходного кода с включенным CGO
RUN CGO_ENABLED=1 go build -o /app_linux ./cmd

# Копируем собранный исполняемый файл из контейнера сборки в рабочую директорию
RUN cp /app_linux /app/app_linux

# Копируем необходимые файлы
COPY config /app/config
COPY web /app/web

# Указываем порт, который нужно пробросить
EXPOSE 7540

# Определение тома
VOLUME /app/db

# Запускаем приложение при старте контейнера
CMD ["/app/app_linux"]
