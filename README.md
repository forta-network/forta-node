![Build](https://github.com/forta-protocol/forta-node/actions/workflows/codedeploy-staging.yml/badge.svg)

# forta-node

Forta node CLI is a Docker container supervisor that runs and manages multiple services and
agents to scan a blockchain network and produce alerts.

## Dependencies

1. [Install Docker](https://docs.docker.com/get-docker/) and start Docker service
2. [Install golang 1.16](https://golang.org/doc/install)

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

It takes a while to build all of the container images. For a faster iteration in local development,
it is sufficient to just build any of the changed containers or the CLI binary. The CLI requires
`forta-protocol/forta-<service>:latest` containers to be available by default and uses the local ones.

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

### Initialize

```shell
$ forta init
```

- This will create a Forta node directory under `~/.forta` by default. You can use the `--dir` flag
to override it.
- Fix the default `config.yml` in Forta dir. You can use the `--config` flag if you want to use a different config file.

See `forta account` command in the CLI help output if you need to work with a specific private key.

### Run

Provide `FORTA_PASSPHRASE` env var or the flag `--passphrase` so your private key can be decrypted on startup.

```shell
$ forta run
```

### View logs

CLI and supervisor logs are made available via stdout. Logs for the rest of the node services and
agents can be inspected by doing:

```shell
$ docker ps # see the running containers from here
$ docker logs -f <container_id>
```

### Stop

```
CTRL-C
```

### See local alerts

Use the port in config (`query.port`) to access the local alerts API:


```shell
$ curl -s http://localhost:8778/alerts
```
