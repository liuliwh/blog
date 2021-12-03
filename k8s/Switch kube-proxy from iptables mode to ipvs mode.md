- [What is kube-proxy](#what-is-kube-proxy)
- [kube-proxy modes](#kube-proxy-modes)
  - [userspace (seldomly used)](#userspace-seldomly-used)
  - [iptables mode (by default)](#iptables-mode-by-default)
    - [Precondition - Service setup](#precondition---service-setup)
    - [Service Chain](#service-chain)
    - [service specific chain](#service-specific-chain)
    - [Service Endpoints chain](#service-endpoints-chain)
    - [KUBE-MARK-MASQ chain](#kube-mark-masq-chain)
    - [Conclusion](#conclusion)
  - [IPVS](#ipvs)
    - [IPVS Service Network Topology](#ipvs-service-network-topology)
    - [How to config kube-proxy to use IPVS while creating cluster](#how-to-config-kube-proxy-to-use-ipvs-while-creating-cluster)
    - [How to switch from iptables mode to IPVS mode for installed cluster](#how-to-switch-from-iptables-mode-to-ipvs-mode-for-installed-cluster)
    - [Test IPVS mode is running](#test-ipvs-mode-is-running)
    - [Extra side note about NodePort and LoadBanalcer service](#extra-side-note-about-nodeport-and-loadbanalcer-service)
# What is kube-proxy
kube-proxy historically acccts as network proxy between pods and services, but since k8s project evolved, kube-proxy more became of a node agent.
# kube-proxy modes
## userspace (seldomly used)
## iptables mode (by default)
iptables is mainly designed for firewall usage. it's wrapper of netfilter. the datastructure is list with it's O(N) complexity.
kube-proxy programs iptables NAT table to DNAT from service clusterIP(VIP) chain -> service specific chain -> service endpoint chain.
eg. 1 service with a deployment of 3 pods
### Precondition - Service setup
```bash
NAME          TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
httpservice   ClusterIP   10.104.218.6   <none>        80/TCP    84m
NAME          ENDPOINTS                                                   AGE
httpservice   192.168.178.7:8080,192.168.178.8:8080,192.168.238.70:8080   85m
sarah@sarah-510-p127c:~/k8s_ansible$ k get pods -owide
NAME                                    READY   STATUS    RESTARTS      AGE   IP               NODE    NOMINATED NODE   READINESS GATES
httpserver-deployment-b86d898db-6dkzl   1/1     Running   1 (36m ago)   85m   192.168.238.70   k-n-1   <none>           <none>
httpserver-deployment-b86d898db-7bwq5   1/1     Running   1 (36m ago)   85m   192.168.178.8    k-n-2   <none>           <none>
httpserver-deployment-b86d898db-mtchv   1/1     Running   1 (36m ago)   85m   192.168.178.7    k-n-2   <none>           <none>
```
The iptables nat table programed by kube-proxy
service chain
### Service Chain
```bash
vagrant@k-n-2:~$ sudo iptables -t nat -L PREROUTING
Chain PREROUTING (policy ACCEPT)
KUBE-SERVICES  all  --  anywhere             anywhere             /* kubernetes service portals */
vagrant@k-n-2:~$ sudo iptables -t nat -L KUBE-SERVICES 
KUBE-SVC-PZZOH4GMXQXQLAH7  tcp  --  anywhere             10.104.218.6         /* default/httpservice:http cluster IP */ tcp dpt:http
```
### service specific chain
```bash
vagrant@k-n-2:~$ sudo iptables -t nat -L KUBE-SVC-PZZOH4GMXQXQLAH7
Chain KUBE-SVC-PZZOH4GMXQXQLAH7 (1 references)
target     prot opt source               destination         
KUBE-MARK-MASQ  tcp  -- !192.168.0.0/16       10.104.218.6         /* default/httpservice:http cluster IP */ tcp dpt:http
KUBE-SEP-SHRHPIF6JJVVZWAJ  all  --  anywhere             anywhere             /* default/httpservice:http */ statistic mode random probability 0.33333333349
KUBE-SEP-RWO7TIJ7ZXLAS5CA  all  --  anywhere             anywhere             /* default/httpservice:http */ statistic mode random probability 0.50000000000
KUBE-SEP-UDTT7ANR2HKPGI6N  all  --  anywhere             anywhere             /* default/httpservice:http */
```
### Service Endpoints chain
```bash
vagrant@k-n-2:~$ sudo iptables -t nat -L KUBE-SEP-SHRHPIF6JJVVZWAJ
Chain KUBE-SEP-SHRHPIF6JJVVZWAJ (1 references)
target     prot opt source               destination         
KUBE-MARK-MASQ  all  --  192.168.178.7        anywhere             /* default/httpservice:http */
DNAT       tcp  --  anywhere             anywhere             /* default/httpservice:http */ tcp to:192.168.178.7:8080
```
### KUBE-MARK-MASQ chain
```bash
vagrant@k-n-2:~$ sudo iptables -t nat -L KUBE-MARK-MASQ
Chain KUBE-MARK-MASQ (23 references)
target     prot opt source               destination         
MARK       all  --  anywhere             anywhere             MARK or 0x4000
```
Purpose: Route packets that arrive at a node from outside the cluster. Or for egress purpose. After mark the pakcage, the IP will not be the POD ip, rather than the node IP.
### Conclusion
It works in the kernal space.
1. but the rules is evaluated one by one until match. When the cluster is large scale and has a lot of services, the IP tables rules are getting huge, so there is performance and scalability limitation.
2. iptables rules are not incremetal, so kube-proxy needs to write out the entire table for every update. When the table is huge, the updates can take minutes, then risks sending traffic to stale endpoints.
## IPVS
IPVS(IP Virtual Server) is designed for load balancing. It is also included into Linux kernel, it uses hash tables with O(1) search complexity.
It is built on top of the netfilter too. It implements L4 load balancing. It support multiple load balancing algorthrims.
### IPVS Service Network Topology
- Create a dummy interface on the node (default interface name is kube-ipvs0)
- Bind service IP (VIP) to the dummy interface
- Create IPVS virtual servers for each Service IP
### How to config kube-proxy to use IPVS while creating cluster
Precondition: service is created.

```bash
sarah@sarah-510-p127c:~/k8s_ansible$ k get service
httpservice-np   NodePort       10.109.208.252   <none>        80:32215/TCP   7s
```
0. Make sure the ipvs modules are loaded
```bash
root@k-n-2:~# lsmod | grep -e ip_vs -e nf_conntrack
nf_conntrack_netlink    45056  0
nfnetlink              16384  4 nf_conntrack_netlink,nf_tables,ip_set
ip_vs_sh               16384  0
ip_vs_wrr              16384  0
ip_vs_rr               16384  0
ip_vs                 155648  6 ip_vs_rr,ip_vs_sh,ip_vs_wrr
nf_conntrack          139264  6 xt_conntrack,nf_nat,xt_nat,nf_conntrack_netlink,xt_MASQUERADE,ip_vs
```
1. Createt a custom kubeadm.yaml
```bash
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: ipvs
```
### How to switch from iptables mode to IPVS mode for installed cluster
0. Make sure the ipvs modules are loaded

1. Change the mode from "" to "ipvs"
```bash
sarah@sarah-510-p127c:~/k8s_ansible$ k -n kube-system edit configmaps kube-proxy 
```
2. kill all kube-proxy pods
```bash
sarah@sarah-510-p127c:~/k8s_ansible$ k -n kube-system delete pods -l k8s-app=kube-proxy
```
3. Verify kube-proxy is started with ipvs proxier
```bash
sarah@sarah-510-p127c:~/k8s_ansible$ k -n kube-system logs kube-proxy-4v5gm | grep ipvs
I1203 15:39:46.292667       1 server_others.go:274] Using ipvs Proxier.
```
### Test IPVS mode is running
- dummy interface (kube-ipvs0) is creating
```bash
10: kube-ipvs0: <BROADCAST,NOARP> mtu 1500 qdisc noop state DOWN group default 
    link/ether 8a:28:ce:01:a5:c4 brd ff:ff:ff:ff:ff:ff
    inet 10.109.208.252/32 scope global kube-ipvs0
       valid_lft forever preferred_lft forever
```
- Service IP is binded to dummy interface (the output is as above)
- Create IPVS virtual servers for each Service IP
If you have ipvsadm installed, you can check the virtual service and real server.
The setup is as below.
```bash
sarah@sarah-510-p127c:~/k8s_ansible$ k get service
NAME             TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
httpservice-np   NodePort    10.104.87.249   <none>        80:30806/TCP   6m24s
sarah@sarah-510-p127c:~/k8s_ansible$ k get pods -owide
NAME                                     READY   STATUS    RESTARTS   AGE   IP               NODE    NOMINATED NODE   READINESS GATES
httpserver-deployment-6586df99f8-74rd5   1/1     Running   0          11m   192.168.238.68   k-n-1   <none>           <none>
httpserver-deployment-6586df99f8-zp7xx   1/1     Running   0          11m   192.168.178.4    k-n-2   <none>           <none>
```
ipvsadm output
```bash
vagrant@k-n-1:~$ sudo ipvsadm --list --numeric --tcp-service 10.104.87.249:80
Prot LocalAddress:Port Scheduler Flags
  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
TCP  10.104.87.249:80 rr
  -> 192.168.178.4:8080           Masq    1      0          0         
  -> 192.168.238.68:8080          Masq    1      0          0   
```
Flush the iptables nat if you switched from iptables to ipvs, you should notice the service chain, service specific chain, service endpoint chain is empty. Please compare with output of iptable mode
```bash
vagrant@k-n-2:~$ sudo iptables -t nat -F
vagrant@k-n-2:~$ sudo iptables -t nat -L PREROUTING
Chain PREROUTING (policy ACCEPT)
target     prot opt source               destination         
cali-PREROUTING  all  --  anywhere             anywhere             /* cali:6gwbT8clXdHdC1b1 */
KUBE-SERVICES  all  --  anywhere             anywhere             /* kubernetes service portals */
vagrant@k-n-2:~$ sudo iptables -t nat -L KUBE-SERVICES
Chain KUBE-SERVICES (2 references)
target     prot opt source               destination         
KUBE-MARK-MASQ  all  -- !192.168.0.0/16       anywhere             /* Kubernetes service cluster ip + port for masquerade purpose */ match-set KUBE-CLUSTER-IP dst,dst
KUBE-NODE-PORT  all  --  anywhere             anywhere             ADDRTYPE match dst-type LOCAL
ACCEPT     all  --  anywhere             anywhere             match-set KUBE-CLUSTER-IP dst,dst
```
Verify the service is reachable with the service IP
```bash
vagrant@k-n-2:~$ curl 10.109.208.252/
ok
```
### Extra side note about NodePort and LoadBanalcer service
kube-proxy creates an IPVS service for the Service's cluster IP. Same as type Cluster.
```bash

```
- also creates an IPVS virtual service for each of the node’s IP addresses 
  
Node IP
```bash
vagrant@k-n-1:~$ ip a s | grep -w inet
    inet 127.0.0.1/8 scope host lo
    inet 192.168.121.103/24 brd 192.168.121.255 scope global dynamic eth0
    inet 192.168.34.101/24 brd 192.168.34.255 scope global eth1
    inet 10.104.87.249/32 scope global kube-ipvs0
    inet 192.168.238.64/32 scope global vxlan.calico
```
IPVS virtual service listening on the node’s IP address (primary/secondary).(xx.xx.xx.xx:30806)
```bash
vagrant@k-n-1:~$ sudo ipvsadm --list --numeric 
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port Scheduler Flags
  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
TCP  192.168.34.101:30806 rr
  -> 192.168.178.4:8080           Masq    1      0          0         
  -> 192.168.238.68:8080          Masq    1      0          0         
TCP  192.168.121.103:30806 rr
  -> 192.168.178.4:8080           Masq    1      0          0         
  -> 192.168.238.68:8080          Masq    1      0          0         
TCP  192.168.238.64:30806 rr
  -> 192.168.178.4:8080           Masq    1      0          0         
  -> 192.168.238.68:8080          Masq    1      0          0   
```
IP Virtual server listening on the Service Cluster IP. (Same as the type=Cluster)
```bash
TCP  10.104.87.249:80 rr
  -> 192.168.178.4:8080           Masq    1      0          0         
  -> 192.168.238.68:8080          Masq    1      0          0   
```