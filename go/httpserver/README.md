# Project layout
refer to https://github.com/golang-standards/project-layout
# Httpserver Go Minimal Requirements
1. 接收客户端 request，并将 request 中带的 header 写入 response header
2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4. 当访问 localhost/healthz 时，应返回 200
# Multistage docker
1. The docker build context is the folder containing file go.mod 
```
cd /path/to/httpserver_with_go.mod
go % ls
build           go.mod          go.sum          httpserver
```
2. build the docker
```
go % docker build -t httpserver:v1 -f build/httpserver/Dockerfile .
```
3. run the docker
go % docker run -p 80:8080 httpserver:v1