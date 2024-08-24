#!/bin/bash

# Start Consul Client agent
consul agent -config-file=./client-config.json > /tmp/consul-client.log 2>&1 &

# Load XDP program here / Start GO client agent

# Run exe
./app/udp_echo_server