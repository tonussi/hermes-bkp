FROM golang:1.16-alpine AS build
WORKDIR /go/src/work
ENV CGO_ENABLED=0

COPY go.* ./

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

COPY . .

RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go get -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@master
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o ./hermes ./cmd/hermes
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/hermes ./cmd/hermes

FROM scratch
COPY --from=build /go/bin/dlv /dlv
COPY --from=build /go/src/work/hermes /hermes
EXPOSE 2345
EXPOSE 8000
EXPOSE 9000
EXPOSE 10000
ENTRYPOINT [ "/dlv" ]
