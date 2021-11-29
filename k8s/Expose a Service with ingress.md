# Description
通过ingress nginx + metallb 发布服务
# Infra
## 实验室集群搭建 (x master + 2 worker)
通过vagrant+ansible搭建,参考连接 https://github.com/liuliwh/k8s-ansible
master1+master2 使用VIP 192.168.34.100, endpoint: VIP:8443
## taint 或者 label 给ingress controller 专用/偏属节点
由于实验条件有限，所以使用 label 的方式将 worker#1 与 worker#2 节点票房为ingress=true
```shell
 k label nodes k-n-1 node-role.kubernetes.io/ingress="true"
 k label nodes k-n-2 node-role.kubernetes.io/ingress="true"
```
## 部署 ingress-nginx
参考baremeta方案，但将 deployment/replicas 设置为2,同时为其设置 nodesector. NodePort 改为 LoadBalancer
diff 如下
```
diff ingress-nginx-baremeta.yaml baremeta.yaml
276c276
<   type: LoadBalancer
---
>   type: NodePort
310d309
<   replicas: 2
402d400
<         node-role.kubernetes.io/ingress: "true"
684c682
<         runAsUser: 2000
\ No newline at end of file
---
>         runAsUser: 2000
```
修改后 manifest: https://github.com/liuliwh/k8s-ansible/blob/main/deployment/ingress-nginx-baremetal.yaml

## 安装 metallb
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.11.0/manifests/namespace.yaml
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.11.0/manifests/metallb.yaml

### 部署metallb 的configmap
按需修改 ip. 由于我的 node 真实 ip 分布为: 
masters: 192.168.34.11-192.168.34.13
workers: 192.168.34.101~
control plan VIP: 192.168.34.100
所以我将metallb的ip pool 范围设置在 :192.168.34.200-192.168.34.205
manifest: 
https://github.com/liuliwh/k8s-ansible/blob/main/deployment/metallb-configmap.yaml
```shell
$ kubectl -n ingress-nginx get svc
NAME                                 TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)                      AGE
ingress-nginx-controller             LoadBalancer   10.96.148.195    192.168.34.200   80:32651/TCP,443:32520/TCP   8m5s
ingress-nginx-controller-admission   ClusterIP      10.104.213.104   <none>           443/TCP                      8m6s
curl http://192.168.34.200/ -H 'Host: myapp.example.com'
<html>
<head><title>404 Not Found</title></head>
<body>
<center><h1>404 Not Found</h1></center>
<hr><center>nginx</center>
</body>
</html>
```
# 部署httpserver 用户级的服务
```
sarah@sarah-510-p127c:~/k8s_ansible$ k apply -f deployment/myhttpserver.yaml 
deployment.apps/httpserver created
service/httpservice created
ingress.networking.k8s.io/httpserver-ingress created
sarah@sarah-510-p127c:~/k8s_ansible$ curl http://192.168.34.200/ -H 'Host: httpserver.localdev.me'
ok
```

## 验证httpserver的功能
1. 源IP保留
```shell
vagrant@k-m-1:~/mani$ k logs -l  app=httpserver
I1129 22:18:09.990947       1 main.go:38] Server Started
I1129 22:18:09.991239       1 httpserver.go:25] Starting http server: :8080
I1129 22:18:09.911241       1 main.go:38] Server Started
I1129 22:18:09.911334       1 httpserver.go:25] Starting http server: :8080
I1129 22:18:46.554741       1 httplog_middleware.go:63] "HTTP" srcIP="192.168.34.1" URI="/" status=200
I1129 22:18:51.574251       1 httplog_middleware.go:63] "HTTP" srcIP="192.168.34.1" URI="/" status=200
```
2. request header 添加上
```bash
sarah@sarah-510-p127c:~/k8s_ansible$ curl http://192.168.34.200/ -H 'Host: httpserver.localdev.me' -v
*   Trying 192.168.34.200:80...
* TCP_NODELAY set
* Connected to 192.168.34.200 (192.168.34.200) port 80 (#0)
> GET / HTTP/1.1
> Host: httpserver.localdev.me
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Mon, 29 Nov 2021 22:45:39 GMT
< Content-Type: text/plain; charset=utf-8
< Content-Length: 3
< Connection: keep-alive
< Version: v1.0.1
< X-Req-Accept: */*
< X-Req-User-Agent: curl/7.68.0
< X-Req-X-Forwarded-For: 192.168.34.1
< X-Req-X-Forwarded-Host: httpserver.localdev.me
< X-Req-X-Forwarded-Port: 80
< X-Req-X-Forwarded-Proto: http
< X-Req-X-Forwarded-Scheme: http
< X-Req-X-Real-Ip: 192.168.34.1
< X-Req-X-Request-Id: 67ccf3ea7ce784f420f9b390b8bc1280
< X-Req-X-Scheme: http
< 
ok
```
3. 日志分级
access-log 必须以 VERBOSE LEVEL >= 3
health check 相关的 access-log, LEVEL >= 4

