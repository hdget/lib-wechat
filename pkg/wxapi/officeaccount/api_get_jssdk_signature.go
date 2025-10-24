package officeaccount

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/hdget/lib-wechat/pkg/wxapi"
	"github.com/hdget/utils/hash"
)

// GetJsSdkSignatureResult signature接口返回结果
type GetJsSdkSignatureResult struct {
	*wxapi.Result
	AppID     string `json:"appId"`
	Ticket    string `json:"ticket"`
	Noncestr  string `json:"noncestr"`
	Url       string `json:"Url"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

// GetJsSdkSignature 生成微信签名
func (impl officeAccountWxApiImpl) GetJsSdkSignature(ticket, url string) (*GetJsSdkSignatureResult, error) {
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

	return &GetJsSdkSignatureResult{
		AppID:     impl.GetAppId(),
		Ticket:    ticket,
		Noncestr:  nonceStr,
		Url:       url,
		Timestamp: now,
		Signature: hashValue,
	}, nil
}
