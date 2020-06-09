# Cgroup demo

- Linux Cgroups
  	- what's the Linux Cgroups
  	- how docker use cgroup
  	- golang example to implement cgroup demo

# Install

如果没有 lssubsys 命令，则需要安装
yum install libcgroup-tools

```bash
-[root@chenqiang-dev namespaces]# lssubsys -a
cpuset
cpu,cpuacct
memory
devices
freezer
net_cls,net_prio
blkio
perf_event
hugetlb
pids

-[root@chenqiang-dev cgroups]# lssubsys -aMi
cpuset /sys/fs/cgroup/cpuset
cpu,cpuacct /sys/fs/cgroup/cpu,cpuacct
memory /sys/fs/cgroup/memory
devices /sys/fs/cgroup/devices
freezer /sys/fs/cgroup/freezer
net_cls,net_prio /sys/fs/cgroup/net_cls,net_prio
blkio /sys/fs/cgroup/blkio
perf_event /sys/fs/cgroup/perf_event
hugetlb /sys/fs/cgroup/hugetlb
pids /sys/fs/cgroup/pids

-[root@chenqiang-dev cgroups]# tree
.
└── cgroup-demo

1 directory, 0 files
-[root@chenqiang-dev cgroups]# mount -t cgroup -o none,name=cgroup-demo ./cgroup-demo
mount: can't find ./cgroup-demo in /etc/fstab
-[root@chenqiang-dev cgroups]# ls
cgroup-demo
-[root@chenqiang-dev cgroups]# mount -t cgroup -o none,name=cgroup-demo cgroup-demo ./cgroup-demo
-[root@chenqiang-dev cgroups]# ls ./cgroup-demo/
cgroup.clone_children  cgroup.event_control  cgroup.procs  cgroup.sane_behavior  notify_on_release  release_agent  tasks
-[root@chenqiang-dev cgroups]# tree
.
└── cgroup-demo
    ├── cgroup.clone_children
    ├── cgroup.event_control
    ├── cgroup.procs
    ├── cgroup.sane_behavior
    ├── notify_on_release
    ├── release_agent
    └── tasks

1 directory, 7 files


-[root@chenqiang-dev cgroup-demo]# cat cgroup.procs  | wc -l
200
-[root@chenqiang-dev cgroup-demo]# ps aux | wc -l
201

-[root@chenqiang-dev cgroup-demo]# ls
cgroup1  cgroup.clone_children  cgroup.procs          notify_on_release  tasks
cgroup2  cgroup.event_control   cgroup.sane_behavior  release_agent
-[root@chenqiang-dev cgroup-demo]# cd cgroup1
-[root@chenqiang-dev cgroup1]# ls
cgroup.clone_children  cgroup.event_control  cgroup.procs  notify_on_release  tasks
-[root@chenqiang-dev cgroup1]# echo $$
1801
-[root@chenqiang-dev cgroup1]# echo $$ >> tasks 
-[root@chenqiang-dev cgroup1]# cat /proc/1801/cgroup 
12:name=cgroup-demo:/cgroup1
11:freezer:/
10:blkio:/
9:hugetlb:/
8:perf_event:/
7:memory:/
6:pids:/
5:cpuacct,cpu:/
4:cpuset:/
3:net_prio,net_cls:/
2:devices:/
1:name=systemd:/user.slice/user-1001.slice/session-2.scope
-[root@chenqiang-dev cgroup1]# cat tasks 
1801
3426


-[root@chenqiang-dev cgroup1]# ls /sys/fs/cgroup/ 
blkio  cpuacct      cpuset   freezer  memory   net_cls,net_prio  perf_event  systemd
cpu    cpu,cpuacct  devices  hugetlb  net_cls  net_prio          pids


-[root@chenqiang-dev cgroup-demo]# mkdir memory
-[root@chenqiang-dev cgroup-demo]# cd memory/
-[root@chenqiang-dev memory]# ls
cgroup.clone_children  cgroup.event_control  cgroup.procs  notify_on_release  tasks
-[root@chenqiang-dev memory]# stress --vm-bytes 200m --vm-keep -m 1
stress: info: [3521] dispatching hogs: 0 cpu, 0 io, 1 vm, 0 hdd
bc <<< "31*0.6"

  PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND                                                     
 3739 root      20   0  212068 204800    132 R  99.7  0.6   0:13.76 stress  
```

## memory cgroups所在目录演示

