## 容器

## cgroups v1 vs v2

cgroups （control group）是一种以 hierarchical（树形层级）方式组织进程的机制（a mechanism to organize processes hierarchically），以及在层级中以受控和 可配置的方式（controlled and configurable manner）分发系统资源 （distribute system resources）。
简单理解就是，cgroup 通过 vfs 将功能暴露给用户，因为用户可以通过创建文件和文件目录来创建，修改和删除cgroups。

### 区别

linux在4.17正式发布了 cgroup v2，v2和v1的区别主要是：
- v2支持线程粒度的资源控制， 这种控制器称为 threaded controllers。
- v2只有单个层级树（挂载点）
- v2 的控制器文件参数和名称都变更了。

docker engine 20.10版本开始支持v2.

#### mount 挂载点区别

cgroup v2 和 v1 挂载目录是相互兼容的，v2是独立在`/sys/fs/cgroup/unified`目录下。

v1 mount目录：
```bash
mount | grep cgroup
tmpfs on /sys/fs/cgroup type tmpfs (ro,nosuid,nodev,noexec,mode=755)
cgroup on /sys/fs/cgroup/systemd type cgroup (rw,nosuid,nodev,noexec,relatime,xattr,release_agent=/usr/lib/systemd/systemd-cgroups-agent,name=systemd)
cgroup on /sys/fs/cgroup/perf_event type cgroup (rw,nosuid,nodev,noexec,relatime,perf_event)
cgroup on /sys/fs/cgroup/cpu,cpuacct type cgroup (rw,nosuid,nodev,noexec,relatime,cpuacct,cpu)
cgroup on /sys/fs/cgroup/blkio type cgroup (rw,nosuid,nodev,noexec,relatime,blkio)
cgroup on /sys/fs/cgroup/hugetlb type cgroup (rw,nosuid,nodev,noexec,relatime,hugetlb)
cgroup on /sys/fs/cgroup/freezer type cgroup (rw,nosuid,nodev,noexec,relatime,freezer)
cgroup on /sys/fs/cgroup/memory type cgroup (rw,nosuid,nodev,noexec,relatime,memory)
cgroup on /sys/fs/cgroup/pids type cgroup (rw,nosuid,nodev,noexec,relatime,pids)
cgroup on /sys/fs/cgroup/net_cls,net_prio type cgroup (rw,nosuid,nodev,noexec,relatime,net_prio,net_cls)
cgroup on /sys/fs/cgroup/devices type cgroup (rw,nosuid,nodev,noexec,relatime,devices)
cgroup on /sys/fs/cgroup/cpuset type cgroup (rw,nosuid,nodev,noexec,relatime,cpuset)
```

v2 mount 目录：
```bash
$ mount | grep cgroup2
cgroup2 on /sys/fs/cgroup/ type cgroup2 (rw,nosuid,nodev,noexec,relatime,nsdelegate)
```

- v2 是在 `/sys/fs/cgroup` 下创建了一个 `unified`的目录，v2 是单一层级树，因此只有一个挂载点。
- v1 根据控制器类型（cpuset/cpu,cpuacct/hugetlb/...），挂载到不同位置。

v1 cgroup controller:
- cpu 子系统，主要限制进程的 cpu 使用率。
- cpuacct 子系统，可以统计 cgroups 中的进程的 cpu 使用报告。
- cpuset 子系统，可以为 cgroups 中的进程分配单独的 cpu 节点或者内存节点。
- memory 子系统，可以限制进程的 memory 使用量。
- blkio 子系统，可以限制进程的块设备 io。
- devices 子系统，可以控制进程能够访问某些设备。
- net_cls 子系统，可以标记 cgroups 中进程的网络数据包，然后可以使用 tc 模块（traffic control）对数据包进行控制。
- freezer 子系统，可以挂起或者恢复 cgroups 中的进程。
- perf_event 子系统，许对cgroup中分组的一组进程进行性能监视。
- hugetlb 子系统，可以限制内存大页。

### cpu 限制原理

cgroup 提供了三种 cpu 资源限制，cpuset, cpuquota, cpushare。
- cpuset 以分配核心的方式对 cpu 资源进行隔离，比如如果我们要限制进程绑定特定的cpu核，可以设置cpuset为对应的核心数。cpuset隔离之后应用之间的相互影响最低
- cpuquota 提供一种比 cpuset 更细粒度的 cpu 隔离手段，按时间分片进行隔离。
-  cpushare, 提供的是一种按照权重比例分配 cpu 时间资源的手段，当cpu空闲的时候，某一个要占用cpu的cgroup可以完全占用剩余cpu时间，充分利用资源。

