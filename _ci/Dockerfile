FROM golang as builder
WORKDIR /project
COPY . /project
ENV GOPROXY=https://goproxy.io,direct
ENV CGO_ENABLED=0
RUN go build -o /go/bin/app /project/main.go

# ------ 分界线 ------

FROM centos:7

# 替换系统镜像源为阿里云
RUN curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo

# 设置时区
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
ENV TZ=Asia/Shanghai

WORKDIR /project
COPY --from=builder /go/bin/app /project/app
COPY --from=builder /project/config /project/config
ENV SERVER_PORT=80
ENV SERVER_MODE=release
ENTRYPOINT ["/project/app"]