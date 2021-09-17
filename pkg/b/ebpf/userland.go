package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"unsafe"

	"github.com/iovisor/gobpf/elf"
)

const Device = "eth0"

func main() {
	flag.Var(&denyList, "deny", "IPv4 to DROP traffic from, repeatable")
	flag.Parse()
	if len(denyList) == 0 {
		fatal_fury("at least one IPv4 address to DROP required (-deny)")
	}

	bina := Load()

	defer func() {
		Unload(bina)
	}()

	fmt.Println("Populate eBPF map with int IPv4 addresses...")
	for i, ip := range denyList {
		UpdateMap(bina, i, ip)
	}

	fmt.Println("Dropping packets, hit CTRL+C to stop")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func Load() *elf.Module {
	bina := elf.NewModule("xdp.o")
	if bina == nil {
		fatal_fury("eBPF program not found")
	}

	var secParams = map[string]elf.SectionParams{}

	if err := bina.Load(secParams); err != nil {
		fatal_fury(err.Error())
	}

	err := bina.AttachXDP(Device, "xdp/bina")
	if err != nil {
		fatal_fury("Failed to attach xdp prog: %v\n", err)
	}

	return bina
}

func Unload(bina *elf.Module) {
	if err := bina.RemoveXDP(Device); err != nil {
		fatal_fury("Failed to remove XDP from %s: %v\n", Device, err)
	}
}

func UpdateMap(bina *elf.Module, index int, ipInt uint32) {
	mp := bina.Map("deny_ips")
	if mp == nil {
		fatal_fury("unable to find dummy_hash map")
	}
	var x uint32
	x = 1
	if err := bina.UpdateElement(mp, unsafe.Pointer(&ipInt), unsafe.Pointer(&x), 0); err != nil {
		fatal_fury("failed trying to update an element: %v\n", err.Error())
	}
	fmt.Printf("\t ADD: %v\n", ipInt)
}

type denyIPs []uint32

var denyList denyIPs

func (d *denyIPs) String() string {
	return fmt.Sprintf("%v", *d)
}

func (d *denyIPs) Set(v string) error {
	ip := net.ParseIP(v)
	if ip == nil {
		return fmt.Errorf("Invalid IP address [%v]", v)
	}
	intIP := binary.LittleEndian.Uint32(ip[12:16])
	*d = append(*d, intIP)
	return nil
}

func fatal_fury(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a)
	fmt.Fprintf(os.Stderr, msg)
	os.Exit(1)
}
