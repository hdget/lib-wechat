package api

import "github.com/hdget/lib/lib-wechat/api"

type Api interface {
	api.Api
	Login(code string) (string, string, error)
}

type apiImpl struct {
	api.Api
}

func New(appId, appSecret string) Api {
	return &apiImpl{
		Api: api.New(appId, appSecret),
	}
}
