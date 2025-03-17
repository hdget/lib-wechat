package api

import (
	"fmt"
	"github.com/hdget/common/intf"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key, value string, expires ...int) error
}

const (
	template = "%d:%s:%s"
)

type cacheImpl struct {
	businessKind  ApiKind
	appId         string
	redisProvider intf.RedisProvider
}

var _ Cache = (*cacheImpl)(nil)

func newCache(businessKind ApiKind, appId string, redisProvider intf.RedisProvider) Cache {
	return &cacheImpl{
		businessKind:  businessKind,
		appId:         appId,
		redisProvider: redisProvider,
	}
}

func (c *cacheImpl) Get(key string) (string, error) {
	bs, err := c.redisProvider.My().Get(c.getFullKey(key))
	return string(bs), err
}

func (c *cacheImpl) Set(key, value string, expires ...int) error {
	if len(expires) == 0 {
		return c.redisProvider.My().Set(c.getFullKey(key), value)
	}
	return c.redisProvider.My().SetEx(c.getFullKey(key), value, expires[0])
}

func (c *cacheImpl) getFullKey(key string) string {
	return fmt.Sprintf(template, c.businessKind, c.appId, key)
}
