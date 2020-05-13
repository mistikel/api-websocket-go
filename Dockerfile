FROM golang:1.14-alpine

WORKDIR /usr/go/src/

COPY . /usr/go/src/

RUN go mod download

EXPOSE 8080

RUN go build -o warpin

CMD ["./warpin"]