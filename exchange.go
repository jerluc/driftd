package main

import (
	"io"
	"net"
	"github.com/jerluc/gobee"
	"golang.org/x/net/ipv6"
	"github.com/jerluc/serial"
)

// Configuration for the Rift packet exchange
type ExchangeConfig struct{
	// The serial device for the exchange to use
	// for all TX/RX operations
	DeviceName string
	// The name of the TUN interface to create
	InterfaceName string
	// The IPv6 CIDR block for routing. Note that
	// this should be a 64-bit prefix, as the
	// remaining 64 bits will be consumed by the
	// serial device's MAC address
	CIDR net.IP
}

// The Rift packet exchange performs all vital
// protocol functions including routing, packet
// encapsulation and stripping, etc.
type Exchange struct{
	cfg ExchangeConfig
	iface io.ReadCloser
	serialPort io.ReadWriteCloser
	inbox <-chan gobee.Frame
	outbox chan<- gobee.Frame
}

// Creates a new Rift packet exchange
func NewExchange(cfg ExchangeConfig) *Exchange {
	iface := CreateTUN(cfg.InterfaceName, cfg.CIDR)
	serialPort := openSerialPort(cfg.DeviceName)
	mailbox := gobee.NewMailbox(serialPort)
	inbox, outbox := mailbox.Inbox(), mailbox.Outbox()
	return &Exchange{
		cfg,
		iface,
		serialPort,
		inbox,
		outbox,
	}
}

// Gracefully shuts down the Rift exchange
func (x *Exchange) Shutdown() {
	Log.Info("Shutting down exchange")
	x.iface.Close()
	x.serialPort.Close()
	close(x.outbox)
}

// Boots up the Rift exchange. Note that this
// function will block indefinitely until shutdown
func (x *Exchange) Start() {
	Log.Info("Booting up exchange")
	go x.incoming()
	go x.outgoing()
	Log.Info("Exchange active")
	select{}
}

func openSerialPort(devName string) io.ReadWriteCloser {
	Log.Debug("Opening serial device at", devName)
	// TODO: Externalize device baud rate?
	serialCfg := &serial.Config{Name: devName, Baud: 9600}
	serialPort, openErr := serial.OpenPort(serialCfg)
	if openErr != nil {
		Log.Fatal("Failed to open serial device:", openErr)
	}
	return serialPort
}

// Creates the "external" IPv6 address for a given
// MAC address
func (x *Exchange) externalIP(mac []byte) []byte {
	return append(x.cfg.CIDR[:8], mac...)
}

// The incoming packet loop continually processes
// RX frames comming from the attached serial device.
// For each received frame, a new raw IPv6 packet is
// constructed by repacking the incoming frame, and
// then eventually forwarded on to the host machine's
// network stack where it is delivered to vanilla
// UDP+IPv6 sockets.
func (x *Exchange) incoming() {
	Log.Debug("Started watching for incoming packets")
	for {
		packet := (<-x.inbox).(*gobee.Rx64Frame)
		Log.Debug("Incoming packet:", packet)
		srcPort := gobee.BytesToUint16(packet.Data[:2])
		dstPort := gobee.BytesToUint16(packet.Data[2:4])
		payload := trimmed(packet.Data[4:])
		srcIP := x.externalIP(packet.Source)

		srcAddr := &net.UDPAddr{
			IP: srcIP,
			Port: int(srcPort),
		}

		dstAddr := &net.UDPAddr{
			// TODO: Is using the loopback address the only way we can do this?
			IP: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Port: int(dstPort),
		}

		Log.Debug("src =", srcAddr, "dst =", dstAddr, "payload =", string(payload))

		rawPkt := BuildRawPacket(srcAddr, dstAddr, payload)
		err := SendRawPacket(rawPkt)
		if err != nil {
			Log.Fatal("Failed to forward packet to destination:", err)
		}
	}
}

// Determines whether or not the IPv6 packet is
// destined for the exchange
func (x *Exchange) isOutgoingPacket(header *ipv6.Header) bool {
	for i, b := range header.Dst[:4] {
		if b != x.cfg.CIDR[i] {
			return false
		}
	}
	return true
}

// The outgoing packet loop continually reads IPv6
// packets coming in from the TUN device itself. If
// the packet is destined for the exchange CIDR, the
// packet is repacked as a TX frame and then sent to
// the underyling serial device
func (x *Exchange) outgoing() {
	Log.Debug("Started watching for outgoing packets")
    for {
		// TODO: Is this okay?
		packet := make([]byte, 127)
        n, err := x.iface.Read(packet)
		if err != nil {
			Log.Fatal("Failed to read from interface:", err)
        }
		Log.Debug("Received packet:", packet[:n])
		header, err := ipv6.ParseHeader(packet[:n])
		if err != nil {
			Log.Error("Failed to parse IPv6 header:", err)
		}

		if x.isOutgoingPacket(header) {
			// IPv6 payload starts after 40 byte fixed-size header
			srcPort := packet[40:42]
			dstPort := packet[42:44]
			dst := header.Dst[8:]
			// TODO: Detect broadcast addressing mode
			payload := packet[48:]
			packet := &gobee.Tx64Frame{
				ID: 0x00,
				Destination: dst,
				Options: 0x00,
				Data: gobee.PackBytes(srcPort, dstPort, payload),
			}
			Log.Debug("Outgoing packet:", packet)
			x.outbox <- packet
		}
    }
}
