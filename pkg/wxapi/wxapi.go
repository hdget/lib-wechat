package wxapi

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

type API interface {
	GetAccessToken() (*GetAccessTokenResult, error)
	GetAppId() string
	GetAppSecret() string
}

type Result struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type wxapiImpl struct {
	AppId     string
	AppSecret string
}

const (
	networkTimeout    = 3 * time.Second
	urlGetAccessToken = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

func New(appId, appSecret string) API {
	return &wxapiImpl{
		AppId:     appId,
		AppSecret: appSecret,
	}
}

func (impl wxapiImpl) GetAppId() string {
	return impl.AppId
}

func (impl wxapiImpl) GetAppSecret() string {
	return impl.AppSecret
}

func (impl wxapiImpl) GetAccessToken() (*GetAccessTokenResult, error) {
	url := fmt.Sprintf(urlGetAccessToken, impl.AppId, impl.AppSecret)

	ret, err := Get[GetAccessTokenResult](url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "get access token, appId: %s", impl.AppId)
	}

	if err = CheckResult(ret.Result, url, nil); err != nil {
		return nil, errors.Wrapf(err, "get access token, appId: %s", impl.AppId)
	}

	return ret, nil
}

func Get[RESULT any](url string, request ...any) (*RESULT, error) {
	var req any
	if len(request) > 0 {
		req = request[0]
	}

	resp, err := resty.New().SetTimeout(networkTimeout).
		SetHeader("Content-Type", "application/json; charset=UTF-8").
		R().SetBody(req).Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "http get request, url: %s, req: %v", url, req)
	}

	var ret RESULT
	err = json.Unmarshal(resp.Body(), &ret)
	if err != nil {
		return nil, errors.Wrapf(err, "parse result, url: %s, req: %v, ret: %s", url, req, convert.BytesToString(resp.Body()))
	}

	return &ret, nil
}

func Post[RESULT any](url string, request ...any) (*RESULT, error) {
	var req any
	if len(request) > 0 {
		req = request[0]
	}

	resp, err := resty.New().SetTimeout(networkTimeout).
		SetHeader("Content-Type", "application/json; charset=UTF-8").
		R().SetBody(req).Post(url)
	if err != nil {
		return nil, errors.Wrapf(err, "http post request, url: %s, req: %v", url, req)
	}

	var ret RESULT
	err = json.Unmarshal(resp.Body(), &ret)
	if err != nil {
		return nil, errors.Wrapf(err, "parse result, url: %s, req: %v, ret: %s", url, req, convert.BytesToString(resp.Body()))
	}

	return &ret, nil
}

// PostResponse http post and get response
func PostResponse(url string, request ...any) ([]byte, error) {
	var req any
	if len(request) > 0 {
		req = request[0]
	}

	resp, err := resty.New().R().Post(url)
	if err != nil {
		return nil, errors.Wrapf(err, "http post request, url: %s, req: %v", url, req)
	}
	return resp.Body(), nil
}

func CheckResult(result *Result, url string, request ...any) error {
	var req any
	if len(request) > 0 {
		req = request[0]
	}

	if result == nil {
		return errors.New("result is nil")
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("wxapi call error, url: %s, req: %v, ret: %v", url, req, result)
	}

	return nil
}
