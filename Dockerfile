FROM golang:1.18.4-alpine3.15 as builder
RUN apk add git

WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY *.go .
COPY cmd cmd/

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" ./cmd/runner

FROM alpine:3.16.0
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /build/runner /
RUN touch .env
ENTRYPOINT [ "/runner" ]
