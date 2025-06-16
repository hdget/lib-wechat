package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib/lib-wechat/api"
	"github.com/pkg/errors"
	"time"
)

type sendTemplateMessageResult struct {
	api.Result
	Msgid int `json:"msgid"`
}

const (
	// 参考：https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
	urlSendTemplateMessage = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
	networkTimeout         = 3 * time.Second
)

// ApiSendTemplateMessage 发送模板消息
func (impl apiImpl) SendTemplateMessage(accessToken string, contents map[string]string) error {
	url := fmt.Sprintf(urlSendTemplateMessage, accessToken)
	resp, err := resty.New().SetTimeout(networkTimeout).SetHeader("Content-Type", "application/json; charset=UTF-8").R().SetBody(contents).Post(url)
	if err != nil {
		return errors.Wrapf(err, "wxoa send template message")
	}

	var result sendTemplateMessageResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return err
	}

	if err = impl.CheckResult(result.Result, url, contents); err != nil {
		return err
	}

	return nil
}
