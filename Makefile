containers:
	docker build -t openzeppelin/zephyr-node -f Dockerfile-zephyr-node .
	docker build -t openzeppelin/zephyr-proxy -f Dockerfile-zephyr-proxy .

main:
	go build -o zephyr-node main.go

proto:
	protoc -I=protocol --go-grpc_out=protocol/. --go_out=protocol/. protocol/*.proto

build: proto main containers

test:
	go test ./...

