package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib/lib-wechat/api"
	"github.com/pkg/errors"
)

type UserInfoResult struct {
	api.Result
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
	// 参考：https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId
	urlGetUnionId = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"
)

// GetUserInfo 通过openId获取用户信息
func (impl apiImpl) GetUserInfo(accessToken, openid string) (*UserInfoResult, error) {
	url := fmt.Sprintf(urlGetUnionId, accessToken, openid)
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "get wxoa user info, appId: %s", impl.GetAppId())
	}

	var result UserInfoResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if err = impl.CheckResult(result.Result, url, nil); err != nil {
		return nil, err
	}

	return &result, nil
}
