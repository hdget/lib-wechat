package api

import "github.com/hdget/lib/lib-wechat/api"

type Api interface {
	api.Api
	GetAccessToken() (*GetAccessTokenResult, error)
	GetJsSdkSignature(ticket, url string) (*GetJsSdkSignatureResult, error)
	GetJsSdkTicket(accessToken string) (*GetJsSdkTicketResult, error)
	SendTemplateMessage(accessToken string, contents map[string]string) error
}

type apiImpl struct {
	api.Api
}

func New(appId, appSecret string) Api {
	return &apiImpl{
		Api: api.New(appId, appSecret),
	}
}
