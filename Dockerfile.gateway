FROM golang:1.19.0-alpine3.16 AS build

RUN apk update && \
    apk add --no-cache gcc musl-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /srv

COPY ./src/go.mod go.mod
COPY ./src/go.sum go.sum

RUN go mod download
COPY data/ ./data
COPY config.json.example ./config.json
COPY src/ .


RUN go build -o ./bin/ipehr-gateway cmd/ipehrgw/main.go

FROM alpine:3.16

WORKDIR /srv

RUN echo "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" > .blockchain.key
COPY --from=build /srv/bin/ /srv
COPY --from=build /srv/config.json /srv/config.json
COPY --from=build /srv/data /srv/data
COPY --from=build /srv/pkg/indexer /srv/inbexer

ENTRYPOINT [ "/srv/ipehr-gateway", "-config=/srv/config.json"]
