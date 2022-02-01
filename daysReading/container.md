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

## linux namespace



## 容器镜像

## 参考文献
- [Control Group v2](https://arthurchiao.art/blog/cgroupv2-zh/#211-%E6%8E%A7%E5%88%B6%E5%99%A8%E4%B8%8E-v1v2-%E7%BB%91%E5%AE%9A%E5%85%B3%E7%B3%BB)