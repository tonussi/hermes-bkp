FROM golang:1.15-alpine AS build-env

WORKDIR /src
COPY ./cmd/hermes ./cmd/hermes
COPY ./pkg ./pkg
COPY go.* ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/hermes ./cmd/hermes

FROM golang:1.15-alpine

COPY --from=build-env /bin/hermes /bin/hermes

EXPOSE 8000
EXPOSE 9000
EXPOSE 10000

ENTRYPOINT [ "/bin/hermes" ]