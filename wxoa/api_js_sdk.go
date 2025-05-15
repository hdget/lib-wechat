package wxoa

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
	"github.com/hdget/utils/hash"
	"github.com/pkg/errors"
	"time"
)

// JsSdkSignatureResult signature接口返回结果
type JsSdkSignatureResult struct {
	api.ApiResult
	AppID     string `json:"appId"`
	Ticket    string `json:"ticket"`
	Noncestr  string `json:"noncestr"`
	Url       string `json:"Url"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

// JsSdkTicketResult 类型
type JsSdkTicketResult struct {
	api.ApiResult
	Value     string `json:"ticket,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"`
}

const (
	redisKeyJsApiTicket = "js_api_ticket"
	urlGetJsSdkTicket   = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
)

// GetJsSdkSignature 生成JsApi SDK微信签名
// nolint:recheck
func (impl *wxoaImpl) GetJsSdkSignature(url string) (*JsSdkSignatureResult, error) {
	ticket, err := impl.getJsSdkTicket()
	if err != nil {
		return nil, err
	}

	signature, err := impl.generateJsSdkSignature(ticket, url)
	if err != nil {
		return nil, err
	}

	if signature == nil || signature.Signature == "" {
		return nil, errors.New("invalid signature")
	}

	return signature, nil
}

// 生成微信签名
func (impl *wxoaImpl) generateJsSdkSignature(ticket, url string) (*JsSdkSignatureResult, error) {
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
		AppID:     impl.AppId,
		Ticket:    ticket,
		Noncestr:  nonceStr,
		Url:       url,
		Timestamp: now,
		Signature: hashValue,
	}, nil
}

func (impl *wxoaImpl) getJsSdkTicket() (string, error) {
	if impl.Cache == nil {
		return "", errors.New("no redis provider")
	}

	cachedTicket, err := impl.Cache.Get(redisKeyJsApiTicket)
	if err != nil {
		return "", errors.Wrap(err, "get wxoa cached ticket")
	}
	if cachedTicket != "" {
		return cachedTicket, nil
	}

	accessToken, err := impl.GetAccessToken()
	if err != nil {
		return "", err
	}

	result, err := impl.generateJsSdkTicket(accessToken)
	if err != nil {
		return "", err
	}

	// 忽略保存ticket错误
	err = impl.Cache.Set(redisKeyJsApiTicket, result.Value, result.ExpiresIn)
	if err != nil {
		return "", errors.Wrap(err, "cache wxoa ticket")
	}

	return result.Value, nil
}

// generateJsSdkTicket jssdk获取凭证
func (impl *wxoaImpl) generateJsSdkTicket(accessToken string) (*JsSdkTicketResult, error) {
	url := fmt.Sprintf(urlGetJsSdkTicket, accessToken)
	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}

	var result JsSdkTicketResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("%s, url: %s", result.ErrMsg, urlGetJsSdkTicket)
	}

	return &result, nil
}
