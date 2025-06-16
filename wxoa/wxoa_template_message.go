package wxoa

type TemplateMessage struct {
	// 必要参数
	ToUser     string `json:"touser"`      // 发送给哪个openId
	TemplateId string `json:"template_id"` // 模板ID
	Data       any    `json:"data"`
	// 非必须的参数
	Url         string                      `json:"url"`           // 跳转链接
	MiniProgram *templateMessageMiniProgram `json:"miniprogram"`   // 跳小程序所需数据
	ClientMsgId string                      `json:"client_msg_id"` // 防重入id。对于同一个openid + client_msg_id, 只发送一条消息,10分钟有效,超过10分钟不保证效果。若无防重入需求，可不填
}

type templateMessageMiniProgram struct {
	AppId    string `json:"appid"`
	PagePath string `json:"pagepath"`
}

type templateMessageRow struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

const (
	defaultColor = "#173177"
)

func NewTemplateMessage(templateId string, contents map[string]string) *TemplateMessage {
	rows := make(map[string]*templateMessageRow)
	for k, v := range contents {
		rows[k] = &templateMessageRow{
			Value: v,
			Color: defaultColor,
		}
	}

	return &TemplateMessage{
		TemplateId: templateId,
		Data:       rows,
	}
}

// LinkUrl 链接外部URL
func (m *TemplateMessage) LinkUrl(url string) *TemplateMessage {
	m.Url = url
	return m
}

// LinkMiniProgram 链接小程序
func (m *TemplateMessage) LinkMiniProgram(appId, pagePath string) *TemplateMessage {
	m.MiniProgram = &templateMessageMiniProgram{
		AppId:    appId,
		PagePath: pagePath,
	}
	return m
}
