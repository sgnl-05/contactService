FROM golang:1.17

RUN  mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY . /usr/src/app

RUN go mod download

EXPOSE 8080

RUN go build -o /conservice

ENTRYPOINT ["/conservice"]
