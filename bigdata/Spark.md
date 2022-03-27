spark提交job会转换成0到多个ShuffleMapStage以及1个RsultStage，然后会再进一步转换成0到多个ShuffleMapTask以及1个ResultTask.

假设有2两个ShuffleMapTask以及一个ResultTask，它们的执行顺序如下：

ShuffleMapTask => ShuffleMapTask => ResultTask

第一个Task的输入来自于外部，如hdfs等。后一个Task的输入是上一个Task的输出。两个task之间的数据传输就是常说的shuffle。



那么一个shuffle流程的流程是怎样的呢？为了方便描述，将前后两个ShuffleMapTask用MapReduce中的task类型类描述。前者作为后者输入，

前者称为MapTask，后者称为ReduceTask.

1、MapTask对数据进行加工

2、MapTask加工的结果，分多个文件写入磁盘

3、对文件进行排序合并

4、





Spark程序优化的一个核心思路是减少shuffle的次数。因为shuffle伴随着数据写盘、跨节点传输数据、排序等性能消耗较大的环节。只有分布式计算中有聚合的逻辑，shuffle就无法避免。因此优化的另一个思路是尽量减少shuffle的数据量。

