- [Description](#description)
- [Infra](#infra)
  - [实验室集群搭建 (x master + 2 worker)](#实验室集群搭建-x-master--2-worker)
  - [taint 或者 label 给ingress controller 专用/偏属节点](#taint-或者-label-给ingress-controller-专用偏属节点)
  - [部署 ingress-nginx](#部署-ingress-nginx)
  - [安装 metallb](#安装-metallb)
    - [部署metallb 的configmap](#部署metallb-的configmap)
  - [集群部署后效果](#集群部署后效果)
    - [集群相关](#集群相关)
    - [网络相关](#网络相关)
- [部署httpserver 用户级的服务](#部署httpserver-用户级的服务)
  - [验证httpserver的功能（Vagrant host)](#验证httpserver的功能vagrant-host)
- [Deploy HTTPS application](#deploy-https-application)
  - [Create cert for the application](#create-cert-for-the-application)
  - [Create secret yaml based on cert](#create-secret-yaml-based-on-cert)
  - [Add tls secret into application ingress](#add-tls-secret-into-application-ingress)
  - [Verify https traffic](#verify-https-traffic)
- [Reference](#reference)
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
参考baremeta方案，但将 deployment 改为 Daemonset,同时为其设置 nodesector. NodePort 改为 LoadBalancer, 以及加上了 externalTrafficPolicy: Local
修改后 manifest: https://github.com/liuliwh/k8s-ansible/blob/main/deployment/ingress-nginx-baremetal.yaml

## 安装 metallb
如果集群已切换为 ipvs 模式， 需要将 configmap.kub-proxy.ipvs.strictARP 设置为 true. 请参考https://metallb.universe.tf/installation/

```bash
# see what changes would be made, returns nonzero returncode if different
kubectl get configmap kube-proxy -n kube-system -o yaml | \
sed -e "s/strictARP: false/strictARP: true/" | \
kubectl diff -f - -n kube-system

# actually apply the changes, returns nonzero returncode on errors only
kubectl get configmap kube-proxy -n kube-system -o yaml | \
sed -e "s/strictARP: false/strictARP: true/" | \
kubectl apply -f - -n kube-system
```

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
## 集群部署后效果
### 集群相关
```bash
sarah@sarah-510-p127c:~$ k get all -n ingress-nginx 
NAME                                            READY   STATUS      RESTARTS   AGE
pod/ingress-nginx-admission-create--1-6kqsg     0/1     Completed   0          18h
pod/ingress-nginx-admission-patch--1-99xj5      0/1     Completed   1          18h
pod/ingress-nginx-controller-5db7ffdb85-5j29z   1/1     Running     0          18h
pod/ingress-nginx-controller-5db7ffdb85-zx96g   1/1     Running     0          18h

NAME                                         TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)                      AGE
service/ingress-nginx-controller             LoadBalancer   10.96.148.195    192.168.34.200   80:32651/TCP,443:32520/TCP   18h
service/ingress-nginx-controller-admission   ClusterIP      10.104.213.104   <none>           443/TCP                      18h

NAME                                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/ingress-nginx-controller   2/2     2            2           18h

NAME                                                  DESIRED   CURRENT   READY   AGE
replicaset.apps/ingress-nginx-controller-5db7ffdb85   2         2         2       18h

NAME                                       COMPLETIONS   DURATION   AGE
job.batch/ingress-nginx-admission-create   1/1           2s         18h
job.batch/ingress-nginx-admission-patch    1/1           3s         18h

sarah@sarah-510-p127c:~$ k get all -n metallb-system 
NAME                              READY   STATUS    RESTARTS   AGE
pod/controller-7dcc8764f4-msfds   1/1     Running   0          18h
pod/speaker-46p9v                 1/1     Running   0          18h
pod/speaker-k68gx                 1/1     Running   0          18h
pod/speaker-nnbh5                 1/1     Running   0          18h

NAME                     DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR            AGE
daemonset.apps/speaker   3         3         3       3            3           kubernetes.io/os=linux   18h

NAME                         READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/controller   1/1     1            1           18h

NAME                                    DESIRED   CURRENT   READY   AGE
replicaset.apps/controller-7dcc8764f4   1         1         1       18h

```
### 网络相关
```bash
sarah@sarah-510-p127c:~$ virsh net-dumpxml k0
<network connections='3' ipv6='yes'>
  <name>k0</name>
  <uuid>1390af10-6b32-49dc-b243-dd75cce9bc7b</uuid>
  <forward mode='nat'>
    <nat>
      <port start='1024' end='65535'/>
    </nat>
  </forward>
  <bridge name='virbr3' stp='on' delay='0'/>
  <mac address='52:54:00:49:98:29'/>
  <ip address='192.168.34.1' netmask='255.255.255.0'>
    <dhcp>
      <range start='192.168.34.2' end='192.168.34.199'/>
    </dhcp>
  </ip>
</network>

sarah@sarah-510-p127c:~$ ip a s virbr3
9: virbr3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
    link/ether 52:54:00:49:98:29 brd ff:ff:ff:ff:ff:ff
    inet 192.168.34.1/24 brd 192.168.34.255 scope global virbr3
       valid_lft forever preferred_lft forever
```
# 部署httpserver 用户级的服务
yaml file: https://github.com/liuliwh/k8s-ansible/blob/main/deployment/myhttpserver.yaml 
```
sarah@sarah-510-p127c:~/k8s_ansible$ k apply -f deployment/myhttpserver.yaml 
deployment.apps/httpserver created
service/httpservice created
ingress.networking.k8s.io/httpserver-ingress created
sarah@sarah-510-p127c:~/k8s_ansible$ curl http://192.168.34.200/ -H 'Host: httpserver.localdev.me'
ok
```

## 验证httpserver的功能（Vagrant host)
1. 源IP保留
```shell
sarah@sarah-510-p127c:~/k8s_ansible$ k logs -l app.kubernetes.io/name=httpserver
I1203 22:18:24.907617       1 main.go:38] Server Started
I1203 22:18:24.908096       1 httpserver.go:25] Starting http server: :8080
I1203 22:19:29.162219       1 httplog_middleware.go:63] "HTTP" srcIP="192.168.121.103" URI="/" status=200
I1203 22:23:33.095631       1 httplog_middleware.go:63] "HTTP" srcIP="192.168.34.200" URI="/" status=200
I1203 22:18:24.889001       1 main.go:38] Server Started
I1203 22:18:24.889106       1 httpserver.go:25] Starting http server: :8080
I1203 22:20:19.026606       1 httplog_middleware.go:63] "HTTP" srcIP="192.168.34.200" URI="/" status=200
I1203 22:24:20.910247       1 httplog_middleware.go:63] "HTTP" srcIP="192.168.34.1" URI="/" status=200
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
# Deploy HTTPS application
## Create cert for the application
```bash
openssl req -x509 -newkey rsa:4096 -keyout tls.key -out tls.crt -days 365 -nodes -subj "/C=US/ST=Florida/L=Miami/O=localdev.me/CN=httpserver.localdev.me"
```
## Create secret yaml based on cert
```bash
sarah@sarah-510-p127c:~/k8s_ansible$ k create secret tls httpserver-localdev-tls --key tls.key --cert tls.crt --dry-run=client -oyaml > httpserver.tls.secret.yaml 
sarah@sarah-510-p127c:~/k8s_ansible$ k delete secrets httpserver-localdev-tls 
secret "httpserver-localdev-tls" deleted
sarah@sarah-510-p127c:~/k8s_ansible$ k create -f httpserver.tls.secret.yaml 
secret/httpserver-localdev-tls created
```
## Add tls secret into application ingress
```bash
spec:
  ingressClassName: nginx
  tls:
    - hosts: 
      - httpserver.localdev.me
      secretName: httpserver-localdev-tls
```
The detail info is in https://github.com/liuliwh/k8s-ansible/blob/main/deployment/myhttpserver.yaml 
## Verify https traffic
由于是自签证书，所以要带-k
```bash
vagrant@k-n-1:~$ curl -H 'Host: httpserver.localdev.me' https://192.168.34.200/  -k
ok
```
# Reference
1. https://kubernetes.github.io/ingress-nginx/deploy/baremetal/#a-pure-software-solution-metallb
2. 
