# API Server

## admission webhook


api-server 通过读取 `mutatingwebhookconfiguration` 和 `validatingwebhookconfiguration` 的 CR 文件的目标地址，然后回调用户自定义的服务。

```
                                            ┌──────────────────────────────────┐
             ┌─────────────────┐            │                                  │
    apply    │                 │    read    │  validatingwebhookconfiguration  │
────────────►│    api-server   │◄───────────┤                                  │
             │                 │            │  mutatingwebhookconfiguration    │
             └────────┬────────┘            │                                  │
                      │                     └──────────────────────────────────┘
                      │
                      │  回调
                      │
                      │
             ┌────────▼────────┐
             │                 │
             │  webhookservice │
             │                 │
             └─────────────────┘
```

api-server 发起的请求是一串json数据格式，header需要设置`content-type`为`application/json`, 我们看看请求的 body :
```json
curl -X POST \
  http://webhook-url \
  -H 'content-type: application/json' \
  -d '{
  "apiVersion": "admission.k8s.io/v1",
  "kind": "AdmissionReview",
  "request": {
    ...
  }
}'
```
返回的结果：
```json
{
    "kind": "AdmissionReview",
    "apiVersion": "admission.k8s.io/v1",
    "response": {
        "uid": "b955fb34-0135-4e78-908e-eeb2f874933f",
        "allowed": true,
        "status": {
            "metadata": {},
            "code": 200
        },
        "patch": "W3sib3AiOiJyZXBsYWNlIiwicGF0aCI6Ii9zcGVjL3JlcGxpY2FzIiwidmFsdWUiOjJ9XQ==",
        "patchType": "JSONPatch"
    }
}
```
这里的 patch 是用base64编码的一个json，我们解码看看，是一个 json patch：
```bash
$ echo "W3sib3AiOiJyZXBsYWNlIiwicGF0aCI6Ii9zcGVjL3JlcGxpY2FzIiwidmFsdWUiOjJ9XQ==" | base64 -d
[
    {
        "op": "replace",
        "path": "/spec/replicas",
        "value": 2
    }
]
```