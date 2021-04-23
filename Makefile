
node:
	docker build -t openzeppelin/zephyr-node -f Dockerfile-zephyr-node .

main:
	go build -o zephyr-node main.go

proto:
	protoc -I=protocol --go-grpc_out=protocol/. --go_out=protocol/. protocol/*.proto

build: proto main node

test:
	go test ./...

