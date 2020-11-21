FROM golang:1.15.5 AS builder

LABEL maintainer="GeXiao <gxgexiao@gmail.com>"

ENV TZ=Asia/Shanghai LANG=C.UTF-8 APP_DIR=/go/src/myapp/

RUN mkdir -p $APP_DIR

WORKDIR $APP_DIR

# make a cached image for go mod
COPY go.mod .
COPY go.sum .

RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=jsoniter -a -o main cmd/service/main.go

FROM alpine:latest AS production

WORKDIR /app

COPY --from=builder /go/src/myapp/main .

ENV TZ=Asia/Shanghai LANG=C.UTF-8

EXPOSE 8080

CMD ["./main"]