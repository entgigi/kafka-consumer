### Build
```
 docker build -t gigiozzz/consumer-kafka:0.0.1 .

 docker push gigiozzz/consumer-kafka:0.0.1
```

### Usage on kube
```
 kubectl -n kafka run kafka-consumer -ti --image=gigiozzz/consumer-kafka:0.0.1 --rm=true   --env="KAFKA=my-cluster-kafka-bootstrap:9092" --env="POD_NAME=test1"
```