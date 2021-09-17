#include <linux/bpf.h>
#include <linux/in.h>
#include <linux/if_ether.h>
#include <linux/if_packet.h>
#include <linux/if_vlan.h>
#include <linux/ip.h>
#include "bpf_helper.h"

static inline int parse_ipv4(void *data, __u64 nh_off, void *data_end, __be32 *src)
{
  struct iphdr *iph = data + nh_off;
  if (iph + 1 > data_end)
    return 0;
  *src = iph->saddr;
  return iph->protocol;
}

struct bpf_map_def SEC("maps/deny_ips") deny_ips = {
	.type = BPF_MAP_TYPE_HASH,
	.key_size = sizeof(__u32),
	.value_size = sizeof(__u32),
	.max_entries = 10,
	.pinning = 0,
	.namespace = "",
};

SEC("xdp/bina")
int xdp_bina(struct xdp_md *ctx)
{
  void *data_end = (void *)(long)ctx->data_end;
  void *data = (void *)(long)ctx->data;
  struct ethhdr *eth = data;
  __be32 src_ip;
  __u16 h_proto;
  __u64 nh_off;
  int ipproto;
  __u32 *exists;

  nh_off = sizeof(*eth);
  if (data + nh_off > data_end)
    goto pass;  

  h_proto = eth->h_proto;
  if (h_proto != __constant_htons(ETH_P_IP))
      goto pass;

  ipproto = parse_ipv4(data, nh_off, data_end, &src_ip);
  __u32 intIP = src_ip;
  exists = bpf_map_lookup_elem(&deny_ips, &intIP);
  if (exists)
      return XDP_DROP;

pass:
  return XDP_PASS;
}
