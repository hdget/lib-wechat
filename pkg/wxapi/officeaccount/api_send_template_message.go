package officeaccount

import (
	"fmt"

	"github.com/hdget/lib-wechat/pkg/wxapi"
	"github.com/pkg/errors"
)

type sendTemplateMessageResult struct {
	*wxapi.Result
	Msgid int `json:"msgid"`
}

const (
	// 参考：https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
	urlSendTemplateMessage = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
)

// SendTemplateMessage 发送模板消息
func (impl officeAccountWxApiImpl) SendTemplateMessage(accessToken string, content map[string]string) error {
	url := fmt.Sprintf(urlSendTemplateMessage, accessToken)

	ret, err := wxapi.Post[sendTemplateMessageResult](url, content)
	if err != nil {
		return errors.Wrap(err, "send template message")
	}

	if err = wxapi.CheckResult(ret.Result, url, content); err != nil {
		return errors.Wrap(err, "send template message")
	}

	return nil
}
