## 服务网格

### istio

1. 针对指定pod的sidecar设置指定端口不走sidecar逻辑

有些时候我们并不需要对服务调用数据库的端口进行客观性和追踪，我们可以通过设置 annotations 设置，比如我们希望出口Port 3306 不走istio，我们可以设置：
```yaml
annotations:
    traffic.sidecar.istio.io/excludeOutboundPorts: 3306
```
值得注意的是`traffic.sidecar.istio.io/excludeOutboundPorts`是和`traffic.sidecar.istio.io/excludeOutboundIPRanges`、`traffic.sidecar.istio.io/includeOutboundIPRanges`共用的，默认`traffic.sidecar.istio.io/includeOutboundIPRanges`是`*`，表示所有IP都走sidecar，因此我们设置`traffic.sidecar.istio.io/excludeOutboundPorts`表示这几个端口不走，如果这时设置`traffic.sidecar.istio.io/includeOutboundPorts`就不生效，因为`traffic.sidecar.istio.io/includeOutboundIPRanges`为`*`已经所有IP端口都走sidecar了，这里确实有些歧义。

2. istio的负载均衡策略在默认情况下和识别的端口名称有关系

我们都知道在k8s集群中，一个带istio sidecar的pod要访问另一个service，需要 service 端口名称设置为"tcp-xxx","http-xxx"或"grpc-xxx"，因为 istio 如果没开启端口协议检查（生产环境不建议开启），会根据 k8s 的 service 的端口名称去识别对应端口是什么协议。

也正因为 istio 的这个特性导致我们在给 service 的 port 设置 name 的时候要注意，如果为了偷懒设置为"tcp-xxx"会导致 istio将下游端口识别为tcp，从而使用 四层负载均衡策略来访问，如果业务时 grpc，那就可能导致下游负载不均衡。