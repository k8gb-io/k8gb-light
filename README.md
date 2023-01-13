# k8gb-light
POC k8gb reduced by CRD.  Besides CRD it also reduces the number of reconciliations and provides new providers as Gateway.

see [🚨Changes](https://github.com/k8gb-io/k8gb-light/issues/1) which contain more detailed information about the changes.


## running locally
### Environments
Locally create three test clusters by displaying `make-deploy-full-local-setup`.  All clusters will have the current version of k8gb installed from your local Git branch.

### Test application
To install the application, use `make deploy-demo`

### Terratests
Run `make terratest`.

If you need to reset the environment, use the commands `make clean-namespaces`, `make deploy-full-local-clusters`, `make redeploy-clusters`



## Ingress

If k8gb is successfully installed, you only need to add the annotation to ingress and load-balancing will be enabled.

`k8gb.io/strategy` is mandatory
`k8gb.io/status` is out only information written by controller back to ingress. The value represents the state of the host application on each cluster.
`k8gb.io/primary-geotag` is used if `k8gb.io/strategy` is `failover`
`k8gb.io/weights` is currently JSON containing key-values for the weights of the individual regions. Only used if `k8gb.io/strategy` is `roundRobin`

Other annotations are `k8gb.io/splitbrain-threshold-seconds` and `k8gb.io/dns-ttl-seconds`


```yaml
kind: Ingress
metadata:
  annotations:
    k8gb.io/dns-ttl-seconds: "351"
    k8gb.io/status: '{"serviceHealth":{"demo.cloud.example.com":"Healthy"},"healthyRecords":{"demo.cloud.example.com":["172.18.0.5","172.18.0.6","172.18.0.3","172.18.0.4"]},"geoTag":"us","hosts":"demo.cloud.example.com"}'
    k8gb.io/strategy: roundRobin
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"networking.k8s.io/v1","kind":"Ingress","metadata":{"annotations":{"k8gb.io/strategy":"roundRobin","x.y.io/ep":"[{\"addresses\":[\"1.2.3.4\"],\"port\":80}]","xxx":"xxx"},"name":"ing","namespace":"demo"},"spec":{"ingressClassName":"nginx","rules":[{"host":"demo.cloud.example.com","http":{"paths":[{"backend":{"service":{"name":"frontend-podinfo","port":{"name":"http"}}},"path":"/","pathType":"Prefix"}]}}]}}
```

```yaml
  endpoints:
  - dnsName: localtargets-demo.cloud.example.com
    recordTTL: 351
    recordType: A
    targets:
    - 172.18.0.5
    - 172.18.0.6
  - dnsName: demo.cloud.example.com
    labels:
      strategy: roundRobin
    recordTTL: 351
    recordType: A
    targets:
    - 172.18.0.5
    - 172.18.0.6
    - 172.18.0.3
    - 172.18.0.4
```