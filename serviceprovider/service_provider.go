package serviceprovider

import (
	"fmt"
	"strings"
	"time"

	"github.com/hdget/common/types"
	"github.com/hdget/lib-wechat/pkg/wxapi/open/provider"
	"github.com/hdget/lib-wechat/serviceprovider/cache"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type API interface {
	//Cache() cache.Cache
	GetAuthUrl(client, redirectUrl string, authType int) (string, error)       // 获取授权链接
	GetAuthorizerAppId(authCode string) (string, error)                        // 通过authCode获取授权应用的appId
	GetAuthorizerInfo(appId string) (*provider.GetAuthorizerInfoResult, error) // 获取授权应用的信息
	HandleEvent(event *Event, handlers map[EventKind]EventHandler) error       // 处理授权事件
}

type serviceProviderImpl struct {
	wxapi provider.WxAPI
	cache cache.Cache
}

const (
	urlPCAuth = "https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s&auth_type=%d"
	urlH5Auth = "https://open.weixin.qq.com/wxaopen/safe/bindcomponent?action=bindcomponent&no_scan=1&component_appid=%s&pre_auth_code=%s&redirect_uri=%s&auth_type=%d#wechat_redirect"
)

func New(appId, appSecret string, redisProvider types.RedisProvider) API {
	return &serviceProviderImpl{
		wxapi: provider.New(appId, appSecret),
		cache: cache.New(appId, redisProvider),
	}
}

func (impl serviceProviderImpl) GetAuthUrl(client, redirectUrl string, authType int) (string, error) {
	componentAccessToken, err := impl.getComponentAccessToken()
	if err != nil {
		return "", err
	}

	preAuthCode, err := impl.wxapi.CreatePreAuthCode(componentAccessToken)
	if err != nil {
		return "", err
	}

	// 校验authCode
	switch cast.ToInt(authType) {
	case 1, 2, 3, 4, 5, 6:
	default:
		return "", fmt.Errorf("invalid auth type: %d", authType)
	}

	// 校验redirectUrl, https://xxx
	if !strings.HasPrefix(redirectUrl, "https") {
		return "", fmt.Errorf("invalid redirect url, redirectUrl: %s", redirectUrl)
	}

	switch strings.ToLower(client) {
	case "pc":
		return fmt.Sprintf(urlPCAuth, impl.wxapi.GetAppId(), preAuthCode, redirectUrl, authType), nil
	case "h5":
		return fmt.Sprintf(urlH5Auth, impl.wxapi.GetAppId(), preAuthCode, redirectUrl, authType), nil
	default:
		return "", fmt.Errorf("unsupported client, client: %s", client)
	}
}

func (impl serviceProviderImpl) getComponentAccessToken() (string, error) {
	componentAccessToken, _ := impl.cache.GetComponentAccessToken()
	if componentAccessToken == "" { // 缓存取不到则通过API接口获取并缓存起来
		componentVerifyTicket, err := impl.getComponentVerifyTicket()
		if err != nil {
			return "", err
		}

		result, err := impl.wxapi.GetComponentAccessToken(componentVerifyTicket)
		if err != nil {
			return "", errors.Wrap(err, "retrieve component access Token")
		}

		// 过期前十分钟过期
		err = impl.cache.SetComponentAccessToken(result.ComponentAccessToken, result.ExpiresIn-600)
		if err != nil {
			return "", err
		}

		return result.ComponentAccessToken, nil
	}

	return componentAccessToken, nil
}

func (impl serviceProviderImpl) GetAuthorizerAppId(authCode string) (string, error) {
	if authCode == "" {
		return "", errors.New("empty auth code")
	}

	componentAccessToken, err := impl.getComponentAccessToken()
	if err != nil {
		return "", err
	}

	// 每次查询一次, accessToken可能会发生变化，需要更新缓存
	authorizationInfo, err := impl.wxapi.QueryAuthorizationInfo(componentAccessToken, authCode)
	if err != nil {
		return "", errors.Wrap(err, "query authorization info")
	}

	err = impl.cache.SetAuthorizerAccessToken(authorizationInfo.AuthorizerAppid, authorizationInfo.AuthorizerAccessToken, authorizationInfo.ExpiresIn)
	if err != nil {
		return "", err
	}

	err = impl.cache.SetAuthorizerRefreshToken(authorizationInfo.AuthorizerAppid, authorizationInfo.AuthorizerRefreshToken)
	if err != nil {
		return "", err
	}

	return authorizationInfo.AuthorizerAppid, nil
}

func (impl serviceProviderImpl) GetAuthorizerInfo(appId string) (*provider.GetAuthorizerInfoResult, error) {
	componentAccessToken, err := impl.getComponentAccessToken()
	if err != nil {
		return nil, err
	}

	authorizerInfo, err := impl.wxapi.GetAuthorizerInfo(componentAccessToken, appId)
	if err != nil {
		return nil, errors.Wrap(err, "get authorizer info")
	}

	return authorizerInfo, nil
}

func (impl serviceProviderImpl) GetAuthorizerAccessToken(authorizerAppid string) (string, error) {
	authorizerAccessToken, _ := impl.cache.GetAuthorizerAccessToken(authorizerAppid)
	if authorizerAccessToken == "" {
		componentAccessToken, err := impl.getComponentAccessToken()
		if err != nil {
			return "", err
		}

		authRefreshToken, err := impl.getAuthorizerRefreshToken(authorizerAppid)
		if err != nil {
			return "", err
		}

		result, err := impl.wxapi.GetAuthorizerAccessToken(componentAccessToken, authorizerAppid, authRefreshToken)
		if err != nil {
			return "", err
		}

		// 缓存authorizer access Token
		if err = impl.cache.SetAuthorizerAccessToken(authorizerAppid, result.AuthorizerAccessToken, result.ExpiresIn); err != nil {
			return "", errors.Wrap(err, "cache authorizer access Token")
		}

		if err = impl.cache.SetAuthorizerRefreshToken(authorizerAppid, result.AuthorizerRefreshToken); err != nil {
			return "", errors.Wrap(err, "cache authorizer access Token")
		}

		return result.AuthorizerAccessToken, nil
	}
	return authorizerAccessToken, nil
}

// getAuthorizerRefreshToken 获取保存的authorizerRefreshToken, 先从缓存中找，找不到从调用WX API接口获取
func (impl serviceProviderImpl) getAuthorizerRefreshToken(appId string) (string, error) {
	refreshToken, _ := impl.cache.GetAuthorizerRefreshToken(appId)
	if refreshToken == "" {
		componentAccessToken, err := impl.getComponentAccessToken()
		if err != nil {
			return "", err
		}

		result, err := impl.wxapi.GetAuthorizerInfo(componentAccessToken, appId)
		if err != nil {
			return "", errors.Wrapf(err, "wxapi get authorizer refresh Token, appId: %s", appId)
		}

		refreshToken = result.Authorization.RefreshToken
		if refreshToken == "" {
			return "", errors.New("empty refresh Token in authorization info")
		}

		err = impl.cache.SetAuthorizerRefreshToken(appId, refreshToken)
		if err != nil {
			return "", errors.Wrap(err, "cache authorizer refresh Token")
		}

		return refreshToken, nil
	}
	return refreshToken, nil
}

func (impl serviceProviderImpl) getComponentVerifyTicket() (string, error) {
	componentVerifyTicket, _ := impl.cache.GetComponentVerifyTicket()
	if componentVerifyTicket == "" {
		// 如果缓存里面没有component verify ticket, 尝试重新推送ticket
		if err := impl.wxapi.StartPushComponentVerifyTicket(); err != nil {
			return "", errors.Wrap(err, "start push component verify ticket")
		}

		// 等待3秒后重新获取
		time.Sleep(3 * time.Second)
		componentVerifyTicket, _ = impl.cache.GetComponentVerifyTicket()
	}
	return componentVerifyTicket, nil
}
