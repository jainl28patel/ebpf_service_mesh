#include "xdp_main.h"

SEC("tc")
int tc_egress(struct __sk_buff *skb) {
    // Pointer to the start of the packet data
    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;

    // Ethernet header
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end) {
        return TC_ACT_OK;  // Drop the packet if it's invalid
    }

    // Check if the packet is IP
    if (eth->h_proto == __constant_htons(ETH_P_IP)) {
        struct iphdr *iph = data + sizeof(struct ethhdr);
        if ((void *)(iph + 1) > data_end) {
            return TC_ACT_OK;
        }

        // check if UDP packet
        if(iph->protocol == IPPROTO_UDP)
        {
            struct udphdr *udh = data + sizeof(struct ethhdr) + sizeof(struct iphdr);
            if((void*)(udh + 1) > data_end) {
                return TC_ACT_OK;
            }
            // Log the IP addresses and protocol
            bpf_printk("TC_EGRESS: src_ip=%u, dest_ip=%u, protocol=%u", 
                        iph->saddr, iph->daddr, iph->protocol);
            bpf_printk("PORT_SOURCE=%u, PORT_DESTINATION=%u", bpf_ntohs(udh->source), bpf_ntohs(udh->dest));

            // bpf_printk("[CHANGED] TC_EGRESS: src_ip=%u, dest_ip=%u, protocol=%u", 
            //             iph->saddr, iph->daddr, iph->protocol);
            // bpf_printk("[CHANGED] PORT_SOURCE=%u, PORT_DESTINATION=%u\n", bpf_ntohs(udh->source), bpf_ntohs(udh->dest));

            // read the map and print value
            struct catalog_key key;
            __builtin_memset(&key, 0, sizeof(key));  // Initialize key to zero
            __builtin_memcpy(key.hostname, "consul", sizeof("consul"));  // Copy the key "hello" into the struct

            bpf_printk("Looking up key: %s\n", key.hostname);

            // Look up value from map
            struct catalog_value *value;
            value = bpf_map_lookup_elem(&service_catalog, &key);

            if (value) {
                // Print the service IP from the map value
                bpf_printk("Service IP: %u\n", bpf_ntohs(value->service_ip));
            } else {
                bpf_printk("No value found for key :%u \n", (value->service_ip));
            }

        }
    }

    return TC_ACT_OK;  // Allow the packet to continue
}

char __license[] SEC("license") = "GPL";