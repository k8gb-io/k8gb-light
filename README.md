# annotation-controller
Demonstrates reconciliation on the top of ingress only 
- No CRD
- No GSLB
- No test
- Not completed
- keeping state within the annotations

`k8gb.io/status` is written back by controller at the end of reconciliation. 
`k8gb.io/strategy` is mandatory
 

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