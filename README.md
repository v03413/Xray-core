# Project X - Socks5

> 一个有意思的项目

## 编译

```bash
# Linux-64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -trimpath -ldflags "-s -w" -ldflags="-s -w" -o socks ./main

```