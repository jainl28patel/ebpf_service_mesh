#pragma once

#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/in.h>
#include <linux/udp.h>
#include <linux/pkt_cls.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

#define TC_ACT_OK 0
#define MAX_MAP_SIZE 1024
#define HOSTNAME_MAX_LEN 256

struct catalog_key {
    char hostname[HOSTNAME_MAX_LEN];     // hostname to resolve
};

struct catalog_value {
    unsigned int service_ip;    // resolved ip provided by consul
};

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__type(key, struct catalog_key);
	__type(value, struct catalog_value);
	__uint(max_entries, MAX_MAP_SIZE);
} service_catalog SEC(".maps");