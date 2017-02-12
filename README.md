# riftd

[![](https://api.travis-ci.org/jerluc/riftd.svg?branch=master)](https://travis-ci.org/jerluc/riftd)
[![GoDoc](https://godoc.org/github.com/jerluc/riftd?status.svg)](https://godoc.org/github.com/jerluc/riftd)

The Rift protocol daemon

## What is Rift?

Rift is a simple 802.15.4 wireless protocol for enabling peer-to-peer IPv6 communication without the
need for any intermediary infrastructure such as routers or gateways. Using this new protocol, we
hope to see a new future for software developers to more easily use direct wireless communications
for applications that would benefit from the increased security, privacy, and simplicity of Rift.

## What is riftd?

riftd is a small daemon process which implements the Rift protocol as an IPv6 TUN device. When
running, riftd routes incoming 802.15.4 packets from nearby peers to local UDP sockets and outgoing
UDP packets in a specific CIDR block to their corresponding remote peer.

As a simple example, given:

  * Peer 0 (MAC `a:b:c:d`) running riftd with CIDR `2001:412:abcd:1::/64`, and a local UDP server
    bound to port 8000
  * Peer 1 (MAC `c:d:e:f`) running riftd with CIDR `2001:412:abcd:1::/64`, and a local UDP client

When Peer 1 sends a UDP packet to `[2001:412:abcd:1:a:b:c:d]:8000`, riftd will route the packet over
the 802.15.4 wireless device to riftd running on Peer 0, who then forwards the UDP packet to the
locally-running UDP server.

## Getting started

### Software requirements

* A Unix-based system (Windows may work, but this has not been tested)
* [Go](https://golang.org)

### Hardware requirements

For testing with real hardware, you will need:

  * One or more [ZigBee/XBee (Series 1) radios](https://www.digi.com/products/models/xb24-api-001)
    configured for API-mode (exact configuration details to come)

### Installation

From Github.com:

```
go install github.com/jerluc/riftd
```

From source:

```
# Clone the source code
git clone https://github.com/jerluc/riftd.git

# Install the riftd binary
cd riftd && make && make install
```

### Usage

```
usage: riftd [<flags>] <command> [<args> ...]

Rift protocol daemon

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.


  run [<flags>]
    Starts the Rift protocol daemon

    --logging="INFO"          Log level
    --iface="rift0"           Network interface name
    --dev=/dev/ttyUSB0        Serial device name
    --cidr=2001:412:abcd:1::  IPv6 64-bit prefix

  version
    Displays riftd version


  configure [<flags>]
    Configures a new device for Rift

    --dev=/dev/ttyUSB0  Serial device name
```

### Shell autocompletion

You can also add shell autocompletion (for Bash or ZSH only) by adding the following to your
`.bash_profile` (or equivalent file) for Bash:

```
eval "$(riftd --completion-script-bash)"
```

or for ZSH:

```
eval "$(riftd --completion-script-zsh)"
```

## License

MIT License

[View full license](LICENSE)
