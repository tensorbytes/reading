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



## Flink应用运行模式

在不同部署模式下，Flink应用和Spark应用运行模式并不相同。



Spark支持4部署模式（资源管理器）和2种运行模式。部署方式分别是standalone、mesos、k8s、yarn；运行方式有2种，分别是client模式以及cluster模式。

部署模式和运行模式的关系是组合关系，也就有4*2=8种情况。

![image-20220220163153038](flink源码学习/image-20220220163153038.png)

Flink有3种部署模式和3种运行模式。部署方式分别是standalone、k8s、yarn， 不支持mesos是因为mesos已经很少人用了；运行模式分别是PerJob模式、Session模式、Application模式。

![image-20220220163447695](flink源码学习/image-20220220163447695.png)

官方说明地址： https://nightlies.apache.org/flink/flink-docs-release-1.14/docs/deployment/overview/#per-job-mode

根据前面提到的组合关系，理论上flink有3*3=9种情况，但是standlone模式下，并不支持perjob模式，因此是8种。



因为flink对运行模式的称呼与spark的不一样，对于spark用户在学习flink的时候，很容易搞混。以官网的flink应用架构图来说明

![image-20220220170009818](flink源码学习/image-20220220170009818.png)

在Spark用户看来，这是对应的Spark里面的standalone集群。flink client很自然对应spark-submit命令，jobmanager对应master,taskManager对应worker。但是在flink的standalone的session模式来说，jobmanager实现了spark standalone中master的资源请求处理、以及Driver主控应用的作用，taskManager实现了worker分配资源以及executor执行计算任务的作用。 而Spark无论哪种情况下，driver与master、executor与worker都是不同的JVM进程，都是独立分开的。









## RPC实现

要理解Flink的源码，需要了解flink各个进程之间是如何通讯的。flink的RPC与Spark相似，用netty仿制akka的api实现rpc的基础类，所有组件的通讯都是通过继承rpc的基础类来实现的。







## Flink Job提交流程

提交程序的主入口类是CLIFrontEnd。















