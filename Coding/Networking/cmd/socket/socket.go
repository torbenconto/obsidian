package socket

import (
	"encoding/binary"
	"errors"
	"math"
	"net"
	"syscall"

	"github.com/josharian/native"
)

func htons(i int) (uint16, error) {
	if i < 0 || i > math.MaxUint16 {
		return 0, errors.New("packet: protocol value out of range")
	}

	var b [2]byte
	binary.BigEndian.PutUint16(b[:], uint16(i))

	return native.Endian.Uint16(b[:]), nil
}

type Socket struct {
	ifi      *net.Interface
	protocol int
	fd       int
}

func (s *Socket) Ifi() *net.Interface {
	return s.ifi
}

func NewSocket(ifi *net.Interface, protocol int) Socket {
	return Socket{
		ifi:      ifi,
		protocol: protocol,
		fd:       0,
	}
}

func (s *Socket) Listen() error {
	proto, err := htons(s.protocol)
	if err != nil {
		return err
	}

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, s.protocol)
	if err != nil {
		return err
	}

	s.fd = fd

	addr := &syscall.SockaddrLinklayer{
		Ifindex:  s.ifi.Index,
		Protocol: proto,
	}

	err = syscall.Bind(fd, addr)

	return err
}

func (s *Socket) Write(data []byte) error {
	_, err := syscall.Write(s.fd, data)
	return err
}

func (s *Socket) Read(data []byte) (int, error) {
	n, _, err := syscall.Recvfrom(s.fd, data, 0)

	return n, err
}

func (s *Socket) Close() error {
	return syscall.Close(s.fd)
}
