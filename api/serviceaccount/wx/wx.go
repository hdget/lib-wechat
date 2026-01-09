package wx

type WxApi interface {
	SendTemplateMessage(accessToken string, msg *TemplateMessage) error // 发送模板消息
}

type wxApiImpl struct {
	appId     string
	appSecret string
}

func New(appId, appSecret string) WxApi {
	return &wxApiImpl{
		appId:     appId,
		appSecret: appSecret,
	}
}
