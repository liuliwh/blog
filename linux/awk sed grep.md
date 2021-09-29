## Prerequisite
```shell script
[vagrant@localhost ~]$ cat ip.log 
192.168.0.1 vagrant
192.168.0.2 tiantian
192.168.0.3 bingbing
192.168.0.4 tingting
192.168.0.4 vagrant
```
### Print the 1st col if the 2nd col containing text 'vagrant'
```shell script
[vagrant@localhost ~]$ awk '$2~/vagrant/ { print $1 }' ./ip.log 
192.168.0.1
192.168.0.4
[vagrant@localhost ~]$ awk 'match($2,/vagrant/,m) { print $1 }'  ./ip.log 
192.168.0.1
192.168.0.4
```
shell script verbose version (not recommended)
```shell script
[vagrant@localhost ~]$ cat findvagrant.sh 
#! /bin/bash
while read line
do
  name=`echo $line | awk '{ print $2 }'`    
  if [[ "$name" =~ "vagrant" ]]; then
    echo $line | awk '{ print $1 }'
  fi
done < ~/ip.log
[vagrant@localhost ~]$ sh findvagrant.sh 
192.168.0.1
192.168.0.4
```
sed solution (not accurate)
```shell script
[vagrant@localhost ~]$ # sed address line: all lines with regex match /vagrant/
[vagrant@localhost ~]$ # sed command: substitude s#old#new#, group backreference \1
[vagrant@localhost ~]$ sed -nr '/vagrant/s#^([^ ]+) (.*)$#\1#gp' ./ip.log 
192.168.0.1
192.168.0.4
```
grep + cut solution (not accurate)
```shell script
[vagrant@localhost ~]$ grep '[0-9.] *vagrant*' ./ip.log | cut -d' ' -f1
192.168.0.1
192.168.0.4
```
### add () if char is a letter char
sed solution with & back reference
```shell script
[vagrant@localhost ~]$ sed -nr 's#[a-z]#(&)#gp' ./ip.log 
192.168.0.1 (v)(a)(g)(r)(a)(n)(t)
192.168.0.2 (t)(i)(a)(n)(t)(i)(a)(n)
192.168.0.3 (b)(i)(n)(g)(b)(i)(n)(g)
192.168.0.4 (t)(i)(n)(g)(t)(i)(n)(g)
192.168.0.4 (v)(a)(g)(r)(a)(n)(t)
```
## Notes
awk is suitable for 
1. the column based file
2. calculation and statistic 
