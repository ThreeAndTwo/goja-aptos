package gojaaptos

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/imroc/req"
	"github.com/threeandtwo/aptclient/client"
)

func NewVMGlobal(chainInfo ChainInfo, accountInfo AccountInfo) (*VMGlobal, error) {
	rt := goja.New()
	client, err := client.NewAptClient(chainInfo.Rpc)
	if err != nil {
		return nil, err
	}
	return &VMGlobal{
		Runtime:     rt,
		ChainInfo:   chainInfo,
		AccountInfo: accountInfo,
		Client:      client,
	}, nil

}

func (gvm *VMGlobal) Init() error {
	registry := require.NewRegistry()
	if !gvm.check() {
		return fmt.Errorf("gvm config error, please check your config")
	}

	vm := gvm.Runtime
	registry.Enable(vm)
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	var err error

	err = vm.Set(string(TokenBalance), gvm.GetTokenBalance)
	if err != nil {
		return err
	}

	// err = vm.Set(string(CALL), gvm.Call)
	// if err != nil {
	// 	return err
	// }

	err = vm.Set(string(GetAddress), gvm.GetAddress)
	if err != nil {
		return err
	}

	err = vm.Set(string(GetPreAddress), gvm.GetPreAddress)
	if err != nil {
		return err
	}

	err = vm.Set(string(GetNextAddress), gvm.GetNextAddress)
	if err != nil {
		return err
	}

	err = vm.Set(string(GetAddressByIndex), gvm.GetAddressByIndex)
	if err != nil {
		return err
	}

	err = vm.Set(string(GetAddressListByIndex), gvm.GetAddressListByIndex)
	if err != nil {
		return err
	}

	err = vm.Set(string(GetCurrentIndex), gvm.GetCurrentIndex)
	if err != nil {
		return err
	}

	err = vm.Set(string(EncryptWithPubKey), gvm.EncryptWithPubKey)
	if err != nil {
		return err
	}

	// http get
	err = vm.Set(string(HttpGetRequest), gvm.HttpGet)
	if err != nil {
		return err
	}

	// http post
	return vm.Set(string(HttpPostRequest), gvm.HttpPost)
}

func (gvm *VMGlobal) check() bool {
	return nil != gvm && nil != gvm.Runtime
}
func (gvm *VMGlobal) getClient() (*client.AptClient, error) {
	return client.NewAptClient(gvm.ChainInfo.Rpc)
}