对于用户态数据隔离来说，cpuset > cpuquota > cpushare，对于内核资源竞争, cpuset > cpushare > cpuquota。

kubelet 参数`--cpu-manager-policy`默认模式none是用 cpuquota 和 cpushare 实现隔离, `static`可以针对整数 Guaranteed pod 中开启 cpuset 来实现绑定特定的CPU核。

kubelet 的 pod 的 limit 的值写入 cpuquota， requests 的值写入 cpushare。

## linux namespace

namespace 类型：
- mount, namespace有自己的挂载信息，即拥有独立的目录层次
- PID, namespace有自己的进程号，使得namespace中的进程PID单独编号
- network, namespace有自己的网络资源，包括网络协议栈、网络设备、路由表、防火墙、端口等
- IPC(Interprocess Communication)，namespace有自己的共享内存、信号量等
- UTS, namepsace有自己的主机信息,如主机名(hostname)、NIS domain name
- User, namespace有自己的用户权限管理机制(比如独立的UID/GID)
- cgroup, namespace有自己单独的cgroup(linux 4.6发布)
- time, namespace有自己的启动时间点信息和单调时间(linux 5.6发布)


### 命名空间常用查看工具

`nsenter`，一个切换namespace的命令

比如要切换到docker容器的命名空间可以，下面是一个查询容器pid，然后切换到其网络命名空间的例子：
```bash
$ docker inspect 6c360c35ccb4 --format="{{.State.Pid}}"
25274
$ nsenter -t 25274 -n iptables -S -t nat
```

`nsenter`跟`-t <PID>`可以指定切换到对应的PID的命名空间，`-n`表示网络命名空间,`-m`表示挂载信息命名空间,`-u`表示主机信息。

例子:
```bash
$ nsenter -t 25274 -m top
Mem: 63850612K used, 1955484K free, 3296692K shrd, 2046232K buff, 37521136K cached
CPU:  10% usr   1% sys   0% nic  88% idle   0% io   0% irq   0% sirq
Load average: 1.75 2.16 2.54 3/5123 10349
  PID  PPID USER     STAT   VSZ %VSZ CPU %CPU COMMAND
    1     0 root     S     711m   1%  11   0% /app/supplement --config=/app/config/config.json

$ nsenter -t 25274 -m ps 
PID   USER     TIME  COMMAND
    1 root      5h59 /app/supplement --config=/app/config/config.json

$ nsenter -t 25274 -m ls  
app    bin    dev    etc    home   lib    media  mnt    opt    proc   root   run    sbin   srv    sys    tmp    usr    var

$ nsenter -t 25274 -u hostname
supplement-55bb4f8755-9jk78
```

`lsns`，一个查看当前已创建的namespace


## 容器镜像(OCI)

OCI 容器规范包括了 image format spec 和 runtime spec。

### image spec

docker 默认是没有开启`manifest`功能，可以通过`export DOCKER_CLI_EXPERIMENTAL=enabled`为测试实验性功能提供了一个临时环境。

```bash
$ docker manifest inspect coredns/coredns:1.1.2
{
   "schemaVersion": 2,
   "mediaType": "application/vnd.docker.distribution.manifest.list.v2+json",
   "manifests": [
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 950,
         "digest": "sha256:d4a3e119061e474af69e4acf3a20a7260fd5b5075906feb07c06c7cdf378429f",
         "platform": {
            "architecture": "amd64",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 946,
         "digest": "sha256:f1a47b6a29b78729005dba14bdbde92c0439e310f2b50f95d190fb0cc4d7de43",
         "platform": {
            "architecture": "arm",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 946,
         "digest": "sha256:8c365ba9ca76753a301d25bd72bb43954972348553c6cb37df184ad3328a1e4d",
         "platform": {
            "architecture": "arm64",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 946,
         "digest": "sha256:bc531416936f219d51b922e7b3f12154667ae85613ab14fb3a195dd0873d03e7",
         "platform": {
            "architecture": "ppc64le",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 946,
         "digest": "sha256:0c57b017e1410ca8083560f6903a63a0e76782fde92670d119bcd93b7fe043dd",
         "platform": {
            "architecture": "s390x",
            "os": "linux"
         }
      }
   ]
}
```
可以看出 manifest 列表中包含了不同系统架构所对应的镜像 manifests，比如有`amd64`,`arm`,`ppc64le` 等。

