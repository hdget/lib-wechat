package provider

import (
	"fmt"

	"github.com/elliotchance/pie/v2"
	"github.com/hdget/lib-wechat/pkg/wxapi"
	"github.com/pkg/errors"
)

type setAuthorizerOptionRequest struct {
	OptionName  string `json:"option_name"`
	OptionValue string `json:"option_value"`
}

type getAuthorizerOptionRequest struct {
	OptionName string `json:"option_name"`
}

type getAuthorizerOptionResult struct {
	*wxapi.Result
	OptionName  string `json:"option_name"`
	OptionValue string `json:"option_value"`
}

const (
	// 授权账号管理 /获取授权方选项信息
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/authorization-management/getAuthorizerOptionInfo.html
	urlGetAuthorizerOption = "https://api.weixin.qq.com/cgi-bin/component/get_authorizer_option?access_token=%s"
)

func (impl serviceProviderWxApiImpl) GetAuthorizerOption(authorizerAccessToken string, optionName string) (string, error) {
	validOptionNames := []string{"location_report", "voice_recognize", "customer_service"}
	if !pie.Contains(validOptionNames, optionName) {
		return "", fmt.Errorf("option name not supported, optionName: %s, valid: %v", optionName, validOptionNames)
	}

	req := &getAuthorizerOptionRequest{
		OptionName: optionName,
	}

	url := fmt.Sprintf(urlGetAuthorizerOption, authorizerAccessToken)

	ret, err := wxapi.Post[getAuthorizerOptionResult](url, req)
	if err != nil {
		return "", errors.Wrap(err, "get authorizer option")
	}

	if err = wxapi.CheckResult(ret.Result, url, req); err != nil {
		return "", errors.Wrap(err, "get authorizer option")
	}

	return ret.OptionValue, nil
}
