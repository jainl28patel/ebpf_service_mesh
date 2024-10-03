#!/bin/bash

# Install all dependencies for XDP and eBPFs
sudo apt -y update
sudo apt install -y iproute2 clang llvm libelf-dev libpcap-dev build-essential libc6-dev-i386
sudo apt install -y linux-tools-$(uname -r)
sudo apt install -y linux-headers-$(uname -r)
sudo apt install -y linux-tools-common linux-tools-generic
sudo apt install -y tcpdump
sudo apt install -y libbpf-dev libxdp-dev
sudo apt install -y xdp-tools
sudo apt install -y unzip
sudo apt install -y wget
sudo apt install -y tar

# Consul Installation
wget https://releases.hashicorp.com/consul/1.19.1/consul_1.19.1_linux_amd64.zip
unzip consul_1.19.1_linux_amd64.zip
sudo mv ./consul /usr/bin

consul agent -config-file=./server-config.json > /tmp/consul-server.log