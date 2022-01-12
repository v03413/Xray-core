package socks

import (
	"github.com/xtls/xray-core/common/protocol"
)

func (a *Account) Equals(another protocol.Account) bool {
	if account, ok := another.(*Account); ok {
		return a.Username == account.Username
	}
	return false
}

func (a *Account) AsAccount() (protocol.Account, error) {
	return a, nil
}

func (c *ServerConfig) HasAccount(username, password, srcIp, cid string) bool {

	return extend.Auth(username, password, srcIp, cid)
}