Image Manifest  组成一个容器 image 的组件描述文档。
如果我们想看详细信息可以加上`verbose`，查看manifest的详细信息：
```bash
$ docker manifest inspect coredns/coredns:1.1.2 --verbose
[
	{
		"Ref": "docker.io/coredns/coredns:1.1.2@sha256:d4a3e119061e474af69e4acf3a20a7260fd5b5075906feb07c06c7cdf378429f",
		"Descriptor": {
			"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
			"digest": "sha256:d4a3e119061e474af69e4acf3a20a7260fd5b5075906feb07c06c7cdf378429f",
			"size": 950,
			"platform": {
				"architecture": "amd64",
				"os": "linux"
			}
		},
		"SchemaV2Manifest": {
			"schemaVersion": 2,
			"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
			"config": {
				"mediaType": "application/vnd.docker.container.image.v1+json",
				"size": 2331,
				"digest": "sha256:8558f8c47fd76c05e4e9182e06a39a7d71453fd2c43a0c4d1a1817ba7c19e6d7"
			},
			"layers": [
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 1990402,
					"digest": "sha256:88286f41530e93dffd4b964e1db22ce4939fffa4a4c665dab8591fbab03d4926"
				},
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 3959848,
					"digest": "sha256:f7a3e79b147dd0a1392288e06fee0b7a23a8f0fb6ec3e2db9a4944d26b8cbe7e"
				},
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 9759024,
					"digest": "sha256:fdb8a64540f06188a7c0c40446c0a1c0948ba890bb2cabf24c1e51056a53f8e7"
				}
			]
		}
	},
	{
		"Ref": "docker.io/coredns/coredns:1.1.2@sha256:f1a47b6a29b78729005dba14bdbde92c0439e310f2b50f95d190fb0cc4d7de43",
		"Descriptor": {
			"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
			"digest": "sha256:f1a47b6a29b78729005dba14bdbde92c0439e310f2b50f95d190fb0cc4d7de43",
			"size": 946,
			"platform": {
				"architecture": "arm",
				"os": "linux"
			}
		},
		"SchemaV2Manifest": {
			"schemaVersion": 2,
			"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
			"config": {
				"mediaType": "application/vnd.docker.container.image.v1+json",
				"size": 2334,
				"digest": "sha256:0bd676e2d04524314467fdd6d76b6144c00f10e6e8d7284bcdb4e8a1a965f182"
			},
			"layers": [
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 1965988,
					"digest": "sha256:0864efeeb5cb8dca4eb53e5d6fd38486daee80fa326fe36d1ad254f8fa6bb310"
				},
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 170,
					"digest": "sha256:3cda69762aee1588fa82aeabf1af6d6ad24f737cce1451fab2e0199849b1e12e"
				},
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 9349981,
					"digest": "sha256:b7e1bed31a1d76ad3c1c29dac00ef297aedf024e9d678252d04044e97c7b6f4e"
				}
			]
		}
	},
	{
		"Ref": "docker.io/coredns/coredns:1.1.2@sha256:8c365ba9ca76753a301d25bd72bb43954972348553c6cb37df184ad3328a1e4d",
		"Descriptor": {
			"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
			"digest": "sha256:8c365ba9ca76753a301d25bd72bb43954972348553c6cb37df184ad3328a1e4d",
			"size": 946,
			"platform": {
				"architecture": "arm64",
				"os": "linux"
			}
		},
		"SchemaV2Manifest": {
			"schemaVersion": 2,
			"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
			"config": {
				"mediaType": "application/vnd.docker.container.image.v1+json",
				"size": 2331,
				"digest": "sha256:c3181d10e56324996c9c8b8ee461320836bc24dbf8bf3695dc35a4371f2a1ce2"
			},
			"layers": [
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 1914748,
					"digest": "sha256:bb473f0ebc12fde1bd45c1bd3c46f2d3aab367b1b7739464771455b9972f7894"
				},
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 176,
					"digest": "sha256:75ff6b7ff3a208b8399e701e7ea1b7edbdc654c8c60d33c6f09a7803e2dda776"
				},
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 9058951,
					"digest": "sha256:6b5a7321f1b5071a9a225fe10075042e57ac20e65b1f4d091ff304abcd292251"
				}
			]
		}
	},
	{
		"Ref": "docker.io/coredns/coredns:1.1.2@sha256:bc531416936f219d51b922e7b3f12154667ae85613ab14fb3a195dd0873d03e7",
		"Descriptor": {
			"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
			"digest": "sha256:bc531416936f219d51b922e7b3f12154667ae85613ab14fb3a195dd0873d03e7",
			"size": 946,
			"platform": {
				"architecture": "ppc64le",
				"os": "linux"
			}
		},
		"SchemaV2Manifest": {
			"schemaVersion": 2,
			"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
			"config": {
				"mediaType": "application/vnd.docker.container.image.v1+json",
				"size": 2335,
				"digest": "sha256:030f75c691c2701fe9071db5c695798d7d11aeb3f32993dcdf7aaa59f9ce9c51"
			},
			"layers": [
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 2008578,
					"digest": "sha256:1e52418956f7d2a8ea35e8e6e3318fd08e005b27457d77868c225e7433bbfa02"
				},
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 176,
					"digest": "sha256:acf472f4e5bb7956ac20bb343b304e1d3de1f79160c0d158cccbe25980022d50"
				},
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 8860774,
					"digest": "sha256:4e2baa5f1f2a49446bfde7bc89f87f35142c4b53ddd4c3e38f96b349c5992003"
				}
			]
		}
	},
	{
		"Ref": "docker.io/coredns/coredns:1.1.2@sha256:0c57b017e1410ca8083560f6903a63a0e76782fde92670d119bcd93b7fe043dd",
		"Descriptor": {
			"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
			"digest": "sha256:0c57b017e1410ca8083560f6903a63a0e76782fde92670d119bcd93b7fe043dd",
			"size": 946,
			"platform": {
				"architecture": "s390x",
				"os": "linux"
			}
		},
		"SchemaV2Manifest": {
			"schemaVersion": 2,
			"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
			"config": {
				"mediaType": "application/vnd.docker.container.image.v1+json",
				"size": 2335,
				"digest": "sha256:ca89c57b8e46e6a1e5ffd6933a246fb1f4fc415d2de4014322b36d8c350a1036"
			},
			"layers": [
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 2065460,
					"digest": "sha256:d45fd9d3c4f188ab1f3a4bf6a9f5202b3f1577dbb998f5f28e82d192e0c1f0e7"
				},
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 176,
					"digest": "sha256:0e5978b6b34b3e943e0fd25dfb50991c0bad82a986cfdaa91c4de756431ba679"
				},
				{
					"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
					"size": 9640536,
					"digest": "sha256:a65042f6f54bb79d348f647b9b69cd34619a86b18eeffa1fc12f6f89ec4ef752"
				}
			]
		}
	}
]
```

