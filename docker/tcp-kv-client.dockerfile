FROM golang:1.15-alpine AS build-env

WORKDIR /src
COPY ./cmd/tcp-kv-client ./cmd/tcp-kv-client
COPY ./pkg/kv ./pkg/kv
COPY go.* ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/tcp-kv-client ./cmd/tcp-kv-client

FROM golang:1.15-alpine

COPY --from=build-env /bin/tcp-kv-client /bin/tcp-kv-client

ENTRYPOINT [ "/bin/tcp-kv-client" ]