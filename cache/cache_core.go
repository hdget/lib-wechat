package cache

import (
	"fmt"
	"github.com/hdget/common/intf"
)

type Kind string // 微信业务类型

// OAP: Office Account Platform
// OP: Open Platform
const (
	KindOAPMiniProgram    Kind = "oap:mp" // 微信小程序
	KindOAPOfficeAccount  Kind = "oap:oa" // 微信公众号
	KindOPServiceProvider Kind = "op:sp"  // 微信开放平台:微信第三方平台
	KindOPWeb             Kind = "op:web" // 微信开放平台:网站应用
)

type Cache interface {
	Get(key string) (string, error)
	Set(key, value string, expires ...int) error
	Del(key string) error
	HGet(key, member string) (string, error)
	HSet(key, member, value string) error
}

type cacheCoreImpl struct {
	Kind          Kind
	AppId         string
	RedisProvider intf.RedisProvider
}

func New(kind Kind, appId string, redisProvider intf.RedisProvider) Cache {
	return &cacheCoreImpl{
		Kind:          kind,
		AppId:         appId,
		RedisProvider: redisProvider,
	}
}

const (
	redisKeyTemplate = "%s:%s:%s"
)

func (c *cacheCoreImpl) Get(key string) (string, error) {
	bs, err := c.RedisProvider.My().Get(c.getFullKey(key))
	return string(bs), err
}

func (c *cacheCoreImpl) Set(key, value string, expires ...int) error {
	if len(expires) == 0 {
		return c.RedisProvider.My().Set(c.getFullKey(key), value)
	}
	return c.RedisProvider.My().SetEx(c.getFullKey(key), value, expires[0])
}

func (c *cacheCoreImpl) Del(key string) error {
	return c.RedisProvider.My().Del(c.getFullKey(key))
}

func (c *cacheCoreImpl) HGet(key, member string) (string, error) {
	return c.RedisProvider.My().HGetString(c.getFullKey(key), member)
}

func (c *cacheCoreImpl) HSet(key, member, value string) error {
	_, err := c.RedisProvider.My().HSet(c.getFullKey(key), member, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *cacheCoreImpl) getFullKey(key string) string {
	return fmt.Sprintf(redisKeyTemplate, c.Kind, c.AppId, key)
}
