## 服务网格

### istio

有些时候我们并不需要对服务调用数据库的端口进行客观性和追踪，我们可以通过设置 annotations 设置，比如我们希望出口Port 3306 不走istio，我们可以设置：
```yaml
annotations:
    traffic.sidecar.istio.io/excludeOutboundPorts: 3306
```
值得注意的是`traffic.sidecar.istio.io/excludeOutboundPorts`是和`traffic.sidecar.istio.io/excludeOutboundIPRanges`、`traffic.sidecar.istio.io/includeOutboundIPRanges`共用的，默认`traffic.sidecar.istio.io/includeOutboundIPRanges`是`*`，表示所有IP都走sidecar，因此我们设置`traffic.sidecar.istio.io/excludeOutboundPorts`表示这几个端口不走，如果这时设置`traffic.sidecar.istio.io/includeOutboundPorts`就不生效，因为`traffic.sidecar.istio.io/includeOutboundIPRanges`为`*`已经所有IP端口都走sidecar了，这里确实有些歧义。