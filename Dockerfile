FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./
COPY config.json ./

RUN go build -o /app/telegramBot

EXPOSE 8080

CMD ["/app/telegramBot"]