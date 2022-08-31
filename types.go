package gojaaptos

import (
	"github.com/dop251/goja"
	"github.com/threeandtwo/aptclient/client"
)

type VMGlobal struct {
	Runtime     *goja.Runtime
	ChainInfo   ChainInfo
	AccountInfo AccountInfo
	Client      *client.AptClient
}

type ChainInfo struct {
	ChainId int64
	Rpc     string
	// Wss     string
}

type AccountInfo struct {
	Key   string
	Index int
}

type VmFunc string

const (
	Balance               VmFunc = "balance"
	TokenBalance          VmFunc = "tokenBalance"
	CALL                  VmFunc = "contractCall"
	GetAddress            VmFunc = "getAddress"
	GetPreAddress         VmFunc = "getPreAddress"
	GetNextAddress        VmFunc = "getNextAddress"
	GetAddressByIndex     VmFunc = "getAddressByIndex"
	GetAddressListByIndex VmFunc = "getAddressListByIndex"
	GetCurrentIndex       VmFunc = "getCurrentIndex"
	PersonalSign          VmFunc = "personalSign"
	HttpGetRequest        VmFunc = "httpGetRequest"
	HttpPostRequest       VmFunc = "httpPostRequest"
	EncryptWithPubKey     VmFunc = "encryptWithPubKey"
)

type AccountToken struct {
	Type string           `json:"type"`
	Data AccountTokenData `json:"data"`
}
type AccountTokenData struct {
	Coin *AccountTokenDataCoin `json:"coin,omitempty"`
}
type AccountTokenDataCoin struct {
	Value string `json:"value"`
}
