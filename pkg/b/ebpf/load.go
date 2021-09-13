package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/iovisor/gobpf/elf"
)

func main() {
	load()
}

func load() {
	bina := elf.NewModule("xdp.o")
	if bina == nil {
		fmt.Fprintf(os.Stderr, "eBPF program not found")
		os.Exit(1)
	}
	var secParams = map[string]elf.SectionParams{}
	if err := bina.Load(secParams); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	device := "eth0"
	err := bina.AttachXDP(device, "xdp/bina")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to attach xdp prog: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := bina.RemoveXDP(device); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove XDP from %s: %v\n", device, err)
		}
	}()

	fmt.Println("Dropping packets, hit CTRL+C to stop")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

}
