FROM alpine AS base
COPY test-bot-txdetector /main
EXPOSE 50051

LABEL "network.forta.settings.agent-logs.enable"="true"

ENTRYPOINT [ "/main" ]
