FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/inputConvert ./cmd/inputConvert/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/inputConvert /app/inputConvert

CMD ["/app/inputConvert"]