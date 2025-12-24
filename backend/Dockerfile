FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Устанавливаем goose для миграций
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем скомпилированное приложение
COPY --from=builder /app/main .

# Копируем goose бинарник
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Копируем миграции
COPY --from=builder /app/db/migrations /app/db/migrations

EXPOSE 3000

CMD ["./main"]

