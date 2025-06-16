package api

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib/lib-wechat/api"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

type startPushComponentTicketRequest struct {
	ComponentAppid     string `json:"component_appid"`
	ComponentAppSecret string `json:"component_secret"`
}

const (
	// 第三方平台调用凭证 /启动票据推送服务
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/startPushTicket.html
	urlStartPushTicket = "https://api.weixin.qq.com/cgi-bin/component/api_start_push_ticket"
)

func (impl apiImpl) StartPushComponentVerifyTicket() error {
	req := &startPushComponentTicketRequest{
		ComponentAppid:     impl.GetAppId(),
		ComponentAppSecret: impl.GetAppSecret(),
	}

	resp, err := resty.New().R().SetBody(req).Post(urlStartPushTicket)
	if err != nil {
		return err
	}

	var result api.Result
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return errors.Wrapf(err, "unmarshal start push component ticket result, data: %s", convert.BytesToString(resp.Body()))
	}

	if err = impl.CheckResult(result, urlStartPushTicket, req); err != nil {
		return err
	}

	return nil
}
