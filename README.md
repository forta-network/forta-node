![Build](https://github.com/forta-network/forta-node/actions/workflows/release-codedeploy-dev.yml/badge.svg)

# forta-node

Forta node CLI is a Docker container supervisor that runs and manages multiple services and detection bots (agents) to scan a blockchain network and produce alerts.

# Running a Node

For information about running a node, see the [Scan Node Quickstart Documentation](https://docs.forta.network/en/latest/scanner-quickstart/)

# Scan Node Development

## Dependencies

1. [Install Docker](https://docs.docker.com/get-docker/) and start Docker service
2. [Install Go](https://golang.org/doc/install)

## Dependencies for local development

### Tools

Install [Protobuf Compiler](https://grpc.io/docs/protoc-installation/).

### Go libraries

```shell
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc 
```
```shell 
$ go install github.com/golang/mock/mockgen@v1.5.0
```

## Build and install

### Full build & install using local version of Go

```shell
$ make install
```

For a faster iteration in local development, it is sufficient to build the common service container only if it has changed. The CLI requires `forta-network/forta-node:latest` containers to be available by default and uses the local ones if other Docker image references were not specified at the compile time.

### CLI-only build using the local version of Go

```shell
$ go build -o forta .
```

### CLI-only build using a specific version of Go

Edit Go image version at build stage inside `Dockerfile.cli` and then:

```shell
$ make main
```

## Run the node

### Run

```shell
$ forta init # if you haven't initialized and configured your Forta directory yet
$ forta run
```

### View logs

CLI logs are made available via stdout. Logs for the rest of the node services and agents can be inspected by doing:

```shell
$ docker ps # see the running containers from here
$ docker logs -f <container_id>
```

### Stop

```
CTRL-C
```

## Bug Bounty

We have a [bug bounty program on Immunefi](https://immunefi.com/bounty/forta). Please report any security issues you find through the Immunefi dashboard, or reach out to [tech@forta.org](mailto:tech@forta.org)
