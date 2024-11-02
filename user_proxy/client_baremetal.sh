#!/bin/bash

# Install all dependencies for XDP and eBPFs
sudo apt -y update
sudo apt install -y iproute2 clang llvm libelf-dev libpcap-dev build-essential libc6-dev-i386
sudo apt install -y linux-tools-$(uname -r)
sudo apt install -y linux-headers-$(uname -r)
sudo apt install -y linux-tools-common linux-tools-generic
sudo apt install -y tcpdump
sudo apt install -y libbpf-dev
sudo apt install -y libxdp-dev
sudo apt install -y xdp-tools
sudo apt install -y unzip
sudo apt install -y wget
sudo apt install -y tar
sudo apt install -y curl

# Consul Installation
wget https://releases.hashicorp.com/consul/1.19.1/consul_1.19.1_linux_amd64.zip
unzip consul_1.19.1_linux_amd64.zip
sudo mv ./consul /usr/bin

# GO Installation
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
rm -rf go1.23.0.linux-amd64.tar.gz

# setting permissions for go on cloudlab
sudo chown -R $(whoami) /users/anakin/.cache/go-build
chmod -R u+rwX /users/anakin/.cache/go-build

# Building ebpf program
make

# Build GO App
export PATH=$PATH:/usr/local/go/bin
cd ./user
go get .
go build .
cd ..

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

# Replace placeholder in JSON configuration file
awk -v var="$SERVICE_NAME" '{gsub(/\$\{SERVICE_NAME\}/, var)}1' client-config.json > client-config-temp-new.json
awk -v var="$SERVER_IP" '{gsub(/\$\{SERVER_IP\}/, var)}1' client-config-temp-new.json > client-config-new.json
rm client-config-temp-new.json

# mount /sys/fs/bpf
sudo mount -t bpf bpf /sys/fs/bpf/ || die "Unable to mount /sys/fs/bpf inside test environment"

# Load eBPF program here / Start GO client agent
sudo ./user/user-agent "$SERVICE_NAME" > /tmp/gologs_"$SERVICE_NAME".logs 2>&1 &

# Start Consul Client agent
consul agent -config-file=./client-config-new.json > /tmp/consul-client.log 2>&1 &

# Run exe / userspace program
./user/udp_echo_server &