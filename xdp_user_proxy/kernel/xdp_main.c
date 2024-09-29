#include "xdp_main.h"
#include "dns.h"

SEC("tc_egress")
int tc_egress_helper(struct __sk_buff *skb) {
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
            bpf_printk("TC_EGRESS: src_ip=%u, dest_ip=%u, protocol=%u, PORT_SOURCE=%u, PORT_DESTINATION=%u", 
                        iph->saddr, iph->daddr, iph->protocol, bpf_ntohs(udh->source), bpf_ntohs(udh->dest));

            if(iph->daddr != 184549503) return TC_ACT_OK;

            
            // Parsing for the DNS header
            // Header size : 12bytes

            void* ptr = data + sizeof(struct ethhdr) + sizeof(struct iphdr) + sizeof(struct udphdr) + sizeof(struct DNS_HEADER);
            if (ptr > data_end) {
                return TC_ACT_OK;
            }

            const char* subdomain = "consul";
            int len = 6;
            int redirect = 0;
            
            for(int no_of_octet = 0; no_of_octet < MAX_OCTET; no_of_octet++)
            {
                if ((void *)(ptr + 1) > data_end) {
                    return TC_ACT_OK;
                }
                unsigned char size = *((char*)(ptr));

                if (size == 0) {
                    break;  // End of the domain name string
                } else {
                    redirect = 0;
                }

                if ((void*)(ptr + 1 + size) > data_end) {
                    return TC_ACT_OK;
                }

                char* octet = (char*)(ptr + 1);

                if(size == len)
                {
                    int same = 1;
                    for (int sz = 0; sz < len; sz++) {
                        if ((void*)(octet + sz + 1) > data_end) {
                            return TC_ACT_OK;  // Avoid buffer overrun
                        }
                        if(octet[sz]!=subdomain[sz]) {
                            same = 0;
                            break;
                        }
                    }
                    if(same==1) {
                        redirect = 1;
                    }
                }
                ptr += 1 + size;  // Move to the next octet
            }

            if(redirect==0) {
                return TC_ACT_OK;
            }

            // redirect the dns request to necessary port
            bpf_printk("REDIRECTING THE REQUEST");

            // Change the destination IP to 127.0.0.1 (0x7f000001 in hex)
            iph->daddr = 16777343;

            // Change the destination port to 8600
            udh->dest = bpf_htons(8600);

            // Recalculate IP checksum
            iph->check = 0;
            __u32 csum = bpf_csum_diff(0, 0, (__be32 *)iph, sizeof(*iph), 0);
            iph->check = ~csum;
            // Prepare the old value for checksum calculation

            bpf_printk("Checksum : %d",iph->check);

            // Recalculate UDP checksum
            udh->check = 0;  // Reset checksum first
            // __u32 udp_csum = bpf_csum_diff(0, 0, (__be32 *)udh, sizeof(*udh), 0);
            // udh->check = ~udp_csum;  // Set the new checksum

            // Action to modify the packet
            // return bpf_redirect(skb->ifindex, BPF_F_INGRESS);  // Redirect to ingress for handling
        }
    }

    return TC_ACT_OK;  // Allow the packet to continue
}


SEC("tc_ingress")
int tc_ingress_helper(struct __sk_buff *skb) {
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
            bpf_printk("TC_INGRESS: src_ip=%u, dest_ip=%u, protocol=%u, PORT_SOURCE=%u, PORT_DESTINATION=%u", 
                        iph->saddr, iph->daddr, iph->protocol, bpf_ntohs(udh->source), bpf_ntohs(udh->dest));

            if(iph->daddr != 184549503) return TC_ACT_OK;

            
            // Parsing for the DNS header
            // Header size : 12bytes

            // void* ptr = data + sizeof(struct ethhdr) + sizeof(struct iphdr) + sizeof(struct udphdr) + sizeof(struct DNS_HEADER);
            // if (ptr > data_end) {
            //     return TC_ACT_OK;
            // }

            // const char* subdomain = "consul";
            // int len = 6;
            // int redirect = 0;
            
            // for(int no_of_octet = 0; no_of_octet < MAX_OCTET; no_of_octet++)
            // {
            //     if ((void *)(ptr + 1) > data_end) {
            //         return TC_ACT_OK;
            //     }
            //     unsigned char size = *((char*)(ptr));

            //     if (size == 0) {
            //         break;  // End of the domain name string
            //     } else {
            //         redirect = 0;
            //     }

            //     if ((void*)(ptr + 1 + size) > data_end) {
            //         return TC_ACT_OK;
            //     }

            //     char* octet = (char*)(ptr + 1);

            //     if(size == len)
            //     {
            //         int same = 1;
            //         for (int sz = 0; sz < len; sz++) {
            //             if ((void*)(octet + sz + 1) > data_end) {
            //                 return TC_ACT_OK;  // Avoid buffer overrun
            //             }
            //             if(octet[sz]!=subdomain[sz]) {
            //                 same = 0;
            //                 break;
            //             }
            //         }
            //         if(same==1) {
            //             redirect = 1;
            //         }
            //     }
            //     ptr += 1 + size;  // Move to the next octet
            // }

            // if(redirect==0) {
            //     return TC_ACT_OK;
            // }

            // // redirect the dns request to necessary port
            // bpf_printk("REDIRECTING THE REQUEST");

            // // Change the destination IP to 127.0.0.1 (0x7f000001 in hex)
            // iph->daddr = 16777343;

            // // Change the destination port to 8600
            // udh->dest = bpf_htons(8600);

            // // Recalculate IP checksum
            // iph->check = 0;
            // __u32 ip_sum = 0;
            // bpf_csum_diff(0, 0, (__be32 *)iph, sizeof(*iph), ip_sum);
            // iph->check = ip_sum;

            // // Set UDP checksum to 0 (optional, or recalculate if necessary)
            // udh->check = 0;

            // // Action to modify the packet
            // return bpf_redirect(skb->ifindex, BPF_F_INGRESS);  // Redirect to ingress for handling
        }
    }

    return TC_ACT_OK;  // Allow the packet to continue
}


char __license[] SEC("license") = "GPL";