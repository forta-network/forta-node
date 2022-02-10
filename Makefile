containers:
	docker build -t forta-protocol/forta-node -f Dockerfile.node .
	docker pull nats:2.3.2

containers-dev:
	DOCKER_BUILDKIT=1 docker build -t forta-protocol/forta-node -f Dockerfile.buildkit.node .
	docker pull nats:2.3.2

main:
	docker build -t build-forta -f Dockerfile.cli .
	docker create --name build-forta build-forta
	docker cp build-forta:/main forta
	docker rm -f build-forta
	chmod 755 forta

proto:
	protoc -I=protocol --go_out=protocol/. protocol/metrics.proto
	protoc -I=protocol --go-grpc_out=protocol/. --go_out=protocol/. protocol/agent.proto
	protoc -I=protocol --go-grpc_out=protocol/. --go_out=protocol/. protocol/publisher.proto
	protoc -I=protocol --go_out=protocol/. protocol/batch.proto

mocks:
	mockgen -source ethereum/client.go -destination ethereum/mocks/mock_client.go
	mockgen -source clients/interfaces.go -destination clients/mocks/mock_clients.go
	mockgen -source feeds/interfaces.go -destination feeds/mocks/mock_feeds.go
	mockgen -source services/registry/registry.go -destination services/registry/mocks/mock_registry.go
	mockgen -source store/registry.go -destination store/mocks/mock_registry.go
	mockgen -source store/updater.go -destination store/mocks/mock_updater.go
	mockgen -source store/ipfs.go -destination store/mocks/mock_ipfs.go
	mockgen -source ethereum/contract_backend.go -destination ethereum/mocks/mock_ethclient.go

test:
	go test -v -count=1 ./...

perf-test:
	go test ./... -tags=perf_test

run:
	go build -o forta . && ./forta --passphrase 123

abigen:
	./scripts/abigen.sh agent_registry
	./scripts/abigen.sh dispatch
	./scripts/abigen.sh scanner_node_version
	./scripts/abigen.sh scanner_registry "--alias _register=underscoreRegister"

build-local: ## Build for local installation from source
	./scripts/build-for-local.sh

build-remote: ## Try the "remote" containers option for build
	./scripts/build-for-release.sh disco-dev.forta.network

.PHONY: install
install: build-local ## Single install target for local installation
	cp forta /usr/local/bin/forta
