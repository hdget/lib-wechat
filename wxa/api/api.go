package api

import (
	"github.com/hdget/lib-wechat/api"
)

type Api interface {
	api.Api
	Code2Session(code string) (*GetSessionResult, error)                               // 小程序静默登录，通过code换取UnionId
	GetUserPhoneNumber(accessToken, code string) (string, error)                       // 获取用户手机号码
	CreateLimitedWxaCode(accessToken, path string, width int) ([]byte, error)          // 生成有限的小程序码
	CreateUnLimitedWxaCode(accessToken, scene, page string, width int) ([]byte, error) // CreateUnLimitedWxaCode 生成小程序码，可接受页面参数较短，生成个数不受限
	GetAccessToken() (*GetAccessTokenResult, error)
}

type apiImpl struct {
	api.Api
}

func New(appId, appSecret string) Api {
	return &apiImpl{
		Api: api.New(appId, appSecret),
	}
}
