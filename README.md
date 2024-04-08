# 17vpn

![image](https://user-images.githubusercontent.com/91862792/172811759-851153ee-8e76-4e77-a45a-a11504dce767.png)


### Pre-Installation

Follow the [confluence](https://17media.atlassian.net/wiki/spaces/H/pages/1027244286/OKTA+Pritunl+VPN) install Pritunl client and import profiles first

### Installation

1. Binary

Download the pre-built binaries from the Releases page. Extract them, move it to your $PATH.

```shell
curl -OL https://github.com/shawnpeng17/17vpn/releases/download/v1.0.0/17vpn_1.0.0_Darwin_x86_64.tar.gz
tar -xzvf 17vpn_1.0.0_Darwin_x86_64.tar.gz
mv 17vpn /usr/local/bin
17vpn
```

2. Source
```shell
# install it to your $GOPATH/bin
go install github.com/shawnpeng17/17vpn@latest 
```

### Usage

```shell
# Initial your OTP key and Pin (first time)
# Enter ID or Server to connect/disconnect
$ 17vpn

# Disconnect all connections
$ 17vpn d
```