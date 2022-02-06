```
1、rpc请求有哪些
2、rpc请求的顺序
3、rpc请求中的client、server在flink中的抽象是什么
4、flink的rpc框架是基于akka还是netty？
5、flink中的job会被抽象graph进行管理，共有几种graph？各个层次的graph的生成是在什么时候？
6、各个层次的graph有什么不一样？为什么要划分这么多层？
7、streamGraph是怎么转换成jobGraph的？
8、JobGraph中为什么要做OperatorChain
9、Operator要满足什么条件才能合并
10、JobVertex、JobEdge、IntermediateDataSet是什么？
11、InputGate、InputChannel、ResultPartition是什么？

12、Flink On Yarn的启动流程
13、Yarn如何实现资源隔离？
14、Flink On Yarn的3种模式区别？
15、3种模式的对应的实现类的类名规律是什么？
16、Yarn模式与StandAlone模式的区别是什么？
17、Flink的HA服务实现中，保证了哪4个服务的高可用？
18、Flink中的大对象存储服务BlobServer
19、心跳机制
20、start-cluster.sh集群启动脚本分析
21、Flink中的JobManager、TaskManager、JobMaster、ResourceManager、Dispatcher、WebMonitorEndpoint的职责
22、StandAlone集群的JobManager启动干了什么？
23、ResourceManager启动干了什么？
24、WebMonitorEndpoint启动干了什么？
25、Dispatcher启动干了什么？
26、TaskExecutor启动干了什么？
27、HA方案下，leader切换，TaskExecutor是如何感知的？
```

















