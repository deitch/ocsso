# openconnect-sso Native

This is a replacement for
[openconnect-sso](https://github.com/vlaci/openconnect-sso/). That wrapper
is great. Its challenge is that it is entirely written in Python, and thus is
dependent on various Python packages' availability and installability on
different platforms. I had problems running it on various macOS, including M1,
So decide to write a replacement.

As of this writing, it only support Cisco AnyConnect. We might add more in the
future.

## Usage

1. Install [openconnect](https://www.infradead.org/openconnect/)
1. Install this binary
1. Run it

For options, just run `ocsso -h`.

Basic usage:

```
ocsso connect --server vpn.remote.com
```

Depending on your platform, you might need to run it with `sudo`, i.e.

```
sudo ocsso connect --server vpn.remote.com
```

## How it works

See [FLOW.md](./FLOW.md)