`SchemaV2Manifest`就是符合 OCI 的 image spec。字段说明：
- schemaVersion，image manifest的schema版本
- mediaType，表示文件类型。

| mediaType字段 | 说明 |
|----:|----:|
|application/vnd.oci.descriptor.v1+json |	Content Descriptor 内容描述文件 |
|application/vnd.oci.layout.header.v1+json |	OCI Layout 布局描述文件 |
|application/vnd.oci.image.index.v1+json |	Image Index 高层次的镜像元信息文件 |
|application/vnd.oci.image.manifest.v1+json |	Image Manifest 镜像元信息文件 |
|application/vnd.oci.image.config.v1+json |	Image Config 镜像配置文件 |
|application/vnd.oci.image.layer.v1.tar	Image | Layer 镜像层文件 |
|application/vnd.oci.image.layer.v1.tar+gzip |	Image Layer 镜像层文件gzip压缩 |
|application/vnd.oci.image.layer.nondistributable.v1.tar |	Image Layer 非内容寻址管理 |
|application/vnd.oci.image.layer.nondistributable.v1.tar+gzip |	Image Layer, gzip压缩 非内容寻址管理 |


Image Configuration， 包含如应用参数、环境等信息。 
```bash
$ docker image inspect coredns/coredns:1.1.2             
[
    {
        "Id": "sha256:8558f8c47fd76c05e4e9182e06a39a7d71453fd2c43a0c4d1a1817ba7c19e6d7",
        "RepoTags": [
            "coredns/coredns:1.1.2"
        ],
        "RepoDigests": [
            "coredns/coredns@sha256:dd2cd70f60ff7895b6a96002a54cbd3d00e88c19ba804aab56c8a5b645cf1e08"
        ],
        "Parent": "",
        "Comment": "",
        "Created": "2018-04-23T13:09:08.592750353Z",
        "Loaded": "2021-04-22T20:46:00.758134827+08:00",
        "Container": "6770f3ff6176c4ec1a7220e2bcbef5617c4573a33c9ff67b99ec07f075264ce6",
        "ContainerConfig": {
            "Hostname": "e1ede117fb1e",
            "Domainname": "",
            "User": "",
            "AttachStdin": false,
            "AttachStdout": false,
            "AttachStderr": false,
            "ExposedPorts": {
                "53/tcp": {},
                "53/udp": {}
            },
            "Tty": false,
            "OpenStdin": false,
            "StdinOnce": false,
            "Env": [
                "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
            ],
            "Cmd": [
                "/bin/sh",
                "-c",
                "#(nop) ",
                "ENTRYPOINT [\"/coredns\"]"
            ],
            "ArgsEscaped": true,
            "Image": "sha256:da5edfc4c1e75bb644fe39bb072e8c99a957365cef3476f7504f1c5fccba111e",
            "Volumes": null,
            "WorkingDir": "",
            "Entrypoint": [
                "/coredns"
            ],
            "OnBuild": null,
            "Labels": {},
            "Annotations": null
        },
        "DockerVersion": "18.03.0-ce",
        "Author": "",
        "Config": {
            "Hostname": "e1ede117fb1e",
            "Domainname": "",
            "User": "",
            "AttachStdin": false,
            "AttachStdout": false,
            "AttachStderr": false,
            "ExposedPorts": {
                "53/tcp": {},
                "53/udp": {}
            },
            "Tty": false,
            "OpenStdin": false,
            "StdinOnce": false,
            "Env": [
                "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
            ],
            "Cmd": null,
            "ArgsEscaped": true,
            "Image": "sha256:da5edfc4c1e75bb644fe39bb072e8c99a957365cef3476f7504f1c5fccba111e",
            "Volumes": null,
            "WorkingDir": "",
            "Entrypoint": [
                "/coredns"
            ],
            "OnBuild": null,
            "Labels": {},
            "Annotations": null
        },
        "Architecture": "amd64",
        "Os": "linux",
        "Size": 46384642,
        "VirtualSize": 46384642,
        "GraphDriver": {
            "Data": {
                "LowerDir": "/var/lib/docker/overlay2/80ea904e4a887745926fab3c86e81196816ef13db4e635321fe812447513a060/diff:/var/lib/docker/overlay2/325bd24eb800ae5af7cf9c8853039d2457f5e778a39dfb6ec376ccc4e54d6f59/diff",
                "MergedDir": "/var/lib/docker/overlay2/b9d45d0687506a2b9bec11b1759174d5e12099e41737c85d2fc454c138b04dcf/merged",
                "UpperDir": "/var/lib/docker/overlay2/b9d45d0687506a2b9bec11b1759174d5e12099e41737c85d2fc454c138b04dcf/diff",
                "WorkDir": "/var/lib/docker/overlay2/b9d45d0687506a2b9bec11b1759174d5e12099e41737c85d2fc454c138b04dcf/work"
            },
            "Name": "overlay2"
        },
        "RootFS": {
            "Type": "layers",
            "Layers": [
                "sha256:5bef08742407efd622d243692b79ba0055383bbce12900324f75e56f589aedb0",
                "sha256:594f5d257cbe4f627ab9f73a0ce21fb806d5773e9b0ce1c22c04403bfc551e77",
                "sha256:3801292bb760b26ad59b779f6ffb8bbc296677eee1343b7a62509e4102d62055"
            ]
        },
        "Metadata": {
            "LastTagTime": "0001-01-01T00:00:00Z"
        }
    }
]
```

## 参考文献
- [Control Group v2](https://arthurchiao.art/blog/cgroupv2-zh/#211-%E6%8E%A7%E5%88%B6%E5%99%A8%E4%B8%8E-v1v2-%E7%BB%91%E5%AE%9A%E5%85%B3%E7%B3%BB)
- [从CPU资源隔离说起](https://zorrozou.github.io/docs/books/cgroup_linux_cpu_control_group.html)
- [Docker CPU 资源限制](https://blog.opskumu.com/docker-cpu-limit.html)