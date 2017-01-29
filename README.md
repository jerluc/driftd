# riftd

[![](https://api.travis-ci.org/jerluc/riftd.svg?branch=master)](https://travis-ci.org/jerluc/riftd)

The Rift protocol daemon

## Getting started

### System requirements

* A Unix-based system (Windows may work, but this has not been tested)
* [Go](https://golang.org)
* For testing with real hardware:
  * One or more [ZigBee/XBee (Series 1) radios](https://www.digi.com/products/models/xb24-api-001) configured for API-mode

### Installation from source

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

## License

MIT License

Copyright (c) 2017 Jeremy Lucas

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
