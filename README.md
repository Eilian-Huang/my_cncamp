# my_cncamp

## 目录
1. [Go语言练习](##1.Go语言练习)  
   1. [基础练习](###1.1基础练习)  
   2. [生产者-消费者模型](###1.2生产者-消费者模型)  
2. [HTTP Server](##2.HTTPServer)
   1. [代码实现](###2.1HTTPServer)
   2. [Docker容器化](###2.2HTTP Server容器化)
   3. [Kubernetes部署](###2.3Kubernetes部署)

## 1. Go语言练习

### 1.1 基础练习

[练习1.1 Code](practice_1_1/main.go)

1. 编写一个小程序 给定一个字符串数组`["I","am","stupid","and","weak"]` 
2. 用 for 循环遍历该数组并修改为`["I","am","smart","and","strong"]`

### 1.2 生产者-消费者模型

[练习1.2 Code](practice_1_2/main.go)

1. 基于 Channel 编写一个简单的单线程生产者消费者模型
2. 队列：队列长度10，队列元素类型为 int
3. 生产者：每1秒往队列中放入一个类型为 int 的元素，队列满时生产者可以阻塞
4. 消费者：每一秒从队列中获取一个元素并打印，队列为空时消费者阻塞

### 1.3 多个生产者和多个消费者模式

## 2. HTTP Server

### 2.1 HTTP Server

[Http Server Code](http_server/main.go)

1. 接收客户端 request，并将 request 中带的 header 写入 response header
2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4. 当访问 localhost/healthz 时，应返回200

### 2.2 HTTP Server容器化

- 构建本地镜像。

```bash
$ GOOS=linux GOARCH=amd64 go build -o bin/linux/httpserver
$ docker build -t http_server_v1.0 .
```

- 编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化（请思考有哪些最佳实践可以引入到 Dockerfile 中来）。

[Dockerfile](http_server/Dockerfile)

1. 最佳实践1：多端构建减小镜像：`FROM ubuntu -> FROM golang:1.16-alpine AS build`
2. 最佳实践2：多条 RUN 命令可通过连接符连接成一条指令集以减少层数。
3. 编写 dockerfile 的时候，应该把变更频率低的编译指令优先构建以便放在镜像底层以有效利用 build cache。

- 将镜像推送至 Docker 官方镜像仓库。

```bash
$ docker tag http_server_v1.0 eilianhuang/cncamp:http_server_v1.0
$ docker push eilianhuang/cncamp:http_server_v1.0
```

- 通过 Docker 命令本地启动 httpserver。

```bash
$ docker run -d --name httpserver -p 800:80 eilianhuang/cncamp:http_server_v1.0
```

- 通过 nsenter 进入容器查看 IP 配置。

```bash
$ docker ps|grep httpserver
$ docker inspect <containerid>|grep -i pid
$ nsenter -t <pid> -n ip a

// 查看http server
$ curl 127.0.0.1:800
$ curl 127.0.0.1:800/healthz
```

### 2.3 Kubernetes部署

> 编写 Kubernetes 部署脚本将 httpserver 部署到 kubernetes 集群
>> 思考维度
>>> - [x] 优雅启动
>>> - [x] 优雅终止 
>>> - [x] 资源需求和 QoS 保证 
>>> - [x] 探活 
>>> - [x] 日常运维需求，日志等级 
>>> - [ ] 配置和代码分离
> 
>> 更加完备的部署spec，将服务发布给集群外部的调用方
>>> - [ ] Service 
>>> - [ ] Ingress
> 
>> 可以考虑的细节
>>> - [ ] 如何确保整个应用的高可用 
>>> - [ ] 如何通过证书保证 httpServer 的通讯安全
