package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	drift "github.com/jerluc/driftd/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	DefaultLogLevel = "INFO"
	DefaultInterface = "drift0"
	DefaultDevice = "/dev/ttyUSB0"
	DefaultLocalIP = "2001:412:abcd:1::"
)

func main() {
	driftd := kingpin.New("driftd", "Drift protocol daemon")
	runCmd := driftd.Command("run", "Starts the Drift protocol daemon")
	versionCmd := driftd.Command("version", "Displays driftd version")
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
	cfgCmd := driftd.Command("configure", "Configures a new device for Drift")
	newDevName := cfgCmd.Flag("dev", "Serial device name").
					Default(DefaultDevice).
					String()

	cmd, parseErr := driftd.Parse(os.Args[1:])
	if parseErr != nil {
		fmt.Println("driftd:", parseErr)
		fmt.Println("Run \"driftd help [cmd]\" for help")
		os.Exit(1)
	}

	switch cmd {
	case runCmd.FullCommand():
		drift.InitLogging(*logLevel)

		xchg := drift.NewExchange(drift.ExchangeConfig{
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
		drift.ConfigureDevice(*newDevName)
	}
}
