containers:
	docker pull nats:2.3.2
	docker build -t openzeppelin/fortify-scanner -f Dockerfile-scanner .
	docker build -t openzeppelin/fortify-query -f Dockerfile-query .
	docker build -t openzeppelin/fortify-json-rpc -f Dockerfile-json-rpc .

docker-login:
	aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 997179694723.dkr.ecr.us-west-2.amazonaws.com

ecr:
	docker tag openzeppelin/fortify-scanner:latest 997179694723.dkr.ecr.us-west-2.amazonaws.com/fortify-scanner:latest
	docker tag openzeppelin/fortify-query:latest 997179694723.dkr.ecr.us-west-2.amazonaws.com/fortify-query:latest
	docker tag openzeppelin/fortify-json-rpc:latest 997179694723.dkr.ecr.us-west-2.amazonaws.com/fortify-json-rpc:latest
	docker push 997179694723.dkr.ecr.us-west-2.amazonaws.com/fortify-scanner:latest
	docker push 997179694723.dkr.ecr.us-west-2.amazonaws.com/fortify-query:latest
	docker push 997179694723.dkr.ecr.us-west-2.amazonaws.com/fortify-json-rpc:latest

main:
	docker build -t build-fortify -f Dockerfile-cli .
	docker create --name build-fortify build-fortify
	docker cp build-fortify:/main fortify
	docker rm -f build-fortify
	chmod 755 fortify

proto:
	protoc -I=protocol --go-grpc_out=protocol/. --go_out=protocol/. protocol/agent.proto
	protoc -I=protocol --go-grpc_out=protocol/. --go_out=protocol/. protocol/query.proto

mocks:
	mockgen -source ethereum/client.go -destination ethereum/mocks/mock_client.go
	mockgen -source clients/interfaces.go -destination clients/mocks/mock_clients.go
	mockgen -source feeds/interfaces.go -destination feeds/mocks/mock_feeds.go
	mockgen -source services/registry/registry.go -destination services/registry/mocks/mock_registry.go

build: proto main containers

test:
	go test -v -count=1 ./...

run:
	go build -o build/fortify . && ./build/fortify --passphrase 123

abigen:
	abigen --abi ./contracts/pool_registry.json --out ./contracts/pool_registry.go --pkg contracts --type PoolRegistry
	abigen --abi ./contracts/alerts.json --out ./contracts/alerts.go --pkg contracts --type Alerts
