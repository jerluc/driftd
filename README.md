# riftd

[![](https://api.travis-ci.org/jerluc/riftd.svg?branch=master)](https://travis-ci.org/jerluc/riftd)

The Rift protocol daemon

## Getting started

### System requirements

* A Unix-based system (Windows may work, but this has not been tested)
* [Rust and Cargo](https://www.rust-lang.org)
* For testing with real hardware:
  * One or more [ZigBee/XBee (Series 1) radios](https://www.digi.com/products/models/xb24-api-001) configured for API-mode

### Installation from source

```
# Clone the source code
git clone https://github.com/jerluc/riftd.git

# Install the riftd binary
cd riftd && cargo install --force
```

### Usage

```
USAGE:
    riftd [OPTIONS] --device <DEVICE_NAME>

FLAGS:
    -h, --help       Prints help information
    -V, --version    Prints version information

OPTIONS:
    -b, --broadcast <BROADCAST>    Sets the node broadcast interval [default: 1000]
    -s, --device <DEVICE_NAME>     Sets the serial device to use
    -p, --poll <POLL_INTERVAL>     Sets the I/O event poll interval [default: 100]
```

E.g. Using a USB device on serial port `/dev/ttyUSB0` and broadcasting
advertising packets every second:

```bash
$ riftd --device /dev/ttyUSB0
```

## Under development

See [TODOS](TODOS.md)

## License

MIT License

Copyright (c) 2016 Jeremy Lucas

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
