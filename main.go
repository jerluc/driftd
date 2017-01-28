package main

import (
    "log"
	"net"
	"os"
	"os/exec"
	"github.com/jerluc/gobee"
	"github.com/jerluc/serial"
	"github.com/songgao/water"
	"golang.org/x/net/ipv6"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	DefaultInterface = "rift0"
	DefaultLocalIP = "2001:412:abcd:1::"
)

var (
	riftd = kingpin.New("riftd", "Rift protocol daemon")
	ifaceName = riftd.Flag("iface", "Network interface name").Default(DefaultInterface).String()
	devname = riftd.Flag("dev", "Serial device name").OverrideDefaultFromEnvar("DEVNAME").Required().String()
	cidr = riftd.Flag("cidr", "IPv6 64-bit prefix").Default(DefaultLocalIP).IP()
)

func packBytes(bas... []byte) []byte {
	var packed []byte
	for _, ba := range bas {
		packed = append(packed, ba...)
	}
	return packed
}

func trimmed(packet []byte) []byte {
	for i, b := range packet {
		if b == 0x00 {
			return packet[:i]
		}
	}
	return packet
}

func incoming(cidr net.IP, inbox <-chan gobee.Frame, iface *water.Interface) {
	for {
		packet := (<-inbox).(*gobee.Rx64Frame)
		log.Println("Incoming packet:", packet)
		srcPort := gobee.BytesToUint16(packet.Data[:2])
		dstPort := gobee.BytesToUint16(packet.Data[2:4])
		payload := trimmed(packet.Data[4:])
		srcIP := append([]byte(cidr[:8]), packet.Source...)

		srcAddr := &net.UDPAddr{
			IP: srcIP,
			Port: int(srcPort),
		}

		dstAddr := &net.UDPAddr{
			// Loopback address
			IP: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Port: int(dstPort),
		}

		log.Println("src =", srcAddr, "dst =", dstAddr, "payload =", string(payload))

		rawPkt := BuildRawPacket(srcAddr, dstAddr, payload)
		err := SendRawPacket(rawPkt)
		if err != nil {
			panic(err)
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

func tun(ifaceName string) *water.Interface {
    iface, err := water.NewTUN(ifaceName)
    if err != nil {
        panic(err)
    }
	return iface
}

func devUp(ifaceName string, cidr net.IP) {
	err := exec.Command("ip", "addr", "add", cidr.String() + "/64", "dev", ifaceName).Run()
	if err != nil {
		panic(err)
	}

	err = exec.Command("ip", "link", "set", ifaceName, "up").Run()
	if err != nil {
		panic(err)
	}
}

func outgoing(iface *water.Interface, outbox chan<- gobee.Frame) {
    for {
		// TODO: Is this okay?
		packet := make([]byte, 127)
        n, err := iface.Read(packet)
		if err != nil {
            panic(err)
        }
		log.Println("Received IPv6 packet: % x", packet[:n])
		header, err := ipv6.ParseHeader(packet[:n])
		if err != nil {
			panic(err)
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
			log.Println("Outgoing packet:", packet)
			outbox <- packet
		}
    }
}

func main() {
	kingpin.MustParse(riftd.Parse(os.Args[1:]))

	serialCfg := &serial.Config{Name: *devname, Baud: 9600}
	serialPort, openErr := serial.OpenPort(serialCfg)
	if openErr != nil {
		panic(openErr)
	}

	iface := tun(*ifaceName)
	devUp(*ifaceName, *cidr)

	mailbox := gobee.NewMailbox(serialPort)
	inbox, outbox := mailbox.Inbox(), mailbox.Outbox()

	go incoming(*cidr, inbox, iface)
	go outgoing(iface, outbox)

	defer func() {
		iface.Close()
		close(outbox)
		serialPort.Close()
	}()
	select{}
}
