FROM golang:1.12 as build

WORKDIR /app

COPY . .

RUN go get -d
RUN go build -o main .

EXPOSE 8083

CMD ["/main"]
