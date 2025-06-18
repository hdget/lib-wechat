package serviceprovider

import (
	"encoding/xml"
)

type xmlAuthAuthorizedEvent struct {
	AuthorizerAppid              string `xml:"AuthorizerAppid"`
	AuthorizationCode            string `xml:"AuthorizationCode"`
	AuthorizationCodeExpiredTime int    `xml:"AuthorizationCodeExpiredTime"`
}

type authorizedEventHandler struct {
	*eventImpl
}

func (impl serviceProviderImpl) newAuthorizedEventProcessor(e *eventImpl) EventProcessor {
	return &authorizedEventHandler{
		eventImpl: e,
	}
}

func (h authorizedEventHandler) Process(data []byte) (string, error) {
	var xmlEvent xmlAuthAuthorizedEvent
	if err := xml.Unmarshal(data, &xmlEvent); err != nil {
		return "", err
	}

	return xmlEvent.AuthorizationCode, nil
}
