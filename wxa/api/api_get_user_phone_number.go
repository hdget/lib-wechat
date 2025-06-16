package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
)

type getUserPhoneNumberRequest struct {
	Code string `json:"code"`
}

type getUserPhoneNumberResult struct {
	api.Result
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
		Watermark       struct {
			AppId     string      `json:"appid"`
			Timestamp interface{} `json:"timestamp"`
		} `json:"watermark"`
	} `json:"phone_info"`
}

const (
	// GetUserPhoneNumber 快速手机号验证
	// 参考: https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-info/phone-number/getPhoneNumber.html
	urlGetUserPhoneNumber = "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s"
)

// GetUserPhoneNumber 通过code获取用户的手机号码
func (impl apiImpl) GetUserPhoneNumber(accessToken, code string) (string, error) {
	req := &getUserPhoneNumberRequest{
		Code: code,
	}

	url := fmt.Sprintf(urlGetUserPhoneNumber, accessToken)
	resp, err := resty.New().R().SetBody(req).Post(url)
	if err != nil {
		return "", err
	}

	var result getUserPhoneNumberResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}

	if err = impl.CheckResult(result.Result, url, req); err != nil {
		return "", err
	}

	return result.PhoneInfo.PhoneNumber, nil
}
