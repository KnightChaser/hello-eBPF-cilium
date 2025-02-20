package main

import (
	"C"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

func main() {
	// Remove resource limits for kernels < 5.11.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("Failed to remove rlimit: %v", err)
	}

	// Load the compiled eBPF ELF and load it into the kernel.
	var objs packetCounterObjects
	if err := loadPacketCounterObjects(&objs, nil); err != nil {
		log.Fatalf("Failed to load packetCounter: %v", err)
	}
	defer objs.Close()

	// Discover the interface index for the interface we want to attach to.
	ifname := "enp3s0"
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		log.Fatalf("Failed to get interface %q: %v", ifname, err)
	}

	// Attach packetCounter to the network interface.
	link, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.CountPackets,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatalf("Failed to attach XDP program to %q: %v", ifname, err)
	}
	defer link.Close()

	log.Printf("Counting packets on the interface %s. Press Ctrl+C to stop.", ifname)

	// Periodically fetch the packet counter from the counter.
	// Exit the program wh it was interrupted.
	tick := time.Tick(time.Second)
	stop := make(chan os.Signal, 5)
	signal.Notify(stop, os.Interrupt)
	for {
		select {
		case <-tick:
			// For every tick, fetch the packet counter.
			key := uint32(0)
			var packetCount uint64
			if err := objs.PacketCountMap.Lookup(&key, &packetCount); err != nil {
				log.Fatalf("Failed to lookup packet count: %v", err)
			} else {
				log.Printf("Packet count: %d", packetCount)
			}
		case <-stop:
			log.Println("Detaching XDP program and exiting...")
			return
		}
	}
}
