package api

import (
	"encoding/json"
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib/lib-wechat/api"
	"github.com/hdget/utils/convert"
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
	api.Result
	OptionName  string `json:"option_name"`
	OptionValue string `json:"option_value"`
}

const (
	// 授权账号管理 /获取授权方选项信息
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/authorization-management/getAuthorizerOptionInfo.html
	urlGetAuthorizerOption = "https://api.weixin.qq.com/cgi-bin/component/get_authorizer_option?access_token=%s"
)

func (impl apiImpl) GetAuthorizerOption(authorizerAccessToken string, optionName string) (string, error) {
	validOptionNames := []string{"location_report", "voice_recognize", "customer_service"}
	if !pie.Contains(validOptionNames, optionName) {
		return "", fmt.Errorf("option name not supported, optionName: %s, valid: %v", optionName, validOptionNames)
	}

	req := &getAuthorizerOptionRequest{
		OptionName: optionName,
	}

	url := fmt.Sprintf(urlGetAuthorizerOption, authorizerAccessToken)
	resp, err := resty.New().R().SetBody(req).Post(url)
	if err != nil {
		return "", err
	}

	var result getAuthorizerOptionResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", errors.Wrapf(err, "unmarshal get authorizer option result, data: %s", convert.BytesToString(resp.Body()))
	}

	if err = impl.CheckResult(result.Result, url, req); err != nil {
		return "", err
	}

	return result.OptionValue, nil
}
