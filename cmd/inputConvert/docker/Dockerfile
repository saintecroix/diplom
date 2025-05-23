FROM golang:1.21-alpine AS builder

RUN apk update && apk add --no-cache ca-certificates
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/data-aggregator .

FROM scratch

COPY --from=builder /app/inputConvert /inputConvert

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["inputConvert"]