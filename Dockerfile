FROM golang:1.22 AS builder
WORKDIR /go/src/app
RUN apt-get update && \
    apt-get -y install libvirt-dev && \
    apt-get clean all
COPY . .
RUN go build

FROM golang:1.22
RUN apt-get update && \
    apt-get -y install libvirt0 && \
    apt-get clean all
COPY --from=builder /go/src/app/libvirtd_exporter /libvirtd_exporter
ENTRYPOINT ["/libvirtd_exporter"]
