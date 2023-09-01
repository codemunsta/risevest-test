FROM golang:1.20-alpine

WORKDIR /src

COPY ./src .

RUN go mod download

RUN go build -o rise .

EXPOSE 3000

CMD [ "./rise" ]
