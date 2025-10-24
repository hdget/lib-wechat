package miniprogram

import (
	"github.com/hdget/common/types"
	"github.com/hdget/lib-wechat/miniprogram/cache"
	"github.com/hdget/lib-wechat/pkg/wxapi/miniprogram"
)

type API interface {
	Login(code string) (string, string, error) // 获取OpenId和UnionId
}

type miniProgramApiImpl struct {
	cache     cache.Cache
	appId     string
	appSecret string
}

func New(appId, appSecret string, redisProvider types.RedisProvider) API {
	return &miniProgramApiImpl{
		cache: cache.New(appId, redisProvider),
	}
}

func (impl miniProgramApiImpl) Login(code string) (string, string, error) {
	result, err := miniprogram.New(impl.appId, impl.appSecret).Code2Session(code)
	if err != nil {
		return "", "", err
	}

	// 保存到缓存中
	err = impl.cache.SetSessionKey(result.SessionKey)
	if err != nil {
		return "", "", err
	}

	return result.OpenId, result.UnionId, err
}
