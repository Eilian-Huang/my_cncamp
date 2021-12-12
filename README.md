# my_cncamp

## 目录
1. [Go语言练习](#1-Go语言练习)  
   1. [基础练习](#11-基础练习)  
   2. [生产者-消费者模型](#12-生产者-消费者模型)  
2. [HTTP Server](#2-http-server)
   1. [代码实现](#21-http-server)
   2. [Docker容器化](#22-http-server容器化)
   3. [Kubernetes部署](#23-kubernetes部署)
   4. [Prometheus监控](#24-为http-server添加监控)

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
[Deployment yaml](http_server/httpserver-deploy.yaml)
> 编写 Kubernetes 部署脚本将 httpserver 部署到 kubernetes 集群
>> 思考维度
>>> - [x] 优雅启动：启动探针和就绪探针
>>> - [x] 优雅终止：preStop
>>> - [x] 资源需求和 QoS 保证：Qos为Burstable
>>> - [x] 探活：存活探针
>>> - [x] 日常运维需求，日志等级：日志等级
>>> - [x] 配置和代码分离：env、secret和configmap
> 
>> 更加完备的部署spec，将服务发布给集群外部的调用方
>>> - [x] Service：NodePort
>>> - [ ] Ingress
> 
>> 可以考虑的细节
>>> - [x] 如何确保整个应用的高可用：多副本、亲和性
>>> - [x] 如何通过证书保证 httpServer 的通讯安全

- 使用yaml文件部署http server到kubernetes
```bash
$ kubectl apply -f httpserver-deploy.yaml

deployment.apps/httpserver created
```

### 2.4 为Http Server添加监控

> - [x] 在Http Server添加0-2秒的随机延时
> - [x] 在Http Server添加延时Metric
> - [x] 在Prometheus界面中查询延时指标数据
> - [x] 创建一个Grafana Dashboard展现延时分配情况

* 安装loki、Grafana和Prometheus
```shell

```
* 暴露Grafana和Prometheus并访问 

将type: ClusterIP改为NodePort
```shell
$ kubectl edit svc loki-grafana

apiVersion: v1
kind: Service
metadata:
  ...
spec:
  ...
  type: NodePort
...

$ kubectl edit svc loki-prometheus-server

apiVersion: v1
kind: Service
metadata:
  ...
spec:
  ...
  type: NodePort
...
```
查看 grafana 的登录密码
```shell
$ kubectl get secret --namespace default loki-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```
查看端口并登陆Grafana和Prometheus
```shell
$ kubectl get pod
loki-grafana                    NodePort    10.107.231.131   <none>        80:31780/TCP   24h
loki-prometheus-server          NodePort    10.106.247.155   <none>        80:30529/TCP   24h
```

* 在Http Server添加0-2秒的随机延时
```go
func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(4).Info("entering root handler")
	timer := metrics.NewTimer()
	// 汇报指标
	defer timer.ObserveTotal()
	// 添加延迟
	delay := randInt(0, 2000)
	time.Sleep(time.Millisecond*time.Duration(delay))
}
```
* 在Http Server添加延时Metric
```go
metrics.Register()
mux.Handle("/metrics", promhttp.Handler())  // prometheus metrics
```
* 在Deployment添加Prometheus发现
```yaml
...
spec:
  ...
  template:
    metadata:
      ...
      # prometheus
      annotations:
        prometheus.io/port: http-metrics
        prometheus.io/scrape: "true"
    spec:
      ...
      containers:
        - name: httpserver
          image: eilianhuang/cncamp:http_server_v1.0
          # prometheus ports
          ports:
            - containerPort: 80
              name: http-metrics
              protocol: TCP
```
* 重新打包docker镜像并修改Deployment镜像tag
```shell
$ docker tag http_server_v2.1 eilianhuang/cncamp:http_server_v2.1
$ docker push eilianhuang/cncamp:http_server_v2.1
```
* 查看是否有metrics打印
```shell
$ kubectl get po -owide
NAME                                           READY   STATUS    RESTARTS      AGE   IP                NODE     NOMINATED NODE   READINESS GATES
httpserver-7965cd4dc-9fzmg                     0/1     Running   0             21s   192.168.216.248   cncamp   <none>           <none>

$ curl 192.168.216.248:80/metrics
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
...
```
* 在Prometheus界面中查询延时指标数据
* 创建一个Grafana Dashboard展现延时分配情况