```bash
-[root@chenqiang-dev memory]# pwd
/sys/fs/cgroup/memory
-[root@chenqiang-dev memory]# ls
cgroup.clone_children       memory.kmem.max_usage_in_bytes      memory.memsw.limit_in_bytes      memory.usage_in_bytes
cgroup.event_control        memory.kmem.slabinfo                memory.memsw.max_usage_in_bytes  memory.use_hierarchy
cgroup.procs                memory.kmem.tcp.failcnt             memory.memsw.usage_in_bytes      notify_on_release
cgroup.sane_behavior        memory.kmem.tcp.limit_in_bytes      memory.move_charge_at_immigrate  release_agent
demo-limit-memory           memory.kmem.tcp.max_usage_in_bytes  memory.numa_stat                 system.slice
docker                      memory.kmem.tcp.usage_in_bytes      memory.oom_control               tasks
memory.failcnt              memory.kmem.usage_in_bytes          memory.pressure_level            user.slice
memory.force_empty          memory.limit_in_bytes               memory.soft_limit_in_bytes
memory.kmem.failcnt         memory.max_usage_in_bytes           memory.stat
memory.kmem.limit_in_bytes  memory.memsw.failcnt                memory.swappiness
-[root@chenqiang-dev memory]# cd demo-limit-memory/

-[root@chenqiang-dev demo-limit-memory]# stress --vm-bytes 200m --vm-keep -m 1
stress: info: [4602] dispatching hogs: 0 cpu, 0 io, 1 vm, 0 hdd
^C
-[root@chenqiang-dev demo-limit-memory]# ls
cgroup.clone_children           memory.kmem.slabinfo                memory.memsw.failcnt             memory.soft_limit_in_bytes
cgroup.event_control            memory.kmem.tcp.failcnt             memory.memsw.limit_in_bytes      memory.stat
cgroup.procs                    memory.kmem.tcp.limit_in_bytes      memory.memsw.max_usage_in_bytes  memory.swappiness
memory.failcnt                  memory.kmem.tcp.max_usage_in_bytes  memory.memsw.usage_in_bytes      memory.usage_in_bytes
memory.force_empty              memory.kmem.tcp.usage_in_bytes      memory.move_charge_at_immigrate  memory.use_hierarchy
memory.kmem.failcnt             memory.kmem.usage_in_bytes          memory.numa_stat                 notify_on_release
memory.kmem.limit_in_bytes      memory.limit_in_bytes               memory.oom_control               tasks
memory.kmem.max_usage_in_bytes  memory.max_usage_in_bytes           memory.pressure_level
-[root@chenqiang-dev demo-limit-memory]# cat memory.limit_in_bytes 
104857600
-[root@chenqiang-dev demo-limit-memory]# echo $$
4557
-[root@chenqiang-dev demo-limit-memory]# cat tasks 
1801
-[root@chenqiang-dev demo-limit-memory]# echo $$ > tasks 
-[root@chenqiang-dev demo-limit-memory]# cat tasks 
1801
4557
4618
-[root@chenqiang-dev demo-limit-memory]# stress --vm-bytes 200m --vm-keep -m 1
stress: info: [4619] dispatching hogs: 0 cpu, 0 io, 1 vm, 0 hdd
stress: FAIL: [4619] (415) <-- worker 4620 got signal 9
stress: WARN: [4619] (417) now reaping child worker processes
stress: FAIL: [4619] (451) failed run completed in 0s
-[root@chenqiang-dev demo-limit-memory]# stress --vm-bytes 50m --vm-keep -m 1  
stress: info: [4621] dispatching hogs: 0 cpu, 0 io, 1 vm, 0 hdd
^C
-[root@chenqiang-dev demo-limit-memory]# stress --vm-bytes 100m --vm-keep -m 1 
stress: info: [4623] dispatching hogs: 0 cpu, 0 io, 1 vm, 0 hdd
stress: FAIL: [4623] (415) <-- worker 4624 got signal 9
stress: WARN: [4623] (417) now reaping child worker processes
stress: FAIL: [4623] (451) failed run completed in 0s
-[root@chenqiang-dev demo-limit-memory]# stress --vm-bytes 90m --vm-keep -m 1  
stress: info: [4625] dispatching hogs: 0 cpu, 0 io, 1 vm, 0 hdd
^C

```

## 看看docker 情况

