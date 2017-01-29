# riftd

[![](https://api.travis-ci.org/jerluc/riftd.svg?branch=master)](https://travis-ci.org/jerluc/riftd)

The Rift protocol daemon

## Getting started

### System requirements

* A Unix-based system (Windows may work, but this has not been tested)
* [Go](https://golang.org)
* For testing with real hardware:
  * One or more [ZigBee/XBee (Series 1) radios](https://www.digi.com/products/models/xb24-api-001) configured for API-mode

### Installating

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
