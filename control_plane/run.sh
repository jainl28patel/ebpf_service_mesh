#!/bin/bash
awk -v var="$SERVICE_NAME" '{gsub(/\$\{SERVICE_NAME\}/, var)}1' server-config.json > server-config-new.json
consul agent -config-file=./server-config-new.json > /tmp/consul-server.log