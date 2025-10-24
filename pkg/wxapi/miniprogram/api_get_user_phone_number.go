package miniprogram

import (
	"fmt"

	"github.com/hdget/lib-wechat/pkg/wxapi"
	"github.com/pkg/errors"
)

type getUserPhoneNumberRequest struct {
	Code string `json:"code"`
}

type getUserPhoneNumberResult struct {
	*wxapi.Result
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
func (impl miniProgramWxApiImpl) GetUserPhoneNumber(accessToken, code string) (string, error) {
	req := &getUserPhoneNumberRequest{
		Code: code,
	}

	url := fmt.Sprintf(urlGetUserPhoneNumber, accessToken)

	ret, err := wxapi.Post[getUserPhoneNumberResult](url, req)
	if err != nil {
		return "", errors.Wrap(err, "get user phone number")
	}

	if err = wxapi.CheckResult(ret.Result, url, req); err != nil {
		return "", errors.Wrap(err, "get user phone number")
	}

	return ret.PhoneInfo.PhoneNumber, nil
}
