
build-safe-node:
	docker build -t openzeppelin/safe-node -f Dockerfile-safe-node .

build-main:
	go build -o safe-node main.go

proto:
	protoc -I=protocol --go-grpc_out=protocol/. --go_out=protocol/. protocol/*.proto

build: proto build-main build-safe-node

test:
	go test ./...

