## 什么是shuffle

spark提交job会根据计算逻辑切分成多个ShuffleMapStage以及1个RsultStage。一个Stage最后会再进一步转换成一个Task集合。ShuffleMapStage转换成ShuffleMapTask集合，ResultStage会转换成ResultTask集合。

假设一个Spark应用有2两个ShuffleMapTask集合以及一个ResultTask集合，它们的执行顺序如下：

ShuffleMapTask集合1 => ShuffleMapTask集合2 => ResultTask集合

第一个集合ShuffleMapTask的输入来自于外部，如hdfs等，第二ShuffleMapTask集合的输入是上一个Task集合的输出。两个ShuffleMapTask集合之间的数据传输就是我们常说的shuffle，当然ShuffleMapTask集合与ResultTask集合之间的数据传输也是shuffle.

之前有个朋友问我，假设ShuffleMapTask集合2与它的的输入来源ShuffleMapTask集合1都在同一个节点上，没有发生跨节点传输数据，那这个数据传输的过程还算shuffle吗？我个人理解是task是否在不同的节点上并不是定义shuffle的必要条件，shuffle本质是两个无法合并两个stage之间的数据传输，全部在一个节点上是任务调度的关系，是shuffle中特殊的情况。

## shuffle的流程（待补充）

那么一个shuffle流程的流程是怎样的呢？为了方便描述，将前后两个ShuffleMapTask用MapReduce中的task类型类描述。前者称为MapTask，后者称为ReduceTask.  无论是Spark还是MapReduce计算框架，都是用了迭代器的方式去触发计算的。

1、MapTask逐行输出计算结果

2、计算结果写入缓存

3、缓存达到阈值，溢写磁盘形成一个小文件

4、重复1-3步骤，直到MapTask输出了全部计算结果，得到一堆小文件

5、利用归并排序算法对小文件进行合并，形成一个以分区id、分区key值排序的大文件



Spark内部把shuffl相关的功能抽象为writer以及reader，executor在执行task的时候，用writer输出数据，用reader读取数据。其底层实现被封装了起来。

这个writer以及reader是通过一个叫ShuffleManager获取的。ShuffleManager实现最初叫HashShuffleManager，后来被SortShuffleManager替代了。HashShuffleManager有两个阶段，分别是优化前，优化后两个阶段。优化前产生的文件数为 map任务数 * reduce任务数个临时文件，优化后为executor 个数 * reduce任务数。<font color="red">SortShuffleManager顾名思义，增加了排序的逻辑，在map端溢写时排序，合并成大文件时会排序,reduce端拉取数据的时候与map端使用同一个数据结构ExternalSorter，因此也会有溢写磁盘以及排序的逻辑。</font>

reader封装了从拉取数据的逻辑，reduce任务只需要考虑对数据集中每条数据如何处理即可。Spark用迭代器来表示一个数据集，reduce任务的实现就是通过遍历迭代器获取每一条记录，然后对记录进行处理。简化后的代码效果大致如下：

```
iterator = rdd.iteator(partitionId)
While iterator.hasNext():
    record = iterator.next()
    process(record)
```

process就是我们编写的业务逻辑处理函数，从不同的executor上拉取map任务输出的中间结果、然后排序合并等等的逻辑都隐藏在 迭代器的  next方法中。



reader在读取数据的时候，需要知道自己需要取的数据落在哪个executor上，因此会先通过MapOutPutTrackerWorker发送rpc请求给到 Driver中的MapOutPutTrackerMaster获取MapStatus对象。 MapStatus对象记录了Map任务的状态信息，这些信息包含了 数据存储所在的executor id，数据文件中各个partition的偏移量。









shuffle过程中，当map task输出数据满足以下条件时，会溢写到磁盘中：

（1）buffer或者map无法扩容时

（2）buffer或者map中的记录数达到阈值时，默认阈值为Integer.Max_value,即2的31次方。该值可以通过配置项修改



在sortShuffleWriter中可以看到map端的数据一定会落盘。



在spark 2.x中，map与reduce之间连接都是 all to all的，在spark 3.x中增加了 external shuffle service实现，优化了shuffle。













