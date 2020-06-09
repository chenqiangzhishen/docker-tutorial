package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

// 挂载了memory subsystem的hierarchy的根目录位置
const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func main() {
	if os.Args[0] == "/proc/self/exe" {
		// 容器进程
		fmt.Printf("current pid %d\n", syscall.Getpid())
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	} else {
		// 得到fork出来进程映射在外部命名空间的pid
		fmt.Printf("%v\n", cmd.Process.Pid)
		// 在系统默认创建挂载了 memory subsystem 的hierarchy上创建cgroup
		os.Mkdir(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit"), 0755)
		// 将容器进程加入到这个cgroup中
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		// 限制cgroup进程使用
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "memory.limit_in_bytes"), []byte("100m"), 0644)
		cmd.Process.Wait()
	}
}

// -[root@chenqiang-dev appuser]# ll /proc/self/exe
// lrwxrwxrwx 1 root root 0 Jun  9 23:55 /proc/self/exe -> /usr/bin/ls
// -[root@chenqiang-dev appuser]# go run demo-ns-cgroup.go /proc/self/exe &
// current pid 1
// stress: info: [6] dispatching hogs: 0 cpu, 0 io, 1 vm, 0 hdd

// 用 top 看看
// PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND
// 25448 root      20   0  212068 204804    132 R 100.0  0.6   0:17.63 stress
