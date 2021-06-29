![Go](https://github.com/OpenZeppelin/fortify-node/workflows/Go/badge.svg)
![Deploy](https://github.com/OpenZeppelin/fortify-node/workflows/Deploy/badge.svg)

## fortify-node

#### Dependencies

1. [Install Docker](https://docs.docker.com/get-docker/)
2. [Install golang 1.16](https://golang.org/doc/install)
3. [Install Protobuf Compiler](https://grpc.io/docs/protoc-installation/)
4. Start Docker Service

#### Libraries to install
```shell
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc 
```
```shell 
$ go install github.com/golang/mock/mockgen@v1.5.0
```

#### Build Node

```shell
$ make build
```

#### Run Node

```shell
$ ./fortify
```

#### Stop Node

```
CTRL-C
```

#### View Logs

```shell
$ docker logs -f CONTAINER_ID
```

#### Get Alerts

```shell
$ curl -s http://localhost:8778/alerts
```


#### Configuration

See config.yml for configuration (Docs TBD)