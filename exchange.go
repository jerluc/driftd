package main

import (
	"io"
	"net"
	"github.com/jerluc/gobee"
	"golang.org/x/net/ipv6"
	"github.com/jerluc/serial"
)

type ExchangeConfig struct{
	DeviceName string
	InterfaceName string
	CIDR net.IP
}

type Exchange struct{
	cfg ExchangeConfig
	iface io.ReadCloser
	serialPort io.ReadWriteCloser
	inbox <-chan gobee.Frame
	outbox chan<- gobee.Frame
}

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

func (x *Exchange) Shutdown() {
	Log.Info("Shutting down exchange")
	x.iface.Close()
	x.serialPort.Close()
	close(x.outbox)
}

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

func (x *Exchange) externalIP(mac []byte) []byte {
	return append(x.cfg.CIDR[:8], mac...)
}

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

func isOutgoingPacket(header *ipv6.Header) bool {
	for i, b := range header.Dst[:4] {
		if b != header.Src[i] {
			return false
		}
	}
	return true
}

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

		if isOutgoingPacket(header) {
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
