FROM golang:1.20-alpine

COPY .env .env

WORKDIR /usr/src/app

RUN go install github.com/cosmtrek/air@latest

COPY . .

RUN go mod tidy

EXPOSE 3000

CMD [ "air", "main.go" ]