func (gvm *VMGlobal) GetTokenBalance(tokenFullQualifiedName, accountAddress string) goja.Value {
	url := fmt.Sprintf("%s/v1/accounts/%s/resources", gvm.ChainInfo.Rpc, gvm.getAddress())
	resp, err := req.Get(url)
	if err != nil {
		gvm.Runtime.Interrupt(`failed requests node for account resources: ` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}
	var accountResources []AccountToken
	err = resp.ToJSON(&accountResources)
	if err != nil {
		gvm.Runtime.Interrupt(`failed decode account resources: ` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}
	balance := big.NewInt(0)
	for _, r := range accountResources {
		if r.Type == fmt.Sprintf("0x1::coin::CoinStore<%s>", tokenFullQualifiedName) {
			balance.SetString(r.Data.Coin.Value, 10)
			break
		}
	}

	return gvm.Runtime.ToValue(balance)
}

func (gvm *VMGlobal) HttpGet(url, params, header string) goja.Value {
	reqHeader, reqParam, _, err := getReqParam(params, header)
	if err != nil {
		gvm.Runtime.Interrupt(err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}

	_req := NewGojaReq(url, reqHeader, reqParam, GET)
	data, err := _req.request()
	if err != nil {
		gvm.Runtime.Interrupt(`http request error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}
	return gvm.Runtime.ToValue(data)
}

func (gvm *VMGlobal) HttpPost(url, params, header string) goja.Value {
	reqHeader, reqParam, isJson, err := getReqParam(params, header)
	if err != nil {
		gvm.Runtime.Interrupt(err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}

	_req := NewGojaReq(url, reqHeader, reqParam, POST)
	_req.isJson = isJson
	resp, err := _req.request()
	if err != nil {
		gvm.Runtime.Interrupt(`http request error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}
	return gvm.Runtime.ToValue(resp)
}

func getReqParam(params, header string) (req.Header, req.Param, bool, error) {
	headerMap := make(map[string]string)
	paramsMap := make(map[string]string)

	if header != "" {
		err := json.Unmarshal([]byte(header), &headerMap)
		if err != nil {
			return nil, nil, false, fmt.Errorf("http params invalidate for header: %s", header)
		}
	}

	if params != "" {
		err := json.Unmarshal([]byte(params), &paramsMap)
		if err != nil {
			return nil, nil, false, fmt.Errorf("http params invalidate for params: %s", params)
		}
	}

	reqHeader, isJson := initHeader(headerMap)
	reqParam := initParam(paramsMap)
	return reqHeader, reqParam, isJson, nil
}

func (gvm *VMGlobal) GetAddress() goja.Value {
	return gvm.getAddress()
}

func (gvm *VMGlobal) GetPreAddress() goja.Value {
	gvm.AccountInfo.Index--
	return gvm.getAddress()
}

func (gvm *VMGlobal) GetNextAddress() goja.Value {
	gvm.AccountInfo.Index++
	return gvm.getAddress()
}

func (gvm *VMGlobal) GetAddressByIndex(index int) goja.Value {
	gvm.AccountInfo.Index = index
	return gvm.getAddress()
}

func (gvm *VMGlobal) getAddressString() (string, error) {

	account := client.NewAptAccount(gvm.AccountInfo.Key, "")
	addr, err := account.GetAptAccount(gvm.AccountInfo.Index)
	if err != nil {
		return "", fmt.Errorf("account invalidated: %s", err)
	}
	return addr.Address, nil
}

func (gvm *VMGlobal) getAddress() goja.Value {
	if gvm.checkAddress() {
		gvm.Runtime.Interrupt(`params invalidate for address, index:` + fmt.Sprintf("%d", gvm.AccountInfo.Index))
		return gvm.Runtime.ToValue(`exception`)
	}

	address, err := gvm.getAddressString()
	if err != nil {
		gvm.Runtime.Interrupt(`params invalidate for address, index:` + fmt.Sprintf("%d", gvm.AccountInfo.Index))
		return gvm.Runtime.ToValue(`exception`)
	}
	return gvm.Runtime.ToValue(address)
}

func (gvm *VMGlobal) GetAddressListByIndex(start, end int) goja.Value {
	gvm.AccountInfo.Index = start
	if gvm.checkAddress() {
		gvm.Runtime.Interrupt(`params invalidate for address, index:` + fmt.Sprintf("%d", gvm.AccountInfo.Index))
		return gvm.Runtime.ToValue(`exception`)
	}

	var arrAddr []string
	for k := start; k < end; k++ {
		gvm.AccountInfo.Index = k

		address, err := gvm.getAddressString()
		if err != nil {
			gvm.Runtime.Interrupt(`params invalidate for address, index:` + fmt.Sprintf("%d", gvm.AccountInfo.Index))
			return gvm.Runtime.ToValue(`exception`)
		}
		arrAddr = append(arrAddr, address)
	}
	addresses := strings.Join(arrAddr, ",")
	return gvm.Runtime.ToValue(addresses)
}

func (gvm *VMGlobal) checkAddress() bool {
	return gvm.AccountInfo.Key == "" || gvm.AccountInfo.Index < 0
}

func (gvm *VMGlobal) GetCurrentIndex() goja.Value {
	return gvm.Runtime.ToValue(gvm.AccountInfo.Index)
}

func (gvm *VMGlobal) EncryptWithPubKey(message string) goja.Value {
	if message == "" {
		gvm.Runtime.Interrupt(`params invalidate for encryptWithPubKey`)
		return gvm.Runtime.ToValue(`exception`)
	}
	account := client.NewAptAccount(gvm.AccountInfo.Key, "")
	indexAccount, err := account.GetAptAccount(gvm.AccountInfo.Index)
	if err != nil {
		gvm.Runtime.Interrupt(`failed get account at index` + err.Error())
		return gvm.Runtime.ToValue(`exception`)

	}

	signerKey, err := hexutil.Decode("0x" + indexAccount.PublicKey)
	if err != nil {
		gvm.Runtime.Interrupt(`decode public key error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}

	pubKey, err := btcec.ParsePubKey(signerKey, btcec.S256())
	if err != nil {
		gvm.Runtime.Interrupt(`parse public key error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}

	encryptData, err := btcec.Encrypt(pubKey, []byte(message))
	if err != nil {
		gvm.Runtime.Interrupt(`encrypt data error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}
	return gvm.Runtime.ToValue(hexutil.Encode(encryptData))
}
