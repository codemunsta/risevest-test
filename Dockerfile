FROM golang:1.20-alpine

COPY .env .env

WORKDIR /usr/src/app

RUN go install github.com/cosmtrek/air@latest

COPY . .

RUN go mod tidy

RUN go build -o rise-vest-test .

EXPOSE 3000

CMD [ "./rise-vest-test" ]