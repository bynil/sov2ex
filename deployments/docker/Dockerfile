FROM golang:1.13.3 AS builder

LABEL maintainer="GeXiao <gexiao1994@gmail.com>"

ENV APP_DIR=$GOPATH/src/myapp/

RUN mkdir -p $APP_DIR

WORKDIR $APP_DIR

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=jsoniter -a -o main cmd/service/main.go

FROM alpine:latest AS production

WORKDIR /app

COPY --from=builder /go/src/myapp/main .

EXPOSE 8080

CMD ["./main"]