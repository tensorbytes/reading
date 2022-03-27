# etcd

## raft算法

### 基础

- leader:  领导者
- Follower: 跟随者
- Candidate:  参与者

节点状态机切换

![Untitled](etcd%20ded24/Untitled.png)

1. 当节点starts up自动进入follower状态
2. follower启动之后，将开启一个选举超时的定时器，当定时器到期后，则进入candidate状态。
3. 进入candidate状态开始选举，如果在下一次选举超时到来之前，没有选出新的leader，则保持在candidate状态重新选举。
4. 当candidate状态收到超过半数的节点选票，则切换状态成为新的leader
5. 当candiate收到来自leader的消息或者更高任期的消息，表示已经有leader，则切换为follower状态
6. 当leader发现更高任期的消息，则切换为follower状态

### leader选举