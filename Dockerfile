FROM golang:1.20-alpine As builder

# 设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn

# 移动到工作目录
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum . 
RUN go mod tidy

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文字 bluebell
RUN go build -o bluebell .

FROM scratch
COPY ./templates /templates
COPY ./static /static
COPY ./conf /conf

WORKDIR /
# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=builder /build/bluebell /

# 声明服务端口
EXPOSE 8081

# 需要运行的命令
ENTRYPOINT [ "/bluebell","conf/config.yaml" ]

