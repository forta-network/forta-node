
build-safe-node:
	docker build -t openzeppelin/safe-node -f Dockerfile-safe-node .

build-main:
	go build -o safe-node main.go

build: build-safe-node build-main

test:
	go test ./...