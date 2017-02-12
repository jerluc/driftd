package lib

import (
	"bytes"
	"encoding/binary"
	"net"
	"syscall"
)

// Computes a UDP checksum per RFC 768
func computeUDPChecksum(buf []byte) uint16 {
	sum := uint32(0)

	for ; len(buf) >= 2; buf = buf[2:] {
		sum += uint32(buf[0])<<8 | uint32(buf[1])
	}
	if len(buf) > 0 {
		sum += uint32(buf[0]) << 8
	}
	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}
	csum := ^uint16(sum)
	
	// From RFC 768:
	// If the computed checksum is zero, it is transmitted as all ones (the
	// equivalent in one's complement arithmetic). An all zero transmitted
	// checksum value means that the transmitter generated no checksum (for
	// debugging or for higher level protocols that don't care).
	if csum == 0 {
		csum = 0xffff
	}
	return csum
}

// Constructs a raw UDP+IPv6 packet
func BuildRawPacket(srcAddr *net.UDPAddr, dstAddr *net.UDPAddr, payload []byte) []byte {
	buf := bytes.NewBuffer([]byte{})
	// Source address
	binary.Write(buf, binary.BigEndian, srcAddr.IP)
	// Destination address
	binary.Write(buf, binary.BigEndian, dstAddr.IP)
	// Length (extra 8 bytes for UDP header)
	binary.Write(buf, binary.BigEndian, uint16(len(payload) + 8))
	// 12-byte zero padding
	binary.Write(buf, binary.BigEndian, uint8(0))
	binary.Write(buf, binary.BigEndian, uint8(0))
	binary.Write(buf, binary.BigEndian, uint8(0))
	// Next header (17 = UDP!)
	binary.Write(buf, binary.BigEndian, byte(17))
	// Source port
	binary.Write(buf, binary.BigEndian, uint16(srcAddr.Port))
	// Destination port
	binary.Write(buf, binary.BigEndian, uint16(dstAddr.Port))
	// Length (extra 8 bytes for UDP header)
	binary.Write(buf, binary.BigEndian, uint16(len(payload) + 8))
	// Checksum placeholder
	binary.Write(buf, binary.BigEndian, uint16(0))
	// Data payload
	binary.Write(buf, binary.BigEndian, payload)

	checksum := computeUDPChecksum(buf.Bytes())

	buf.Reset()

	// Version, traffic class, flow label
	binary.Write(buf, binary.BigEndian, byte(96))
	binary.Write(buf, binary.BigEndian, byte(0))
	binary.Write(buf, binary.BigEndian, byte(0))
	binary.Write(buf, binary.BigEndian, byte(0))
	// Length (extra 8 bytes for UDP header)
	binary.Write(buf, binary.BigEndian, uint16(len(payload) + 8))
	// Next header (17 = UDP!)
	binary.Write(buf, binary.BigEndian, byte(17))
	// Hop limit
	binary.Write(buf, binary.BigEndian, byte(64))
	// Source address
	binary.Write(buf, binary.BigEndian, srcAddr.IP)
	// Destination address
	binary.Write(buf, binary.BigEndian, dstAddr.IP)

	// Source port
	binary.Write(buf, binary.BigEndian, uint16(srcAddr.Port))
	// Destination port
	binary.Write(buf, binary.BigEndian, uint16(dstAddr.Port))
	// Length (extra 8 bytes for UDP header)
	binary.Write(buf, binary.BigEndian, uint16(len(payload) + 8))
	// Checksum placeholder
	binary.Write(buf, binary.BigEndian, uint16(checksum))
	// Data payload
	binary.Write(buf, binary.BigEndian, payload)

	return buf.Bytes()
}

// Sends a raw packet encoded as a byte array to
// the networking stack. Note that this may be
// platform dependent!
func SendRawPacket(packet []byte) error {
	fd, err := syscall.Socket(syscall.AF_INET6, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	return syscall.Sendto(fd, packet, 0, &syscall.SockaddrInet6{})
}
