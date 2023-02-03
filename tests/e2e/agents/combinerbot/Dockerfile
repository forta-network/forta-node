# syntax = docker/dockerfile:latest

FROM alpine AS base
COPY test-bot-combiner /main
EXPOSE 50051

ENTRYPOINT [ "/main" ]
