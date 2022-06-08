# syntax = docker/dockerfile:latest

FROM alpine AS base

FROM golang:1.16.4-alpine AS go-builder

RUN apk add --no-cache bash gcc g++

WORKDIR /go/app
COPY go.mod .
COPY go.sum .

RUN --mount=type=cache,target=/root/.cache/go go mod download

COPY . /go/app
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go test -c -o main -race -covermode=atomic -coverpkg \
	$(go list ./... | grep -v tests | tr "\n" ",") github.com/forta-network/forta-node/tests/e2e/cmd/node

FROM base

RUN apk add --no-cache bash gcc g++ tini

COPY --from=go-builder /go/app/main /usr/bin/forta-node
COPY --from=go-builder /go/app/tests/e2e/cmd/node/run.sh /run.sh
EXPOSE 8089 8090

RUN mkdir -p /.forta/coverage

RUN apk add --no-cache iptables

# receives '/forta-node <command>' from runner and supervisor
# https://docs.docker.com/engine/reference/builder/#understand-how-cmd-and-entrypoint-interact
ENTRYPOINT [ "tini", "-v", "--", "/run.sh" ]
