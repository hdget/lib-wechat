package wxmp

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
	"github.com/pkg/errors"
)

type UserInfo struct {
	OpenID    string `json:"openId"`
	UnionID   string `json:"unionId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
	Language  string `json:"language"`
	Watermark struct {
		Timestamp int64  `json:"timestamp"`
		AppId     string `json:"appid"`
	} `json:"watermark"`
}

type MobileInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	Watermark       struct {
		AppId     string      `json:"appid"`
		Timestamp interface{} `json:"timestamp"`
	} `json:"watermark"`
}

// sessionResult wechat miniprogram login session
type sessionResult struct {
	api.ApiResult
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
	UnionId    string `json:"unionid"`
}

type GetUserPhoneNumberResult struct {
	api.ApiResult
	PhoneInfo MobileInfo `json:"phone_info"`
}

const (
	redisKeySessionKey    = "session_key"
	wxmpSessionKeyExpires = 3600 // session key过期时间3600秒
	urlCode2Session       = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	urlGetUserPhoneNumber = "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s"
)

// Login 微信小程序登录，获取openid,unionid
func (impl *wxmpImpl) Login(code string) (string, string, error) {
	result, err := impl.code2session(code)
	if err != nil {
		return "", "", err
	}
	return result.OpenId, result.UnionId, nil
}

func (impl *wxmpImpl) DecryptUserInfo(encryptedData, iv string) (*UserInfo, error) {
	if impl.Cache == nil {
		return nil, errors.New("no redis provider")
	}

	sessKey, err := impl.Cache.Get(redisKeySessionKey)
	if err != nil {
		return nil, errors.Wrap(err, "session key not found, you should invoke wx.login() firstly")
	}

	cipherText, err := decrypt(sessKey, encryptedData, iv)
	if err != nil {
		return nil, errors.Wrap(err, "decrypt encrypted data")
	}

	var userInfo UserInfo
	err = json.Unmarshal(cipherText, &userInfo)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal to UserInfo")
	}

	if userInfo.Watermark.AppId != impl.AppId {
		return nil, errAppIDNotMatch
	}

	return &userInfo, nil
}

// DecryptMobileInfo 基础库2.21.2以前，换取手机号信息的方式需要先wx.login(),然后解密信息来获取手机号码
func (impl *wxmpImpl) DecryptMobileInfo(encryptedData, iv string) (*MobileInfo, error) {
	if impl.Cache == nil {
		return nil, errors.New("no redis provider")
	}

	sessKey, err := impl.Cache.Get(redisKeySessionKey)
	if err != nil {
		return nil, errors.Wrap(err, "session key not found, you should invoke wx.login() firstly")
	}

	cipherText, err := decrypt(sessKey, encryptedData, iv)
	if err != nil {
		return nil, errors.Wrap(err, "decrypt encrypted data")
	}

	var mobileInfo MobileInfo
	err = json.Unmarshal(cipherText, &mobileInfo)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal to UserInfo")
	}

	return &mobileInfo, nil
}

// GetUserPhoneNumber 从基础库2.21.2开始，换取手机号信息的方式进行了安全升级，新方式不再需要提前调用wx.login进行登录。
func (impl *wxmpImpl) GetUserPhoneNumber(code string) (*MobileInfo, error) {
	accessToken, err := impl.GetAccessToken()
	if err != nil {
		return nil, errors.Wrap(err, "get access token")
	}

	body := struct {
		Code string `json:"code"`
	}{
		Code: code,
	}

	resp, err := resty.New().R().SetBody(body).Post(fmt.Sprintf(urlGetUserPhoneNumber, accessToken))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get user phone number, status_code: %d", resp.StatusCode())
	}

	var result GetUserPhoneNumberResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("%s, url: %s", result.ErrMsg, urlGetUserPhoneNumber)
	}

	return &result.PhoneInfo, nil
}

func (impl *wxmpImpl) code2session(code string) (*sessionResult, error) {
	if impl.Cache == nil {
		return nil, errors.New("no redis provider")
	}

	// url to get sessionKey, openId and unionId from Weixin server
	// do http get request against Wechat server
	url := fmt.Sprintf(urlCode2Session, impl.AppId, impl.AppSecret, code)

	// 登录凭证校验
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "wxmp code to result")
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("wxmp code to session, status_code: %d", resp.StatusCode())
	}

	var result sessionResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("%s, url: %s", result.ErrMsg, url)
	}

	if result.SessionKey == "" {
		return nil, errors.New("empty result key")
	}

	// 保存到缓存中
	err = impl.Cache.Set(redisKeySessionKey, result.SessionKey, wxmpSessionKeyExpires)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
