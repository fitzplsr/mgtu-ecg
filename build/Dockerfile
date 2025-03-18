# Используем мультиархитектурный образ для сборки
FROM --platform=$BUILDPLATFORM golang:1.24.0-alpine as builder

# Копируем исходный код
COPY . /github.com/fitzplsr/mgtu-ecg/
WORKDIR /github.com/fitzplsr/mgtu-ecg/

# Скачиваем зависимости
RUN go mod download
RUN go clean --modcache

# Собираем приложение с учетом целевой ОС и архитектуры
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -mod=readonly -o ./.bin ./cmd/main/main.go

# Используем мультиархитектурный образ для запуска
FROM --platform=$TARGETPLATFORM alpine AS runner

WORKDIR /docker-main/

# Копируем собранный бинарник и конфиги
COPY --from=builder /github.com/fitzplsr/mgtu-ecg/.bin .
COPY --from=builder /github.com/fitzplsr/mgtu-ecg/config config/

# Открываем порт и запускаем приложение
EXPOSE 4000
ENTRYPOINT ["./.bin"]