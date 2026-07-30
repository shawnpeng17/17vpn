# 17vpn

<img width="520" height="295" alt="image" src="https://github.com/user-attachments/assets/07a1511e-cc1a-4ce1-9be8-888077d2a198" />

### Pre-Installation

Follow the [confluence](https://17media.atlassian.net/wiki/spaces/H/pages/1027244286/OKTA+Pritunl+VPN) install Pritunl client and import profiles first

### Installation

1. Binary

Download the pre-built binaries from the Releases page. Extract them, move it to your $PATH.

```shell
curl -OL https://github.com/shawnpeng17/17vpn/releases/download/v2.0.0/17vpn_2.0.0_Darwin_arm_64.tar.gz
tar -xzvf 17vpn_2.0.0_Darwin_arm_64.tar.gz
mv 17vpn /usr/local/bin
17vpn
```

2. Source
```shell
# install it to your $GOPATH/bin
go install github.com/shawnpeng17/17vpn@v2.0.0 
```

### Usage

```shell
# Initial your OTP key and Pin (first time)
# Enter ID or Server to connect/disconnect
$ 17vpn

# Disconnect all connections
$ 17vpn d
```
