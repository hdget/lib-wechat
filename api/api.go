package api

import (
	"github.com/hdget/common/intf"
	"github.com/hdget/common/types"
	"github.com/pkg/errors"
)

type Api struct {
	Kind      ApiKind
	AppId     string
	AppSecret string
	Cache     Cache
}

// ApiResult 微信的API结果
type ApiResult struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type ApiCommon interface {
	GetAccessToken() (string, error)
}

func New(kind ApiKind, appId, appSecret string, providers ...intf.Provider) (*Api, error) {
	if kind == ApiKindUnknown || appId == "" || appSecret == "" {
		return nil, errors.New("invalid parameter")
	}

	b := &Api{
		Kind:      kind,
		AppId:     appId,
		AppSecret: appSecret,
	}

	for _, provider := range providers {
		if provider.GetCapability().Category == types.ProviderCategoryRedis {
			b.Cache = newCache(kind, appId, provider.(intf.RedisProvider))
		}
	}

	return b, nil
}
