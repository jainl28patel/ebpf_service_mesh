#!/bin/bash

# Replace placeholder in JSON configuration file
awk -v var="$SERVICE_NAME" '{gsub(/\$\{SERVICE_NAME\}/, var)}1' client-config.json > client-config-new.json

# mount /sys/fs/bpf
mount -t bpf bpf /sys/fs/bpf/ || die "Unable to mount /sys/fs/bpf inside test environment"

# Load XDP program here / Start GO client agent
./user/user-agent > /tmp/gologs_"$SERVICE_NAME".logs 2>&1 &

# Start Consul Client agent
consul agent -config-file=./client-config-new.json > /tmp/consul-client.log 2>&1 &

# Configuring DNS
apt install -y dnsmasq >/dev/null                          # move to custom image
apt install -y dnsutils >/dev/null
echo -e "# Forward requests for the consul domain to 127.0.0.1:8600\nserver=/consul/127.0.0.1#8600\n\n# Use Docker's default DNS resolver for all other domains\n# This is usually set to Docker's DNS address, typically 127.0.0.11\nserver=127.0.0.11" > /etc/dnsmasq.conf
echo "nameserver 127.0.0.1" > /etc/resolv.conf
service dnsmasq restart

# Run exe
./app/udp_echo_server