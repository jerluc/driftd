package lib

import (
	"io"
	"net"
	"os/exec"
	"github.com/songgao/water"
)

// Creates a named TUN device and configures
// the device to handle packets with the given
// 64-bit CIDR prefix
func CreateTUN(ifaceName string, cidr net.IP) io.ReadCloser {
	iface := tun(ifaceName)
	devUp(ifaceName, cidr)
	return iface
}

// Creates the TUN interface with the given name
func tun(ifaceName string) *water.Interface {
	Log.Debugf("Creating TUN device %s", ifaceName)
	config := water.Config{
        DeviceType: water.TUN,
    }
    config.Name = ifaceName
    iface, err := water.New(config)
    if err != nil {
		Log.Fatal("Failed to create TUN device:", err)
    }
	return iface
}

// Configures the TUN device to handle the given
// 64-bit CIDR prefix and then finally brings the
// network device up
func devUp(ifaceName string, cidr net.IP) {
	Log.Debugf("Setting %s CIDR to %s/64", ifaceName, cidr)
	ipCmd, ipCmdLookupErr := exec.LookPath("ip")
	if ipCmdLookupErr != nil {
		Log.Fatal("Failed to find `ip` command on your $PATH:", ipCmdLookupErr)
	}
	// TODO: Would be great to find a better way to configure the TUN device
	err := exec.Command(ipCmd, "addr", "add", cidr.String() + "/64", "dev", ifaceName).Run()
	if err != nil {
		Log.Fatal("Failed to set CIDR for interface:", err)
	}

	Log.Debugf("Bringing up device %s", ifaceName)
	err = exec.Command(ipCmd, "link", "set", ifaceName, "up").Run()
	if err != nil {
		Log.Fatal("Failed to bring up interface:", err)
	}
}