```bash
-[appuser@chenqiang-dev workstation]$  docker run -itd -m 128m qzschen/nginx-hello
-[appuser@chenqiang-dev workstation]$ docker ps
CONTAINER ID        IMAGE                                                      COMMAND                  CREATED             STATUS              PORTS               NAMES
cc09e736e62a        qzschen/nginx-hello   "/bin/sh -c 'nginx -…"   4 seconds ago       Up 2 seconds        80/tcp              nostalgic_pike
-[appuser@chenqiang-dev workstation]$ sudo su
-[root@chenqiang-dev workstation]# cd /sys/fs/cgroup/memory/d
demo-limit-memory/ docker/            
-[root@chenqiang-dev workstation]# cd /sys/fs/cgroup/memory/d
demo-limit-memory/ docker/            
-[root@chenqiang-dev workstation]# cd /sys/fs/cgroup/memory/d
demo-limit-memory/ docker/            
-[root@chenqiang-dev workstation]# cd /sys/fs/cgroup/memory/docker/cc09e736e62abfb2d4c326f78cfdcb1e0b885d2d30e86ee4ab48f4319d3c6bd5/
-[root@chenqiang-dev cc09e736e62abfb2d4c326f78cfdcb1e0b885d2d30e86ee4ab48f4319d3c6bd5]# ls
cgroup.clone_children           memory.kmem.slabinfo                memory.memsw.failcnt             memory.soft_limit_in_bytes
cgroup.event_control            memory.kmem.tcp.failcnt             memory.memsw.limit_in_bytes      memory.stat
cgroup.procs                    memory.kmem.tcp.limit_in_bytes      memory.memsw.max_usage_in_bytes  memory.swappiness
memory.failcnt                  memory.kmem.tcp.max_usage_in_bytes  memory.memsw.usage_in_bytes      memory.usage_in_bytes
memory.force_empty              memory.kmem.tcp.usage_in_bytes      memory.move_charge_at_immigrate  memory.use_hierarchy
memory.kmem.failcnt             memory.kmem.usage_in_bytes          memory.numa_stat                 notify_on_release
memory.kmem.limit_in_bytes      memory.limit_in_bytes               memory.oom_control               tasks
memory.kmem.max_usage_in_bytes  memory.max_usage_in_bytes           memory.pressure_level
-[root@chenqiang-dev cc09e736e62abfb2d4c326f78cfdcb1e0b885d2d30e86ee4ab48f4319d3c6bd5]# cat memory.limit_in_bytes 
134217728
-[root@chenqiang-dev cc09e736e62abfb2d4c326f78cfdcb1e0b885d2d30e86ee4ab48f4319d3c6bd5]# bc <<< "134217728/1024/1024"
128
```

## CPU cfs

cfs_period_us用来配置时间周期长度，cfs_quota_us用来配置当前cgroup在设置的周期长度内所能使用的CPU时间数，两个文件配合起来设置CPU的使用上限。两个文件的单位都是微秒（us），cfs_period_us的取值范围为1毫秒（ms）到1秒（s），cfs_quota_us的取值大于1ms即可，如果cfs_quota_us的值为-1（默认值），表示不受cpu时间的限制。下面是几个例子：
1.限制只能使用1个CPU（每250ms能使用250ms的CPU时间）
```
    # echo 250000 > cpu.cfs_quota_us /* quota = 250ms */
    # echo 250000 > cpu.cfs_period_us /* period = 250ms */
 ```

2.限制使用2个CPU（内核）（每500ms能使用1000ms的CPU时间，即使用两个内核）
```
    # echo 1000000 > cpu.cfs_quota_us /* quota = 1000ms */
    # echo 500000 > cpu.cfs_period_us /* period = 500ms */
```

3.限制使用1个CPU的20%（每50ms能使用10ms的CPU时间，即使用一个CPU核心的20%）
```
    # echo 10000 > cpu.cfs_quota_us /* quota = 10ms */
    # echo 50000 > cpu.cfs_period_us /* period = 50ms */
```

# cpu.shares

shares用来设置CPU的相对值，并且是针对所有的CPU（内核），默认值是1024，假如系统中有两个cgroup，分别是A和B，A的shares值是1024，B的shares值是512，那么A将获得1024/(1204+512)=66%的CPU资源，而B将获得33%的CPU资源。shares有两个特点：

- 如果A不忙，没有使用到66%的CPU时间，那么剩余的CPU时间将会被系统分配给B，即B的CPU使用率可以超过33%

- 如果添加了一个新的cgroup C，且它的shares值是1024，那么A的限额变成了1024/(1204+512+1024)=40%，B的变成了20%

