package types

import "github.com/hdget/lib-wechat/api"

// JsSdkSignatureResult signature接口返回结果
type JsSdkSignatureResult struct {
	api.Result
	AppID     string `json:"appId"`
	Ticket    string `json:"ticket"`
	Noncestr  string `json:"noncestr"`
	Url       string `json:"Url"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

// JsSdkTicketResult 类型
type JsSdkTicketResult struct {
	api.Result
	Value     string `json:"ticket,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"`
}
