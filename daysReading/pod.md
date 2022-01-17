## Pod 字段小tips

### 加上当地时区或指定时区

k8s的 pod 默认是没有设置时区信息的，如果需要时区需要挂载`/etc/localtime`文件，在本地有`/etc/localtime`之后改时区可以通过环境变量`TZ`修改：

```yaml
...
    spec:
      containers:
        - image: xxxx
          env:
            - name: TZ
                value: Asia/Shanghai
          volumeMounts:
            name: host-time
            readOnly: true
      volumes:
      - hostPath:
          path: /etc/localtime
        name: host-time
```

这样设置，直接加载宿主机的 localtime，可以保持跟 宿主机的时区一致，同时也通过`TZ`环境变量提供了修改时区的能力。
`TZ`包含哪些时区可以通过`ls /usr/share/zoneinfo/`查看。