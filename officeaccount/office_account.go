package officeaccount

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/hdget/common/types"
	"github.com/hdget/lib-wechat/officeaccount/cache"
	"github.com/hdget/lib-wechat/pkg/wxapi/officeaccount"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type API interface {
	CreateJsSdkSignature(url string) (*JsSdkSignatureResult, error)           // 生成jsapi_ticket签名
	VerifyServer(token string, queryParams map[string]string) (string, error) // 校验微信公众号服务器
	HandleMessage(data []byte) ([]byte, error)                                // 接收普通消息,接收事件消息以及被动回复消息
	SendTemplateMessage(toUser string, msg *TemplateMessage) error            // 发送模板消息
	//GetUserInfo(openid string) (*api.UserInfoResult, error)         // 获取用户信息
}

type wxoaApiImpl struct {
	wxapi officeaccount.WxAPI
	cache cache.Cache
}

func New(appId, appSecret string, redisProvider types.RedisProvider) API {
	return &wxoaApiImpl{
		wxapi: officeaccount.New(appId, appSecret),
		cache: cache.New(appId, redisProvider),
	}
}

// VerifyServer 公众号接入时校验
// 参考： https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Access_Overview.html
func (impl *wxoaApiImpl) VerifyServer(token string, queryParams map[string]string) (string, error) {
	signature, timestamp, nonce, echostr := queryParams["signature"], queryParams["timestamp"], queryParams["nonce"], queryParams["echostr"]
	if signature == "" || timestamp == "" || nonce == "" {
		return "", fmt.Errorf("invalid request, urlParams: %v", queryParams)
	}

	si := []string{token, timestamp, nonce}
	sort.Strings(si)              //字典序排序
	str := strings.Join(si, "")   //组合字符串
	s := sha1.New()               //返回一个新的使用SHA1校验的hash.Hash接口
	_, _ = io.WriteString(s, str) //WriteString函数将字符串数组str中的内容写入到s中
	calculatedSignature := fmt.Sprintf("%x", s.Sum(nil))

	if signature != calculatedSignature {
		return "", errors.New("signature not matched")
	}

	return echostr, nil
}

func (impl *wxoaApiImpl) SendTemplateMessage(toUser string, msg *TemplateMessage) error {
	accessToken, err := impl.getAccessToken()
	if err != nil {
		return err
	}

	var contents map[string]string
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &contents,
		TagName: "json", // 使用json标签作为key
	})
	if err != nil {
		return err
	}

	msg.ToUser = toUser
	err = decoder.Decode(msg)
	if err != nil {
		return err
	}

	return impl.wxapi.SendTemplateMessage(accessToken, contents)
}

func (impl *wxoaApiImpl) getAccessToken() (string, error) {
	accessToken, _ := impl.cache.GetAccessToken()
	if accessToken == "" {
		result, err := impl.wxapi.GetAccessToken()
		if err != nil {
			return "", err
		}

		err = impl.cache.SetAccessToken(result.AccessToken, result.ExpiresIn)
		if err != nil {
			return "", err
		}

		accessToken = result.AccessToken
	}

	return accessToken, nil
}
