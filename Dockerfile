FROM golang:alpine

WORKDIR /app
COPY . .

RUN apk add --no-cache git
RUN go get -u github.com/go-telegram-bot-api/telegram-bot-api/v5
RUN go build -o main .

CMD ["./main"]