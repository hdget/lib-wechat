package miniprogram

import "github.com/hdget/lib-wechat/pkg/wxapi"

type WxAPI interface {
	wxapi.API
	Code2Session(code string) (*GetSessionResult, error)                               // 小程序静默登录，通过code换取UnionId
	GetUserPhoneNumber(accessToken, code string) (string, error)                       // 获取用户手机号码
	CreateLimitedWxaCode(accessToken, path string, width int) ([]byte, error)          // 生成有限的小程序码
	CreateUnLimitedWxaCode(accessToken, scene, page string, width int) ([]byte, error) // CreateUnLimitedWxaCode 生成小程序码，可接受页面参数较短，生成个数不受限
}

type miniProgramWxApiImpl struct {
	wxapi.API
}

func New(appId, appSecret string) WxAPI {
	return &miniProgramWxApiImpl{
		API: wxapi.New(appId, appSecret),
	}
}
