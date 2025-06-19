FROM golang:1.23.1-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/inputConvert ./cmd/inputConvert/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/inputConvert /app/inputConvert

CMD ["/app/inputConvert"]