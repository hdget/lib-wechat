package wxoa

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/wxoa/types"
	"github.com/pkg/errors"
)

const (
	//https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId
	urlGetUnionId = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"
)

// GetUserInfo 通过openId获取用户信息
func (impl *wxoaImpl) GetUserInfo(openid string) (*types.UserInfoResult, error) {
	accessToken, err := impl.GetAccessToken()
	if err != nil {
		return nil, errors.Wrap(err, "get access token")
	}

	url := fmt.Sprintf(urlGetUnionId, accessToken, openid)
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "get wxoa user info, appId: %s", impl.AppId)
	}

	var result types.UserInfoResult
	err = impl.ParseApiResult(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if result.Openid == "" {
		return nil, fmt.Errorf("invalid userinfo result, result: %v", result)
	}

	return &result, nil
}
