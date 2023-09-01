FROM golang:1.20-alpine

COPY .env .env

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

EXPOSE 3000

CMD [ "go", "run", "main.go" ]
