# hello-eBPF-cilium

A simple example of using **eBPF (extended Berkeley Packet Filter)** with Cilium (powered by Go). Remade from the [official tutorial](https://ebpf-go.dev/guides/getting-started/#the-go-application), fixing some incompatibilities (*explanations that don't work*) with Ubuntu 22.04 LTS environment. It uses **XDP(eXpress Data Path)** to simply count the packet flown through the network interface (Refer to `main.py`) and show the number of accumulative captured packets every second.

### Getting started
1. Prepare `vmlinux.h` header file at the same directory of this project. For example, you can generate the header file like below.
```sh
bpftool btf dump file /sys/kernel/btf/vmlinux format c > vmlinux.h
```
2. Hit `go generate` to build the Golang objects based on the C eBPF code. After that, you can use eBPF stuff originally written in the C code(`packetCounter.c`) file like a Go object. (`packetcounter_bpfeb.go/o` and `packetcounter_bpfel.go/o` will be automatically created!)
3. Hit `go build .`. Now, `main.go` is compiled, and the subsequent binary file will be dropped. 
4. Execute and observe how the program is going on!
```
knightchaser@passc0de:~/ghRepo/hello-eBPF-cilium$ sudo ./example_cilium_xdp 
2025/02/20 15:15:04 Counting packets on the interface enp3s0. Press Ctrl+C to stop.
2025/02/20 15:15:05 Packet count: 75
2025/02/20 15:15:06 Packet count: 138
2025/02/20 15:15:07 Packet count: 267
2025/02/20 15:15:08 Packet count: 276
2025/02/20 15:15:09 Packet count: 290
^C2025/02/20 15:15:10 Detaching XDP program and exiting...
```

