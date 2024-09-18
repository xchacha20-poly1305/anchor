# Anchor

[![Go Reference](https://pkg.go.dev/badge/github.com/xchacha20-poly1305/anchor.svg)](https://pkg.go.dev/github.com/xchacha20-poly1305/anchor)

A protocol that allow switching proxy info in Lan.

Forked from [SagerNet/SagerConnect](https://github.com/SagerNet/SagerConnect)

# Warn

I know nothing about L4 proxies and transparent proxies, and this is my first time to
something about them. So this software is for educational and research purposes only.
Please do not use it for illegal activities or in production environments.

# Developing

## Build

```shell
CGO_ENABLED=0 make
```

## Format

```shell
make fmt
```

## Protocol

All use Big Endian.

### Query:

| Version | Device Name Length | Device Name |
 |---------|--------------------|-------------|
| 1       | 1                  | 0 to 128    |

- **Version**: always constant `0x01`

### Response:

| Version | Dns Port | Device Name Length | Device Name | Debug | Bypass Lan | Socks Port | Socks User Length | Socks User | Socks Password Length | Socks Password |
 |---------|----------|--------------------|-------------|-------|------------|------------|-------------------|------------|-----------------------|----------------|
| 1       | 2        | 1                  | 0 to 128    | 1     | 1          | 2          | 1                 | 0 to 255   | 1                     | 0 to 255       |

- **Version**: always constant `0x01`

# Credits

Inspired by: [SagerNet/SagerConnect](https://github.com/SagerNet/SagerConnect)

- [SagerNet/sing-tun](https://github.com/SagerNet/sing-tun)
- [SagerNet/sing-box](https://github.com/SagerNet/sing-box)