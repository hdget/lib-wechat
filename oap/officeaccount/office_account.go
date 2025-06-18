package officeaccount

import (
	"crypto/sha1"
	"fmt"
	"github.com/hdget/common/intf"
	"github.com/hdget/lib-wechat/oap/officeaccount/api"
	"github.com/hdget/lib-wechat/oap/officeaccount/cache"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"io"
	"sort"
	"strings"
)

type Lib interface {
	CreateJsSdkSignature(url string) (*JsSdkSignatureResult, error)           // 生成jsapi_ticket签名
	VerifyServer(token string, queryParams map[string]string) (string, error) // 校验微信公众号服务器
	HandleMessage(data []byte) ([]byte, error)                                // 接收普通消息,接收事件消息以及被动回复消息
	SendTemplateMessage(toUser string, msg *TemplateMessage) error            // 发送模板消息
	//GetUserInfo(openid string) (*api.UserInfoResult, error)         // 获取用户信息
}

type wxoaImpl struct {
	api   api.Api
	cache cache.Cache
}

func New(appId, appSecret string, redisProvider intf.RedisProvider) Lib {
	return &wxoaImpl{
		api:   api.New(appId, appSecret),
		cache: cache.New(appId, redisProvider),
	}
}

// VerifyServer 公众号接入时校验
// 参考： https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Access_Overview.html
func (impl *wxoaImpl) VerifyServer(token string, queryParams map[string]string) (string, error) {
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

func (impl *wxoaImpl) SendTemplateMessage(toUser string, msg *TemplateMessage) error {
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

	return impl.api.SendTemplateMessage(accessToken, contents)
}

func (impl *wxoaImpl) getAccessToken() (string, error) {
	accessToken, _ := impl.cache.GetAccessToken()
	if accessToken == "" {
		result, err := impl.api.GetAccessToken()
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
