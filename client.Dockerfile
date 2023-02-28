FROM golang:1.20.1 as builder

ENV GO111MODULE=on
ENV GOPROXY=direct

WORKDIR /go/src/gitlab.com/yvv4git/task-wow
COPY . .

RUN export GOOS=linux GOARCH=amd64 && \
    go build -tags netgo -ldflags '-w -extldflags "-static"' -o /go/bin/client ./cmd/client/client.go


FROM scratch

WORKDIR /root
COPY --from=builder /go/bin/client ./client

ENTRYPOINT ["./client"]
