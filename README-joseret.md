CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./dist/k8s-endpoints-sync-controller src/main/main.go

docker build -t  gcr.io/jr-network-infra-1-4978/joseret/k8s-endpoints-sync-controller:latest  .
docker push  gcr.io/jr-network-infra-1-4978/joseret/k8s-endpoints-sync-controller:latest