# Project layout
refer to https://github.com/golang-standards/project-layout
# Httpserver Go Minimal Requirements
1. 接收客户端 request，并将 request 中带的 header 写入 response header
2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4. 当访问 localhost/healthz 时，应返回 200
# How to build 
1. build container
```
make build
```
2. build and push container
```
make push
```
3. run
```
make run
```
# Ready to use image on hub
```
docker pull sarah4test/httpserver
```
# 通过 nsenter 进入容器查看 IP 配置
0. Precondition: run the docker container
```
vagrant@cncamp:~$ docker run -p 8080:8080 -d sarah4test/httpserver
3d5d51e4187288ffe1e2b6a6a4eea9d7b3ee6c333f008c4df869dfb0b721f0fb
```
1. Find the pid by using docker ps and docker inspect
<b>注意: host ps 与docker ps 查出来的pid不一样，有对应关系</b>
```
vagrant@cncamp:~$ docker inspect --format "{{.State.Pid}}" $(docker ps -q --filter ancestor=sarah4test/httpserver)
4172
```
2. nsenter -t $PID -n command
```
vagrant@cncamp:~$ sudo nsenter -t 4172 -n ip a s
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
13: eth0@if14: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever
```
# TODO
1. graceful shutdown 