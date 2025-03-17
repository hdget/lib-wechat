package wxoa

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
	"github.com/pkg/errors"
)

type UserInfoResult struct {
	api.ApiResult
	Subscribe      int8   `json:"subscribe"`
	Openid         string `json:"openid"`
	Language       string `json:"language"`
	SubscribeTime  int64  `json:"subscribe_time"`
	UnionId        string `json:"unionid"`
	Remark         string `json:"remark"`
	GroupId        int    `json:"groupid"`
	TagIdList      []int  `json:"tagid_list"`
	SubscribeScene string `json:"subscribe_scene"`
	QrScene        int    `json:"qr_scene"`
	QrSceneStr     string `json:"qr_scene_str"`
}

const (
	//https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId
	urlGetUnionId = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"
)

// GetUserInfo 通过openId获取用户信息
func (impl *wxoaImpl) GetUserInfo(openid string) (*UserInfoResult, error) {
	accessToken, err := impl.GetAccessToken()
	if err != nil {
		return nil, errors.Wrap(err, "get access token")
	}

	url := fmt.Sprintf(urlGetUnionId, accessToken, openid)
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "get wxoa user info, appId: %s", impl.AppId)
	}

	var result UserInfoResult
	err = impl.ParseApiResult(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if result.Openid == "" {
		return nil, fmt.Errorf("invalid userinfo result, result: %v", result)
	}

	return &result, nil
}
