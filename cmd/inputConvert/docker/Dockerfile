FROM golang:1.23.1-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /app

# Копируем только модули сначала для кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Убираем проблемные replace-директивы
RUN sed -i '/replace/d' go.mod

# Явно скачиваем проблемную зависимость
RUN go get github.com/xuri/excelize/v2

# Собираем приложение с правильным путём
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/inputConvert ./cmd/inputConvert/main.go

FROM scratch
COPY --from=builder /app/inputConvert /inputConvert
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/inputConvert"]