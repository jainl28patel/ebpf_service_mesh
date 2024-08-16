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

# Build the project
make

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

ls -la /app

# Load Kernel XDP Program
# xdp-loader load -m skb lo ./kernel/xdp.o

# Run User space program

# Run the application
