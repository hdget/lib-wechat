package api

import (
	"fmt"
	"github.com/hdget/common/intf"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key, value string, expires ...int) error
	Del(key string) error
	HGet(key, member string) (string, error)
	HSet(key, member, value string) error
}

const (
	redisKeyTemplate = "%s:%s:%s"
)

type cacheImpl struct {
	apiKind       ApiKind
	appId         string
	redisProvider intf.RedisProvider
}

var _ Cache = (*cacheImpl)(nil)

func newCache(businessKind ApiKind, appId string, redisProvider intf.RedisProvider) Cache {
	return &cacheImpl{
		apiKind:       businessKind,
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

func (c *cacheImpl) Del(key string) error {
	return c.redisProvider.My().Del(c.getFullKey(key))
}

func (c *cacheImpl) HGet(key, member string) (string, error) {
	return c.redisProvider.My().HGetString(c.getFullKey(key), member)
}

func (c *cacheImpl) HSet(key, member, value string) error {
	_, err := c.redisProvider.My().HSet(c.getFullKey(key), member, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *cacheImpl) getFullKey(key string) string {
	return fmt.Sprintf(redisKeyTemplate, c.apiKind, c.appId, key)
}
