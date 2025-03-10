The `Content-Type` header must be set correctly, when using the polaris-scheduler or polaris-cluster-agent APIs, otherwise Gin will silently fail to parse the request body.

```sh
curl -XPOST localhost:8080/pods -H "Content-Type: application/json" -d '{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "namespace": "test",
        "name": "myapp-01",
        "labels": {
            "name": "myapp-01"
        }
    },
    "spec": {
        "containers": [
            {
                "name": "myapp",
                "image": "gcr.io/google-containers/pause:3.2",
                "resources": {
                    "limits": {
                        "polaris-slo-cloud.github.io/fake-milli-cpu": "4000",
                        "polaris-slo-cloud.github.io/fake-memory": "8Gi"
                    }
                }
            }
        ],
        "tolerations": [
            {
                "key": "fake-kubelet/provider",
                "operator": "Exists",
                "effect": "NoSchedule"
            }
        ]
    }
}'
```
