FROM golang:1.18

WORKDIR /go/src/grpc-service

COPY ./go.mod ./go.sum

RUN go mod download && do mod tidy && go mod verify

COPY . .

RUN go build -o app ./...

EXPOSE 8080

CMD ["app"]