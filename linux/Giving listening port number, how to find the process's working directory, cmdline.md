## Giving listening port number, how to find the process's working directory, cmdline? 
### Key steps: 
1. find the process info by port number. (ss +/ grep)
```shell script
[root@localhost ~]# ss -ltunp '(  sport = 22  )'
Netid State      Recv-Q Send-Q                                                Local Address:Port                                                               Peer Address:Port              
tcp   LISTEN     0      128                                                               *:22                                                                            *:*                   users:(("sshd",pid=628,fd=3))
tcp   ESTAB      0      0                                                         10.0.2.15:22                                                                     10.0.2.2:49679               users:(("sshd",pid=1985,fd=3),("sshd",pid=1982,fd=3))
tcp   ESTAB      0      0                                                         10.0.2.15:22                                                                     10.0.2.2:56080               users:(("sshd",pid=2504,fd=3),("sshd",pid=2501,fd=3))
tcp   LISTEN     0      128                                                            [::]:22                                                                         [::]:*                   users:(("sshd",pid=628,fd=4))
# more refined, ss filter expression
[root@localhost ~]# ss -atupn state ESTABLISHED '( dport = 22 or sport = 22 )' 
Netid Recv-Q Send-Q                                                     Local Address:Port                                                                    Peer Address:Port              
tcp   0      0                                                              10.0.2.15:22                                                                          10.0.2.2:49679               users:(("sshd",pid=1985,fd=3),("sshd",pid=1982,fd=3))
tcp   0      0                                                              10.0.2.15:22                                                                          10.0.2.2:56080               users:(("sshd",pid=2504,fd=3),("sshd",pid=2501,fd=3))
```
2. find the process id. (awk)
```shell script
awk 'match($NF,/pid=([0-9]+)/,m) { print m[1] }'
```
3. ls -l /proc/{}/cwd
### Final Answers
```shell script
[root@localhost ~]# ss -ltunp '(  sport = 22  )' | awk 'match($NF,/pid=([0-9]+)/,m) { print m[1] }' | uniq | xargs -i ls -l /proc/{}/cwd /proc/{}/exe
lrwxrwxrwx. 1 root root 0 Sep  2 21:43 /proc/628/cwd -> /
lrwxrwxrwx. 1 root root 0 Sep  1 21:12 /proc/628/exe -> /usr/sbin/sshd
```
### Other notes
1. ss -p, if not the root user, then pid might not be able to print
2. use uniq to filter the duplicate of ipv4/ipv6. or, we can use 
```shell script
# option 4 for ipv4 
ss -ltunp4 '(  sport = 22  )'
# option 6 for ipv6
ss -ltunp6 '(  sport = 22  )'
```
3. use the awk match for regex grouping
4. ls doesn't read arguments from stdin, so we need to use 
```shell script
xargs -i
```
**5. The file structure of /proc/{pid}/**
```shell script
[root@localhost ~]# ss -ltunp '(  sport = 22  )' | awk 'match($NF,/pid=([0-9]+)/,m) { print m[1] }' | uniq | xargs -i ls -l /proc/{}/cwd /proc/{}/exe /proc/{}/cmdline
-r--r--r--. 1 root root 0 Sep  1 21:12 /proc/628/cmdline
lrwxrwxrwx. 1 root root 0 Sep  2 21:43 /proc/628/cwd -> /
lrwxrwxrwx. 1 root root 0 Sep  1 21:12 /proc/628/exe -> /usr/sbin/sshd
[root@localhost ~]# cat /proc/628/cmdline
/usr/sbin/sshd-D-u0[root@localhost ~]# 
```
