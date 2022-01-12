第8章 可靠性探究



AR、ISR、OSR、HW、LEO指的是什么？

```
AR指所有副本
ISR全称 In Sync Replicas指所有与leader副本保持一致的副本集合为ISR
OSR全称 Out of Sync Replicas,指落后leader副本的副本集合
```

HW与LEO的关系？

```
当所有的副本都处于与leader一致的情况下，LEO=HW。
消费者只能消费低于HW的消息，当有副本（ISR中的副本）滞后时，HW=滞后最多的副本的LEO
```

一个副本什么情况下会被判定为失效副本？

```
（1）一段时间内没有向leader发起同步请求，例如频繁Full GC
（2）follwer副本进程同步过慢，在一段时间内都无法追赶上leader副本 (replica.lag.time.max.ms 参数)
（3）新增副本，在赶上leader之前都是失效状态
（4）副本下线（broker宕机）后又上线（恢复）
```



为什么不用落后的消息记录数来判定为失效副本呢？

```
0.9版本前有一个参数replica.lag.max.messages，默认值为4000。kafka也用该参数判定副本是否失效。当一个follower副本滞后leader副本消息数超过replica.lag.max.messages的大小时，则会被判定为同步失效（该参数与replica.lag.time.max.ms参数判定出的失效副本取并集组成OSR集合）。

使用滞留消息记录数来判定同步失效的问题在于很难给出一个合适的值，若设置太大，则这个参数本身就没有太多意义（大多数时间内处于失效状态），值过小则会导致ISR的频繁伸缩，可能会导致性能问题（具体案例可参考 https://blog.csdn.net/prestigeding/article/details/122094581  ）

此处需要补充ISR频繁伸缩为啥会引发性能下降



```







