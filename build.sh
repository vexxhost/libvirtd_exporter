#!/bin/bash
# Small script to build libvirtd-exporter using docker and export the resulting binary

set -e
docker build -t libvirtd-exporter .

# Start container with result
docker run  --entrypoint sh -itd --name libvirtd-exporter libvirtd-exporter

# Extract build atifact to host
docker cp libvirtd-exporter:/libvirtd_exporter .

docker stop libvirtd-exporter
docker rm -f libvirtd-exporter
