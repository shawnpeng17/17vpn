# 17vpn

![image](https://user-images.githubusercontent.com/91862792/172811759-851153ee-8e76-4e77-a45a-a11504dce767.png)

## Installation

Follow the [confluence](https://17media.atlassian.net/wiki/spaces/H/pages/1027244286/OKTA+Pritunl+VPN) install Pritunl client and import profiles

### Binary

Download the pre-built binaries from the Releases page. Extract them, move it to your $PATH.

```shell
curl -OL https://github.com/shawnpeng17/17vpn/releases/download/v1.0.0/17vpn_1.0.0_Darwin_x86_64.tar.gz
tar -xzvf 17vpn_1.0.0_Darwin_x86_64.tar.gz
mv 17vpn /usr/local/bin
17vpn
```

### Source
```shell
git clone git@github.com:shawnpeng17/17vpn.git
cd 17vpn
go install
```
