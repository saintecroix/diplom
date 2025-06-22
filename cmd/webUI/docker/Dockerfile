FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/webUI ./cmd/webUI/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/webUI /app/webUI
COPY --from=builder /app/cmd/webUI/internal/web/static /app/cmd/webUI/internal/web/static

EXPOSE 8080

CMD ["/app/webUI"]