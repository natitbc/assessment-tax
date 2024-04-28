FROM golang:1.21.9

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /go/bin/app

CMD ["/go/bin/app"]