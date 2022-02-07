#### api-server 

提供kubernetes资源的restful api

##### GVK vs GVR

##### GVK

GROUP VERION KIND

##### GVR
GVR 常用于组合成 RESTful API 请求路径

这种 GVK 与 GVR 的映射叫做 RESTMapper

|         | Syntax      | Description |
|---------| ----------- | ----------- |
|实体类型   | Resource      | Kind       |
|实现方式	| http   | Controller        |
|资源定位	| URL PATH   | GroupVersionKind        |
