containers:
	docker build -t openzeppelin/fortify-scanner -f Dockerfile-scanner .
	docker build -t openzeppelin/fortify-query -f Dockerfile-query .

main:
	go build -o fortify main.go

proto:
	protoc -I=protocol --go-grpc_out=protocol/. --go_out=protocol/. protocol/agent.proto
	protoc -I=protocol --go-grpc_out=protocol/. --go_out=protocol/. protocol/query.proto

build: proto main containers

test:
	go test ./...

