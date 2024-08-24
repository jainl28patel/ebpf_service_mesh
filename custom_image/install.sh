#!/bin/bash

# Install all dependencies for XDP and eBPFs
apt -y update
apt install -y iproute2 clang llvm libelf-dev libpcap-dev build-essential libc6-dev-i386
apt install -y linux-tools-$(uname -r)
apt install -y linux-headers-$(uname -r)
apt install -y linux-tools-common linux-tools-generic
apt install -y tcpdump
apt install -y libbpf-dev libxdp-dev
apt install -y xdp-tools
apt install -y unzip
apt install -y wget

wget https://releases.hashicorp.com/consul/1.19.1/consul_1.19.1_linux_amd64.zip
unzip consul_1.19.1_linux_amd64.zip
mv ./consul /usr/bin