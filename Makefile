containers:
	docker build -t forta-network/forta-node -f Dockerfile.node .
	docker pull nats:2.3.2

containers-dev:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o forta-node cmd/node/main.go
	DOCKER_BUILDKIT=1 docker build --no-cache --network=host -t forta-network/forta-node -f Dockerfile.buildkit.dev.node .
	docker pull nats:2.3.2

main:
	docker build -t build-forta -f Dockerfile.cli .
	docker create --name build-forta build-forta
	docker cp build-forta:/main forta
	docker rm -f build-forta
	chmod 755 forta

mocks:
	mockgen -source clients/interfaces.go -destination clients/mocks/mock_clients.go
	mockgen -source clients/ratelimiter/rate_limiter.go -destination clients/ratelimiter/mocks/mock_rate_limiter.go
	mockgen -source services/registry/registry.go -destination services/registry/mocks/mock_registry.go
	mockgen -source store/registry.go -destination store/mocks/mock_registry.go
	mockgen -source services/storage/ipfs.go -destination services/storage/mocks/mock_ipfs.go
	mockgen -source store/scanner_release.go -destination store/mocks/mock_scanner_release.go

test:
	go test -v -count=1 ./... -coverprofile=coverage.out

.PHONY: coverage
coverage:
	go tool cover -func=coverage.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}'

coverage-func:
	go tool cover -func=coverage.out

coverage-html:
	go tool cover -html=coverage.out -o=coverage.html

perf-test:
	go test ./... -tags=perf_test

MOCKREG = $$(pwd)/tests/e2e/misccontracts/contract_mock_registry

.PHONY: e2e-test-mock
e2e-test-mock:
	solc --bin --abi -o $(MOCKREG) --include-path . --base-path $(MOCKREG) --overwrite --input-file $(MOCKREG)/MockRegistry.sol
	abigen --out $(MOCKREG)/mock_registry.go --pkg contract_mock_registry --type MockRegistry --abi $(MOCKREG)/MockRegistry.abi --bin $(MOCKREG)/MockRegistry.bin

.PHONY: e2e-test-deps
e2e-test-deps:
	./tests/e2e/deps-start.sh

.PHONY: e2e-test
e2e-test:
	./tests/e2e/build.sh

	cd tests/e2e && E2E_TEST=1 go test -v -count=1 .

run:
	go build -o forta . && ./forta --passphrase 123

build-local: ## Build for local installation from source
	./scripts/build-for-local.sh

build-remote: ## Try the "remote" containers option for build
	./scripts/build-for-release.sh disco-dev.forta.network

.PHONY: install
install: build-local ## Single install target for local installation
	cp forta /usr/local/bin/forta

.PHONY: update-core
update-core:
	go get github.com/forta-network/forta-core-go && go mod tidy
