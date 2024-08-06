package main

import (
	"fmt"
	"net"

	"github.com/torbenconto/obsidian/Coding/Networking/cmd/arp"
)

func main() {
	a := arp.NewARP()
	f, err := a.Resolve(net.IPv4(0x0A, 0, 0, 1))

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(net.HardwareAddr(f.SourceMAC[:]))
}
