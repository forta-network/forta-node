containers:
	docker pull nats:2.3.2
	docker build -t forta-protocol/forta-scanner -f Dockerfile.scanner .
	docker build -t forta-protocol/forta-query -f Dockerfile.query .
	docker build -t forta-protocol/forta-json-rpc -f Dockerfile.json-rpc .

main:
	docker build -t build-forta -f Dockerfile.cli .
	docker create --name build-forta build-forta
	docker cp build-forta:/main forta
	docker rm -f build-forta
	chmod 755 forta

proto:
	protoc -I=protocol --go_out=protocol/. protocol/metrics.proto
	protoc -I=protocol --go-grpc_out=protocol/. --go_out=protocol/. protocol/agent.proto
	protoc -I=protocol --go-grpc_out=protocol/. --go_out=protocol/. protocol/query.proto
	protoc -I=protocol --go_out=protocol/. protocol/batch.proto

mocks:
	mockgen -source ethereum/client.go -destination ethereum/mocks/mock_client.go
	mockgen -source clients/interfaces.go -destination clients/mocks/mock_clients.go
	mockgen -source feeds/interfaces.go -destination feeds/mocks/mock_feeds.go
	mockgen -source services/registry/registry.go -destination services/registry/mocks/mock_registry.go
	mockgen -source store/registry.go -destination store/mocks/mock_registry.go

test:
	go test -v -count=1 ./...

run:
	go build -o forta . && ./forta --passphrase 123

abigen:
	abigen --abi ./contracts/agent_registry.json --out ./contracts/agent_registry.go --pkg contracts --type AgentRegistry
	abigen --abi ./contracts/dispatch.json --out ./contracts/dispatch.go --pkg contracts --type Dispatch
	abigen --abi ./contracts/alerts.json --out ./contracts/alerts.go --pkg contracts --type Alerts

build-local: ## Build for local installation from source
	./scripts/build-for-local.sh

build-remote: ## Try the "remote" containers option for build
	./scripts/build-for-release.sh disco-dev.forta.network

.PHONY: install
install: build-local ## Single install target for local installation
	cp forta /usr/bin/forta
