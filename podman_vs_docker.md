本文将介绍三个符合 CRI 标准的容器工具 Podman、 Buildah 和 Skopeo。这三个工具都是基于 *nix 传统的 fork-exec 模型，解决了由于 Docker 守护程序导致的启动和安全问题，提高了容器的性能和安全。

---

# 专注于OCI容器管理的工具

## 1、Docker

Docker 是目前最流行的 Linux 容器解决方案，即使 Docker 是目前管理 Linux 容器的一个非常方便的工具，但它也有两个缺点：

- Docker 需要在你的系统上运行一个守护进程。

- Docker 是以 root 身份在你的系统上运行该守护程序。

这些缺点的存在可能有一定的安全隐患，为了解决这些问题，下一代容器化工具 Podman 出现了 。

## 2、Podman

Podman 是一个开源的容器运行时项目，可在大多数 Linux 平台上使用。Podman 提供与 Docker 非常相似的功能。

Podman 可以管理和运行任何符合 OCI（Open Container Initiative）规范的容器和容器镜像。Podman 提供了一个与 Docker 兼容的命令行前端来管理 Docker 镜像。

- Podman 官网地址：https://podman.io/
- Podman 项目地址：https://github.com/containers/libpod

## 3、Podman 目前已支持大多数发行版本通过软件包来进行安装

```bash
$ sudo yum -y install podman
```

## 4、Podman 常用命令

```bash
$ podman run -dt -p 8080:80/tcp qzschen/nginx-hello
$ podman ps -a
$ podman inspect -l | grep IPAddress
$ sudo podman logs --latest
$ sudo podman top <container_id>
$ sudo podman stop --latest
$ sudo podman rm --latest
$ sudo podman container checkpoint <container_id>
$ sudo podman container restore <container_id>
```

## 5、迁移容器

Podman 支持将容器从一台机器迁移到另一台机器。

首先，在源机器上对容器设置检查点，并将容器打包到指定位置。

```bash
$ sudo podman container checkpoint <container_id> -e /tmp/checkpoint.tar.gz
$ scp /tmp/checkpoint.tar.gz <destination_system>:/tmp
```

其次，在目标机器上使用源机器上传输过来的打包文件对容器进行恢复。

`$ sudo podman container restore -i /tmp/checkpoint.tar.gz`

## 6、配置别名

如果习惯了使用 Docker 命令，可以直接给 Podman 配置一个别名来实现无缝转移。你只需要在 .bashrc 下加入以下行内容即可：

```bash
$ echo "alias docker=podman" >> .bashrc
$ source .bashrc

```

## 7、Podman 如何实现开机重启容器

由于 Podman 不再使用守护进程管理服务，所以不能通过守护进程去实现自动重启容器的功能。那如果要实现开机自动重启容器，又该如何实现呢？

其实方法很简单，现在大多数系统都已经采用 Systemd 作为守护进程管理工具。这里我们就可以使用 Systemd 来实现 Podman 开机重启容器，这里我们以启动一个 Nginx 容器为例子。

首先，我们先运行一个 Nginx 容器。

`$ sudo podman run -t -d -p 80:80 --name nginx qzschen/nginx-hello`

然后，在建立一个 Systemd 服务配置文件。

```bash
$ vim /etc/systemd/system/nginx_container.service

[Unit]
Description=Podman Nginx Service
After=network.target
After=network-online.target

[Service]
Type=simple
ExecStart=/usr/bin/podman start -a nginx
ExecStop=/usr/bin/podman stop -t 10 nginx
Restart=always

[Install]
WantedBy=multi-user.target
```

接下来，启用这个 Systemd 服务。

```bash
$ sudo systemctl daemon-reload
$ sudo systemctl enable nginx_container.service
$ sudo systemctl start nginx_container.service
```

---

# 专注于构建 OCI 容器镜像的工具

Podman 只是 OCI 容器生态系统计划中的一部分，主要专注于帮助用户维护和修改符合 OCI 规范的容器镜像。其它的组件还有 Buildah、Skopeo 等。

虽然 Podman 也可以支持用户构建 Docker 镜像，但是构建速度比较慢。并且默认情况下使用 VFS 存储驱动程序会消耗大量磁盘空间。

Buildah 是一个专注于构建 OCI 容器镜像的工具，Buildah 构建速度非常快并使用覆盖存储驱动程序，可以节约大量的空间。

Buildah 基于 fork-exec 模型，不以守护进程运行。Buildah 支持 Dockerfile 中的所有命令。你可以直接使用 Dockerfiles 来构建镜像，并且不需要任何 root 权限。Buildah 也支持用自己的语法文件构建镜像，可以允许将其他脚本语言集成到构建过程中。

Buildah 项目地址 https://github.com/containers/buildah

---

# 专注于镜像管理工具

Skopeo 是一个镜像管理工具，允许我们通过 Push、Pull和复制镜像来处理 Docker 和符合 OCI 规范的镜像。

Skopeo 项目地址 https://github.com/containers/skopeo