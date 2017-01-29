package main

import (
	"io"
	"net"
	"os/exec"
	"github.com/songgao/water"
)

func CreateTUN(ifaceName string, cidr net.IP) io.ReadCloser {
	iface := tun(ifaceName)
	devUp(ifaceName, cidr)
	return iface
}


func tun(ifaceName string) *water.Interface {
	Log.Debugf("Creating TUN device %s", ifaceName)
    iface, err := water.NewTUN(ifaceName)
    if err != nil {
		Log.Fatal("Failed to create TUN device:", err)
    }
	return iface
}

func devUp(ifaceName string, cidr net.IP) {
	Log.Debugf("Setting %s CIDR to %s/64", ifaceName, cidr)
	// TODO: Would be great to find a better way to configure the TUN device
	err := exec.Command("ip", "addr", "add", cidr.String() + "/64", "dev", ifaceName).Run()
	if err != nil {
		Log.Fatal("Failed to set CIDR for interface:", err)
	}

	Log.Debugf("Bringing up device %s", ifaceName)
	err = exec.Command("ip", "link", "set", ifaceName, "up").Run()
	if err != nil {
		Log.Fatal("Failed to bring up interface:", err)
	}
}
