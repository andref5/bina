#include <linux/bpf.h>
#include <linux/in.h>
#include <linux/if_ether.h>
#include <linux/if_packet.h>
#include <linux/if_vlan.h>
#include <linux/ip.h>

static inline int parse_ipv4(void *data, __u64 nh_off, void *data_end, __be32 *src)
{
  struct iphdr *iph = data + nh_off;
  if (iph + 1 > data_end)
    return 0;
  *src = iph->saddr;
  return iph->protocol;
}

int xdp_bina(struct xdp_md *ctx)
{
  void *data_end = (void *)(long)ctx->data_end;
  void *data = (void *)(long)ctx->data;
  struct ethhdr *eth = data;
  __be32 src_ip;
  __u16 h_proto;
  __u64 nh_off;
  int ipproto;

  nh_off = sizeof(*eth);
  if (data + nh_off > data_end)
    goto pass;  

  h_proto = eth->h_proto;
  if (h_proto != __constant_htons(ETH_P_IP))
      goto pass;

  ipproto = parse_ipv4(data, nh_off, data_end, &src_ip);
  unsigned int intIP = 184554668; // int value of 172.20.0.11
  if (src_ip == intIP)
      return XDP_DROP;

pass:
  return XDP_PASS;
}

