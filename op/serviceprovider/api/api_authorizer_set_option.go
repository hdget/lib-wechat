package api

import (
	"encoding/json"
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

const (
	// 授权账号管理 /设置授权方选项信息 限制：1000次/天/平台
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/authorization-management/setAuthorizerOptionInfo.html
	urlSetAuthorizerOption = "https://api.weixin.qq.com/cgi-bin/component/set_authorizer_option?access_token=%s"
)

func (impl apiImpl) SetAuthorizerOption(authorizerAccessToken string, optionName string, optionValue string) error {
	validOptionNames := []string{"location_report", "voice_recognize", "customer_service"}
	validOptionValues := []string{"0", "1"}
	switch optionName {
	case "location_report":
		validOptionValues = []string{"0", "1", "2"}
	}

	if !pie.Contains(validOptionNames, optionName) {
		return fmt.Errorf("option name not supported, optionName: %s, valid: %v", optionName, validOptionNames)
	}

	if !pie.Contains(validOptionValues, optionValue) {
		return fmt.Errorf("option value not supported, optionValue: %s, valid: %v", optionValue, validOptionValues)
	}

	req := &setAuthorizerOptionRequest{
		OptionName:  optionName,
		OptionValue: optionValue,
	}

	url := fmt.Sprintf(urlSetAuthorizerOption, authorizerAccessToken)
	resp, err := resty.New().R().SetBody(req).Post(url)
	if err != nil {
		return err
	}

	var result api.Result
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return errors.Wrapf(err, "unmarshal set authorizer option result, data: %s", convert.BytesToString(resp.Body()))
	}

	if err = impl.CheckResult(result, url, req); err != nil {
		return err
	}

	return nil
}
