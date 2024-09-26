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
apt install -y tar
apt install -y dnsmasq
apt install -y dnsutils

# Consul Installation
wget https://releases.hashicorp.com/consul/1.19.1/consul_1.19.1_linux_amd64.zip
unzip consul_1.19.1_linux_amd64.zip
mv ./consul /usr/bin

# GO Installation
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
rm -rf go1.23.0.linux-amd64.tar.gz