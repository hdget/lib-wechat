package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib/lib-wechat/api"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

type getAuthorizerInfoRequest struct {
	ComponentAppid  string `json:"component_appid"`
	AuthorizerAppid string `json:"authorizer_appid"`
}

type GetAuthorizerInfoResult struct {
	api.Result
	Authorizer    *AppInfo `json:"authorizer_info"`
	Authorization struct {
		RefreshToken string      `json:"authorizer_refresh_token"`
		FuncInfo     []*FuncInfo `json:"func_info"`
	} `json:"authorization_info"`
}

type AppInfo struct {
	NickName        string `json:"nick_name"`
	HeadImg         string `json:"head_img"`
	ServiceTypeInfo struct {
		Id int `json:"id"`
	} `json:"service_type_info"`
	VerifyTypeInfo struct {
		Id int `json:"id"`
	} `json:"verify_type_info"`
	UserName     string `json:"user_name"`
	Alias        string `json:"alias"`
	QrcodeUrl    string `json:"qrcode_url"`
	BusinessInfo struct {
		OpenPay   int `json:"open_pay"`
		OpenShake int `json:"open_shake"`
		OpenScan  int `json:"open_scan"`
		OpenCard  int `json:"open_card"`
		OpenStore int `json:"open_store"`
	} `json:"business_info"`
	Idc             int    `json:"idc"`
	PrincipalName   string `json:"principal_name"`
	Signature       string `json:"signature"`
	MiniProgramInfo struct {
		Network struct {
			RequestDomain   []string      `json:"RequestDomain"`
			WsRequestDomain []interface{} `json:"WsRequestDomain"`
			UploadDomain    []interface{} `json:"UploadDomain"`
			DownloadDomain  []interface{} `json:"DownloadDomain"`
			BizDomain       []interface{} `json:"BizDomain"`
			UDPDomain       []interface{} `json:"UDPDomain"`
		} `json:"network"`
		Categories []struct {
			First  string `json:"first"`
			Second string `json:"second"`
		} `json:"categories"`
		VisitStatus int `json:"visit_status"`
	} `json:"MiniProgramInfo"`
	BasicConfig struct {
		IsPhoneConfigured bool `json:"is_phone_configured"` // 小程序注册方式
		IsEmailConfigured bool `json:"is_email_configured"` // 小程序注册方式
	} `json:"basic_config"`
	RegisterType  int `json:"register_type"`  // 小程序注册方式
	AccountStatus int `json:"account_status"` // 帐号状态，该字段小程序也返回
	ChannelsInfo  int `json:"channels_info"`  // 视频号账号类型；如果该授权账号为视频号则返回该字段
}

type FuncInfo struct {
	FuncscopeCategory struct {
		Id int `json:"id"`
	} `json:"funcscope_category"`
}

const (
	// 授权账号管理 /获取授权账号详情 限制：1000次/天/授权方
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/authorization-management/getAuthorizerInfo.html
	urlGetAuthorizerInfo = "https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info?access_token=%s"
)

// GetAuthorizerInfo 获取授权方的基本信息，包括头像、昵称、账号类型、认证类型、原始ID等信息
func (impl apiImpl) GetAuthorizerInfo(componentAccessToken, authorizerAppid string) (*GetAuthorizerInfoResult, error) {
	req := &getAuthorizerInfoRequest{
		ComponentAppid:  impl.GetAppId(),
		AuthorizerAppid: authorizerAppid,
	}

	url := fmt.Sprintf(urlGetAuthorizerInfo, componentAccessToken)
	resp, err := resty.New().R().SetBody(req).Post(url)
	if err != nil {
		return nil, err
	}

	var result GetAuthorizerInfoResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal authorizer info, data: %s", convert.BytesToString(resp.Body()))
	}

	if err = impl.CheckResult(result.Result, url, req); err != nil {
		return nil, err
	}

	return &result, nil
}
