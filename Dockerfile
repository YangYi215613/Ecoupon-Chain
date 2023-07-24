# 基于Ubuntu 20.04镜像作为基础镜像
FROM ubuntu:20.04

# 增加国内源
# RUN sed -i 's/archive.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
# RUN sed -i 's/security.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
RUN sed -i 's/archive.ubuntu.com/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list
RUN sed -i 's/security.ubuntu.com/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list

# 安装依赖项
RUN apt-get update && apt-get install gcc libc6-dev git lrzsz wget vim -y

# 设置环境变量
ENV GOLANG_VERSION 1.18
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH
ENV GOPROXY https://goproxy.io,direct

# RUN cd &GOPATH && mkdir src pkg bin
RUN mkdir -p $GOPATH/src/Ecoupon-Chain $GOPATH/bin $GOPATH/pkg

# 下载和安装Go语言
RUN wget -q https://go.dev/dl/go1.18.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GOLANG_VERSION.linux-amd64.tar.gz && \
    rm go$GOLANG_VERSION.linux-amd64.tar.gz

# 拷贝文件
COPY . $GOPATH/src/Ecoupon-Chain

# 设置工作目录
WORKDIR $GOPATH/src/Ecoupon-Chain

RUN go mod download && go mod tidy