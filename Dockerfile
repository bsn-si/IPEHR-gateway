FROM golang:1.18-alpine3.16 AS build
WORKDIR /srv
COPY src/ .
COPY config.json config.json
RUN CGO_ENABLED=0 go build -o ./bin/ipehr-gateway cmd/ipehrgw/main.go

FROM alpine:3.16
WORKDIR /srv
COPY data/ /data
COPY --from=build /srv/bin/ /srv
COPY --from=build /srv/config.json /srv
CMD ["./ipehr-gateway", "-config=./config.json"]
