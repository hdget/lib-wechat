package wxop

import (
	"encoding/xml"
	"strings"
)

type authEvent struct {
	InfoType string `xml:"InfoType"`
}

type xmlAuthAuthorizedEvent struct {
	AuthorizerAppid              string `xml:"AuthorizerAppid"`
	AuthorizationCode            string `xml:"AuthorizationCode"`
	AuthorizationCodeExpiredTime int    `xml:"AuthorizationCodeExpiredTime"`
}

type xmlAuthUnAuthorizedEvent struct {
	AuthorizerAppid string `xml:"AuthorizerAppid"`
}

// AuthCallback 授权回调函数
type AuthCallback func(appId string, params ...any) error

func (impl wxopImpl) HandleAuthEvent(signature, timestamp, nonce, body string, callbacks map[string]AuthCallback) error {
	data, err := impl.crypt.DecryptMsg(signature, timestamp, nonce, body)
	if err != nil {
		return err
	}

	var event authEvent
	if err = xml.Unmarshal(data, &event); err != nil {
		return err
	}

	switch infoType := strings.ToLower(event.InfoType); infoType {
	case "component_verify_ticket":
		if err = impl.apiProcessComponentVerifyTicket(data); err != nil {
			return err
		}
	case "authorized", "updateauthorized": // 授权成功, 换取authorizer info
		authorizationCode, err := impl.parseAuthorizationCode(data)
		if err != nil {
			return err
		}

		info, err := impl.apiQueryAuthorizationInfo(authorizationCode)
		if err != nil {
			return err
		}

		if cb, exist := callbacks[infoType]; exist {
			return cb(info.AuthorizerAppid, info)
		}
	case "unauthorized": // 取消授权
		unAuthorizationAppId, err := impl.parseUnAuthorizationAppId(data)
		if err != nil {
			return err
		}

		if cb, exist := callbacks[infoType]; exist {
			return cb(unAuthorizationAppId)
		}
	}
	return nil
}

func (impl wxopImpl) parseAuthorizationCode(data []byte) (string, error) {
	var event xmlAuthAuthorizedEvent
	if err := xml.Unmarshal(data, &event); err != nil {
		return "", err
	}
	return event.AuthorizationCode, nil
}

func (impl wxopImpl) parseUnAuthorizationAppId(data []byte) (string, error) {
	var event xmlAuthUnAuthorizedEvent
	if err := xml.Unmarshal(data, &event); err != nil {
		return "", err
	}
	return event.AuthorizerAppid, nil
}
