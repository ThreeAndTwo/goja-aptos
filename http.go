package goja_onchain_vm

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req"
	"strings"
)

type ReqType string

const (
	POST ReqType = "post"
	GET  ReqType = "get"
)

type gojaReq struct {
	url     string
	header  req.Header
	param   req.Param
	reqType ReqType
	isJson  bool
}

func NewGojaReq(url string, header req.Header, param req.Param, reqType ReqType) *gojaReq {
	return &gojaReq{url: url, header: header, param: param, reqType: reqType}
}

func initHeader(header map[string]string) (req.Header, bool) {
	authHeader := req.Header{}
	hasJson := false

	for k, v := range header {
		authHeader[k] = v
		if hasJsonInHeader(k, v) {
			hasJson = true
		}
	}
	return authHeader, hasJson
}

func hasJsonInHeader(key, value string) bool {
	return strings.ToLower(key) == "content-type" && strings.Contains(strings.ToLower(value), "json")
}

func initParam(params map[string]string) req.Param {
	reqParams := req.Param{}
	for k, v := range params {
		reqParams[k] = v
	}
	return reqParams
}

func (gReq *gojaReq) check() bool {
	return gReq.url == ""
}

func (gReq *gojaReq) request() (string, error) {
	if gReq.check() {
		return "", fmt.Errorf("url should be null")
	}

	switch gReq.reqType {
	case POST:
		return gReq.post()
	case GET:
		return gReq.get()
	default:
		return gReq.get()
	}
}

func (gReq *gojaReq) post() (string, error) {
	var reqResp = &req.Resp{}
	var err error
	if gReq.isJson {
		jsonParam, _ := json.Marshal(gReq.param)
		reqResp, err = req.Post(gReq.url, jsonParam, gReq.header)
	} else {
		reqResp, err = req.Post(gReq.url, gReq.param, gReq.header)
	}

	if err != nil {
		return "", err
	}

	msHeader, _ := json.Marshal(reqResp.Response().Header)
	return reqResp.String() + "||" + string(msHeader), err
}

func (gReq *gojaReq) get() (string, error) {
	resp, err := req.Get(gReq.url, gReq.header)
	if err != nil {
		return "", err
	}

	msHeader, err := json.Marshal(resp.Response().Header)
	if err != nil {
		return "", err
	}
	return resp.String() + "||" + string(msHeader), nil
}
