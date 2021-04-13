FROM golang AS builder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.io,direct

WORKDIR /project

COPY go.mod go.sum /project/
RUN go mod download
COPY . /project
RUN go get -u github.com/swaggo/swag/cmd/swag && \
    rm -rf ./docs && \
    swag init && \
    go build -o /go/bin/app /project/main.go

# ------ 分界线 ------

FROM alpine
MAINTAINER xusheng:20691718@qq.com
EXPOSE 80

# 替换系统镜像源为阿里云
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# 设置时区
RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /project
COPY --from=builder /go/bin/app /project/app
COPY --from=builder /project/config /project/config
ENTRYPOINT ["/project/app"]