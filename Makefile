containers:
	docker build -t forta-network/forta-node -f Dockerfile.node .
	docker pull nats:2.3.2

containers-dev:
	DOCKER_BUILDKIT=1 docker build --network=host -t forta-network/forta-node -f Dockerfile.buildkit.node .
	docker pull nats:2.3.2

main:
	docker build -t build-forta -f Dockerfile.cli .
	docker create --name build-forta build-forta
	docker cp build-forta:/main forta
	docker rm -f build-forta
	chmod 755 forta

mocks:
	mockgen -source clients/interfaces.go -destination clients/mocks/mock_clients.go
	mockgen -source services/registry/registry.go -destination services/registry/mocks/mock_registry.go
	mockgen -source store/registry.go -destination store/mocks/mock_registry.go
	mockgen -source services/storage/ipfs.go -destination services/storage/mocks/mock_ipfs.go

test:
	go test -v -count=1 ./...

perf-test:
	go test ./... -tags=perf_test

.PHONY: e2e-test-deps
e2e-test-deps:
	./tests/e2e/deps-start.sh

.PHONY: e2e-test
e2e-test:
	rm -rf tests/e2e/.forta/coverage
	mkdir -p tests/e2e/.forta/coverage
	rm -rf tests/e2e/.forta-local/coverage
	mkdir -p tests/e2e/.forta-local/coverage

	./tests/e2e/build.sh

	cd tests/e2e && E2E_TEST=1 go test -v -count=1 .

	rm -rf coverage
	cp -r tests/e2e/.forta/coverage .
	cp -r tests/e2e/.forta-local/coverage/* coverage/
	./scripts/total-coverage.sh e2e

run:
	go build -o forta . && ./forta --passphrase 123

build-local: ## Build for local installation from source
	./scripts/build-for-local.sh

build-remote: ## Try the "remote" containers option for build
	./scripts/build-for-release.sh disco-dev.forta.network

.PHONY: install
install: build-local ## Single install target for local installation
	cp forta /usr/local/bin/forta
