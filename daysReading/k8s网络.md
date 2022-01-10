

# 1、容器网络



## 1.1 namespace

> Linux的namespace（名字空间）的作用就是“隔离内核资源”。在Linux的世界里，文件系统挂载点、主机名、POSIX进程间通信消息队列、进程PID数字空间、IP地址、user ID数字空间等全局系统资源被namespace分割，装到一个个抽象的独立空间里。而隔离上述系统资源的namespace分别是Mount namespace、UTS namespace、IPC namespace、PIDnamespace、network namespace和user namespace

上述面文字摘抄自《Kubernetes网络权威指南》，namespace是linux虚拟化技术的基石。linux在创建进程的时候可以直接指定该进程的namespace。目前可指定的namespace 有上述文字描述的6种。



docker 本质上是一个指定了各种namespace的进程。默认情况下会为其分配一个属于该容器的Network Namespace，使其与其他容器、宿主机隔离，也可以在创建容器的时候，使用`–net=host`参数直接使用宿主机的网络栈（即宿主机的network namespace，普通的进程就是与宿主机使用相同的network namespace）。示例如下：

```shell
# nginx容器直接监听的就是宿主机的 80 端口。
docker run –d –net=host --name nginx-host nginx
```



关于network namespace相关的命令

```shell
# 查看network namespace
ip netns list

# 创建network namespace
ip netns add netns1

# 删除network namespace (只要仍旧有进程使用,该ns仍旧存在)
ip netns delete netns1

# 进入netns1这个network namespace查询网卡信息的命令
ip netns exec netns1 ip link list
```



## 1.2 veth pair

veth是虚拟以太网卡（Virtual Ethernet）的缩写。veth设备总是成对的，又称为veth pair。veth pair一端发送的数据会在另外一端接收。根据这一特性，veth pair常被用于跨networknamespace之间的通信，即分别将veth pair的两端放在不同的namespace里，即可使不同namespace间的进程可以通信了。

仅有veth pair设备，容器是无法访问外部网络的。为什么呢？因为从容器发出的数据包，实际上是直接进了veth pair设备的协议栈。如果容器需要访问网络，则需要使用网桥等技术将veth pair设备接收的数据包通过某种方式转发出去。



**veth pari 创建和使用的例子**

```shell
# 创建veth pair 分别是 veth0、veth1
ip link add veth0 type veth peer name veth1
```

然后用`ip link` 查看创建出来的veth, 初始状态是down

```shell
6: veth1@veth0: <BROADCAST,MULTICAST,M-DOWN> mtu 1500 qdisc noop state DOWN mode DEFAULT group default qlen 1000
    link/ether 0e:e7:43:99:a0:fe brd ff:ff:ff:ff:ff:ff
7: veth0@veth1: <BROADCAST,MULTICAST,M-DOWN> mtu 1500 qdisc noop state DOWN mode DEFAULT group default qlen 1000
    link/ether 66:d0:66:86:02:61 brd ff:ff:ff:ff:ff:ff
```

用ip link命令将这两块网卡的状态设置为UP

```shell
ip link set dev veth0 up
ip link set dev veth1 up
```

设置IP

```shell
ifconfig veth0 10.20.30.40/24
ifconfig veth0 10.20.30.41/24
```

将veth1放到某个network namespace中, 加入netns1后，状态会变回down,需要重新设置ip

```
ip link set veth1 netns netns1
```



## 1.3 bridge

