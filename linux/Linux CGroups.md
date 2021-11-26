- [What is CGroups?](#what-is-cgroups)
- [Lab](#lab)
  - [Create a child cgroup](#create-a-child-cgroup)
  - [Use cgroup to limit CPU resource usage](#use-cgroup-to-limit-cpu-resource-usage)
  - [Use cgroup to limit Memory](#use-cgroup-to-limit-memory)
- [Equivalent of Docker usage](#equivalent-of-docker-usage)
# What is CGroups?
CGroups(Control Groups) allows processes to be organized into **hierarchical groups** whose usage of various types of resources can then be limited and monitored. 
The cgroup interface is provided through pseudo-filesystem **cgroupfs**. 
# Lab
## Create a child cgroup
A child cgroup can be created by creating a sub-directory *makedir $CGROUP_NAME*, the cgroup driver will generate the related files.
```bash
sh1 $pwd
/sys/fs/cgroup/cpu
sh1 $sudo mkdir cputest
sh1 $cd cputest/
sh1 $ll
total 0
drwxr-xr-x 2 root root 0 Oct 30 16:58 ./
dr-xr-xr-x 6 root root 0 Oct 28 19:46 ../
-rw-r--r-- 1 root root 0 Oct 30 16:58 cgroup.clone_children
-rw-r--r-- 1 root root 0 Oct 30 16:58 cgroup.procs
-r--r--r-- 1 root root 0 Oct 30 16:58 cpuacct.stat
-rw-r--r-- 1 root root 0 Oct 30 16:58 cpuacct.usage
-r--r--r-- 1 root root 0 Oct 30 16:58 cpuacct.usage_all
-r--r--r-- 1 root root 0 Oct 30 16:58 cpuacct.usage_percpu
-r--r--r-- 1 root root 0 Oct 30 16:58 cpuacct.usage_percpu_sys
-r--r--r-- 1 root root 0 Oct 30 16:58 cpuacct.usage_percpu_user
-r--r--r-- 1 root root 0 Oct 30 16:58 cpuacct.usage_sys
-r--r--r-- 1 root root 0 Oct 30 16:58 cpuacct.usage_user
-rw-r--r-- 1 root root 0 Oct 30 16:58 cpu.cfs_period_us
-rw-r--r-- 1 root root 0 Oct 30 16:58 cpu.cfs_quota_us
-rw-r--r-- 1 root root 0 Oct 30 16:58 cpu.shares
-r--r--r-- 1 root root 0 Oct 30 16:58 cpu.stat
-rw-r--r-- 1 root root 0 Oct 30 16:58 cpu.uclamp.max
-rw-r--r-- 1 root root 0 Oct 30 16:58 cpu.uclamp.min
-rw-r--r-- 1 root root 0 Oct 30 16:58 notify_on_release
-rw-r--r-- 1 root root 0 Oct 30 16:58 tasks
```
## Use cgroup to limit CPU resource usage
```bash
cd /sys/fs/cgroup/cpu/cputest
echo $PID > cgroup.procs
echo 10000 > cpu.cfs_quota_us # default is -1 means no limit
```
## Use cgroup to limit Memory
```bash
cd /sys/fs/cgroup/memory/memtest
echo $PID > cgroup.procs
echo 10240000 > memory.limit_in_bytes # default is -1 means no limit
```
# Equivalent of Docker usage
```bash
sh1 $docker run -it --rm --cpu-quota 10000 ubuntu:20.04
sh2 $docker ps
CONTAINER ID   IMAGE                   COMMAND         CREATED          STATUS          PORTS                    NAMES
dbbcc2bb1758   ubuntu:20.04            "bash"          28 seconds ago   Up 27 seconds                            busy_tharp
sh2 $docker inspect dbbcc2bb1758 | grep Pid
            "Pid": 8276,
            "PidMode": "",
            "PidsLimit": null,
sh2 $pwd
/sys/fs/cgroup/cpu/docker     
sh2 $cat dbbcc2bb1758cd9499406d9a6261f1629d1ad76add26be315a94b1d99180aea7/cgroup.procs 
8276
sh2 $cat dbbcc2bb1758cd9499406d9a6261f1629d1ad76add26be315a94b1d99180aea7/cpu.cfs_quota_us 
10000
```