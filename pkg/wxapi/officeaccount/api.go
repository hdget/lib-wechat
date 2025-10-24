package officeaccount

import (
	"github.com/hdget/lib-wechat/pkg/wxapi"
)

type WxAPI interface {
	wxapi.API
	GetJsSdkSignature(ticket, url string) (*GetJsSdkSignatureResult, error)
	GetJsSdkTicket(accessToken string) (*GetJsSdkTicketResult, error)
	SendTemplateMessage(accessToken string, contents map[string]string) error
}

type officeAccountWxApiImpl struct {
	wxapi.API
}

func New(appId, appSecret string) WxAPI {
	return &officeAccountWxApiImpl{
		API: wxapi.New(appId, appSecret),
	}
}
