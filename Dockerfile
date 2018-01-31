FROM golang:1.9-stretch AS build

RUN mkdir -p /go/src/github.com/wkozyra95/go-graphql-starter && \
  mkdir -p /build && \
  go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/github.com/wkozyra95/go-graphql-starter

RUN cd /go/src/github.com/wkozyra95/go-graphql-starter && \
    dep ensure && \
    go build -i -o /build/go-graphql-starter
  
FROM debian:9
COPY --from=build /build /root/backend

ENTRYPOINT /root/backend/go-graphql-starter
