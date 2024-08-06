package arp

import (
	"net"
	"syscall"

	"github.com/torbenconto/obsidian/Coding/Networking/cmd/socket"
)

type ARP struct {
	sock *socket.Socket
}

func NewARP() ARP {
	var iface, _, err = getIface()
	if err != nil {
		panic(err)
	}

	s := socket.NewSocket(&iface, syscall.ETH_P_ARP)
	return ARP{
		sock: &s,
	}
}

func (a *ARP) Request(ip net.IP) error {
	var _, localIP, err = getIface()
	if err != nil {
		return err
	}

	broadcastMacBytes := [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	emptyMacBytes := [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	frame := NewARPFrame([2]byte{0x00, 0x01}, a.sock.Ifi().HardwareAddr, net.HardwareAddr(emptyMacBytes[:]), net.HardwareAddr(broadcastMacBytes[:]), localIP[:], ip)

	a.sock.Listen()

	return a.sock.Write(frame.ToBytes())
}

// Read reads a single ARP frame with a return opcode
func (a *ARP) Read() (*ARPFrame, error) {
	buf := make([]byte, 128)
	for {
		n, err := a.sock.Read(buf)
		if err != nil {
			return &ARPFrame{}, err
		}

		f := ParseFrame(buf[:n])

		if f.Opcode == [2]byte{0x00, 0x02} {
			return f, nil
		}

		continue
	}
}

func (a *ARP) Resolve(ip net.IP) (*ARPFrame, error) {
	err := a.Request(ip)
	if err != nil {
		return &ARPFrame{}, err
	}

	f, err := a.Read()
	return f, err
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
