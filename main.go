package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	rift "github.com/jerluc/riftd/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	DefaultLogLevel = "INFO"
	DefaultInterface = "rift0"
	DefaultDevice = "/dev/ttyUSB0"
	DefaultLocalIP = "2001:412:abcd:1::"
)

func main() {
	riftd := kingpin.New("riftd", "Rift protocol daemon")
	runCmd := riftd.Command("run", "Starts the Rift protocol daemon")
	versionCmd := riftd.Command("version", "Displays riftd version")
	logLevel := runCmd.Flag("logging", "Log level").
					Default(DefaultLogLevel).
					Enum("DEBUG", "INFO", "NOTICE", "WARNING", "ERROR", "CRITICAL")
	ifaceName := runCmd.Flag("iface", "Network interface name").
					Default(DefaultInterface).
					String()
	devName := runCmd.Flag("dev", "Serial device name").
					Default(DefaultDevice).
					String()
	cidr := runCmd.Flag("cidr", "IPv6 64-bit prefix").
					Default(DefaultLocalIP).
					IP()
	cfgCmd := riftd.Command("configure", "Configures a new device for Rift")
	newDevName := cfgCmd.Flag("dev", "Serial device name").
					Default(DefaultDevice).
					String()

	cmd, parseErr := riftd.Parse(os.Args[1:])
	if parseErr != nil {
		fmt.Println("riftd:", parseErr)
		fmt.Println("Run \"riftd help [cmd]\" for help")
		os.Exit(1)
	}

	switch cmd {
	case runCmd.FullCommand():
		rift.InitLogging(*logLevel)

		xchg := rift.NewExchange(rift.ExchangeConfig{
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
		rift.ConfigureDevice(*newDevName)
	}
}
