#!/bin/bash

# Replace placeholder in JSON configuration file
awk -v var="$SERVICE_NAME" '{gsub(/\$\{SERVICE_NAME\}/, var)}1' client-config.json > client-config-new.json

# Start Consul Client agent
consul agent -config-file=./client-config-new.json > /tmp/consul-client.log 2>&1 &

# Load XDP program here / Start GO client agent

# Run exe
./app/udp_echo_server