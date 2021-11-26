- [Linux namespace is used for isolatiion](#linux-namespace-is-used-for-isolatiion)
- [7(8 if Linux kernal 5.6+) Namespace types](#78-if-linux-kernal-56-namespace-types)
  - [Peers or Hierarchies](#peers-or-hierarchies)
  - [The namespaces that supports being manipulated by setns](#the-namespaces-that-supports-being-manipulated-by-setns)
  - [The limits on # of namespaces **per user** can be created](#the-limits-on--of-namespaces-per-user-can-be-created)
  - [The mount namespace (-m)](#the-mount-namespace--m)
  - [PID and IPC namespaces (-ip)](#pid-and-ipc-namespaces--ip)
  - [Network and UTS namespaces (-un)](#network-and-uts-namespaces--un)
  - [User namespaces (-Ur)](#user-namespaces--ur)
  - [All together](#all-together)
- [Common used commands](#common-used-commands)
  - [list the ns by type](#list-the-ns-by-type)
  - [find ns by pid](#find-ns-by-pid)
  - [Enter the ns by pid](#enter-the-ns-by-pid)
    - [When to use it?](#when-to-use-it)
# Linux namespace is used for isolatiion
# 7(8 if Linux kernal 5.6+) Namespace types
## Peers or Hierarchies
- Peers:
MNT (FS), UTS (Hostname domains), IPC (Seldom use), NET (Independant network)
- Hierarchies:
PID, User, Cgroup (hide system limits)
## The namespaces that supports being manipulated by setns
```bash
$ ll /proc/self/ns | awk '{print $1,$9,$10,$11}'
total   
dr-x--x--x ./  
dr-xr-xr-x ../  
lrwxrwxrwx cgroup -> cgroup:[4026531835]
lrwxrwxrwx ipc -> ipc:[4026531839]
lrwxrwxrwx mnt -> mnt:[4026531840]
lrwxrwxrwx net -> net:[4026531992]
lrwxrwxrwx pid -> pid:[4026531836]
lrwxrwxrwx pid_for_children -> pid:[4026531836]
lrwxrwxrwx user -> user:[4026531837]
lrwxrwxrwx uts -> uts:[4026531838]
```
## The limits on # of namespaces **per user** can be created
```bash
$ ll /proc/sys/user/
total 0
dr-xr-xr-x 1 root root 0 Oct 29 13:57 ./
dr-xr-xr-x 1 root root 0 Oct 28 19:46 ../
-rw-r--r-- 1 root root 0 Oct 29 13:57 max_cgroup_namespaces
-rw-r--r-- 1 root root 0 Oct 29 13:57 max_inotify_instances
-rw-r--r-- 1 root root 0 Oct 29 13:57 max_inotify_watches
-rw-r--r-- 1 root root 0 Oct 29 13:57 max_ipc_namespaces
-rw-r--r-- 1 root root 0 Oct 29 13:57 max_mnt_namespaces
-rw-r--r-- 1 root root 0 Oct 29 13:57 max_net_namespaces
-rw-r--r-- 1 root root 0 Oct 29 13:57 max_pid_namespaces
-rw-r--r-- 1 root root 0 Oct 29 13:57 max_user_namespaces
-rw-r--r-- 1 root root 0 Oct 29 13:57 max_uts_namespaces
$ # the default value is half of kernel/threads-max
$ cat /proc/sys/kernel/threads-max 
95147
$ cat /proc/sys/user/max_mnt_namespaces 
47573
```
## The mount namespace (-m)
The mount namespace stores the mount table.
```bash
sh1 $ # check filesystems supported by kernal
sh1 $cat /proc/filesystems 
nodev	sysfs # pseudo-filesystem for exporting kernal objects
nodev	tmpfs # virtual memory filesystem
nodev	proc # process information pseudo-filesystem for exporting kernal objects
...
sh1 $ # check mountinfo by pid 
sh1 $cat /proc/self/mountinfo 
26 32 0:24 / /sys rw,nosuid,nodev,noexec,relatime shared:7 - sysfs sysfs rw
27 32 0:5 / /proc rw,nosuid,nodev,noexec,relatime shared:14 - proc proc rw
...
```
> For container <br> When our sandboxed environment runs in a new Mount Namespace, it can mount filesystems not present on the host.
```bash
sh1 $ # create a new mount namespace and start a Bash shell inside that namespace.
sh1 $ sudo unshare -m /bin/bash
$mount -t tmpfs tmpfs /mnt
$cat /proc/self/mountinfo | grep mnt
801 723 0:75 / /mnt rw,relatime - tmpfs tmpfs rw
$ # tmpfs at /mnt as private. means no shared nor master options in mountinfo
```
## PID and IPC namespaces (-ip)
The PID namespace allows a process and its children to run in a new process tree that maps back to the host process tree. The new PID namespace starts with PID 1 which will map to a much higher PID in the host’s native PID namespace. 

The Inter-Process Communication (IPC) Namespace limits the processes ability to share memory.

```bash
sh1 $ # create a new PID and IPC namespaces with fork
sh1 $sudo unshare -mipf /bin/bash
 $ # verify current PID is 1
 $echo $$
1
 $ # check ps 
 $ps aux | head -n5
USER         PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root           1  0.0  0.0 167352 11200 ?        Ss   Oct28   0:02 /sbin/init
root           2  0.0  0.0      0     0 ?        S    Oct28   0:00 [kthreadd]
root           3  0.0  0.0      0     0 ?        I<   Oct28   0:00 [rcu_gp]
root           4  0.0  0.0      0     0 ?        I<   Oct28   0:00 [rcu_par_gp]
 $ # ? the PID==1 is systemd(init)? why? because ps reads from /proc, which still fork the native mount namespace. so we need to use --mount-proc option, or mount a new /proc to this new PID namespace
 ```
 Let's fix above problem
 ```bash
 sh1 $ # unshare with --mount-proc
 sh1 $sudo unshare --mount-proc -mipf /bin/bash
 $ps aux | head -n5
USER         PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root           1  0.0  0.0  10864  5140 pts/1    S    19:25   0:00 /bin/bash
root          10  0.0  0.0  11492  3408 pts/1    R+   19:26   0:00 ps aux
root          11  0.0  0.0   8092   528 pts/1    S+   19:26   0:00 head -n5
$mount | grep proc
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
systemd-1 on /proc/sys/fs/binfmt_misc type autofs (rw,relatime,fd=28,pgrp=0,timeout=0,minproto=5,maxproto=5,direct,pipe_ino=1348)
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
$echo $$
1
```
## Network and UTS namespaces (-un)
The Network Namespace allows a new network stack to exist in the sandbox. This means our sandboxed environment can have its own network interfaces, routing tables, DNS lookup servers, IP addresses, subnets…​ 
The UTS Namespace exists solely for storing the system’s domain name and hostname. 
We are going to use application *slirp4netns* to create a tap inside the container's network namespace and attached to it. slirp4netns provides user-mode networking ("slirp") for unprivileged network namespaces.
```bash
sh1 $ # on host, configure to allow rootless container to ping through the physic network
sh1 $ sysctl -w "net.ipv4.ping_group_range=0 2000000"
sh1 $unshare -unUr /bin/bash
$echo $$
5039
$# configure the tap on the host
sh2 $slirp4netns -c 5039 tap0
sent tapfd=5 for tap0
received tapfd=5
Starting slirp
* MTU:             1500
* Network:         10.0.2.0
* Netmask:         255.255.255.0
* Gateway:         10.0.2.2
* DNS:             10.0.2.3
* Recommended IP:  10.0.2.100
WARNING: 127.0.0.1:* on the host is accessible as 10.0.2.2 (set --disable-host-loopback to prohibit connecting to 127.0.0.1:*)
$ip a 
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host 
       valid_lft forever preferred_lft forever
3: tap0: <BROADCAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UNKNOWN group default qlen 1000
    link/ether 1a:38:5d:66:5c:46 brd ff:ff:ff:ff:ff:ff
    inet 10.0.2.100/24 brd 10.0.2.255 scope global tap0
       valid_lft forever preferred_lft forever
    inet6 fe80::1838:5dff:fe66:5c46/64 scope link 
       valid_lft forever preferred_lft forever

```
## User namespaces (-Ur)
The user namespace allows non-root users to run containers without a heavy root daemon.
If a new User Namespace is created, then that same user can now spawn any type of namespace.
```bash
sh1 $whoami
vagrant
sh1 $unshare -Ur /bin/bash
$whoami
root
$cat /proc/self/uid_map 
         0       1000          1
$id
uid=0(root) gid=0(root) groups=0(root),65534(nogroup)
$ # the files owned by the real root soutside of the sandbox is marked as nobody/nogroup
$ll /bin/bash
-rwxr-xr-x 1 nobody nogroup 1183448 Jun 18  2020 /bin/bash*
$ # write some files in the sandbox
$date > /tmp/test
$ll /tmp/test
-rw-rw-r-- 1 root root 32 Oct 29 19:44 /tmp/test
$ # check in the native namespace, the owner is mapped to host user(vagrant)
sh2 $ll /tmp/test 
-rw-rw-r-- 1 vagrant vagrant 32 Oct 29 19:44 /tmp/test
```
## All together
By create a new user namespace, this same user can now spawn any type of namespace. so no need root previlliage.
```bash
sh1 $unshare -mipfunUr /bin/bash
$mount -t proc proc /proc
$mount -t tmpfs none /tmp
$mount -t sysfs none /sys
$ #network related
$hostname sandbox
$echo "nameserver 10.0.2.3" > /tmp/resolv.conf
$mount --bind /tmp/resolv.conf /etc/resolv.conf
$ping google.com
PING google.com (142.251.45.46) 56(84) bytes of data.
64 bytes from dfw25s47-in-f14.1e100.net (142.251.45.46): icmp_seq=1 ttl=255 time=29.1 ms
```
# Common used commands
## list the ns by type
```bash
sh1 $sudo lsns --type net
        NS TYPE NPROCS   PID USER     NETNSID NSFS                           COMMAND
4026531992 net     123     1 root  unassigned                                /sbin/init
4026532232 net       1   551 root  unassigned                                /usr/sbin/haveged --Foreground --verbose=1 -w 1024
4026532297 net       1  1548 65532          0 /run/docker/netns/bfcc0e7b6227 /httpserver
4026532419 net       1  6602 root           1 /run/docker/netns/26e9fd31ed65 bash
```
## find ns by pid
```bash
sh1 $sudo ls -al /proc/1548/ns/
total 0
dr-x--x--x 2 65532 65532 0 Oct 28 19:49 .
dr-xr-xr-x 9 65532 65532 0 Oct 28 19:49 ..
lrwxrwxrwx 1 65532 65532 0 Oct 29 19:06 cgroup -> 'cgroup:[4026531835]'
lrwxrwxrwx 1 65532 65532 0 Oct 29 19:06 ipc -> 'ipc:[4026532294]'
lrwxrwxrwx 1 65532 65532 0 Oct 29 19:06 mnt -> 'mnt:[4026532292]'
lrwxrwxrwx 1 65532 65532 0 Oct 28 19:49 net -> 'net:[4026532297]'
lrwxrwxrwx 1 65532 65532 0 Oct 28 20:03 pid -> 'pid:[4026532295]'
lrwxrwxrwx 1 65532 65532 0 Oct 30 13:00 pid_for_children -> 'pid:[4026532295]'
lrwxrwxrwx 1 65532 65532 0 Oct 29 19:06 user -> 'user:[4026531837]'
lrwxrwxrwx 1 65532 65532 0 Oct 29 19:06 uts -> 'uts:[4026532293]'
```
## Enter the ns by pid
### When to use it?
Sometimes the container image(distroless) doesn't provide sh:
1. we can use nsenter to enter the namespace to check
2. we can run another container which provide debug utils, let the new container run in the same namespace.
```bash
sh1 $sudo nsenter --target 1548 --net ip a s
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
7: eth0@if8: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever
```