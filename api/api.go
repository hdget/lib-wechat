package api

import (
	"encoding/json"
	"fmt"
	"github.com/hdget/common/intf"
	"github.com/hdget/common/types"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

type Api struct {
	Kind      ApiKind
	AppId     string
	AppSecret string
	Cache     Cache
}

// Result 微信的API结果
type Result struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
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

func (a *Api) ParseApiResult(data []byte, result any) error {
	err := json.Unmarshal(data, result)
	if err != nil {
		// 如果unmarshal错误,尝试解析错误信息
		var ret Result
		_ = json.Unmarshal(data, &ret)

		if ret.ErrCode != 0 {
			return fmt.Errorf("api error, err: %s", ret.ErrMsg)
		} else {
			return errors.Wrapf(err, "unmarshal api result, body: %s", convert.BytesToString(data))
		}
	}
	return nil
}
