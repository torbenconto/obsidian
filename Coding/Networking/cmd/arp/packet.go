package arp

import (
	"bytes"
	"net"
)

type EthernetHeader struct {
	TargetMAC [6]byte
	SourceMAC [6]byte
	EtherType [2]byte
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

func NewARPFrame(op [2]byte, SourceMAC, TargetMAC, NetMAC net.HardwareAddr, SourceIP, TargetIP net.IP) *ARPFrame {
	return &ARPFrame{
		EthernetHeader: EthernetHeader{
			TargetMAC: [6]byte(NetMAC),
			SourceMAC: [6]byte(SourceMAC),
			EtherType: [2]byte{0x08, 0x06},
		},
		HardwareType: [2]byte{0x00, 0x01},
		ProtocolType: [2]byte{0x08, 0x00},
		HardwareSize: 6,
		ProtocolSize: 4,
		Opcode:       op,
		SourceMAC:    [6]byte(SourceMAC),
		TargetMAC:    [6]byte(TargetMAC),
		SourceIP:     [4]byte(SourceIP.To4()),
		TargetIP:     [4]byte(TargetIP.To4()),
	}
}

func (a *ARPFrame) ToBytes() []byte {
	buf := new(bytes.Buffer)

	// Ethernet header
	buf.Write(a.EthernetHeader.TargetMAC[:])
	buf.Write(a.EthernetHeader.SourceMAC[:])
	buf.Write(a.EthernetHeader.EtherType[:])

	// ARP header
	buf.Write(a.HardwareType[:])
	buf.Write(a.ProtocolType[:])
	buf.WriteByte(a.HardwareSize)
	buf.WriteByte(a.ProtocolSize)
	buf.Write(a.Opcode[:])
	buf.Write(a.SourceMAC[:])
	buf.Write(a.SourceIP[:])
	buf.Write(a.TargetMAC[:])
	buf.Write(a.TargetIP[:])

	return buf.Bytes()
}

func ParseFrame(frame []byte) *ARPFrame {
	if len(frame) < 42 {
		return nil // Not enough data to parse
	}

	a := &ARPFrame{}
	offset := 0

	// Ethernet header
	copy(a.EthernetHeader.TargetMAC[:], frame[offset:offset+6])
	offset += 6
	copy(a.EthernetHeader.SourceMAC[:], frame[offset:offset+6])
	offset += 6
	copy(a.EthernetHeader.EtherType[:], frame[offset:offset+2])
	offset += 2

	// ARP header
	copy(a.HardwareType[:], frame[offset:offset+2])
	offset += 2
	copy(a.ProtocolType[:], frame[offset:offset+2])
	offset += 2
	a.HardwareSize = frame[offset]
	offset += 1
	a.ProtocolSize = frame[offset]
	offset += 1
	copy(a.Opcode[:], frame[offset:offset+2])
	offset += 2
	copy(a.SourceMAC[:], frame[offset:offset+6])
	offset += 6
	copy(a.SourceIP[:], frame[offset:offset+4])
	offset += 4
	copy(a.TargetMAC[:], frame[offset:offset+6])
	offset += 6
	copy(a.TargetIP[:], frame[offset:offset+4])

	return a
}
