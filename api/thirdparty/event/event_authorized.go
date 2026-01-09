package event

import (
	"encoding/xml"
)

type xmlAuthAuthorizedEvent struct {
	AuthorizerAppid              string `xml:"AuthorizerAppid"`
	AuthorizationCode            string `xml:"AuthorizationCode"`
	AuthorizationCodeExpiredTime int    `xml:"AuthorizationCodeExpiredTime"`
}

type authorizedEventHandler struct {
}

func newAuthorizedEventHandler() Handler {
	return &authorizedEventHandler{}
}

func (h authorizedEventHandler) Handle(data []byte) (string, error) {
	var e xmlAuthAuthorizedEvent
	if err := xml.Unmarshal(data, &e); err != nil {
		return "", err
	}

	return e.AuthorizationCode, nil
}
