package wxoa

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

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

type templateMessageSendResult struct {
	api.ApiResult
	Msgid int `json:"msgid"`
}

const (
	urlSendTemplateMessage = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
	defaultColor           = "#173177"
	networkTimeout         = 3 * time.Second
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

// SendTemplateMessage 发送模板消息
func (impl *wxoaImpl) SendTemplateMessage(toUser string, m *TemplateMessage) error {
	if m.TemplateId == "" {
		return errors.New("invalid template message")
	}

	accessToken, err := impl.GetAccessToken()
	if err != nil {
		return err
	}

	// send to whom
	m.ToUser = toUser

	// new http client
	httpClient := resty.New()
	httpClient.SetTimeout(networkTimeout)

	url := fmt.Sprintf(urlSendTemplateMessage, accessToken)
	resp, err := httpClient.SetHeader("Content-Type", "application/json; charset=UTF-8").R().SetBody(m).Post(url)
	if err != nil {
		msgContent, _ := json.Marshal(m)
		return errors.Wrapf(err, "wxoa send template message, message: %s", msgContent)
	}

	if resp.StatusCode() != http.StatusOK {
		msgContent, _ := json.Marshal(m)
		return errors.Wrapf(err, "wxoa send template message, message: %s, statusCode: %d", msgContent, resp.StatusCode())
	}

	var result templateMessageSendResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return err
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("%s, url: %s", result.ErrMsg, url)
	}

	return nil
}
