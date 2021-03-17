FROM golang:1.15-alpine AS build-env

WORKDIR /src
COPY ./cmd/tcp-kv-hashicorp-raft ./cmd/tcp-kv-hashicorp-raft
COPY ./pkg/kv ./pkg/kv
COPY go.* ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/tcp-kv-hashicorp-raft ./cmd/tcp-kv-hashicorp-raft

FROM golang:1.15-alpine

COPY --from=build-env /bin/tcp-kv-hashicorp-raft /bin/tcp-kv-hashicorp-raft

EXPOSE 8000
EXPOSE 9000
EXPOSE 10000

ENTRYPOINT [ "/bin/tcp-kv-hashicorp-raft" ]