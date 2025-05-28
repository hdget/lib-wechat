package wxop

import (
	"github.com/hdget/common/intf"
	"github.com/hdget/lib-wechat/api"
	"github.com/pkg/errors"
)

// LoadAuthorizerRefreshToken 定义从持久化存储中加载refreshToken的函数
type LoadAuthorizerRefreshToken func(authorizerAppid string) (string, error)

type ApiWxop interface {
	HandleAuthEvent(signature, timestamp, nonce, body string, callbacks map[string]AuthCallback) error                  // 处理授权事件
	GetAuthUrl(client, redirectUrl, authCode string) (string, error)                                                    // 获取授权链接
	GetAuthorizerAccessToken(authorizerAppid string, fnLoadAuthRefreshToken LoadAuthorizerRefreshToken) (string, error) // 获取授权访问令牌
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
