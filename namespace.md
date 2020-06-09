
Linux Namespace

# docker namespaces

## UTS Namespaces

Cloneflags: syscall.CLONE_NEWUTS

### uts_ns.go

```go
package main

import (
        "log"
        "os"
        "os/exec"
        "syscall"
)

func main() {
        cmd := exec.Command("sh")
        cmd.SysProcAttr = &syscall.SysProcAttr{
                Cloneflags: syscall.CLONE_NEWUTS,
        }
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.Env = []string{"PS1=-[ns-demo]- # "}
        if err := cmd.Run(); err != nil {
                log.Fatal(err)
        }
}
```

### steps

```bash
 # cd training/docker-demo/namespaces/
-[root@chenqiang-dev namespaces]# go run uts_ns.go 
-[ns-demo]- # pstree -pl

-[ns-demo]- # echo $$
31636
-[ns-demo]- # readlink /proc/31636/ns/uts
uts:[4026532212]
-[ns-demo]- # readlink /proc/31631/ns/uts
uts:[4026531838]
-[ns-demo]- # hostname
chenqiang-dev.novalocal
You have new mail in /var/mail/root
-[ns-demo]- # hostname -b demo
-[ns-demo]- # hostname
demo
-[ns-demo]- # exit
exit
-[root@chenqiang-dev namespaces]# hostname
chenqiang-dev.novalocal

-[root@chenqiang-dev namespaces]#  ls -lth /var/spool/mail/
-[root@chenqiang-dev namespaces]# cat /dev/null > /var/spool/mail/root
```
---

## IPC Namespaces

### ipc_ns.go

```go
package main

import (
        "log"
        "os"
        "os/exec"
        "syscall"
)

func main() {
        cmd := exec.Command("sh")
        cmd.SysProcAttr = &syscall.SysProcAttr{
                Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC,
        }
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.Env = []string{"PS1=-[ns-demo]- # "}
        if err := cmd.Run(); err != nil {
                log.Fatal(err)
        }
}
```

### steps

```bash
-[ns-demo]- # echo $$
3930
-[ns-demo]- # ipcs -q

------ Message Queues --------
key        msqid      owner      perms      used-bytes   messages    

-[ns-demo]- # ipcmk -Q
Message queue id: 0
-[ns-demo]- # ipcs -q            

------ Message Queues --------
key        msqid      owner      perms      used-bytes   messages    
0x892abf0c 0          root       644        0            0           

-[ns-demo]- # exit
exit
-[root@chenqiang-dev namespaces]# ipcs -q

------ Message Queues --------
key        msqid      owner      perms      used-bytes   messages    

```

## PID Namespace

syscall.CLONE_NEWPID

## steps

```bash
-[ns-demo]- # pstree -pl 
-[ns-demo]- # echo $$
```

### MOUNT Namepace

syscall.CLONE_NEWNS

### steps

```bash
-[root@chenqiang-dev namespaces]# go run mnt_ns.go 
-[ns-demo]- # ls /proc
-[ns-demo]- # mount -t proc proc /proc
-[ns-demo]- # ls /proc
```

### USER Namepaces

Centos中对user ns的支持需要开启并重启：
`grubby --args="user_namespace.enable=1" --update-kernel="$(grubby --default-kernel)"`