# syntax = docker/dockerfile:latest

FROM alpine AS base

FROM golang:1.19 AS go-builder
WORKDIR /go/app
COPY go.mod .
COPY go.sum .

RUN --mount=type=cache,target=/root/.cache/go go mod download

COPY . /go/app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/app/main \
	/go/app/tests/e2e/agents/txdetectoragent/main.go

FROM base
COPY --from=go-builder /go/app/main /main
EXPOSE 50051

ENTRYPOINT [ "/main" ]
