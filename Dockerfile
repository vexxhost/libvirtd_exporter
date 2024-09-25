ARG RELEASE=bookworm

FROM golang:1.22-${RELEASE} AS builder
WORKDIR /go/src/app
RUN apt-get update && \
    apt-get -y install libvirt-dev && \
    apt-get clean all
COPY . .
RUN go build

FROM debian:${RELEASE}
RUN apt-get update && \
    apt-get -y install libvirt0 && \
    apt-get clean all
COPY --from=builder /go/src/app/libvirtd_exporter /usr/local/bin/libvirtd_exporter
ENTRYPOINT ["/usr/local/bin/libvirtd_exporter"]
