containers:
	docker pull nats:2.3.2
	docker build -t forta-network/forta-scanner -f Dockerfile-scanner .
	docker build -t forta-network/forta-query -f Dockerfile-query .
	docker build -t forta-network/forta-json-rpc -f Dockerfile-json-rpc .

docker-login:
	aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 997179694723.dkr.ecr.us-west-2.amazonaws.com

ecr:
	docker tag forta-network/forta-scanner:latest 997179694723.dkr.ecr.us-west-2.amazonaws.com/forta-scanner:latest
	docker tag forta-network/forta-query:latest 997179694723.dkr.ecr.us-west-2.amazonaws.com/forta-query:latest
	docker tag forta-network/forta-json-rpc:latest 997179694723.dkr.ecr.us-west-2.amazonaws.com/forta-json-rpc:latest
	docker push 997179694723.dkr.ecr.us-west-2.amazonaws.com/forta-scanner:latest
	docker push 997179694723.dkr.ecr.us-west-2.amazonaws.com/forta-query:latest
	docker push 997179694723.dkr.ecr.us-west-2.amazonaws.com/forta-json-rpc:latest

main:
	docker build -t build-forta -f Dockerfile-cli .
	docker create --name build-forta build-forta
	docker cp build-forta:/main forta
	docker rm -f build-forta
	chmod 755 forta

proto:
	protoc -I=protocol --go-grpc_out=protocol/. --go_out=protocol/. protocol/agent.proto
	protoc -I=protocol --go-grpc_out=protocol/. --go_out=protocol/. protocol/query.proto
	protoc -I=protocol --go_out=protocol/. protocol/batch.proto

mocks:
	mockgen -source ethereum/client.go -destination ethereum/mocks/mock_client.go
	mockgen -source clients/interfaces.go -destination clients/mocks/mock_clients.go
	mockgen -source feeds/interfaces.go -destination feeds/mocks/mock_feeds.go
	mockgen -source services/registry/registry.go -destination services/registry/mocks/mock_registry.go

build: proto main containers

test:
	go test -v -count=1 ./...

run:
	go build -o build/forta . && ./build/forta --passphrase 123

abigen:
	abigen --abi ./contracts/agent_registry.json --out ./contracts/agent_registry.go --pkg contracts --type AgentRegistry
	abigen --abi ./contracts/scanner_registry.json --out ./contracts/scanner_registry.go --pkg contracts --type ScannerRegistry
	abigen --abi ./contracts/alerts.json --out ./contracts/alerts.go --pkg contracts --type Alerts
