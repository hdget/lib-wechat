package officeaccount

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/hdget/utils/hash"
	"github.com/pkg/errors"
)

// JsSdkSignatureResult signature接口返回结果
type JsSdkSignatureResult struct {
	AppID     string `json:"appId"`
	Ticket    string `json:"ticket"`
	Noncestr  string `json:"noncestr"`
	Url       string `json:"Url"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

// CreateJsSdkSignature 生成JsApi SDK微信签名
// nolint:recheck
func (impl *wxoaApiImpl) CreateJsSdkSignature(url string) (*JsSdkSignatureResult, error) {
	ticket, err := impl.getJsSdkTicket()
	if err != nil {
		return nil, err
	}

	signature, err := impl.createJsSdkSignature(ticket, url)
	if err != nil {
		return nil, err
	}

	if signature == nil || signature.Signature == "" {
		return nil, errors.New("invalid signature")
	}

	return signature, nil
}

func (impl *wxoaApiImpl) getJsSdkTicket() (string, error) {
	ticket, _ := impl.cache.GetJsSdkTicket()
	if ticket == "" {
		accessToken, err := impl.getAccessToken()
		if err != nil {
			return "", err
		}

		ticketResult, err := impl.wxapi.GetJsSdkTicket(accessToken)
		if err != nil {
			return "", err
		}

		err = impl.cache.SetJsSdkTicket(ticketResult.Value, ticketResult.ExpiresIn)
		if err != nil {
			return "", err
		}

		ticket = ticketResult.Value
	}
	return ticket, nil
}

// 生成微信签名
func (impl *wxoaApiImpl) createJsSdkSignature(ticket, url string) (*JsSdkSignatureResult, error) {
	now := time.Now().Unix()
	nonceStr := hash.GenerateRandString(32)
	s := fmt.Sprintf(
		"jsapi_ticket=%s&noncestr=%s&timestamp=%d&Url=%s",
		ticket,
		nonceStr,
		now,
		url,
	)

	// 获取signature
	h := sha1.New()
	_, err := h.Write([]byte(s))
	if err != nil {
		return nil, err
	}
	hashValue := fmt.Sprintf("%x", h.Sum(nil))

	return &JsSdkSignatureResult{
		AppID:     impl.wxapi.GetAppId(),
		Ticket:    ticket,
		Noncestr:  nonceStr,
		Url:       url,
		Timestamp: now,
		Signature: hashValue,
	}, nil
}
