FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go mod download

COPY internal ./internal

COPY cmd ./cmd

RUN go build -o /app/telego-bot ./cmd

EXPOSE 8080

CMD ["/app/telego-bot"]