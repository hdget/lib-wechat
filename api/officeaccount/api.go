package officeaccount

import (
	"github.com/hdget/lib-wechat/api"
	"github.com/hdget/lib-wechat/api/officeaccount/wx"
)

type API interface {
	api.Api
	GetJsSdkSignature(ticket, url string) (*wx.GetJsSdkSignatureResult, error)
	GetJsSdkTicket(accessToken string) (*wx.GetJsSdkTicketResult, error)
}

type officeAccountApiImpl struct {
	api.Api
	wx.WxApi
}

func New(appId, appSecret string) API {
	return &officeAccountApiImpl{
		Api:   api.New(appId, appSecret),
		WxApi: wx.New(appId, appSecret),
	}
}

func (impl officeAccountApiImpl) GetJsSdkSignature(ticket, url string) (*wx.GetJsSdkSignatureResult, error) {
	return impl.WxApi.GetJsSdkSignature(ticket, url)
}

func (impl officeAccountApiImpl) GetJsSdkTicket(accessToken string) (*wx.GetJsSdkTicketResult, error) {
	return impl.WxApi.GetJsSdkTicket(accessToken)
}