从上面两个特点可以看出：

- 在闲的时候，shares基本上不起作用，只有在CPU忙的时候起作用，这是一个优点。

- 由于shares是一个绝对值，需要和其它cgroup的值进行比较才能得到自己的相对限额，而在一个部署很多容器的机器上，cgroup的数量是变化的，所以这个限额也是变化的，自己设置了一个高的值，但别人可能设置了一个更高的值，所以这个功能没法精确的控制CPU使用率。

# cpu.stat

包含了下面三项统计结果

- nr_periods： 表示过去了多少个cpu.cfs_period_us里面配置的时间周期

- nr_throttled： 在上面的这些周期中，有多少次是受到了限制（即cgroup中的进程在指定的时间周期中用光了它的配额）

- throttled_time: cgroup中的进程被限制使用CPU持续了多长时间(纳秒)

# demo again

```bash
-[appuser@chenqiang-dev cgroup-demo]$ cd /sys/fs/cgroup/cpu,cpuacct/
-[appuser@chenqiang-dev cpu,cpuacct]$ ls
cgroup.clone_children  cpuacct.stat          cpu.cfs_quota_us   cpu.stat           system.slice
cgroup.event_control   cpuacct.usage         cpu.rt_period_us   docker             tasks
cgroup.procs           cpuacct.usage_percpu  cpu.rt_runtime_us  notify_on_release  user.slice
cgroup.sane_behavior   cpu.cfs_period_us     cpu.shares         release_agent
-[appuser@chenqiang-dev cpu,cpuacct]$ sudo mkdir demo
-[appuser@chenqiang-dev cpu,cpuacct]$ ls
cgroup.clone_children  cgroup.sane_behavior  cpuacct.usage_percpu  cpu.rt_period_us   cpu.stat  notify_on_release  tasks
cgroup.event_control   cpuacct.stat          cpu.cfs_period_us     cpu.rt_runtime_us  demo      release_agent      user.slice
cgroup.procs           cpuacct.usage         cpu.cfs_quota_us      cpu.shares         docker    system.slice
-[appuser@chenqiang-dev cpu,cpuacct]$ cd demo/
-[appuser@chenqiang-dev demo]$ ls
cgroup.clone_children  cgroup.procs  cpuacct.usage         cpu.cfs_period_us  cpu.rt_period_us   cpu.shares  notify_on_release
cgroup.event_control   cpuacct.stat  cpuacct.usage_percpu  cpu.cfs_quota_us   cpu.rt_runtime_us  cpu.stat    tasks
-[appuser@chenqiang-dev demo]$ cat cpu.cfs_period_us 
100000
-[appuser@chenqiang-dev demo]$ cat cpu.cfs_quota_us 
-1
-[appuser@chenqiang-dev demo]$ cat cgroup.procs 
-[appuser@chenqiang-dev demo]$ cat tasks 
-[appuser@chenqiang-dev demo]$ sudo bash -c "echo 50000 > cpu.cfs_quota_us"
-[appuser@chenqiang-dev demo]$ cat cpu.cfs_quota_us 
50000
-[appuser@chenqiang-dev demo]$ echo $$
3737
-[appuser@chenqiang-dev demo]$ sudo bash -c "echo $$ > cgroups.procs"
bash: cgroups.procs: Permission denied
-[appuser@chenqiang-dev demo]$ sudo bash -c "echo $$ > cgroup.procs" 
-[appuser@chenqiang-dev demo]$ cat cgroup.procs 
3737
4144
-[appuser@chenqiang-dev demo]$ while :; do echo test>/dev/null;done
^C
-[appuser@chenqiang-dev demo]$ ls
cgroup.clone_children  cgroup.procs  cpuacct.usage         cpu.cfs_period_us  cpu.rt_period_us   cpu.shares  notify_on_release
cgroup.event_control   cpuacct.stat  cpuacct.usage_percpu  cpu.cfs_quota_us   cpu.rt_runtime_us  cpu.stat    tasks
-[appuser@chenqiang-dev demo]$ cat cpu.s
cpu.shares  cpu.stat    
-[appuser@chenqiang-dev demo]$ cat cpu.s
cpu.shares  cpu.stat    
-[appuser@chenqiang-dev demo]$ cat cpu.stat 
nr_periods 245
nr_throttled 216
throttled_time 3119108009
-[appuser@chenqiang-dev demo]$ 
```

![IMAGE](assets/cpu_50_usage)