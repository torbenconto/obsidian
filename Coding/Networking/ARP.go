package main

import (
	"fmt"
	"net"
	"syscall"
)

var NETMAC = [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
var EMPTYMAC = [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

type EthernetHeader struct {
	DestinationMAC [6]byte
	SourceMAC      [6]byte
	EtherType      [2]byte
}

type ARPFrame struct {
	EthernetHeader EthernetHeader
	HardwareType   [2]byte
	ProtocolType   [2]byte
	HardwareSize   byte
	ProtocolSize   byte
	Opcode         [2]byte
	SourceMAC      [6]byte
	SourceIP       [4]byte
	TargetMAC      [6]byte
	TargetIP       [4]byte
}

func (a *ARPFrame) Send(iface *net.Interface) {
	// Create Socket                                                    0x806
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ARP)
	if err != nil {
		panic(err)
	}

	defer syscall.Close(fd)

	addr := &syscall.SockaddrLinklayer{
		Ifindex:  iface.Index,
		Protocol: syscall.ETH_P_ARP,
	}

	err = syscall.Bind(fd, addr)
	if err != nil {
		panic(err)
	}

	fmt.Println(a.ToBytes())

	_, err = syscall.Write(fd, a.ToBytes())
	if err != nil {
		panic(err)
	}
}

func (a *ARPFrame) ToBytes() []byte {
	buf := make([]byte, 42) // Ethernet header (14 bytes) + ARP header (28 bytes)
	offset := 0

	// Ethernet header
	copy(buf[offset:], a.EthernetHeader.DestinationMAC[:])
	offset += 6
	copy(buf[offset:], a.EthernetHeader.SourceMAC[:])
	offset += 6
	copy(buf[offset:], a.EthernetHeader.EtherType[:])
	offset += 2

	// ARP header
	copy(buf[offset:], a.HardwareType[:])
	offset += 2
	copy(buf[offset:], a.ProtocolType[:])
	offset += 2
	buf[offset] = a.HardwareSize
	offset += 1
	buf[offset] = a.ProtocolSize
	offset += 1
	copy(buf[offset:], a.Opcode[:])
	offset += 2
	copy(buf[offset:], a.SourceMAC[:])
	offset += 6
	copy(buf[offset:], a.SourceIP[:])
	offset += 4
	copy(buf[offset:], a.TargetMAC[:])
	offset += 6
	copy(buf[offset:], a.TargetIP[:])

	return buf
}

func getIface() (net.Interface, [4]byte, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return net.Interface{}, [4]byte{}, err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			if ipNet.IP.To4() != nil && !ipNet.IP.IsLoopback() {
				return iface, [4]byte(ipNet.IP.To4()), nil
			}
		}
	}

	return net.Interface{}, [4]byte{}, err
}

func main() {

	var iface, localIP, err = getIface()
	if err != nil {
		panic(err)
	}

	var SRCMAC = [6]byte(iface.HardwareAddr)
	var SRCIP = localIP

	a := ARPFrame{
		EthernetHeader: EthernetHeader{
			DestinationMAC: NETMAC,
			SourceMAC:      SRCMAC,
			EtherType:      [2]byte{0x08, 0x06},
		},

		HardwareType: [2]byte{0x00, 0x01},
		ProtocolType: [2]byte{0x08, 0x00},
		HardwareSize: 6,
		ProtocolSize: 4,
		Opcode:       [2]byte{0x00, 0x01},
		SourceMAC:    SRCMAC,
		TargetMAC:    EMPTYMAC,
		SourceIP:     SRCIP,
		TargetIP:     [4]byte{0xA, 0x0, 0x0, 0x1},
	}

	a.Send(&iface)

	fmt.Println(a)

}
