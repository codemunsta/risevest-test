FROM golang:1.20-alpine

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

EXPOSE 3000

CMD [ "go", "run", "src/main.go" ]
