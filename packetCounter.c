//go:build ignore

#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

struct {
  __uint(type, BPF_MAP_TYPE_ARRAY);
  __uint(max_entries, 1);
  __type(key, __u32);
  __type(value, __u64);
} packet_count_map SEC(".maps");

SEC("xdp")
int count_packets(struct xdp_md *ctx) {
  __u32 key = 0;
  __u64 *value = bpf_map_lookup_elem(&packet_count_map, &key);
  if (value) {
    __sync_fetch_and_add(value, 1);
  }
  return XDP_PASS;
}

char _license[] SEC("license") = "GPL";
