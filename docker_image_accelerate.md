在 docker 中，只能够为 docker 官方镜像仓库 hub.docker.com 提供 mirror 加速，自建的 docker 镜像仓库是不能配置 mirror 加速的；

下面是配置镜像加速的方法：

# 修改镜像仓库 mirror 地址

1、修改 /etc/docker/daemon.json 文件（如果没有，则创建）：

```bash
vim /etc/docker/daemon.json
```

2、添加 registry-mirrors 字段

```json
{
 "registry-mirrors": ["https://registry.cn-hangzhou.aliyuncs.com"]
}

```

3、重启 docker/kubelet

```bash
systemctl daemon-reload
systemctl restart docker
systemctl start kubelet # 假设您安装了 kubenetes
```

# 查看修改结果

执行 `docker info`
![IMAGE](assets/docker_info)

# 其他镜像源

```bash
# Docker中国 mirror
# export REGISTRY_MIRROR="https://registry.docker-cn.com"
# 腾讯云 docker hub mirror
# export REGISTRY_MIRROR="https://mirror.ccs.tencentyun.com"
# 华为云镜像
# export REGISTRY_MIRROR="https://05f073ad3c0010ea0f4bc00b7105ec20.mirror.swr.myhuaweicloud.com"
# DaoCloud 镜像
# export REGISTRY_MIRROR="http://f1361db2.m.daocloud.io"
# 阿里云 docker hub mirror
export REGISTRY_MIRROR=https://registry.cn-hangzhou.aliyuncs.com

```

# 配置举例

```bash
$ cat /etc/docker/daemon.json
{
    "log-level": "warn",
    "selinux-enabled": false,
    "insecure-registries": ["0.0.0.0/0"],
    "registry-mirrors": ["https://registry.cn-hangzhou.aliyuncs.com"],
    "max-concurrent-downloads": 10,
    "max-concurrent-uploads": 10,
    "log-driver": "json-file",
    "log-opts": {
        "max-size": "200m",
        "max-file": "7"
    },
    "storage-driver": "overlay2",
    "storage-opts": [
        "overlay2.override_kernel_check=true"
    ],
    "live-restore": true,
    "exec-opts": [
        "native.cgroupdriver=cgroupfs"
    ]
}
```