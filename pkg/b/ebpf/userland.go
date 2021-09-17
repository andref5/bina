package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/iovisor/gobpf/elf"
)

const Device = "eth0"

func main() {
	bina := Load()

	defer func() {
		Unload(bina)
	}()

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

func fatal_fury(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a)
	fmt.Fprintf(os.Stderr, msg)
	os.Exit(1)
}
