
build-node:
	docker build -t openzeppelin/zephyr-node -f Dockerfile-zephyr-node .

build-main:
	go build -o zephyr-node main.go

proto:
	protoc -I=protocol --go-grpc_out=protocol/. --go_out=protocol/. protocol/*.proto

build: proto build-main build-node

test:
	go test ./...

