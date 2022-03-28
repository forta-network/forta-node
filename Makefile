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

mocks:
	mockgen -source clients/interfaces.go -destination clients/mocks/mock_clients.go
	mockgen -source services/registry/registry.go -destination services/registry/mocks/mock_registry.go
	mockgen -source store/registry.go -destination store/mocks/mock_registry.go

test:
	go test -v -count=1 ./...

perf-test:
	go test ./... -tags=perf_test

e2e-test:
	echo 'mode: atomic' > tests/e2e/.forta/coverage.txt
	go test -v -count=1 -tags=e2e_test -race \
		-covermode=atomic -coverprofile=coverage.tmp \
		-coverpkg $$(go list ./... | grep -v tests | tr "\n" ",") \
		github.com/forta-protocol/forta-node/tests/e2e
	tail -n +2 coverage.tmp >> tests/e2e/.forta/coverage.txt

run:
	go build -o forta . && ./forta --passphrase 123

build-local: ## Build for local installation from source
	./scripts/build-for-local.sh

build-remote: ## Try the "remote" containers option for build
	./scripts/build-for-release.sh disco-dev.forta.network

.PHONY: install
install: build-local ## Single install target for local installation
	cp forta /usr/local/bin/forta
