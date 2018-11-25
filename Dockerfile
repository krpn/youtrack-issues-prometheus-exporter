FROM golang:1.11-alpine3.8 as builder

RUN \
    cd / && \
    apk update && \
    apk add --no-cache git ca-certificates make tzdata curl gcc libc-dev && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN \
    mkdir -p src/github.com/krpn && \
    cd src/github.com/krpn && \
    git clone https://github.com/krpn/youtrack-issues-prometheus-exporter && \
    cd youtrack-issues-prometheus-exporter && \
    dep ensure -v && \
    go test ./... && \
    cd cmd/youtrack-issues-prometheus-exporter && \
    CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o youtrack-issues-prometheus-exporter


FROM alpine:3.8
COPY --from=builder /go/src/github.com/krpn/youtrack-issues-prometheus-exporter/cmd/youtrack-issues-prometheus-exporter/youtrack-issues-prometheus-exporter /
RUN apk add --no-cache ca-certificates tzdata
EXPOSE 8080
ENTRYPOINT ["/youtrack-issues-prometheus-exporter"]