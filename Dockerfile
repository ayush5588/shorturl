FROM golang:1.20-alpine

WORKDIR $GOPATH/src/github.com/ayush5588/shorturl

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    DOMAIN_NAME="http://localhost:8080/" \
    REDIS_HOSTNAME="redis"

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o shorturl

EXPOSE 8080

CMD [ "./shorturl" ]
