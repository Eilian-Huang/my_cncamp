# my_cncamp

## 1. Go语言练习1.1

[练习1.1 Code](practice_1_1/main.go)

1. 编写一个小程序 给定一个字符串数组`["I","am","stupid","and","weak"]` 
2. 用 for 循环遍历该数组并修改为`["I","am","smart","and","strong"]`

## 2. Go语言练习1.2：生产者-消费者模型

[练习1.2 Code](practice_1_2/main.go)

1. 基于 Channel 编写一个简单的单线程生产者消费者模型
2. 队列：队列长度10，队列元素类型为 int
3. 生产者：每1秒往队列中放入一个类型为 int 的元素，队列满时生产者可以阻塞
4. 消费者：每一秒从队列中获取一个元素并打印，队列为空时消费者阻塞

## 3. HTTP Server

[Http Server Code](http_server/main.go)

1. 接收客户端 request，并将 request 中带的 header 写入 response header
2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4. 当访问 localhost/healthz 时，应返回200

## 4. HTTP Server容器化

- 构建本地镜像。

```bash
$ GOOS=linux GOARCH=amd64 go build -o bin/linux/httpserver
$ docker build -t my_cncamp/http_server:huangsiyi_v0.0 .
```

- 编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化（请思考有哪些最佳实践可以引入到 Dockerfile 中来）。

[Dockerfile](http_server/Dockerfile)

1. 最佳实践1：多端构建减小镜像：`FROM ubuntu -> FROM golang:1.16-alpine AS build`
2. 最佳实践2：多条 RUN 命令可通过连接符连接成一条指令集以减少层数。
3. 编写 dockerfile 的时候，应该把变更频率低的编译指令优先构建以便放在镜像底层以有效利用 build cache。

- 将镜像推送至 Docker 官方镜像仓库。

```bash
$ docker push my_cncamp/http_server:huangsiyi_v0.0
```

- 通过 Docker 命令本地启动 httpserver。

```bash
docker run -d --name httpserver -p 80:800 my_cncamp/http_server:huangsiyi_v0.0
```

- 通过 nsenter 进入容器查看 IP 配置。

```bash
docker ps|grep httpserver
docker inspect <containerid>|grep -i pid
nsenter -t <pid> -n ip a

// 查看http server
curl 127.0.0.1/800
curl 127.0.0.1/healthz/800
```
