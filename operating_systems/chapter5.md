#### 进程API

##### UNIX系统进程创建
#####  创建新进程

一对系统调用fork() 和 exec()
调用wait() 等待子进程创建完成

### 为什么要这样设计API？
分离fork()及exec()的做法在构建UNIX shell的时候非常有用，这给了shell在fork之后exec之前运行代码的机会