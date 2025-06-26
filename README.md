# Anchor

[![Go Reference](https://pkg.go.dev/badge/github.com/xchacha20-poly1305/anchor.svg)](https://pkg.go.dev/github.com/xchacha20-poly1305/anchor)

A protocol that allow switching proxy info in Lan. And a tun2socks implementation based on sing-tun.

Forked from [SagerNet/SagerConnect](https://github.com/SagerNet/SagerConnect)

# Warn

I know nothing about L4 proxies and transparent proxies, and this is my first time to
something about them. So this software is for educational and research purposes only.
Please do not use it for illegal activities or in production environments.

# Packages

## anchor

Protocol serialization implementation.

Here is the specification of anchor protocol:

* All field use Big Endian.

* Generally speaking, service should listen at port `45947`, which generated randomly.

### Query:

| Version | Device Name Length | Device Name |
 |---------|--------------------|-------------|
| 1       | 1                  | 0 to 128    |

- **Version**: always constant `0x02`.

- **Device Name Length** and **Device Name**: device name should not more than 128.

### Response:

| Version | Dns Port | Device Name Length | Device Name | Socks Port |
|---------|----------|--------------------|-------------|------------|
| 1       | 2        | 1                  | 0 to 128    | 2          |

- **Version**: always constant `0x02`.

- **Dns Port**: can be zero, which means server not providing dns port.

- **Device Name Length** and **Device Name**: device name should not more than 128.

- **Socks Port**: the socks service's port.

## tun2dialer

A router routes connections from tunnel to dialer. You can easily implement "tun2socks" via it.

## route

Route functions for tun2dialer.

## log

Implements `github.com/sagernet/sing/common/logger.Logger` by [zap](https://go.uber.org/zap).

## sockshttp

An inbound that mixed socks and HTTP.

## anchorservice

A server that provides anchor protocol service.

## cmd

Anchor service detector + tun2socks instance, which also includes simple routes function.

Config example: [example.config.json](./example.config.json). The detail description of
each field see [sing-box](https://sing-box.sagernet.org/configuration/inbound/tun/).

Besides, you can just run it as a tun2socks software using command like
`anchor -a <ip> -socks1080 -c <config>`.

```
Usage of anchor:
  -a string
    	remote ip address (skip scan)
  -c string
    	Configuration file path
  -d int
    	selected device index (skip select) (default -1)
  -dns int
    	remote dns port (skip scan) (default 6450)
  -i	skip waiting when the first device is found
  -l string
    	Log level (default "warn")
  -o string
    	Log output. (default "stderr")
  -socks int
    	remote socks port (skip scan) (default 2080)
  -v	Show version

```

## Build

```shell
CGO_ENABLED=0 make
```

# Credits

Inspired by: [SagerNet/SagerConnect](https://github.com/SagerNet/SagerConnect)

- [SagerNet/sing-tun](https://github.com/SagerNet/sing-tun)
- [SagerNet/sing-box](https://github.com/SagerNet/sing-box)
