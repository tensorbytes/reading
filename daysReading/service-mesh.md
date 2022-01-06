## 服务网格

### istio

有些时候我们并不需要对服务调用数据库的端口进行客观性和追踪，我们可以通过设置 annotations 设置，比如我们希望出口Port 3306 不走istio，我们可以设置：
```yaml
annotations:
    traffic.sidecar.istio.io/excludeOutboundPorts: 3306
```
