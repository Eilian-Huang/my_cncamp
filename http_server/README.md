# HTTP Server容器化

- 构建本地镜像。

```
GOOS=linux GOARCH=amd64 go build -o bin/linux/httpserver
docker build -t my_cncamp/http_server:huangsiyi_v0.0 .
```

- 编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化（请思考有哪些最佳实践可以引入到 Dockerfile 中来）。

[Dockerfile](./Dockerfile)

1. 最佳实践1：多端构建减小镜像：`FROM ubuntu -> FROM golang:1.16-alpine AS build`
2. 最佳实践2：多条 RUN 命令可通过连接符连接成一条指令集以减少层数。
3. 编写 dockerfile 的时候，应该把变更频率低的编译指令优先构建以便放在镜像底层以有效利用 build cache。

- 将镜像推送至 Docker 官方镜像仓库。
    
`docker push my_cncamp/http_server:huangsiyi_v0.0`

- 通过 Docker 命令本地启动 httpserver。
    
`docker run -d --name httpserver -p 80:800 my_cncamp/http_server:huangsiyi_v0.0`

- 通过 nsenter 进入容器查看 IP 配置。

```bash
docker ps|grep httpserver
docker inspect <containerid>|grep -i pid
nsenter -t <pid> -n ip a

// 查看http server
curl 127.0.0.1/800
curl 127.0.0.1/healthz/800
```