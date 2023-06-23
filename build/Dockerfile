FROM golang:1.20

WORKDIR /usr/src/app

COPY ./go.mod ./go.sum ./

RUN go mod download && go mod verify

COPY . .

WORKDIR /usr/src/app/cmd/hurma

RUN go build -v -o /usr/local/bin/app

EXPOSE 8080

CMD ["/usr/local/bin/app"]