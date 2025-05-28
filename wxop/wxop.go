package wxop

import (
	"github.com/hdget/common/intf"
	"github.com/hdget/lib-wechat/api"
	"github.com/pkg/errors"
)

type ApiWxop interface {
	api.Common
	HandleAuthEvent(signature, timestamp, nonce, body string, callbacks map[string]AuthCallback) error // 处理授权事件
	GetComponentVerifyTicket() (string, error)                                                         // 获取验证票据, 有效期12个小时
	GetComponentAccessToken() (string, error)                                                          // 获取第三方平台接口的调用凭据, 有效期为2小时
	GetPreAuthCode() (string, error)                                                                   // 获取预授权码
	GetAuthUrl(client, redirectUrl, authCode string) (string, error)                                   // 获取授权链接
}

type wxopImpl struct {
	*api.Api
	crypt *WXBizMsgCrypt
}

func New(appId, appSecret, token, encodingAESKey string, providers ...intf.Provider) (ApiWxop, error) {
	crypt, err := NewWXBizMsgCrypt(appId, token, encodingAESKey)
	if err != nil {
		return nil, err
	}

	a, err := api.New(api.ApiKindWxop, appId, appSecret, providers...)
	if err != nil {
		return nil, err
	}

	if a.Cache == nil {
		return nil, errors.New("redis provider not provided")
	}

	return &wxopImpl{
		Api:   a,
		crypt: crypt,
	}, nil
}
