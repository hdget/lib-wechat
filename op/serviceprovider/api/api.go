package api

import "github.com/hdget/lib-wechat/api"

type Api interface {
	api.Api
	ApiAuthorizer
	ApiCredential
	ApiWxa
}

type ApiAuthorizer interface {
	GetAuthorizerInfo(componentAccessToken, authorizerAppid string) (*GetAuthorizerInfoResult, error) // 授权账号管理 /获取授权账号详情
	GetAuthorizerOption(appId string, optionName string) (string, error)                              // 授权账号管理 /获取授权方选项信息
	SetAuthorizerOption(authorizerAccessToken string, optionName string, optionValue string) error    // 授权账号管理 /设置授权方选项信息
}

type ApiCredential interface {
	StartPushComponentVerifyTicket() error                                                                                                         // 第三方平台调用凭证/启动票据推送服务
	CreatePreAuthCode(componentAccessToken string) (string, error)                                                                                 // 第三方平台调用凭证/获取预授权码
	GetAuthorizerAccessToken(componentAccessToken string, authorizerAppid, authorizerRefreshToken string) (*GetAuthorizerAccessTokenResult, error) // 第三方平台调用凭证/获取授权账号调用令牌
	QueryAuthorizationInfo(componentAccessToken, authCode string) (*AuthorizationInfo, error)                                                      // 第三方平台调用凭证/获取刷新令牌
	GetComponentAccessToken(componentVerifyTicket string) (*GetComponentAccessTokenResult, error)                                                  // 第三方平台调用凭证 /获取令牌
}

type ApiWxa interface {
	WxaCode2Session(componentAppId, componentAccessToken string, appId, code string) (*WxaCode2SessionResult, error) // 小程序登录
}

type apiImpl struct {
	api.Api
}

func New(appId, appSecret string) Api {
	return &apiImpl{
		Api: api.New(appId, appSecret),
	}
}
