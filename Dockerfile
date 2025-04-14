FROM golang:1.24.2

# Устанавливаем рабочую директорию
WORKDIR /avito-pvz

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Переходим в директорию с main.go и собираем приложение
RUN cd cmd/app && go build -o avito-pvz .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./cmd/app/avito-pvz"]