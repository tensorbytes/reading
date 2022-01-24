# k8s network

## The path of network traffic in Kubernetes

k8s网络模型要求：
- pod 和 pod 之间自由通信，不需要 NAT 做地址转换
- 每个 pod 都有自己的 ip 地址 (ip-per pod)
- 集群同一节点上的 pod 和 应用程序 可以自由通信

只要满足这几点，k8s的网络模型可以采用任意的技术。

