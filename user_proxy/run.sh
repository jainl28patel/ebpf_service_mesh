#!/bin/bash

# Replace placeholder in JSON configuration file
awk -v var="$SERVICE_NAME" '{gsub(/\$\{SERVICE_NAME\}/, var)}1' client-config.json > client-config-new.json

# mount /sys/fs/bpf
mount -t bpf bpf /sys/fs/bpf/ || die "Unable to mount /sys/fs/bpf inside test environment"

# Load XDP program here / Start GO client agent
./user/user-agent "$SERVICE_NAME" > /tmp/gologs_"$SERVICE_NAME".logs 2>&1 &

# Start Consul Client agent
consul agent -config-file=./client-config-new.json > /tmp/consul-client.log 2>&1 &

# provide permission
chmod +x metric_update.py

# Run exe
./app/udp_echo_server