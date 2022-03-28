# syntax = docker/dockerfile:latest

FROM alpine AS base

FROM golang:1.16.4 AS go-builder
WORKDIR /go/app
COPY go.mod .
COPY go.sum .

RUN --mount=type=cache,target=/root/.cache/go go mod download

COPY . /go/app
RUN go test -c -o main -race -covermode=atomic -coverprofile=coverage.tmp -coverpkg \
	$(go list ./... | grep -v tests | tr "\n" ",") github.com/forta-protocol/forta-node/tests/e2e/cmd/node

FROM base
COPY --from=go-builder /go/app/main /main
EXPOSE 8089 8090

RUN mkdir -p /.forta

# append coverage output to coverage.txt in mounted forta dir
ENTRYPOINT ./main -test.coverprofile=coverage.tmp && tail -n +2 coverage.tmp >> /.forta/coverage.txt
