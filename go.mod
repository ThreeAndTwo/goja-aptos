module github.com/ThreeAndTwo/goja-onchain-vm

go 1.16

require (
	github.com/btcsuite/btcd v0.22.1
	github.com/dop251/goja v0.0.0-20211203105952-bf6af58bbcc8
	github.com/dop251/goja_nodejs v0.0.0-20210225215109-d91c329300e7
	github.com/ethereum/go-ethereum v1.10.16
	github.com/imroc/req v0.3.0
	github.com/threeandtwo/aptclient v0.1.0
)

replace (
	github.com/alethio/web3-multicall-go v0.0.15 => github.com/ThreeAndTwo/web3-multicall-go v0.0.17
	github.com/dgrijalva/jwt-go v3.2.0+incompatible => github.com/golang-jwt/jwt/v4 v4.4.0
	github.com/docker/docker v1.4.2-0.20180625184442-8e610b2b55bf => github.com/docker/docker v20.10.13+incompatible
)
