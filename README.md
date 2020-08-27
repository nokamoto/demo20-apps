# demo20-apps
golang gRPC servers of demo applications

## Run
```bash
# initialize
$ kubectl apply -k deployments/migration

# run
$ skaffold dev

# clean up
$ kubectl delete -k deployments/migration
```
