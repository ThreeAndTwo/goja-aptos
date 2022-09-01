package gojaaptos

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dop251/goja"
)

func TestGetAddress(t *testing.T) {
	type fields struct {
		VmGlobal *VMGlobal
	}
	vmg, err := NewVMGlobal(ChainInfo{
		Rpc:     "https://fullnode.devnet.aptoslabs.com",
		ChainId: 2,
	}, AccountInfo{
		Key: "0x3bf53a2dc48aedf452c8962950013b325747ece60bc7de6e6a9a70e9d04bb4a8",
		// Key:   "goat energy okay cube kangaroo army picnic wolf grit stairs draft sting",
		Index: 0,
	})
	if err != nil {
		t.Error(err)
	}
	vmg.Init()
	vm := vmg.Runtime
	_, err = vm.RunString(`function run(){return getAddress()}`)
	runFunc, ok := goja.AssertFunction(vm.Get("run"))
	if !ok {
		t.Error(fmt.Errorf("config params error, via mismatch %s", err))
	}

	value, err := runFunc(goja.Undefined())
	if err != nil {
		t.Error(err)
	}
	addr, _ := json.Marshal(value)
	expectAddr := `"0x4e7a58adca88cfa5c99dcf92e29f248d1acabd3efc8d5183e3148849c07f7659"`
	if string(addr) != expectAddr {
		t.Errorf("expect %s, got %s", expectAddr, addr)
	}
}

func TestGetBalance(t *testing.T) {
	type fields struct {
		VmGlobal *VMGlobal
	}
	vm, err := NewVMGlobal(ChainInfo{
		Rpc:     "https://fullnode.devnet.aptoslabs.com",
		ChainId: 2,
	}, AccountInfo{
		Key:   "0x3bf53a2dc48aedf452c8962950013b325747ece60bc7de6e6a9a70e9d04bb4a8",
		Index: 0,
	})
	if err != nil {
		t.Error(err)
	}
	addr, err := vm.getAddressString()
	if err != nil {
		t.Error(err)
	}
	balance := vm.GetTokenBalance("0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::coins::USDT", addr)
	if balance.String() != "256501" {
		t.Errorf("expect %d, got %s", 256501, balance)
	}
}
