# k8s network

## The path of network traffic in Kubernetes

k8s网络模型要求：
- pod 和 pod 之间自由通信，不需要 NAT 做地址转换
- 每个 pod 都有自己的 ip 地址 (ip-per pod)
- 集群同一节点上的 pod 和 应用程序 可以自由通信

只要满足这几点，k8s的网络模型可以采用任意的技术。

## k8s的服务发现

我们都知道 k8s 的服务发现是用 dns 来做的，但为什么可以用 dns 来做服务发现呢？

DNS 来做服务发现的好处是具有极强的扩展性，通过 forward 用层次结构的方式可以搭建出服务全球网络域名发现，所以在扩展性上是毋庸置疑的。但是用 DNS 做服务发现有一个致命的缺点，那就是其缓存的功能。缓存对于域名发现来说是非常棒的功能，因为对网站的域名变更通常是非常低频的，因此通过设置缓存可以极大的降低 DNS 服务器的负载。
但应用间的服务发现对实时性通常会有极强需求，当一个服务 IP 地址变更的时候，服务通常希望他的调用方能立刻感知到，如果因为 DNS 缓存导致调用方访问到未更新的错误IP，那么将是灾难性的。那么 k8s 是怎么解决这个问题的呢？ k8s 通过 service 引入一个全局的 clusterIP 来解决这个问题。clusterIP 是一个VIP（虚拟IP），与 service 共生，同时是全局唯一的。具体应用的 IP 地址，会通过一系列的 iptables 规则来转发到真实的目标服务IP(endpoints)，而 cluseterIP 和目标服务IP 的关系则是由每个边缘端节点上的 kube-proxy 来实时更新。这样，k8s 就通过 dns 配合虚拟IP创建了一个基于 DNS 服务发现的基础架构。你只需要创建一个 service 对象资源，然后直接用域名去访问它就可以实时地被转发到目标资源，从而通过 DNS 这种优雅的方式实现了服务发现。

## 现有的主流 cni 模型

现有的主流 cni 模型可以分为两类：
- 基于 overlay 模式
- 基于 路由 模式

### 基于 overlay 模型

基于 overlay 的有：
- flannel 的 UDP 模式
- flannel 的 vxlan 模式
- calica 的 IPIP 模式



### 基于路由模型

基于路由的有：
- flannel 的 hostgw 模式
- calico 的 BGP 路由模式