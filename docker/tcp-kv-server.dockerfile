FROM golang:1.15-alpine AS build-env

WORKDIR /src
COPY ./cmd/tcp-kv-server ./cmd/tcp-kv-server
COPY ./pkg/kv ./pkg/kv
COPY go.* ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/tcp-kv-server ./cmd/tcp-kv-server

FROM golang:1.15-alpine

COPY --from=build-env /bin/tcp-kv-server /bin/tcp-kv-server

EXPOSE 8001

ENTRYPOINT [ "/bin/tcp-kv-server" ]