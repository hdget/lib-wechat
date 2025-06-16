package api

import (
	"fmt"
)

type Api interface {
	CheckResult(result Result, url string, request ...any) error
	GetAppId() string
	GetAppSecret() string
}

type apiImpl struct {
	AppId     string
	AppSecret string
}

// Result 微信的API结果
type Result struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func New(appId, appSecret string) Api {
	return &apiImpl{
		AppId:     appId,
		AppSecret: appSecret,
	}
}

func (c *apiImpl) GetAppId() string {
	return c.AppId
}

func (c *apiImpl) GetAppSecret() string {
	return c.AppSecret
}

func (c *apiImpl) CheckResult(result Result, url string, request ...any) error {
	if result.ErrCode != 0 {
		if len(request) > 0 {
			return fmt.Errorf("invoke wx api, err: %s, url: %s, request: %v", result.ErrMsg, url, request[0])
		} else {
			return fmt.Errorf("invoke wx api, err: %s, url: %s", result.ErrMsg, url)
		}
	}
	return nil
}
