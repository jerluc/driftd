package main

import (
	"os"
	"os/signal"
	"syscall"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	DefaultLogLevel = "INFO"
	DefaultInterface = "rift0"
	DefaultDevice = "/dev/ttyUSB0"
	DefaultLocalIP = "2001:412:abcd:1::"
)

var (
	riftd = kingpin.New("riftd", "Rift protocol daemon")
	runCmd = riftd.Command("run", "Starts the Rift protocol daemon")
	versionCmd = riftd.Command("version", "Displays riftd version")
	logLevel = runCmd.Flag("logging", "Log level").Default(DefaultLogLevel).String()
	ifaceName = runCmd.Flag("iface", "Network interface name").Default(DefaultInterface).String()
	devName = runCmd.Flag("dev", "Serial device name").Default(DefaultDevice).ExistingFile()
	cidr = runCmd.Flag("cidr", "IPv6 64-bit prefix").Default(DefaultLocalIP).IP()
	cfgCmd = riftd.Command("configure", "Configures a new device for Rift")
	newDevName = cfgCmd.Flag("dev", "Serial device name").Default(DefaultDevice).ExistingFile()
)

func main() {
	switch kingpin.MustParse(riftd.Parse(os.Args[1:])) {
	case runCmd.FullCommand():
		InitLogging(*logLevel)

		xchg := NewExchange(ExchangeConfig{
			DeviceName: *devName,
			InterfaceName: *ifaceName,
			CIDR: *cidr,
		})

		// Install signal handler to gracefully shutdown
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
		go func() {
			<-c
			xchg.Shutdown()
			os.Exit(0)
		}()

		xchg.Start()
	case versionCmd.FullCommand():
		PrintVersionInfo()
	case cfgCmd.FullCommand():
		ConfigureDevice(*newDevName)
	}
}
