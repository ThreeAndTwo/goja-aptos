package gojaaptos

import (
	"testing"
)

func TestGetAddress(t *testing.T) {
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
	expectAddr := "0x4e7a58adca88cfa5c99dcf92e29f248d1acabd3efc8d5183e3148849c07f7659"
	if addr != expectAddr {
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
